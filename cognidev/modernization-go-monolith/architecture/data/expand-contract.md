---
id: expand-contract
title: Expand–contract schema change
family: data
currency: standard
role: invariant
applies_when:
  - always
teaches: Two versions of the code always share the schema during a rollout — add, migrate, then remove, never rename in place.
---

# Expand–contract schema change

## What it is

Schema changes in three releases: expand (add the new column/table, write both), migrate (backfill, read new), contract (drop old). Required because rolling deploys run old and new replicas against one database.

## Currency (2026)

**STANDARD — invariant.**

## Fit signals (from `.cognidev/feature/context.md`)

- Any service with migrations and rolling updates.

## What it changes in the generated services

- Per-service `migrations/` with forward-only SQL (golang-migrate or goose); no destructive change in the same release as its code change.

## Go 2026 implementation

`migrations/0002_add_status_code.up.sql` adds; `0003_drop_status.up.sql` ships a release later.

## Anti-patterns

- Renaming a column in one migration.
- Down-migrations as a rollback strategy in production.

## Interacts with

- [[database-per-service]] · [[progressive-delivery]]
