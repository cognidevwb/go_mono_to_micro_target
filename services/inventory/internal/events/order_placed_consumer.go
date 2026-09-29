package events

import (
	"context"

	pkgevents "github.com/acme/shop/pkg/events"
	"github.com/acme/shop/services/inventory/internal/contracts/orders"
	"github.com/acme/shop/services/inventory/internal/inventory"
)

// OrderPlacedDurable names the JetStream consumer (letters, digits, - and _ only).
const OrderPlacedDurable = "inventory-service-order_placed"

// SubscribeOrderPlaced consumes `order.placed` from the broker and hands the
// payload to the ported handler. Redeliveries are dropped by message id; the
// broker acks only after the handler returns without error.
func SubscribeOrderPlaced(ctx context.Context, bus pkgevents.Bus, seen pkgevents.Seen, svc *inventory.Service) error {
	h := pkgevents.Idempotent(OrderPlacedDurable, seen, func(_ context.Context, e pkgevents.Envelope) error {
		var payload any
		if err := e.Decode(&payload); err != nil {
			return err
		}
		svc.OnOrderPlaced(payload)
		return nil
	})
	return bus.Subscribe(ctx, orders.TopicOrderPlaced, OrderPlacedDurable, h)
}
