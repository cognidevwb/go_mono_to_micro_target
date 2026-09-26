---
id: outbox-relay-topology
title: Outbox relay topology
family: data
currency: standard
role: invariant
applies_when:
  - answer:outbox != none
  - answer:event_backbone != none
teaches: Decide who drains the outbox — one relay per service, leader-elected or SKIP LOCKED — and make sure every replica is allowed to run it.
---

# Outbox relay topology

## What it is

The relay is a loop that must run exactly somewhere. With N replicas, either every replica polls using `FOR UPDATE SKIP LOCKED` (safe concurrent draining) or one leader relays (Kubernetes Lease). Scale-to-zero platforms need a minimum replica or the outbox never drains.

## Currency (2026)

**STANDARD — invariant** wherever an outbox exists.

## Fit signals (from `.cognidev/feature/context.md`)

- Autoscaled services that publish; `orchestrator` = azure-containerapps with scale-to-zero.

## What it changes in the generated services

- Relay runs in-process in every replica with SKIP LOCKED; `minReplicas: 1` for publishing services in ACA/HPA.

## Go 2026 implementation

```go
rows, err := tx.Query(ctx, `SELECT id, subject, payload FROM outbox WHERE sent_at IS NULL ORDER BY created_at LIMIT 100 FOR UPDATE SKIP LOCKED`)
```

## Anti-patterns

- A relay only in replica 0 with no leader election.
- Scale-to-zero on a publisher.

## Interacts with

- [[transactional-outbox]] · [[autoscaling]]
