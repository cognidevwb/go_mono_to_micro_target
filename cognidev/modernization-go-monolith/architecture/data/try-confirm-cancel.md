---
id: try-confirm-cancel
title: Try–Confirm–Cancel
family: data
currency: standard
role: reference
applies_when: []
why_not: TCC requires every participant to expose Try/Confirm/Cancel; this platform generates compensating sagas because a strangler extraction has participants (and a legacy monolith) that cannot be reshaped that way. Adopt it per participant where isolation is worth it.
teaches: Reserve, then confirm or cancel — isolation instead of compensation, if every participant can expose the three operations.
---

# Try–Confirm–Cancel

## What it is

Each participant exposes Try (reserve), Confirm, Cancel; the coordinator confirms all or cancels all.

## Currency (2026)

**STANDARD** in payments/booking.

## Fit signals (from `.cognidev/feature/context.md`)

- Participants that can hold reservations with timeouts.

## What it changes in the generated services

Not generated as a fleet pattern.

## Go 2026 implementation

`POST /v1/reservations` (try) → `POST /v1/reservations/{id}/confirm` | `DELETE /v1/reservations/{id}`

## Anti-patterns

- Reservations without expiry.

## Interacts with

- [[saga-orchestration]] · [[saga-isolation]]
