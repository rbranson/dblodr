package main

import (
	"fmt"
	"os"

	"github.com/rbranson/dblodr/cmd"

	_ "github.com/rbranson/dblodr/workloads/exec"
	_ "github.com/rbranson/dblodr/workloads/openconns"
	_ "github.com/rbranson/dblodr/workloads/simpleinsert1"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
