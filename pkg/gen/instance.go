package gen

import (
	"context"
	"database/sql"
	"sync/atomic"

	"go.uber.org/zap"
)

// Instance is a single instance, or unit of concurrency, for a load run.
type Instance struct {
	// ID uniquely identifies the instance
	ID int

	// DB is the database for this instance
	DB *sql.DB

	// Ctx is the context specific to this instance
	Ctx context.Context

	// Logger is the logger to use for this instance
	Logger *zap.Logger

	// Stats tracks the stats of the instance
	Stats *Stats

	// Invokes is the number of times the instance has been invoked
	Invokes atomic.Int64
}
