---
id: aggregator
title: Gateway aggregation
family: communication
currency: standard
role: default
applies_when:
  - calls > 0
teaches: When one screen needs three services, compose once at the edge with bounded fan-out — not three round-trips from the client.
---

# Gateway aggregation

## What it is

A thin composition endpoint that calls several services concurrently and merges the results, so a client does one request. Read-only, no business rules.

## Currency (2026)

**STANDARD.** Azure Architecture Center (Gateway Aggregation).

## Fit signals (from `.cognidev/feature/context.md`)

- Endpoints in `decomposition-facts.json` whose handler reaches more than one service (`routes_spanning`).

## What it changes in the generated services

- Such routes get an aggregation handler in the gateway or the owning service, calling clients in parallel with `errgroup`.

## Go 2026 implementation

```go
g, ctx := errgroup.WithContext(r.Context())
var cust Customer; var orders []Order
g.Go(func() (err error) { cust, err = customers.Get(ctx, id); return })
g.Go(func() (err error) { orders, err = ordersC.ForCustomer(ctx, id); return })
if err := g.Wait(); err != nil { problem(w, err); return }
```

## Anti-patterns

- Writes through an aggregator.
- Unbounded fan-out (`errgroup.SetLimit` it).

## Interacts with

- [[api-composition]] · [[api-gateway]] · [[resilience]]
