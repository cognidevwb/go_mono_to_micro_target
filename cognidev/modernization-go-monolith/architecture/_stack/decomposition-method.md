---
id: decomposition-method
title: How the cut is computed
family: method
currency: standard
role: reference
applies_when:
  - always
teaches: Boundaries come from data ownership and co-writes measured in the shards; the model only refines names and edges it can be checked against.
---

# How the cut is computed

## What it is

The method behind `map-services-go`: the monolith is cut by **data ownership**,
not by folder. Each persisted struct (gorm model / sqlc row type) is an entity;
entities that are composed (`Order.Lines []OrderLine`) or written together in
one method (co-writes) cluster into an aggregate; each aggregate names a
bounded context; every other file follows the entities it declares, the types
it references, or the package directory it sits in.

## Currency (2026)

**STANDARD.** Data-ownership-first decomposition is the microservices.io /
Team Topologies consensus; the determinism (no model in the default path) is
this playbook's own rule.

## Fit signals

Always applies. The inputs are `.cognidev/understand/entities.json`,
`entity-access.json` (write sets, co-writes), `dataflow.json` (consistency
boundaries — a `db.Transaction` spanning contexts), `graph.json` (receiver-bound
call edges), `events.json`, `externals.json`, `datastructures.json`,
`ambient-context.json`.

## What it changes in the generated services

- `service-map.json` — every file → its service.
- `plan.json` — services, owned tables, aggregates, `calls`, `sagas`, `order`.
- `decomposition-facts.json` — the measured edges the optional refine is checked against.
- A Go package is one compile unit: files in one directory land in one service
  unless they measurably straddle a boundary, and a straddling file is reported.

## Go 2026 implementation

Order of evidence, strongest first: owned entity → type reference → own name
(`OrderService` → `Order`) → path segment (`internal/orders/`) → sibling files
in the same package directory.

## Anti-patterns

- Cutting by top-level folder (`internal/`, `cmd/`, `pkg/`) — structural, not domain.
- Asking a model for the whole decomposition in one call, unchecked.
- Treating a package-level `var` or an in-process bus as free after the split.

## Interacts with

- [[right-sizing]] · [[decompose-by-business-capability]] · [[database-per-service]] · [[saga-orchestration]]
