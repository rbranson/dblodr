package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rbranson/dblodr/pkg/gen"
	"github.com/rbranson/dblodr/pkg/workload"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute load against the target database",
	Run: func(cmd *cobra.Command, args []string) {
		exitErrorWithUsage(cmd, errors.New("a subcommand is required"))
	},
}

var (
	concurrency   int
	iterTimeout   time.Duration
	statFrequency time.Duration
	iterations    int64
)

func addDatabaseFlags(c *cobra.Command) {
	c.PersistentFlags().StringVarP(&dsn, "dsn", "c", "", "DSN or Connection String")
	c.PersistentFlags().StringVarP(&driverName, "driver", "d", "", "Driver Name (e.g. 'mysql')")
	c.MarkFlagRequired("dsn")
	c.MarkFlagRequired("driver")
}

func init() {
	addDatabaseFlags(runCmd)
	runCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "t", 1, "Concurrency")
	runCmd.PersistentFlags().DurationVar(&iterTimeout, "iter-timeout", 0, "Timeout for an individual iteration of the run (0 is no timeout)")
	runCmd.PersistentFlags().DurationVar(&statFrequency, "stat-frequency", 1*time.Second, "Frequency to print run stats")
	runCmd.PersistentFlags().Int64VarP(&iterations, "iterations", "n", 0, "Iterations to run (0 is unlimited)")
}

func Run(w *workload.Workload) {
	db := dbFromFlags()
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewDevelopment(zap.WithCaller(false))
	defer logger.Sync()

	logger.Info("Starting Run", zap.String("workload", w.Use))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		logger.Info("Received an interrupt, stopping all workers...")
		cancel()
	}()

	if w.Prerun != nil {
		err := w.Prerun(ctx, db)
		if err != nil {
			exitError(fmt.Errorf("error in prerun: %w", err))
		}
	}

	gstats := &gen.Stats{}
	var wg sync.WaitGroup
	var invokes atomic.Int64

	// Initialize the stats so they show up initially
	gstats.ResetCounter("iterations")
	gstats.ResetCounter("errors")

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		stats := &gen.Stats{Parent: gstats}
		inst := &gen.Instance{
			ID:     i,
			DB:     db,
			Ctx:    ctx,
			Logger: logger,
			Stats:  stats,
		}

		go func() {
			defer wg.Done()
			worker(inst, w, &invokes)
		}()
	}

	statsCh := make(chan struct{})
	go emitStats(ctx, logger, gstats, statsCh)

	wg.Wait()
	cancel()
	<-statsCh
}

func worker(inst *gen.Instance, w *workload.Workload, invokes *atomic.Int64) {
	for {
		if inst.Ctx.Err() != nil {
			return
		}
		if iterations > 0 && invokes.Add(1) > iterations {
			return
		}

		inst.Invokes.Add(1)

		ctx := inst.Ctx
		cancel := func() {}
		if iterTimeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, iterTimeout)
		}
		defer cancel()

		start := time.Now()
		err, isFatal := w.Run(ctx, inst, inst.DB)
		dur := time.Since(start)

		inst.Stats.CounterIncr("iterations", 1)
		inst.Stats.CounterIncrDur("runtime_ns", dur, time.Nanosecond)

		if err != nil {
			if isFatal {
				exitError(fmt.Errorf("fatal error during run: %w", err))
			}
			if errors.Is(err, context.DeadlineExceeded) {
				inst.Stats.CounterIncr("timeouts", 1)
				continue
			}
			if errors.Is(err, context.Canceled) {
				return
			}
			inst.Stats.CounterIncr("errors", 1)
			inst.Logger.Warn("error in run",
				zap.Int("instance", inst.ID),
				zap.Int64("invocation", inst.Invokes.Load()),
				zap.Error(err))
		}
	}
}

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%.0fns", float64(d)/float64(time.Nanosecond))
	case d < time.Millisecond:
		return fmt.Sprintf("%.3fµs", float64(d)/float64(time.Microsecond))
	case d < time.Second:
		return fmt.Sprintf("%.3fms", float64(d)/float64(time.Millisecond))
	case d < time.Minute:
		return fmt.Sprintf("%.3fs", float64(d)/float64(time.Second))
	case d < time.Hour:
		return fmt.Sprintf("%.3fm", float64(d)/float64(time.Minute))
	default:
		return fmt.Sprintf("%.3fh", float64(d)/float64(time.Hour))
	}
}

type statsFormatter struct {
	lastInvokes int64
	lastDur     int64
}

func (f *statsFormatter) getOutput(counters map[string]int64) map[string]string {
	out := map[string]string{}
	for k, v := range counters {
		out[k] = fmt.Sprintf("%d", v)
	}

	invokes, hasInvokes := counters["iterations"]
	dur, hasDur := counters["runtime_ns"]

	if hasDur {
		delete(out, "runtime_ns")
	}

	if hasInvokes && hasDur {
		thisInvokes := invokes - f.lastInvokes
		if thisInvokes > 0 {
			thisDur := time.Duration(dur-f.lastDur) * time.Nanosecond
			avg := thisDur / time.Duration(thisInvokes)
			out["avg_runtime"] = formatDuration(avg)
		} else {
			out["avg_runtime"] = "NaN"
		}
	}

	if hasInvokes {
		f.lastInvokes = invokes
	}
	if hasDur {
		f.lastDur = dur
	}
	return out

}

func logStats(logger *zap.Logger, counters map[string]string) {
	keys := make([]string, 0, len(counters))
	for k := range counters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fields := []zap.Field{}
	for _, k := range keys {
		fields = append(fields, zap.String(k, counters[k]))
	}
	logger.Info("Run Stats", fields...)
}

func emitStats(ctx context.Context, logger *zap.Logger, stats *gen.Stats, done chan struct{}) {
	defer close(done)

	var fmtr statsFormatter
	t := time.NewTicker(statFrequency)
	defer t.Stop()

	emit := func() {
		counters := fmtr.getOutput(stats.ReadCounters())
		logStats(logger, counters)
	}

	for {
		select {
		case <-t.C:
			emit()
		case <-ctx.Done():
			emit()
			return
		}
	}
}
