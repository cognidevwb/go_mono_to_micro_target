---
id: backend-for-frontend
title: Backend for frontend
family: communication
currency: rising
role: reference
applies_when: []
why_not: A single Go gateway is what this run scaffolds; per-client backends are a later refinement you add by hand.
teaches: One backend per client type when their needs diverge; the token-handling BFF keeps tokens out of the browser.
---

# Backend for frontend

## What it is

A dedicated edge per client (web, mobile, partner) shaping responses for that client, and — in the security variant — holding OAuth tokens server-side so the SPA only has a cookie.

## Currency (2026)

**RISING** (security variant), per the IETF browser-based-apps BCP.

## Fit signals (from `.cognidev/feature/context.md`)

- Several client types with divergent payloads; an SPA holding tokens in the browser.

## What it changes in the generated services

Not generated. A second gateway module per client plus a client list question is the next slice.

## Go 2026 implementation

A Go BFF is the same `ReverseProxy` gateway with session cookies (`gorilla/securecookie` or `alexedwards/scs`) and server-side token exchange via `golang.org/x/oauth2`.

## Anti-patterns

- A BFF per team rather than per client.
- Shared BFF that becomes a new monolith.

## Interacts with

- [[api-gateway]] · [[security]]
