---
id: undefined-the-type-the-monolith-shared-in-package
runtime: go
since: 0
category: build-failure
tags: [port, packages, contracts, typed-client, cross-context]
symptoms:
  - "undefined: catalog."
  - "undefined:"
  - "could not import"
  - "imported and not used"
remedy: { kind: seam }
reference: "go/types error codes · UndeclaredImportedName / UnusedImport (pkg.go.dev/golang.org/x/tools/internal/typesinternal)"
---

# `undefined:` — the type the monolith shared by importing a package

## The signature

```
services/orders/internal/orders/service.go:21:13: undefined: catalog.Service
services/orders/internal/orders/service.go:22:13: undefined: inventory.Service
services/orders/internal/orders/service.go:9:2: "github.com/acme/shop/internal/catalog" imported and not used
```

## What it means

Inside the monolith a context used another context's types directly:
`orders.Service` held a `*catalog.Service` and called `s.catalog.PriceOf(id)`.
The port moved `orders` into its own module and correctly did NOT bring
`catalog` with it — so every reference to that package is now undefined.

The wrong fix is the one the compiler suggests: vendor `catalog` into the
orders module. That compiles, and it is the hollow decomposition — two services
each carrying a copy of the other's domain, reading two databases that used to
be one.

## What to do

- Each `undefined: <ctx>.Service` field becomes a typed client interface owned
  by the CONSUMER (`internal/clients/catalog_client.go`), with only the methods
  it actually calls.
- Each `undefined: <ctx>.<Entity>` becomes the provider's contract DTO in
  `pkg/contracts/<ctx>/` — fields the consumer reads, JSON tags, no gorm tags.
- Calls that passed a `tx *gorm.DB` are not portable as written — see
  [two gorm.DBs are not one transaction](two-gorm-dbs-are-not-one-transaction.md).
