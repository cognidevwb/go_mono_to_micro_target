---
id: change-data-capture
title: Change data capture
family: data
currency: rising
role: choice
applies_when:
  - answer:migration_approach = parallel-run
teaches: Tail the database log to publish changes — for tables you cannot touch and for parallel runs.
---

# Change data capture

## What it is

Read Postgres logical replication (Debezium, or `github.com/jackc/pglogrepl` in Go) to turn committed row changes into events, without changing the writer.

## Currency (2026)

**RISING.**

## Fit signals (from `.cognidev/feature/context.md`)

- Parallel run; legacy monolith tables that must feed new services during strangling.

## What it changes in the generated services

- A CDC relay deployment reading the monolith's tables and publishing to NATS/Kafka.

## Go 2026 implementation

Debezium Server with a NATS JetStream sink, or a small Go relay on `pglogrepl`.

## Anti-patterns

- CDC of internal tables as a public contract.

## Interacts with

- [[transactional-outbox]] · [[parallel-run]]
