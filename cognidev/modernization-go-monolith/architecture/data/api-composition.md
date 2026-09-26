---
id: api-composition
title: API composition
family: data
currency: standard
role: default
applies_when:
  - calls > 0
teaches: Answer a cross-service query by calling the owners and joining in memory — try this before building a read model.
---

# API composition

## What it is

A query needing data from several services is served by a composer that calls each owner's API and joins the results. The replacement for the monolith's cross-context `Preload`/JOIN.

## Currency (2026)

**STANDARD — default query pattern.**

## Fit signals (from `.cognidev/feature/context.md`)

- Cross-context reads in `graph.json` (orders reading catalog prices, customers' status).

## What it changes in the generated services

- The consumer's typed client + an in-memory join in its handler; batch endpoints (`?ids=`) on providers where a loop exists.

## Go 2026 implementation

```go
prices, err := catalog.Prices(ctx, productIDs) // one call, not one per line
```

## Anti-patterns

- N+1 calls inside a loop.
- Composition over services with incompatible SLAs on a hot route.

## Interacts with

- [[aggregator]] · [[materialized-view]] · [[cqrs]]
