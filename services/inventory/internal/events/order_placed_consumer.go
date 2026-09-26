// services · consume `order.placed` (raised by orders) over the broker: a
// durable subscription per consumer, an idempotent handler keyed on the
// message id (processed_messages, seen via pkg/events.Idempotent) — replaces
// the in-process bus Subscribe (inventory-service)
package events

import (
	"context"

	pkgevents "github.com/acme/shop/pkg/events"
	"github.com/acme/shop/services/inventory/internal/contracts/orders"
	"github.com/acme/shop/services/inventory/internal/inventory"
)

// DurableOrderPlacedConsumer names inventory's durable subscription to
// order.placed.
const DurableOrderPlacedConsumer = "inventory-service.order-placed"

// SubscribeOrderPlaced durably subscribes to order.placed and hands each
// fresh delivery to svc.OnOrderPlaced, dropping redeliveries the consumer
// already handled.
func SubscribeOrderPlaced(ctx context.Context, bus pkgevents.Bus, seen pkgevents.Seen, svc *inventory.Service) error {
	handler := pkgevents.Idempotent(DurableOrderPlacedConsumer, seen, func(_ context.Context, e pkgevents.Envelope) error {
		var payload any
		if err := e.Decode(&payload); err != nil {
			return err
		}
		svc.OnOrderPlaced(payload)
		return nil
	})
	return bus.Subscribe(ctx, orders.TopicOrderPlaced, DurableOrderPlacedConsumer, handler)
}
