---
id: chassis
title: Microservice chassis
family: observability
currency: standard
role: invariant
applies_when:
  - always
teaches: One small shared startup path — config, slog, OTel, health, graceful shutdown — so every main.go is the same ten lines.
---

# Microservice chassis

## What it is

The cross-cutting startup every service needs, written once. In Go it is deliberately small: a `pkg/chassis` package (or generated per-service `internal/platform`) providing config parsing, the slog logger, OpenTelemetry setup, health handlers and a signal-aware run loop.

## Currency (2026)

**STANDARD — invariant.** microservices.io Microservice Chassis.

## Fit signals (from `.cognidev/feature/context.md`)

- Always.

## What it changes in the generated services

- Each `cmd/<ctx>/main.go` calls the chassis; nothing domain-specific lives in it.

## Go 2026 implementation

```go
func Run(ctx context.Context, name string, h http.Handler) error {
    ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM); defer stop()
    shutdown, err := telemetry.Setup(ctx, name); if err != nil { return err }; defer shutdown(context.Background())
    srv := &http.Server{Addr: ":8080", Handler: otelhttp.NewHandler(h, name), ReadHeaderTimeout: 5 * time.Second}
    go func() { <-ctx.Done(); _ = srv.Shutdown(context.Background()) }()
    if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) { return err }
    return nil
}
```

## Anti-patterns

- A chassis that grows domain helpers.
- Frameworks (go-kit, go-micro) pulled in for what is 100 lines.

## Interacts with

- [[observability]] · [[health-checks]] · [[secrets-config]]
