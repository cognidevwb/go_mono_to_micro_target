---
id: anti-corruption-layer
title: Anti-corruption layer
family: decomposition
currency: standard
role: invariant
applies_when:
  - calls > 0
teaches: Translate at every boundary — the consumer maps the provider’s contract into its own types, so one service’s model never leaks into another.
---

# Anti-corruption layer

## What it is

A translation layer on the consumer side of every cross-context call. The consumer never uses the provider's structs as its own model; it maps the contract DTO into local types at the edge. Between a new service and the remaining monolith it also shields the new model from legacy shapes (bare-string statuses, float money).

## Currency (2026)

**STANDARD — invariant where calls cross.** Evans/DDD; Azure Architecture Center.

## Fit signals (from `.cognidev/feature/context.md`)

- Any cross-context call edge in `plan.json` (`calls`).
- Legacy smells in the monolith model (`Payment.Status string`).

## What it changes in the generated services

- Each `internal/clients/<ctx>_client.go` maps `pkg/contracts/<ctx>` DTOs into local types; no contract type appears in the domain package.

## Go 2026 implementation

```go
func (c *CatalogClient) PriceOf(ctx context.Context, id uint) (float64, error) {
    var dto catalogv1.Price
    if err := c.get(ctx, fmt.Sprintf("/v1/products/%d/price", id), &dto); err != nil { return 0, fmt.Errorf("catalog price %d: %w", id, err) }
    return dto.Amount, nil
}
```

## Anti-patterns

- Importing the provider's domain package to reuse its struct.
- One giant shared `models` module.

## Interacts with

- [[branch-by-abstraction]] · [[contract-testing]] · [[decompose-by-subdomain]]
