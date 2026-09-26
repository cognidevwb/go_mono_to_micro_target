---
id: dapr
title: Dapr
family: communication
currency: rising
role: reference
applies_when: []
why_not: This run talks to NATS/Kafka and HTTP directly from Go; Dapr is an additional sidecar runtime nothing here installs.
teaches: A sidecar runtime giving portable pub/sub, state and invocation APIs — useful across heterogeneous stacks, extra runtime for a Go-only fleet.
---

# Dapr

## What it is

Dapr runs beside each service and exposes building blocks (pub/sub, state, service invocation, workflow) over HTTP/gRPC; the Go SDK (`github.com/dapr/go-sdk`) wraps them.

## Currency (2026)

**RISING**, CNCF graduated.

## Fit signals (from `.cognidev/feature/context.md`)

- Polyglot estates; platform teams standardising on Dapr components.

## What it changes in the generated services

Not generated — this run wires NATS/Kafka clients and net/http directly.

## Go 2026 implementation

`client, _ := dapr.NewClient(); client.PublishEvent(ctx, "pubsub", "order.placed", payload)`

## Anti-patterns

- Adding a sidecar runtime to get pub/sub a Go client library already gives you.

## Interacts with

- [[event-driven]] · [[sidecar]]
