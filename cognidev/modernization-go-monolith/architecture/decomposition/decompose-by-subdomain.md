---
id: decompose-by-subdomain
title: Decompose by subdomain (DDD)
family: decomposition
currency: rising
role: choice
applies_when:
  - answer:boundary_strategy = subdomain
teaches: Split along DDD subdomains — core, supporting, generic — so the investment goes where the business differentiates.
---

# Decompose by subdomain (DDD)

## What it is

Domain-Driven Design's refinement of capability decomposition: classify each area as **core** (differentiating — pricing, order orchestration), **supporting** (necessary, specific — inventory) or **generic** (buy it — auth, email) and cut finer where the core is. Bounded contexts carry their own ubiquitous language; the same word (`Product`) may mean different structs in catalog and inventory.

## Currency (2026)

**RISING.** Event storming and context mapping have gone mainstream with the modular-monolith movement (Vernon, Khononov *Learning DDD*).

## Fit signals (from `.cognidev/feature/context.md`)

- One package's model used with two meanings (`Product` with price in catalog, with `OnHand` in inventory).
- The user filled `cognidev/domain.md` with a glossary.

## What it changes in the generated services

- Finer services than `domain`; generic subdomains are flagged as buy/replace candidates rather than carved.
- Each context gets its own copy of shared-name structs; a translation lives in its client ([[anti-corruption-layer]]).

## Go 2026 implementation

No framework — subdomains are packages with their own types. Duplicate the struct per context rather than share it through `pkg/`:

```go
// services/inventory/internal/inventory/model.go
type StockItem struct{ ProductID uint; OnHand, Reserved int }
// services/catalog/internal/catalog/model.go
type Product struct{ ID uint; Name string; Price float64 }
```

## Anti-patterns

- Fine-grained cuts without a team per subdomain — more hops, no autonomy gained.
- One canonical `Product` struct in `pkg/` shared by all — a distributed monolith through the type system.

## Interacts with

- [[decompose-by-business-capability]] · [[anti-corruption-layer]] · [[right-sizing]]
