package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dblodr",
	Short: "dblodr is a simple SQL database load generator",
}

func Execute() error {
	return rootCmd.Execute()
}

var (
	dsn        string
	driverName string
)

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
}

func AddWorkloadCommand(c *cobra.Command) {
	initCmd.AddCommand(c)
	runCmd.AddCommand(c)
}

func exitErrorWithUsage(c *cobra.Command, err error) {
	fmt.Printf("Fatal: %s\n", err)
	_ = c.Usage()
	os.Exit(1)
}

func exitError(err error) {
	fmt.Printf("Fatal: %s\n", err)
	os.Exit(2)
}

func dbFromFlags() *sql.DB {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		exitError(fmt.Errorf("opening database: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		exitError(fmt.Errorf("pinging database: %w", err))
	}

	return db
}
