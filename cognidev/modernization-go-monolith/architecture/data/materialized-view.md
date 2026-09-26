---
id: materialized-view
title: Materialized view / read model
family: data
currency: rising
role: choice
applies_when:
  - answer:read_write_model = cqrs
teaches: When composition is too slow or too chatty, keep a local copy of the data you read, fed by the owner’s events.
---

# Materialized view / read model

## What it is

A consumer keeps its own table of another service's data it needs to read (product names and prices in orders), updated from that service's events.

## Currency (2026)

**RISING.**

## Fit signals (from `.cognidev/feature/context.md`)

- Hot cross-service reads; CQRS chosen.

## What it changes in the generated services

- Consumer table + event handler per replicated fact; staleness documented per field.

## Go 2026 implementation

`ON CONFLICT (product_id) DO UPDATE SET price = EXCLUDED.price WHERE product_prices.version < EXCLUDED.version`

## Anti-patterns

- Writing to the replica as if owned.
- No version check — out-of-order events overwrite newer data.

## Interacts with

- [[cqrs]] · [[idempotent-consumer]]
