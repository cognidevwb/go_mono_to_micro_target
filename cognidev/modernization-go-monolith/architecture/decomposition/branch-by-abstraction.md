---
id: branch-by-abstraction
title: Branch by abstraction
family: decomposition
currency: standard
role: choice
applies_when:
  - answer:migration_approach = strangler
teaches: Put an interface in front of the in-process call, swap the implementation for a client behind it — the in-code half of the strangler.
---

# Branch by abstraction

## What it is

The in-process sibling of the strangler fig. Replace a concrete dependency (`*catalog.Service`) with an interface the caller owns, keep the in-process implementation, add the HTTP-client implementation, and switch per environment. The monolith keeps compiling at every step.

## Currency (2026)

**STANDARD.** Fowler/Hammant; Go's implicit interfaces make it nearly free.

## Fit signals (from `.cognidev/feature/context.md`)

- Struct fields holding another context's concrete service (`catalog *catalog.Service`).
- Cross-context calls in `graph.json`.

## What it changes in the generated services

- The port rewrites such fields to consumer-owned interfaces; `internal/clients/<ctx>_client.go` implements them.
- Seam: `// CW-SEAM[kind=cross-context-call]` on each converted call site.

## Go 2026 implementation

```go
type Catalog interface{ PriceOf(ctx context.Context, id uint) (float64, error) }
// in the monolith: *catalog.Service satisfies it (add ctx)
// in the service:  *clients.CatalogClient satisfies it
```

## Anti-patterns

- A "toggle" left in production forever with both paths live.
- An interface mirroring every method of the provider instead of what the caller uses.

## Interacts with

- [[strangler-fig]] · [[anti-corruption-layer]] · [[sync-transport]]
