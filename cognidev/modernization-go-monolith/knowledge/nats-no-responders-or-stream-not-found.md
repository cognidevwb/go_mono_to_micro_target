---
id: nats-no-responders-or-stream-not-found
runtime: go
since: 0
category: runtime-failure
tags: [nats, jetstream, events, outbox, broker, startup-order]
symptoms:
  - "nats: no responders available for request"
  - "nats: stream not found"
  - "nats: no stream matches subject"
  - "nats: no heartbeat received"
remedy: { kind: environment }
reference: "nats.go / jetstream package errors (pkg.go.dev/github.com/nats-io/nats.go/jetstream#pkg-variables)"
---

# The in-process bus became a broker, and nobody created the stream

## The signature

```
level=ERROR msg="publish order.placed" err="nats: no responders available for request"
```

or on the consumer side:

```
level=ERROR msg="consume" subject=order.placed err="nats: stream not found"
```

## What it means

In the monolith `platform.Bus.Publish("order.placed", id)` called the inventory
handler synchronously, in the same goroutine — delivery could not fail. After
the cut that topic is a JetStream subject, and a JetStream publish needs a
stream whose subjects cover it. `no responders` from a JetStream publish means
no stream is listening on that subject; `stream not found` means the consumer
started before anything created it.

Both are ordering problems, not code problems: whichever service starts first
decides whether the stream exists.

## What to do

- Streams are declared idempotently at startup by the PUBLISHER
  (`js.CreateOrUpdateStream`), one stream per publishing context
  (`ORDERS` → subjects `order.>`), and consumers are durable.
- The outbox relay retries a failed publish; the row stays until the broker
  acks. A publish error must never be swallowed in a handler.
- Locally, compose starts `nats -js` before any service (`depends_on` with a
  healthcheck on `:8222/healthz`).
