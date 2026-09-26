---
id: govulncheck-fails-the-pipeline
runtime: go
since: 0
category: build-failure
tags: [govulncheck, cve, advisory, ci, security, grype]
symptoms:
  - "Your code is affected by"
  - "Vulnerability #1: GO-"
  - "Found in: "
  - "Fixed in: "
remedy: { kind: build-file }
reference: "govulncheck (pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) · exits 3 when a called vulnerable symbol is found"
---

# The advisory the monolith carried now fails every service's pipeline

## The signature

```
=== Symbol Results ===

Vulnerability #1: GO-2025-3553
    Excessive memory allocation during header parsing in github.com/golang-jwt/jwt
  Found in: github.com/golang-jwt/jwt/v5@v5.2.1
  Fixed in: github.com/golang-jwt/jwt/v5@v5.2.2

Your code is affected by 1 vulnerability from 1 module.
```

and the CI step exits `3`.

## What it means

The decomposition does not choose dependency versions — it inherits them from
the monolith's `go.sum` and then multiplies them: one vulnerable module becomes
five services each shipping it. The monolith may never have run `govulncheck`,
so the generated per-service CI is the first thing ever to look.

`govulncheck` reports only vulnerabilities in code that is actually reached, so
the same advisory can fail one service and not another — that difference is
real, not flakiness.

## What to do

- Treat it as a pre-existing condition, not something the cut introduced — it
  is reported as such in the gate, never charged to the decomposition.
- Upgrade in the monolith's baseline `go.mod` once (`go get <module>@<fixed>`),
  then `go work sync` so every service moves together.
- Container images are scanned separately by Syft + Grype; a clean `govulncheck`
  does not cover the base image.
