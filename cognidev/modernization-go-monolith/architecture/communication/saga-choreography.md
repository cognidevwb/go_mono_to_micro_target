---
id: saga-choreography
title: Saga — choreography
family: communication
currency: standard
role: choice
applies_when:
  - answer:coordination = choreography
  - answer:event_backbone != none
teaches: No coordinator — each service reacts to the previous fact and emits its own; simple to start, hard to trace past three steps.
---

# Saga — choreography

## What it is

Each participant subscribes to the event that precedes its step and publishes its outcome: `order.created` → inventory reserves → `stock.reserved` → payments charges → `payment.charged` → orders confirms. Compensation is also event-driven (`payment.failed` → inventory releases).

## Currency (2026)

**STANDARD** for short, independent fan-out; teams switch to orchestration when flows grow.

## Fit signals (from `.cognidev/feature/context.md`)

- User chose choreography; saga steps ≤ 3; each participant's reaction independent.

## What it changes in the generated services

- No `internal/saga/`; each participant gets a consumer + a publisher; a `saga_id` travels in every payload for tracing.

## Go 2026 implementation

```go
func (h *Handler) OnOrderCreated(ctx context.Context, e OrderCreated) error {
    return h.tx(ctx, func(tx pgx.Tx) error {
        r, err := h.reserve(ctx, tx, e.Lines)
        if err != nil { return h.outbox.Add(ctx, tx, "stock.rejected", StockRejected{SagaID: e.SagaID}) }
        return h.outbox.Add(ctx, tx, "stock.reserved", StockReserved{SagaID: e.SagaID, Reservation: r})
    })
}
```

## Anti-patterns

- Cyclic event chains nobody can draw.
- No correlation id — the flow cannot be traced.

## Interacts with

- [[saga-orchestration]] · [[event-driven]] · [[observability]]
