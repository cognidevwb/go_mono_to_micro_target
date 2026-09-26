---
id: event-sourcing
title: Event sourcing
family: data
currency: niche
role: choice
applies_when:
  - answer:read_write_model = event-sourced
teaches: Store the events, derive the state — powerful for audit-heavy domains, a heavy commitment for a port.
---

# Event sourcing

## What it is

The aggregate's state is the fold of its events, stored append-only; projections build read models.

## Currency (2026)

**NICHE.**

## Fit signals (from `.cognidev/feature/context.md`)

- Audit/regulatory domains; the user explicitly chose it.

## What it changes in the generated services

- An `events` table per aggregate with optimistic append on stream version; projectors. Not offered in the questionnaire until a live run proves it.

## Go 2026 implementation

`INSERT INTO events (stream_id, version, type, data) VALUES ($1,$2,$3,$4)` with a unique (stream_id, version).

## Anti-patterns

- Event sourcing a CRUD context during a decomposition.

## Interacts with

- [[cqrs]] · [[materialized-view]]
