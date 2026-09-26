// Command orders is the orders-service entrypoint: config from the
// environment, JSON logs, OpenTelemetry, liveness/readiness, the outbox relay
// and graceful shutdown. The service's own routes are mounted by app.Routes.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/acme/shop/pkg/httpx"
	"github.com/acme/shop/pkg/otelx"
	"github.com/acme/shop/pkg/outbox"
	"github.com/acme/shop/services/orders/internal/app"
	"github.com/acme/shop/services/orders/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("orders-service stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", cfg.Service))

	shutdown, err := otelx.Setup(ctx, cfg.Service)
	if err != nil {
		return err
	}
	defer func() { _ = shutdown(context.Background()) }()

	deps, err := app.New(ctx, cfg)
	if err != nil {
		return err
	}
	defer deps.Close()

	handler, err := app.Routes(deps)
	if err != nil {
		return err
	}
	go outbox.Relay{DB: deps.SQL, Bus: deps.Bus}.Run(ctx)

	health := httpx.NewHealth()
	health.AddReadiness("database", deps.Ready)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", health.Live)
	mux.HandleFunc("GET /readyz", health.Ready)
	mux.Handle("/", otelhttp.NewHandler(handler, cfg.Service))

	return httpx.Run(ctx, &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	})
}
