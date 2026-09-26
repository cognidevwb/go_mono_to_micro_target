// Command gateway is the edge: it authenticates the caller once, then
// reverse-proxies each path prefix to the service that owns it. Anything no
// service has taken over yet goes to LEGACY_UPSTREAM (the strangler fig).
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

	"github.com/acme/shop/gateway/internal/proxy"
	"github.com/acme/shop/pkg/httpx"
	"github.com/acme/shop/pkg/otelx"
)

func main() {
	if err := run(); err != nil {
		slog.Error("gateway stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "gateway"))

	shutdown, err := otelx.Setup(ctx, "gateway")
	if err != nil {
		return err
	}
	defer func() { _ = shutdown(context.Background()) }()

	handler, err := proxy.New(proxy.FromEnv(proxy.Routes()), os.Getenv("LEGACY_UPSTREAM"))
	if err != nil {
		return err
	}
	auth, err := proxy.Authenticate(ctx, os.Getenv("AUTH_ISSUER"), os.Getenv("AUTH_AUDIENCE"))
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	health := httpx.NewHealth()
	mux.HandleFunc("GET /healthz", health.Live)
	mux.HandleFunc("GET /readyz", health.Ready)
	mux.Handle("/", otelhttp.NewHandler(auth(handler), "gateway"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return httpx.Run(ctx, &http.Server{Addr: ":" + port, Handler: mux, ReadHeaderTimeout: 5 * time.Second})
}
