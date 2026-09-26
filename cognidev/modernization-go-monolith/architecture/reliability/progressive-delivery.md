---
id: progressive-delivery
title: Progressive delivery
family: reliability
currency: standard
role: choice
applies_when:
  - answer:migration_approach = strangler
teaches: Move traffic to an extracted service gradually — weighted routes and flags, with a fast way back.
---

# Progressive delivery

## What it is

Canary a route at the gateway (1% → 10% → 100%) or behind a feature flag (OpenFeature Go SDK), watching error rate and latency; roll back by weight, not redeploy.

## Currency (2026)

**STANDARD.**

## Fit signals (from `.cognidev/feature/context.md`)

- Strangler migration.

## What it changes in the generated services

- Weighted targets in the gateway route table; Argo Rollouts / ACA revision traffic splitting.

## Go 2026 implementation

```go
if rand.IntN(100) < weight(route) { newProxy.ServeHTTP(w, r) } else { monolith.ServeHTTP(w, r) }
```

## Anti-patterns

- Big-bang cutover of a route.

## Interacts with

- [[strangler-fig]] · [[api-gateway]]
