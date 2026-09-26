// Package app composes the orders-service: its own database, the event bus,
// and the routes. Deps is what the ported code is wired against.
package app

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/acme/shop/pkg/events"
	"github.com/acme/shop/pkg/outbox"
	"github.com/acme/shop/services/orders/internal/config"
)

// Deps are the service's shared dependencies. DB is THIS service's database —
// no other service's tables are reachable through it.
type Deps struct {
	Config config.Config
	DB     *gorm.DB
	SQL    *sql.DB
	Bus    events.Bus
}

// New opens the service database, migrates the service's own models and the
// outbox tables, and connects the event bus.
func New(ctx context.Context, cfg config.Config) (*Deps, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database handle: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxConns)
	if _, err := sqlDB.ExecContext(ctx, outbox.Schema); err != nil {
		return nil, fmt.Errorf("outbox schema: %w", err)
	}
	if ms := models(); len(ms) > 0 {
		if err := db.WithContext(ctx).AutoMigrate(ms...); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}
	bus, err := events.Connect(cfg.BrokerURL, cfg.Service)
	if err != nil {
		return nil, err
	}
	return &Deps{Config: cfg, DB: db, SQL: sqlDB, Bus: bus}, nil
}

// Ready reports whether the database answers.
func (d *Deps) Ready(ctx context.Context) error { return d.SQL.PingContext(ctx) }

// Close releases the bus and the database.
func (d *Deps) Close() {
	if d.Bus != nil {
		_ = d.Bus.Close()
	}
	if d.SQL != nil {
		_ = d.SQL.Close()
	}
}
