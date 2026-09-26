---
id: modular-monolith-first
title: Modular monolith first
family: decomposition
currency: rising
role: choice
applies_when:
  - services <= 2
teaches: With one or two real contexts, enforce package boundaries in one binary instead of paying for a network.
---

# Modular monolith first

## What it is

Keep one deployable, but make the boundaries real: each context a package with an exported API, no cross-context access to tables, boundaries enforced by lint. Go's `internal/` and the import graph make this cheap, and a module can graduate to a service later without rewriting callers.

## Currency (2026)

**RISING.** The 2024–2026 correction to microservices-by-default (Fowler, ThoughtWorks Radar, Shopify).

## Fit signals (from `.cognidev/feature/context.md`)

- The map found ≤2 data-owning contexts.
- Few coupling edges; no named force (team, scale, isolation) per context.

## What it changes in the generated services

- Recommends not carving; if the run continues it still generates the services, and the briefing says the split is not justified by the measurements.
- `depguard` rules forbidding cross-context imports except through each context's API package.

## Go 2026 implementation

Boundary enforcement in `.golangci.yml`:

```yaml
linters-settings:
  depguard:
    rules:
      orders:
        files: ["**/internal/orders/**"]
        deny: [{ pkg: "github.com/acme/shop/internal/inventory/store", desc: "use inventory's API" }]
```

## Anti-patterns

- Calling it modular while every package still reads every table.
- Skipping it because "we will split later anyway".

## Interacts with

- [[right-sizing]] · [[branch-by-abstraction]]
