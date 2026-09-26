package orders

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"

	pkgsaga "github.com/acme/shop/pkg/saga"
	"github.com/acme/shop/services/orders/internal/clients/catalog"
	"github.com/acme/shop/services/orders/internal/clients/customers"
	"github.com/acme/shop/services/orders/internal/clients/inventory"
	"github.com/acme/shop/services/orders/internal/clients/payments"
	"github.com/acme/shop/services/orders/internal/events"
	"github.com/acme/shop/services/orders/internal/outbox"
	"github.com/acme/shop/services/orders/internal/platform"
	ordersaga "github.com/acme/shop/services/orders/internal/saga"
)

// TopicOrderPlaced is published after an order commits.
const TopicOrderPlaced = "order.placed"

// Status is an order's typed lifecycle state.
type Status string

const (
	StatusPending Status = "pending"
	StatusPlaced  Status = "placed"
	StatusFailed  Status = "failed"
)

// Service owns the orders context and orchestrates the others.
type Service struct {
	db        *gorm.DB
	catalog   *catalog.Service
	customers *customers.Service
	inventory *inventory.Service
	payments  *payments.Service
	bus       *platform.Bus
}

// NewService wires the order service.
func NewService(db *gorm.DB, c *catalog.Service, cu *customers.Service,
	i *inventory.Service, p *payments.Service, bus *platform.Bus) *Service {
	return &Service{db: db, catalog: c, customers: cu, inventory: i, payments: p, bus: bus}
}

// PlaceOrder prices and records the order in orders' own database, then runs
// the CreateOrder saga to reserve stock and charge payment: once each
// context owns its own tables, one database transaction can no longer span
// all four. Stock is reserved before payment is charged, since a charge is
// the harder step to undo, and a failed step compensates the completed ones
// in reverse.
func (s *Service) PlaceOrder(req PlaceOrderRequest) (*Order, error) {
	ctx := context.Background()

	active, err := s.customers.IsActive(req.CustomerID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, errors.New("customer inactive")
	}

	order := &Order{CustomerID: req.CustomerID, Status: string(StatusPending)}
	for _, item := range req.Items {
		price, err := s.catalog.PriceOf(item.ProductID)
		if err != nil {
			return nil, err
		}
		order.Lines = append(order.Lines, OrderLine{ProductID: item.ProductID, Quantity: item.Quantity, UnitPrice: price})
		order.Total += price * float64(item.Quantity)
	}
	if err := s.db.WithContext(ctx).Create(order).Error; err != nil {
		return nil, err
	}

	items := req.Items
	reserved := make(map[uint]int, len(items))
	sg := ordersaga.CreateOrder(map[string]pkgsaga.Step{
		"inventory": {
			Name: "inventory",
			Do: func(ctx context.Context) error {
				for _, item := range items {
					ok, err := s.inventory.Reserve(s.db, item.ProductID, item.Quantity)
					if err != nil {
						return err
					}
					if !ok {
						return errors.New("out of stock")
					}
					reserved[item.ProductID] += item.Quantity
				}
				return nil
			},
			Compensate: func(ctx context.Context) error {
				// inventory-service exposes no release endpoint yet; log so
				// an operator can reconcile the reservation by hand.
				slog.ErrorContext(ctx, "order saga: stock reservation could not be released automatically",
					"order", order.ID, "reserved", reserved)
				return nil
			},
		},
		"payments": {
			Name: "payments",
			Do: func(ctx context.Context) error {
				return s.payments.Charge(s.db, order.ID, order.Total)
			},
			Compensate: func(ctx context.Context) error {
				// payments-service exposes no refund endpoint yet; log so an
				// operator can reconcile the charge by hand.
				slog.ErrorContext(ctx, "order saga: charge could not be refunded automatically",
					"order", order.ID, "total", order.Total)
				return nil
			},
		},
	})

	if _, err := sg.Run(ctx); err != nil {
		order.Status = string(StatusFailed)
		if saveErr := s.db.WithContext(ctx).Save(order).Error; saveErr != nil {
			slog.ErrorContext(ctx, "order saga: failed to persist failed status", "order", order.ID, "err", saveErr)
		}
		return nil, err
	}

	order.Status = string(StatusPlaced)
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(order).Error; err != nil {
			return err
		}
		return events.Publish(ctx, outbox.TxExecer{Tx: tx}, TopicOrderPlaced, events.OrderPlaced{
			OrderID:    order.ID,
			CustomerID: order.CustomerID,
			Total:      order.Total,
		})
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

// ForCustomer lists a customer's orders.
func (s *Service) ForCustomer(customerID uint) ([]Order, error) {
	var out []Order
	err := s.db.Preload("Lines").Where("customer_id = ?", customerID).Find(&out).Error
	return out, err
}
