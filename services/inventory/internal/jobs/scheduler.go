// Package jobs runs inventory's scheduled work once across replicas, under a
// Postgres advisory-lock lease.
package jobs

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

// leaseKey is the advisory lock the restock job holds while it runs.
const leaseKey int64 = 0x1e5701c4

// Run calls work every interval until ctx ends. Each tick only the replica that
// wins the advisory lock runs it; the others skip that tick.
func Run(ctx context.Context, db *sql.DB, interval time.Duration, work func() error) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			runLeased(ctx, db, work)
		}
	}
}

func runLeased(ctx context.Context, db *sql.DB, work func() error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		slog.WarnContext(ctx, "job lease connection", "err", err)
		return
	}
	defer conn.Close()
	var got bool
	if err := conn.QueryRowContext(ctx, `SELECT pg_try_advisory_lock($1)`, leaseKey).Scan(&got); err != nil || !got {
		return
	}
	defer func() { _, _ = conn.ExecContext(context.WithoutCancel(ctx), `SELECT pg_advisory_unlock($1)`, leaseKey) }()
	if err := work(); err != nil {
		slog.WarnContext(ctx, "restock failed", "err", err)
	}
}
