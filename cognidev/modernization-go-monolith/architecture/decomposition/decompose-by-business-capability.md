---
id: decompose-by-business-capability
title: Decompose by business capability
family: decomposition
currency: standard
role: default
applies_when:
  - answer:boundary_strategy = domain
teaches: A service is a thing the business does and owns data for — catalog, orders, payments — never a technical layer.
---

# Decompose by business capability

## What it is

Draw each service around a business capability — something the business does, owns data for and could staff a team around: *catalog*, *customers*, *inventory*, *payments*, *orders*. In a Go monolith the capability is usually already a package (`internal/orders`), and the gorm models it migrates are its tables.

## Currency (2026)

**STANDARD — the default first cut.** microservices.io, Team Topologies; the shape every mainstream decomposition guide starts from.

## Fit signals (from `.cognidev/feature/context.md`)

- Packages that each own ≥1 persisted struct (`entities.json`).
- Few co-writes between them (`entity-access.json` `co_writes`).
- Route groups already split by capability (`/products`, `/orders`).

## What it changes in the generated services

- One service per data-owning package; package-less helpers (`internal/platform`) are NOT services — they become per-service copies or `pkg/`.
- `plan.json.services[].entities` lists the owned structs; `sources` every file the map assigned.

## Go 2026 implementation

The cut is computed, not written: `map-services-go` clusters entities by composition + co-write, names each cluster after the package its structs were declared in, and places the rest by type reference.

## Anti-patterns

- Cutting by layer (`handlers` service, `models` service) — every request crosses every service.
- A service per struct — nano-services with a network hop per join.

## Interacts with

- [[decompose-by-subdomain]] · [[database-per-service]] · [[right-sizing]]
