// Package outbox re-exports the shared transactional outbox for orders: rows are
// written in the same transaction as the change (see internal/events) and the
// shared relay publishes them, marking them sent only after the broker accepts.
package outbox

import (
	"context"
	"database/sql"

	pkgevents "github.com/acme/shop/pkg/events"
	pkgoutbox "github.com/acme/shop/pkg/outbox"
)

// Schema is the DDL for the outbox and processed_messages tables.
const Schema = pkgoutbox.Schema

// Relay publishes unsent outbox rows until ctx ends.
func Relay(ctx context.Context, db *sql.DB, bus pkgevents.Bus) {
	pkgoutbox.Relay{DB: db, Bus: bus}.Run(ctx)
}
