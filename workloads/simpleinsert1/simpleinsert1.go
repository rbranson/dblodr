package simpleinsert1

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/rbranson/dblodr/pkg/gen"
	"github.com/rbranson/dblodr/pkg/workload"
	"github.com/rbranson/dblodr/workloads"
	"github.com/spf13/cobra"
)

var (
	columns         string
	dataLen         int
	preCommitDelay  time.Duration
	preCommitJitter float64
	tableName       string
)

func randStringArray(len int, size int) []string {
	strs := make([]string, len)
	for i := 0; i < len; i++ {
		strs[i] = gen.RandAlphanum(size)
	}
	return strs
}

func parseColumns(s string) []string {
	return strings.Split(s, ",")
}

var wl = &workload.Workload{
	Use:  "simpleinsert1",
	Desc: "Insert into a simple table",

	Cmd: func(c *cobra.Command) {
		c.Flags().StringVar(&columns, "columns", "data1,data2", "Comma-separated list of varchar columns to insert into (e.g. 'data1,data2')")
		c.Flags().IntVar(&dataLen, "data-len", 20, "Length of the random data to insert")
		c.Flags().DurationVar(&preCommitDelay, "precommit-delay", 0*time.Second, "Delay between the insert and commit")
		c.Flags().Float64Var(&preCommitJitter, "precommit-jitter", 0.0, "Jitter for the precommit delay")
		c.Flags().StringVar(&tableName, "table-name", "tbl1", "Table to insert into")
	},

	Init: func(ctx context.Context, db *sql.DB) error {
		return nil
	},

	Run: func(ctx context.Context, inst *gen.Instance, db *sql.DB) (error, bool) {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("in begin tx: %w", err), false
		}
		defer tx.Rollback()

		cols := parseColumns(columns)
		columnList := strings.Join(cols, ",")
		data := randStringArray(len(cols), dataLen)
		args := make([]any, len(data))
		for i, v := range data {
			args[i] = v
		}

		q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnList, strings.Repeat("?,", len(cols)-1)+"?")
		_, err = tx.ExecContext(ctx, q, args...)
		if err != nil {
			return fmt.Errorf("in exec for query: %q, error: %w", q, err), false
		}

		if err := gen.InterruptibleSleep(ctx, gen.JitterDuration(preCommitDelay, preCommitJitter)); err != nil {
			return err, false
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("in commit: %w", err), false
		}

		return nil, false
	},
}

func init() {
	workloads.Register(wl)
}
