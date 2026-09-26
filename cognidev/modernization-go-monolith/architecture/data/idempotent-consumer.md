---
id: idempotent-consumer
title: Idempotent consumer
family: data
currency: standard
role: invariant
applies_when:
  - answer:event_backbone != none
teaches: At-least-once is reality — record each processed message id in the same transaction as its effect.
---

# Idempotent consumer

## What it is

Every consumer records the message id it processed in a `processed_messages` table inside the same local transaction as the side effect; a redelivery finds the id and acks without re-applying.

## Currency (2026)

**STANDARD — invariant.** JetStream and Kafka are at-least-once by default.

## Fit signals (from `.cognidev/feature/context.md`)

- Every subscriber in `events.json` that lands in a service.

## What it changes in the generated services

- `processed_messages(consumer, msg_id)` primary key in each consuming service's migration; the consumer wrapper checks it.

## Go 2026 implementation

```go
tag, err := tx.Exec(ctx, `INSERT INTO processed_messages (consumer, msg_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, "inventory", m.Headers().Get(jetstream.MsgIDHeader))
if tag.RowsAffected() == 0 { return m.Ack() } // already applied
```

## Anti-patterns

- Deduping in memory (lost on restart, per replica).
- Acking before the transaction commits.

## Interacts with

- [[transactional-outbox]] · [[event-driven]]
