---
id: read-replica
title: Read replica
family: data
currency: standard
role: choice
applies_when:
  - read_heavy_services >= 1
  - answer:load_profile = high-volume
teaches: Same model, another server — route read-only queries of read-heavy services to a replica and accept its lag.
---

# Read replica

## What it is

A second connection string for queries that tolerate replication lag.

## Currency (2026)

**STANDARD.**

## Fit signals (from `.cognidev/feature/context.md`)

- Read-heavy services and high-volume load.

## What it changes in the generated services

- `DATABASE_READ_URL`; a second pool used by query handlers only.

## Go 2026 implementation

`readPool, _ := pgxpool.New(ctx, cfg.DatabaseReadURL)`

## Anti-patterns

- Reading your own write from the replica.

## Interacts with

- [[read-through-cache]] · [[cqrs]]
