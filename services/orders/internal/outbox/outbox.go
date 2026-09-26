// Package outbox adapts orders' own *gorm.DB transaction to pkg/outbox's
// Execer, so an outgoing event can be enqueued in the SAME local transaction
// as the business row that caused it. Delivery is pkg/outbox.Relay, started
// from cmd/orders/main.go against the table this Execer writes into.
package outbox

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// TxExecer satisfies pkg/outbox.Execer for a *gorm.DB held mid-transaction.
type TxExecer struct{ Tx *gorm.DB }

// ExecContext runs query against the wrapped transaction.
func (e TxExecer) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	res := e.Tx.WithContext(ctx).Exec(query, args...)
	return execResult(res.RowsAffected), res.Error
}

// execResult reports the rows an outbox insert affected; gorm's raw Exec
// gives no last-insert-id, which pkg/outbox.Enqueue never asks for.
type execResult int64

func (r execResult) LastInsertId() (int64, error) { return 0, nil }
func (r execResult) RowsAffected() (int64, error) { return int64(r), nil }
