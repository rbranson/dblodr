package workload

import (
	"context"
	"database/sql"

	"github.com/rbranson/dblodr/pkg/gen"
	"github.com/spf13/cobra"
)

type Workload struct {
	// Use is a unique short name for the workload (i.e. 'myinsertwork1')
	Use string

	// Desc is a description for the workload (i.e. 'Generates insert workload')
	Desc string

	// Init is called during init command
	Init func(context.Context, *sql.DB) error

	// Run is called to execute the workload. Return fatal as "true" if the
	// run should exit and not retry.
	Run func(context.Context, *gen.Instance, *sql.DB) (err error, fatal bool)

	// Prerun is an optional function that is called once during a run before the
	// actual workload is ran.
	Prerun func(context.Context, *sql.DB) error

	// Cmd is called at register to provide the command line for the workload,
	// which is useful for providing extra flags if needed.
	Cmd func(c *cobra.Command)
}
