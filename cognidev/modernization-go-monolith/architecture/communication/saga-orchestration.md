---
id: saga-orchestration
title: Saga — orchestration
family: communication
currency: standard
role: default
applies_when:
  - answer:coordination = orchestration
  - answer:event_backbone != none
  - sagas > 0
teaches: One coordinator owns the multi-service write: each step a local transaction, each failure a compensating step, the riskiest step last.
---

# Saga — orchestration

## What it is

Replaces a `db.Transaction` that spanned several contexts. The orchestrator (the service that owned the flow — orders for `PlaceOrder`) persists saga state, sends a command per step, waits for the reply event, and on failure runs compensations in reverse. For the fixture: reserve stock (inventory) → charge payment (payments) → confirm order (orders); on charge failure release the reservation and mark the order failed.

## Currency (2026)

**STANDARD — the default** for flows over ~3 steps or business-critical ones (microservices.io, Azure).

## Fit signals (from `.cognidev/feature/context.md`)

- `plan.json.sagas` / `decomposition-facts.json` sagas: a flow writing in two or more services, with its ordered steps.
- `dataflow.json` consistency boundaries whose effective writes span contexts.

## What it changes in the generated services

- `internal/saga/<saga>.go` in the orchestrator: a state struct persisted in its own table (`saga_state`), a step table, compensations.
- Commands and replies ride the outbox; seam `// CW-SEAM[kind=saga]` replaces the cross-context `db.Transaction`.

## Go 2026 implementation

```go
type PlaceOrderSaga struct{ ID uuid.UUID; OrderID uint; Step string; Reservation string }
func (s *Orchestrator) OnStockReserved(ctx context.Context, e StockReserved) error {
    return s.tx(ctx, func(tx pgx.Tx) error {
        if err := s.store.Advance(ctx, tx, e.SagaID, "charging"); err != nil { return err }
        return s.outbox.Add(ctx, tx, "payments.charge", ChargePayment{SagaID: e.SagaID, Amount: e.Total})
    })
}
```
No framework is required; for long-running flows see [[durable-execution]].

## Anti-patterns

- A saga with no persisted state — a crash forgets where it was.
- Compensations that can fail silently.
- Charging payment before the reversible steps.

## Interacts with

- [[saga-isolation]] · [[transactional-outbox]] · [[idempotent-consumer]] · [[saga-choreography]]
