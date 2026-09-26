// Package events is how the payments-service publishes: through the
// transactional outbox, in the same transaction as the state change.
package events

import (
	"context"

	pkgevents "github.com/acme/shop/pkg/events"
	"github.com/acme/shop/pkg/outbox"
)

// Source names this service on every envelope it emits.
const Source = "payments-service"

// Publish enqueues payload as subject inside tx; the outbox relay delivers it.
func Publish(ctx context.Context, tx outbox.Execer, subject string, payload any) error {
	e, err := pkgevents.New(subject, Source, payload)
	if err != nil {
		return err
	}
	return outbox.Enqueue(ctx, tx, subject, e)
}
