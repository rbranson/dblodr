package cmd

import (
	"errors"

	"github.com/rbranson/dblodr/pkg/workload"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the target database for a load profile",
	Run: func(c *cobra.Command, args []string) {
		exitErrorWithUsage(c, errors.New("a subcommand is required"))
	},
}

func init() {
	addDatabaseFlags(initCmd)
}

func Initialize(w *workload.Workload) {

}
