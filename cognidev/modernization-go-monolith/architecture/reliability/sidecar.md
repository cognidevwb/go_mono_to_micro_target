---
id: sidecar
title: Sidecar / ambassador
family: reliability
currency: declining
role: reference
applies_when: []
why_not: Sidecarless (ambient / eBPF) meshes are displacing it, and nothing in this run scaffolds a sidecar.
teaches: Helpers deployed beside each pod — fading for meshes, still used for agents and Dapr.
---

# Sidecar / ambassador

## What it is

A companion container per pod handling cross-cutting concerns.

## Currency (2026)

**DECLINING** for meshes (ambient/eBPF).

## Fit signals (from `.cognidev/feature/context.md`)

- Legacy mesh estates.

## What it changes in the generated services

Not generated.

## Go 2026 implementation

None.

## Anti-patterns

- Sidecars for what a Go library does in-process.

## Interacts with

- [[service-mesh]] · [[dapr]]
