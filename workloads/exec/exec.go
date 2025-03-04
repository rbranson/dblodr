package exec

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rbranson/dblodr/pkg/gen"
	"github.com/rbranson/dblodr/pkg/workload"
	"github.com/rbranson/dblodr/workloads"
	"github.com/spf13/cobra"
)

var (
	statement string
)

var wl = &workload.Workload{
	Use:  "exec",
	Desc: "Executes a SQL statement",

	Cmd: func(c *cobra.Command) {
		c.Flags().StringVarP(&statement, "sql", "q", "", "SQL statement to execute")
	},

	Run: func(ctx context.Context, inst *gen.Instance, db *sql.DB) (error, bool) {
		_, err := db.ExecContext(ctx, statement)
		if err != nil {
			return fmt.Errorf("in exec: %w", err), false
		}
		return err, false
	},
}

func init() {
	workloads.Register(wl)
}
