---
id: cqrs
title: CQRS
family: data
currency: standard
role: choice
applies_when:
  - answer:read_write_model = cqrs
teaches: Separate the write model from a read model shaped for queries — per service, where reads are heavy or their shape fights the writes.
---

# CQRS

## What it is

Commands mutate the aggregate; queries read a projection kept current by events. <a id="crud-baseline"></a>The default (CRUD) is one model per service.

## Currency (2026)

**STANDARD, often over-applied.**

## Fit signals (from `.cognidev/feature/context.md`)

- User chose `read_write_model = cqrs`; read-heavy services in the measured route surface.

## What it changes in the generated services

- `internal/<pkg>/commands/` and `internal/<pkg>/queries/`; sqlc query files for the read side; a projector consumer updating read tables.

## Go 2026 implementation

```go
type PlaceOrder struct{ CustomerID uint; Lines []Line }
func (h *Handlers) PlaceOrder(ctx context.Context, c PlaceOrder) (uint, error)
func (q *Queries) OrdersForCustomer(ctx context.Context, id uint) ([]OrderView, error) // sqlc-generated
```

## Anti-patterns

- Fleet-wide CQRS.
- A read model without the outbox feeding it.

## Interacts with

- [[materialized-view]] · [[api-composition]] · [[transactional-outbox]]
