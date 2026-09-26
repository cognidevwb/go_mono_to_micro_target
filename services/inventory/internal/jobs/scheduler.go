// services · the scheduled work inventory inherits (StartRestockJob): run it
// exactly once across replicas via a Postgres advisory-lock lease, started
// with the service's own context so ctx.Done() (SIGTERM) stops it
// (inventory-service)
package jobs

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

// restockLockKey serializes the restock job across replicas; any fixed
// bigint works as long as it is unique to this job.
const restockLockKey = 727100001

// Task is one unit of scheduled work.
type Task func() error

// Run executes task every interval, holding a Postgres advisory lock so only
// one replica runs it at a time, until ctx is done.
func Run(ctx context.Context, db *sql.DB, interval time.Duration, task Task) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				runLeased(ctx, db, task)
			}
		}
	}()
}

// runLeased runs task only if this replica acquires the advisory lock.
func runLeased(ctx context.Context, db *sql.DB, task Task) {
	conn, err := db.Conn(ctx)
	if err != nil {
		slog.WarnContext(ctx, "job: acquire connection", "err", err)
		return
	}
	defer conn.Close()

	var acquired bool
	if err := conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", restockLockKey).Scan(&acquired); err != nil {
		slog.WarnContext(ctx, "job: advisory lock", "err", err)
		return
	}
	if !acquired {
		return
	}
	defer func() {
		if _, err := conn.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", restockLockKey); err != nil {
			slog.WarnContext(ctx, "job: advisory unlock", "err", err)
		}
	}()

	if err := task(); err != nil {
		slog.WarnContext(ctx, "job failed", "err", err)
	}
}
