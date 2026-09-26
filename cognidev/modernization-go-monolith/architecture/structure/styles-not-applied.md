---
id: styles-not-applied
title: Architecture styles this playbook does not apply
family: structure
currency: reference
role: reference
applies_when:
  - always
teaches: Omissions are decisions — plugin architectures, space-based grids, serverless-per-function and a shared-library "platform" are not targets here.
---

# Architecture styles this playbook does not apply

## What it is

The styles a general architecture survey lists that this playbook deliberately
does not generate, each with its reason, so their absence reads as a decision.

## Currency (2026)

Reference only.

## Fit signals

Always listed.

## What it changes in the generated services

Nothing — that is the point.

| Style | Why not here |
| --- | --- |
| Go `plugin` package / plugin architecture | Linux/macOS only, same toolchain and dependency versions required on both sides; not a deployment boundary |
| Function-per-endpoint serverless | Cold starts are cheap in Go, but a monolith cut into 60 functions has no owner per table |
| Space-based / in-memory data grid | No Go-native grid is mainstream; the monolith has one `*gorm.DB` |
| Pipe-and-filter / ETL as the target | Batch pipelines are a different playbook |
| N-tier as the target | It is the shape being left behind |
| A shared "platform" library every service imports | Re-creates the monolith as a dependency; only `pkg/contracts` is shared |
| WASM (WASI) services | wasip2 support is maturing; not a 2026 production default |

## Go 2026 implementation

None generated.

## Anti-patterns

Treating an omitted style as an oversight and adding it without a named force.

## Interacts with

- [[flat-package]] · [[hexagonal]] · [[right-sizing]]
