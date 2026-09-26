// Package events defines the payloads orders-service publishes through the
// outbox. order.placed is announced once PlaceOrder's saga commits.
package events

// OrderPlaced is orders-service's payload for the order.placed event.
type OrderPlaced struct {
	OrderID    uint    `json:"orderId"`
	CustomerID uint    `json:"customerId"`
	Total      float64 `json:"total"`
}
