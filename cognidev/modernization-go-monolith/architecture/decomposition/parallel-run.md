---
id: parallel-run
title: Parallel run
family: decomposition
currency: niche
role: choice
applies_when:
  - answer:migration_approach = parallel-run
teaches: Run old and new on the same input, compare, then switch — proof for logic where a wrong answer costs money.
---

# Parallel run

## What it is

For high-risk logic (pricing, payments), send each request to both the monolith and the new service, serve the monolith's answer, and record differences. Switch when the diff rate is zero over a meaningful window.

## Currency (2026)

**NICHE.** GitHub Scientist-style; expensive, and justified only where correctness is costly.

## Fit signals (from `.cognidev/feature/context.md`)

- The user chose `migration_approach = parallel-run`.
- Money-bearing flows (a payments context, an order total).

## What it changes in the generated services

- Gateway shadows traffic to the new service (reads) or both sides consume the same events (writes, via CDC/outbox).
- A `diffs` table / metric per compared endpoint.

## Go 2026 implementation

```go
go func(r *http.Request) { // shadow, never affects the response
    ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), 2*time.Second)
    defer cancel()
    compare(ctx, legacyResp, newClient.Do(r.Clone(ctx)))
}(r)
```

## Anti-patterns

- Shadowing non-idempotent writes into a live payment provider.
- Running it forever instead of switching.

## Interacts with

- [[strangler-fig]] · [[change-data-capture]] · [[idempotency-key]]
