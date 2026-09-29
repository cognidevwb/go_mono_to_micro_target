package orders

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"gorm.io/gorm"

	pkgevents "github.com/acme/shop/pkg/events"
	"github.com/acme/shop/pkg/httpx"
	pkgsaga "github.com/acme/shop/pkg/saga"
	"github.com/acme/shop/services/orders/internal/clients/catalog"
	"github.com/acme/shop/services/orders/internal/clients/customers"
	"github.com/acme/shop/services/orders/internal/clients/inventory"
	"github.com/acme/shop/services/orders/internal/clients/payments"
	"github.com/acme/shop/services/orders/internal/events"
	"github.com/acme/shop/services/orders/internal/platform"
	"github.com/acme/shop/services/orders/internal/saga"
)

// TopicOrderPlaced is published after an order commits.
const TopicOrderPlaced = "order.placed"

// Service owns the orders context and orchestrates the others.
type Service struct {
	db        *gorm.DB
	catalog   *catalog.Service
	customers *customers.Service
	inventory *inventory.Service
	payments  *payments.Service
	bus       *platform.Bus
	http      *http.Client
}

// NewService wires the order service.
func NewService(db *gorm.DB, c *catalog.Service, cu *customers.Service,
	i *inventory.Service, p *payments.Service, bus *platform.Bus) *Service {
	return &Service{db: db, catalog: c, customers: cu, inventory: i, payments: p, bus: bus, http: httpx.NewClient()}
}

// PlaceOrder runs the CreateOrder saga: inventory is reserved first, and the
// payment — the hardest step to undo — is charged last. The order and its lines
// are written inside the local transaction that only commits once the charge
// succeeded, so a failed saga leaves no order behind.
func (s *Service) PlaceOrder(req PlaceOrderRequest) (*Order, error) {
	ctx := context.Background()
	active, err := s.customers.IsActive(req.CustomerID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, errors.New("customer inactive")
	}
	order := &Order{CustomerID: req.CustomerID, Status: "pending"}
	sagaID := pkgevents.NewID()
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reserved []LineItem
		release := func(ctx context.Context) error {
			var errs []error
			for _, it := range reserved {
				if err := s.releaseStock(ctx, sagaID, it.ProductID, it.Quantity); err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		}
		flow := saga.CreateOrder(map[string]pkgsaga.Step{
			"inventory": {
				Name: "inventory",
				Do: func(ctx context.Context) error {
					for _, item := range req.Items {
						ok, err := s.inventory.Reserve(tx, item.ProductID, item.Quantity)
						if err != nil {
							_ = release(ctx)
							return err
						}
						if !ok {
							_ = release(ctx)
							return errors.New("out of stock")
						}
						reserved = append(reserved, item)
						price, err := s.catalog.PriceOf(item.ProductID)
						if err != nil {
							_ = release(ctx)
							return err
						}
						order.Lines = append(order.Lines, OrderLine{ProductID: item.ProductID, Quantity: item.Quantity, UnitPrice: price})
						order.Total += price * float64(item.Quantity)
					}
					return nil
				},
				Compensate: release,
			},
			"payments": {
				Name: "payments",
				Do: func(ctx context.Context) error {
					if err := tx.Create(order).Error; err != nil {
						return err
					}
					if err := s.payments.Charge(tx, order.ID, order.Total); err != nil {
						return err
					}
					order.Status = "placed"
					return tx.Save(order).Error
				},
			},
		})
		if _, err := flow.Run(ctx); err != nil {
			return err
		}
		return events.PublishOrderPlaced(ctx, tx.Statement.ConnPool, order.ID)
	})
	if err != nil {
		var se *pkgsaga.StepError
		if errors.As(err, &se) && len(se.Compensation) > 0 {
			slog.Error("saga compensation failed", "saga", se.Saga, "step", se.Step, "err", errors.Join(se.Compensation...))
		}
		return nil, err
	}
	return order, nil
}

// releaseStock undoes one reservation through inventory-service's idempotent
// compensation route (keyed by saga id and product).
func (s *Service) releaseStock(ctx context.Context, sagaID string, productID uint, qty int) error {
	base := os.Getenv("INVENTORY_SERVICE_URL")
	if base == "" {
		base = "http://inventory-service:8080"
	}
	payload, err := json.Marshal(map[string]any{"saga_id": sagaID, "product_id": productID, "qty": qty})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/v1/inventory/create_order/compensate", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("release stock for product %d: status %d", productID, resp.StatusCode)
	}
	return nil
}

// ForCustomer lists a customer's orders.
func (s *Service) ForCustomer(customerID uint) ([]Order, error) {
	var out []Order
	err := s.db.Preload("Lines").Where("customer_id = ?", customerID).Find(&out).Error
	return out, err
}
