# Migration report

## Summary

This migration decomposed the monolith into **5 service(s)** (catalog-service, customers-service, inventory-service, payments-service, orders-service). 23/23 develop tasks completed cleanly. Built with go1.26.8.

## Measured against the monolith, before the cut

The monolith was green before anything was written: 1 target(s) built and 0 test(s) passed on go1.26.4. That is the denominator for everything below.

| | Monolith, before | Services (with gateway and kernel), after — monolith excluded |
|---|---|---|
| Targets built | 1 of 1 | 7 of 7 |
| Tests ran | 0 | 56 |
| Tests failed | 0 | 0 |

The monolith had no test run recorded, so nothing here compares behaviour. The build result is a compilation proof only.


## Manual intervention required

**Files with unfilled markers:**

- `CLAUDE.md (1)`


## Build proof

Every module is built from its own directory with `go build ./...` (never through a `go.work`), so a pass means each service compiles on its own.

| Module | `go build ./...` |
| --- | --- |
| `gateway` | ✓ |
| `pkg` | ✓ |
| `services/catalog` (service `catalog`) | ✓ |
| `services/customers` (service `customers`) | ✓ |
| `services/inventory` (service `inventory`) | ✓ |
| `services/orders` (service `orders`) | ✓ |
| `services/payments` (service `payments`) | ✓ |

The monolith's own module (`.`) is still in the tree beside the services. It is not a service, so it is neither built nor tested here and none of its tests are counted below — its build and suite are the before-state.

**Independence:** no service module `replace`s another service's module with a local path. ✓

## Test proof

`go test -json ./...`: **56** test(s) ran, **0** failed.

| Module | Tests | Failed |
| --- | --- | --- |
| `gateway` | 2 | 0 |
| `pkg` | 6 | 0 |
| `services/catalog` | 6 | 0 |
| `services/customers` | 9 | 0 |
| `services/inventory` | 5 | 0 |
| `services/orders` | 15 | 0 |
| `services/payments` | 13 | 0 |

## Static analysis

**go vet:** 0 finding(s).


**golangci-lint:** 60 issue(s).

- `gateway/../monolith-rt6/gateway/internal/proxy/proxy.go:61:15` — ST1023: should omit type http.Handler from declaration; it will be inferred from the right-hand side (golangci-lint/staticcheck)
- `pkg/../monolith-rt6/pkg/httpx/server_test.go:38:17` — Error return value of `resp.Body.Close` is not checked (golangci-lint/errcheck)
- `pkg/../monolith-rt6/pkg/httpx/server_test.go:49:17` — Error return value of `resp.Body.Close` is not checked (golangci-lint/errcheck)
- `pkg/../monolith-rt6/pkg/outbox/outbox.go:102:18` — Close should use defer (golangci-lint/sqlclosecheck)
- `services/catalog/../monolith-rt6/services/catalog/internal/platform/auth.go:1:1` — package-comments: should have a package comment (golangci-lint/revive)
- `services/catalog/internal/app/app.go:35:1` — exported: exported function New should have comment or be unexported (golangci-lint/revive)
- `services/catalog/internal/app/models.go:1:1` — package-comments: package comment should be of the form "Package app ..." (golangci-lint/revive)
- `services/catalog/internal/app/wire.go:1:1` — package-comments: package comment should be of the form "Package app ..." (golangci-lint/revive)
- `services/catalog/internal/catalog/catalog_test.go:1:1` — package-comments: should have a package comment (golangci-lint/revive)
- `services/customers/../monolith-rt6/services/customers/internal/platform/auth.go:1:1` — package-comments: should have a package comment (golangci-lint/revive)
- `services/customers/internal/app/app.go:35:1` — exported: exported function New should have comment or be unexported (golangci-lint/revive)
- `services/customers/internal/app/models.go:1:1` — package-comments: package comment should be of the form "Package app ..." (golangci-lint/revive)
- `services/customers/internal/app/wire.go:1:1` — package-comments: package comment should be of the form "Package app ..." (golangci-lint/revive)
- `services/customers/internal/customers/customers_test.go:1:1` — package-comments: should have a package comment (golangci-lint/revive)
- `services/inventory/../monolith-rt6/services/inventory/internal/contracts/orders/contracts.go:1:1` — package-comments: package comment should be of the form "Package orders ..." (golangci-lint/revive)
- `services/inventory/../monolith-rt6/services/inventory/internal/contracts/orders/contracts.go:11:7` — exported: exported const TopicOrderPlaced should have comment or be unexported (golangci-lint/revive)
- `services/inventory/../monolith-rt6/services/inventory/internal/platform/auth.go:1:1` — package-comments: should have a package comment (golangci-lint/revive)
- `services/inventory/internal/app/app.go:35:1` — exported: exported function New should have comment or be unexported (golangci-lint/revive)
- `services/inventory/internal/app/models.go:1:1` — package-comments: package comment should be of the form "Package app ..." (golangci-lint/revive)
- `services/inventory/internal/app/wire.go:1:1` — package-comments: package comment should be of the form "Package app ..." (golangci-lint/revive)
- `services/inventory/internal/httpapi/stock_item_handlers.go:18:49` — unused-parameter: parameter 'svc' seems to be unused, consider removing or renaming it as _ (golangci-lint/revive)
- `services/inventory/internal/inventory/inventory_test.go:1:1` — package-comments: should have a package comment (golangci-lint/revive)
- `services/inventory/internal/jobs/scheduler.go:36:18` — Error return value of `conn.Close` is not checked (golangci-lint/errcheck)
- `services/orders/../monolith-rt6/services/orders/internal/contracts/catalog/contracts.go:1:1` — package-comments: package comment should be of the form "Package catalog ..." (golangci-lint/revive)
- `services/orders/../monolith-rt6/services/orders/internal/contracts/payments/contracts.go:1:1` — package-comments: package comment should be of the form "Package payments ..." (golangci-lint/revive)
- `services/orders/../monolith-rt6/services/orders/internal/platform/auth.go:1:1` — package-comments: should have a package comment (golangci-lint/revive)
- `services/orders/internal/app/app.go:35:1` — exported: exported function New should have comment or be unexported (golangci-lint/revive)
- `services/orders/internal/app/models.go:1:1` — package-comments: package comment should be of the form "Package app ..." (golangci-lint/revive)
- `services/orders/internal/app/wire.go:1:1` — package-comments: package comment should be of the form "Package app ..." (golangci-lint/revive)
- `services/orders/internal/clients/catalog/client.go:1:1` — package-comments: package comment should be of the form "Package catalog ..." (golangci-lint/revive)
- … and 30 more (see `.cognidev/feature/verify.json`)
