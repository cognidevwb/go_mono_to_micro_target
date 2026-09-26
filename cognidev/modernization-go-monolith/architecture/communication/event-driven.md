---
id: event-driven
title: Event-driven messaging
family: communication
currency: standard
role: default
applies_when:
  - answer:event_backbone != none
teaches: The in-process bus becomes a broker: every topic a context publishes becomes a durable subject, published through the outbox.
---

# Event-driven messaging

## What it is

Services publish facts (`order.placed`) and others react asynchronously. In the monolith `platform.Bus` delivered synchronously in-process; after the cut each topic is a NATS JetStream subject (or Kafka topic), delivery is at-least-once, and the publisher does not know the consumers.

## Currency (2026)

**STANDARD.** NATS JetStream is the lightweight Go-native default (the server is written in Go); Kafka via franz-go (`github.com/twmb/franz-go`) for high-throughput streaming.

## Fit signals (from `.cognidev/feature/context.md`)

- `events.json`: publishers and subscribers of each event, and which services they landed in; an event whose publisher and subscriber land in different services crosses the cut.

## What it changes in the generated services

- `internal/events/<event>.go` per published event (subject name + payload struct in `pkg/contracts/<ctx>`).
- Publisher declares its stream at startup; each subscriber has a durable consumer.
- Seam `// CW-SEAM[kind=event]` where `bus.Publish` was called.

## Go 2026 implementation

```go
js, _ := jetstream.New(nc)
_, _ = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{Name: "ORDERS", Subjects: []string{"order.>"}})
cons, _ := js.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{Durable: "inventory", FilterSubject: "order.placed", AckPolicy: jetstream.AckExplicitPolicy})
_, _ = cons.Consume(func(m jetstream.Msg) { if handle(m) == nil { m.Ack() } })
```

## Anti-patterns

- Publishing directly after the DB commit (dual write) — use the [[transactional-outbox]].
- Events carrying the full entity as a shared schema.
- Core NATS (fire-and-forget) for events that must not be lost.

## Interacts with

- [[transactional-outbox]] · [[idempotent-consumer]] · [[saga-choreography]]
