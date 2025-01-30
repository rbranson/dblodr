package workloads

import (
	"github.com/rbranson/dblodr/cmd"
	"github.com/rbranson/dblodr/pkg/workload"
	"github.com/spf13/cobra"
)

var Workloads = &WorkloadSet{}

type WorkloadSet struct {
	workloads []*workload.Workload
}

// Register adds a Workload w to the set
func (s *WorkloadSet) Register(w *workload.Workload) {
	s.workloads = append(s.workloads, w)
}

// Registers a global workload
func Register(w *workload.Workload) {
	Workloads.Register(w)

	c := &cobra.Command{
		Use:   w.Use,
		Short: w.Desc,
		Run: func(c *cobra.Command, args []string) {
			switch c.Parent().Use {
			case "init":
				cmd.Initialize(w)
			case "run":
				cmd.Run(w)
			default:
				panic("unknown parent use, inconceivable")
			}
		},
	}
	w.Cmd(c)

	cmd.AddWorkloadCommand(c)
}
