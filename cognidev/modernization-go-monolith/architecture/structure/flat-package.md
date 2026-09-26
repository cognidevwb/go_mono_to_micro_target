---
id: flat-package
title: Flat package per context
family: structure
currency: standard
role: default
applies_when:
  - answer:service_internal_structure = flat-package
teaches: Keep the ported package as the unit — one package per context inside the service, wiring in cmd/, adapters beside it; Go rewards fewer layers.
---

# Flat package per context

## What it is

The idiomatic Go default: a service is its `cmd/<ctx>/main.go` plus the ported
context package (`internal/orders`) holding the model, the service type and the
HTTP handlers together, with `internal/clients`, `internal/outbox` and
`internal/saga` beside it. The package is the unit of encapsulation; unexported
identifiers are the boundary. It is what the monolith already had, so the port
moves the package verbatim.

## Currency (2026)

**STANDARD.** "Package by feature, keep it flat" is the dominant Go guidance
(go.dev/doc/modules/layout; Ben Johnson's Standard Package Layout).

## Fit signals

A right-sized service: one context, a handful of structs, the monolith already
package-per-context (`internal/catalog`, `internal/orders`).

## What it changes in the generated services

- Ported files keep their package name and directory:
  `internal/orders/{model,service,handlers}.go`.
- The dependency rule (domain code does not import `net/http` clients directly)
  is enforced by `depguard` in `.golangci.yml`, not by extra modules.
- Plan field: `"structure": "flat-package"`.

## Go 2026 implementation

```go
// internal/orders/service.go — the ported type, collaborators become interfaces
type Service struct {
    db        *gorm.DB
    catalog   CatalogClient   // was *catalog.Service
    inventory InventoryClient // was *inventory.Service
}
type CatalogClient interface{ PriceOf(ctx context.Context, productID uint) (float64, error) }
```

The consumer declares the interface it needs ("accept interfaces, return
structs"), so the typed client in `internal/clients` satisfies it and tests use a fake.

## Anti-patterns

- Splitting into `models/`, `services/`, `handlers/` packages for symmetry —
  every type becomes exported and the package boundary stops meaning anything.
- A `domain` package importing `gorm` tags everywhere when the service is CRUD
  — acceptable here; the hexagonal card is for when it is not.

## Interacts with

- [[hexagonal]] — the alternative for large, long-lived services.
- [[microservices-conventions]] — the surrounding layout.
