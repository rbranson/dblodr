package openconns

import (
	"context"
	"database/sql"
	"time"

	"github.com/rbranson/dblodr/pkg/gen"
	"github.com/rbranson/dblodr/pkg/workload"
	"github.com/rbranson/dblodr/workloads"
)

var wl = &workload.Workload{
	Use:  "openconns",
	Desc: "Open a bunch of connections that do nothing",
	Run: func(ctx context.Context, inst *gen.Instance, db *sql.DB) (error, bool) {
		conn, err := db.Conn(ctx)
		if err != nil {
			return err, false
		}
		defer conn.Close()
		if err := gen.InterruptibleSleep(ctx, 24*time.Hour); err != nil {
			return err, false
		}
		return nil, false
	},
}

func init() {
	workloads.Register(wl)
}
