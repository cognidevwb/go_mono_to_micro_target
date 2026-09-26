---
id: domain-event-aggregate
title: Domain events and aggregates
family: data
currency: standard
role: invariant
applies_when:
  - always
teaches: An aggregate is the consistency boundary; what it records as having happened is the event other services react to.
---

# Domain events and aggregates

## What it is

An aggregate (Order with its OrderLines) is changed in one local transaction; the facts it produces (`OrderPlaced`) are recorded as domain events and mapped to integration events (the published contract) at the boundary.

## Currency (2026)

**STANDARD — invariant.**

## Fit signals (from `.cognidev/feature/context.md`)

- `plan.json.services[].aggregates` (the root the clustering computed); `events.json`.

## What it changes in the generated services

- Integration event structs in `pkg/contracts/<ctx>/events.go`; the aggregate method returns the events it raised; the application service adds them to the outbox.

## Go 2026 implementation

```go
func (o *Order) Place() []any { o.Status = StatusPlaced; return []any{OrderPlaced{OrderID: o.ID, Total: o.Total}} }
```

## Anti-patterns

- Publishing the gorm struct as the event.
- Aggregates spanning services.

## Interacts with

- [[transactional-outbox]] · [[event-driven]]
