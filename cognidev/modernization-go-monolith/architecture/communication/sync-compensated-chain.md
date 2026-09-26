---
id: sync-compensated-chain
title: Synchronous compensated call chain
family: communication
currency: standard
role: choice
applies_when:
  - answer:event_backbone = none
  - sagas > 0
teaches: With no broker, the orchestrator calls each step over HTTP with an idempotency key and calls the undo endpoints itself on failure.
---

# Synchronous compensated call chain

## What it is

The saga without a broker: the orchestrator calls reserve, then charge, synchronously, persisting its progress locally between calls, and on failure calls the compensating endpoints (`DELETE /v1/reservations/{id}`) with retries. Weaker than an event saga — the orchestrator must be up for compensation — but honest for a small estate that refused a broker.

## Currency (2026)

**STANDARD** for small estates without messaging.

## Fit signals (from `.cognidev/feature/context.md`)

- `event_backbone = none` and the map still found a multi-service write.

## What it changes in the generated services

- Orchestrator gets `internal/saga/` with a local progress table and a retry sweeper goroutine; every step endpoint accepts `Idempotency-Key`.

## Go 2026 implementation

```go
key := sagaID.String() + ":reserve"
res, err := inv.Reserve(ctx, key, lines)
if err != nil { return s.fail(ctx, sagaID, err) }
if _, err := pay.Charge(ctx, sagaID.String()+":charge", total); err != nil {
    return errors.Join(err, inv.Release(ctx, key, res))
}
```

## Anti-patterns

- No persisted progress — a crash mid-chain leaks reservations.
- Non-idempotent step endpoints retried.

## Interacts with

- [[idempotency-key]] · [[saga-isolation]] · [[resilience]]
