// Package outbox is the transactional outbox: an event is written in the SAME
// database transaction as the state change that caused it, and a relay
// publishes it afterwards. No dual write, no lost event on a crash between the
// commit and the publish.
package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/acme/shop/pkg/events"
)

// Schema is the DDL each service's migrations carry.
const Schema = `CREATE TABLE IF NOT EXISTS outbox (
    id           BIGSERIAL PRIMARY KEY,
    subject      TEXT        NOT NULL,
    envelope     JSONB       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS processed_messages (
    consumer   TEXT        NOT NULL,
    message_id TEXT        NOT NULL,
    seen_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (consumer, message_id)
);`

// Execer is satisfied by *sql.Tx and *sql.DB.
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// Enqueue writes e inside tx. Call it in the transaction that changes state.
func Enqueue(ctx context.Context, tx Execer, subject string, e events.Envelope) error {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshal envelope: %w", err)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO outbox (subject, envelope) VALUES ($1, $2)`, subject, body)
	return err
}

// Relay publishes unpublished rows in id order.
type Relay struct {
	DB       *sql.DB
	Bus      events.Bus
	Interval time.Duration
	Batch    int
}

// Run polls until ctx ends.
func (r Relay) Run(ctx context.Context) {
	interval := r.Interval
	if interval <= 0 {
		interval = time.Second
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if n, err := r.Flush(ctx); err != nil {
				slog.WarnContext(ctx, "outbox relay", "err", err, "published", n)
			}
		}
	}
}

// Flush publishes one batch and returns how many rows it published.
func (r Relay) Flush(ctx context.Context) (int, error) {
	batch := r.Batch
	if batch <= 0 {
		batch = 100
	}
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, `SELECT id, subject, envelope FROM outbox
        WHERE published_at IS NULL ORDER BY id LIMIT $1 FOR UPDATE SKIP LOCKED`, batch)
	if err != nil {
		return 0, err
	}
	type row struct {
		id      int64
		subject string
		env     events.Envelope
	}
	var pending []row
	for rows.Next() {
		var rw row
		var body []byte
		if err := rows.Scan(&rw.id, &rw.subject, &body); err != nil {
			_ = rows.Close()
			return 0, err
		}
		if err := json.Unmarshal(body, &rw.env); err != nil {
			_ = rows.Close()
			return 0, err
		}
		pending = append(pending, rw)
	}
	if err := rows.Close(); err != nil {
		return 0, err
	}
	done := 0
	for _, rw := range pending {
		if err := r.Bus.Publish(ctx, rw.subject, rw.env); err != nil {
			break
		}
		if _, err := tx.ExecContext(ctx, `UPDATE outbox SET published_at = now() WHERE id = $1`, rw.id); err != nil {
			return done, err
		}
		done++
	}
	return done, tx.Commit()
}
