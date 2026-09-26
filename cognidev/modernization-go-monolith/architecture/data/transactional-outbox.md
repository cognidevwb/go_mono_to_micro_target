---
id: transactional-outbox
title: Transactional outbox
family: data
currency: standard
role: invariant
applies_when:
  - answer:outbox != none
  - answer:event_backbone != none
teaches: Write the row and the message in one local transaction, relay after commit — the only honest way out of the dual write.
---

# Transactional outbox

## What it is

A handler that changes state and publishes must not do those as two operations. The message is inserted into an `outbox` table in the SAME database transaction as the state change; a relay goroutine reads unsent rows, publishes to NATS/Kafka, and marks them sent after the broker acks. Delivery is at-least-once, so consumers dedupe.

## Currency (2026)

**STANDARD — invariant.** microservices.io; in Go it is a table plus ~100 lines, or Watermill's SQL forwarder (`github.com/ThreeDotsLabs/watermill-sql`).

## Fit signals (from `.cognidev/feature/context.md`)

- The monolith's `bus.Publish` after `db.Transaction` returns (the fixture's `PlaceOrder` publishes `order.placed` after commit).
- Every saga step message.

## What it changes in the generated services

- `internal/outbox/outbox.go` + `outbox` table in the service's first migration; `Add(ctx, tx, subject, payload)` takes the caller's `pgx.Tx`/`*gorm.DB` tx.
- A relay started from `main` with the service context.

## Go 2026 implementation

```go
func (o *Outbox) Add(ctx context.Context, tx pgx.Tx, subject string, v any) error {
    b, err := json.Marshal(v); if err != nil { return err }
    _, err = tx.Exec(ctx, `INSERT INTO outbox (id, subject, payload) VALUES ($1,$2,$3)`, uuid.New(), subject, b)
    return err
}
// relay: SELECT ... FOR UPDATE SKIP LOCKED LIMIT 100 → js.Publish(ctx, subj, payload, jetstream.WithMsgID(id)) → UPDATE sent_at
```
`WithMsgID` lets JetStream's duplicate window drop a republish.

## Anti-patterns

- `tx.Commit()` then `js.Publish()` — the dual write.
- Publishing inside the transaction to the broker (couples availability, still not atomic).
- Assuming exactly-once.

## Interacts with

- [[outbox-relay-topology]] · [[idempotent-consumer]] · [[saga-orchestration]] · [[change-data-capture]]
