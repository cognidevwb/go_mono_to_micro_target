package events

import (
	"context"

	"github.com/acme/shop/pkg/outbox"
)

// SubjectOrderPlaced is the topic an order announces itself on.
const SubjectOrderPlaced = "order.placed"

// OrderPlaced is the versioned payload of the order.placed event.
type OrderPlaced struct {
	Version int  `json:"version"`
	OrderID uint `json:"orderId"`
}

// PublishOrderPlaced writes order.placed to the outbox inside tx, so it is
// relayed only if the order commits.
func PublishOrderPlaced(ctx context.Context, tx outbox.Execer, orderID uint) error {
	return Publish(ctx, tx, SubjectOrderPlaced, OrderPlaced{Version: 1, OrderID: orderID})
}
