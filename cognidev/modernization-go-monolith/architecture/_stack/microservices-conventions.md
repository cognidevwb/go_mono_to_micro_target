---
id: microservices-conventions
title: Go microservice conventions
family: method
currency: standard
role: reference
applies_when:
  - always
teaches: Every service looks the same when you open it — cmd/, internal/, migrations/, one go.mod — so a reader learns the fleet once.
---

# Go microservice conventions

## What it is

The per-service conventions that make five services one system: the same
directory shape, the same startup order, the same health, config, error and
logging conventions. None of them is clever; all of them are enforced.

## Currency (2026)

**STANDARD.** The `cmd/` + `internal/` layout is the Go project's own guidance
(go.dev/doc/modules/layout); `internal/` is compiler-enforced privacy.

## Fit signals

Always applies.

## What it changes in the generated services

```
services/<ctx>/
  go.mod                         module <base_module>/services/<ctx>
  cmd/<ctx>/main.go              wiring only: config → telemetry → db → router → run
  internal/<pkg>/                the PORTED monolith package, package name kept
  internal/clients/<to>_client.go  one typed client per upstream context
  internal/saga/                 orchestrator only
  internal/outbox/               when the service publishes
  internal/events/               event contracts it publishes/consumes
  migrations/0001_init.sql       owned tables only
  Dockerfile
pkg/contracts/<ctx>/             DTOs other services may import — nothing else
```

Rules the generated tree follows:

1. **Config from the environment** (`DATABASE_URL`, `NATS_URL`, `<CTX>_URL`,
   `OTEL_EXPORTER_OTLP_ENDPOINT`) — one struct, parsed once in `main`.
2. **Errors wrap** with `fmt.Errorf("reserve stock: %w", err)`; HTTP maps them
   to RFC 9457 problem JSON at the edge only.
3. **Every exported call takes `context.Context` first**; no goroutine outlives
   the context `main` created.
4. **`/healthz` (liveness) and `/readyz` (readiness)** on every service.
5. **URL-segment versioning** (`/v1/...`) on every route that leaves the service.
6. **No service imports another service's module.** Only `pkg/contracts`.

## Go 2026 implementation

`go vet`, `staticcheck` and `golangci-lint` (with `depguard` forbidding
`services/*` imports across modules) run in each service's CI.

## Anti-patterns

- A `utils/` or `common/` module every service imports — the monolith again.
- `init()` side effects registering models or routes (the monolith's
  `platform.Register` pattern) — wiring belongs in `main`.
- Package-level mutable state (`var priceCache = map[...]`) — per-replica,
  inconsistent, and invisible to the next reader.

## Interacts with

- [[target-architecture-2026]] · [[flat-package]] · [[hexagonal]] · [[chassis]]
