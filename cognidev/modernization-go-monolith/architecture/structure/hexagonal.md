---
id: hexagonal
title: Hexagonal (ports and adapters)
family: structure
currency: standard
role: choice
applies_when:
  - answer:service_internal_structure = hexagonal
teaches: Domain package imports nothing infrastructural; adapters implement its ports — worth it only for a large, long-lived, infrastructure-heavy service.
---

# Hexagonal (ports and adapters)

## What it is

Ports and adapters in Go: `internal/domain` holds entities and the interfaces
(ports) the domain needs — `OrderRepository`, `PaymentGateway`; `internal/app`
holds use cases; `internal/adapters/{postgres,http,nats}` implement the ports.
The dependency rule points inward and is enforced by the import graph: the
domain package imports only the standard library.

## Currency (2026)

**STANDARD, situational.** Widely used in larger Go codebases (ThreeDotsLabs
"Wild Workouts"); over-applied to CRUD services.

## Fit signals

Many aggregates in one context, several infrastructure adapters (two stores,
a broker, an external gateway), a long expected life — and a monolith package
that already mixed gorm tags into domain rules.

## What it changes in the generated services

```
internal/domain/        entities + ports (interfaces), no gorm/pgx/http imports
internal/app/           use cases: PlaceOrder, Cancel
internal/adapters/postgres/  repository implementations (gorm or sqlc)
internal/adapters/http/      handlers → app
internal/adapters/clients/   typed clients implementing domain ports
```

- The porter still moves the monolith package verbatim into `internal/<pkg>`;
  the split into ports is develop-loop work, marked `// CW-SEAM[kind=port]`.
- `depguard`: `internal/domain` may import only stdlib; a violation fails lint.
- Plan field: `"structure": "hexagonal"`.

## Go 2026 implementation

```go
// internal/domain/order.go
type StockReserver interface {
    Reserve(ctx context.Context, productID uint, qty int) (Reservation, error)
    Release(ctx context.Context, r Reservation) error
}
```

## Anti-patterns

- Hexagonal for a three-struct CRUD service — three packages of indirection
  around one table.
- Mapper layers copying identical structs between every layer.

## Interacts with

- [[flat-package]] · [[anti-corruption-layer]] · [[testing]]
