---
id: target-architecture-2026
title: The 2026 Go target architecture
family: method
currency: standard
role: reference
applies_when:
  - always
teaches: One pinned stack, stated once — Go 1.26, a module per service in one go.work, slog + OpenTelemetry, NATS JetStream, distroless nonroot images.
---

# The 2026 Go target architecture

## What it is

The single, version-pinned stack every extracted service is generated on. A
decomposition that lets each service pick its own router, logger and driver
ships five dialects of one system; this card is the one dialect. Every choice
below is the default of a questionnaire answer — change the answer, not the file.

## Currency (2026)

**STANDARD.** Go 1.26 is the current release (Go supports the two newest
minors); `log/slog` (1.21), per-iteration loop variables (1.22), range-over-func
(1.23) and the `toolchain` directive (1.21) are all baseline. Source: go.dev/doc/devel/release.

## Fit signals (from `.cognidev/feature/context.md`)

Always applies — it is the substrate the other patterns are expressed in.

## What it changes in the generated services

| Concern | Pinned default | Questionnaire key |
| --- | --- | --- |
| Toolchain | Go 1.26 in `go.work`, every `go.mod`, `golang:1.26` build image | `go_version` |
| Layout | `go.work` → `services/<ctx>` (one module each), `pkg/` shared kernel, `gateway/` | — |
| HTTP | keep the monolith's router (gin / echo / chi) or chi v5 | `web_framework` |
| Data | keep gorm, or pgx v5 + sqlc; Postgres 17, database per service | `data_access`, `database` |
| Events | NATS JetStream (nats.go `jetstream`) · Kafka (franz-go) · none | `event_backbone` |
| Outbox | an `outbox` table per service + relay goroutine | `outbox` |
| Logging | `log/slog` JSON handler, trace ids injected | — |
| Telemetry | OpenTelemetry Go SDK, OTLP → collector | `observability` |
| Auth | JWT (golang-jwt/jwt v5 / coreos go-oidc) at gateway + per service | `auth` |
| Gateway | Go `httputil.ReverseProxy` gateway, or Envoy / Traefik | `gateway` |
| Images | multi-stage → `gcr.io/distroless/static-debian12:nonroot`, `CGO_ENABLED=0` | — |
| Deploy | docker compose (local) · Kubernetes + Helm · Azure Container Apps | `orchestrator`, `cloud` |
| CI | GitHub Actions / Azure DevOps: vet, test -race, golangci-lint, govulncheck, Syft SBOM + Grype | `ci` |

## Go 2026 implementation

```go
// services/orders/cmd/orders/main.go (shape)
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
shutdownOTel, err := telemetry.Setup(ctx, "orders-service")
pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
srv := &http.Server{Addr: ":8080", Handler: router(pool), ReadHeaderTimeout: 5 * time.Second}
g, gctx := errgroup.WithContext(ctx)
g.Go(srv.ListenAndServe)
g.Go(func() error { <-gctx.Done(); return srv.Shutdown(context.Background()) })
```

## Anti-patterns

- **A different stack per service** "because it is independent now". Independent
  deployability is not a licence for five loggers and three routers.
- **Upgrading Go inside the decomposition** without saying so — see the
  `loop-variable-semantics-changed-at-go-1-22` knowledge card.
- **Alpine + CGO by accident** — a `-race` or sqlite dependency silently needs
  cgo; the distroless static image then fails at exec.

## Interacts with

- [[microservices-conventions]] — the per-service conventions on this stack.
- [[decomposition-method]] — how the map decides what lands on it.
- [[chassis]] — the shared startup code every `main.go` calls.
