---
id: autoscaling
title: Autoscaling — HPA and KEDA
family: reliability
currency: standard
role: choice
applies_when:
  - answer:orchestrator != compose
teaches: Scale HTTP services on CPU/RPS and consumers on queue depth — CPU is the wrong signal for a JetStream consumer.
---

# Autoscaling — HPA and KEDA

## What it is

HPA on CPU (requests are the denominator) for HTTP services; KEDA's NATS JetStream or Kafka scaler on consumer lag for event consumers; Container Apps scale rules the same way. Go's `GOMAXPROCS` respects cgroup CPU limits since Go 1.25.

## Currency (2026)

**STANDARD.**

## Fit signals (from `.cognidev/feature/context.md`)

- `load_profile`, `service_size` answers; consumers in `events.json`.

## What it changes in the generated services

- HPA per HTTP service; KEDA ScaledObject per consumer; ACA scale rules; pool budget recomputed.

## Go 2026 implementation

KEDA `nats-jetstream` trigger: `stream: ORDERS, consumer: inventory, lagThreshold: "100"`.

## Anti-patterns

- Scale-to-zero on outbox publishers.
- Autoscaling without a connection-pool budget.

## Interacts with

- [[connection-pool-budget]] · [[outbox-relay-topology]] · [[deployment]]
