---
id: use-of-internal-package-not-allowed
runtime: go
since: 0
category: build-failure
tags: [modules, internal, go-work, shared-kernel, port]
symptoms:
  - "use of internal package"
  - "not allowed"
remedy: { kind: build-file }
reference: "Go spec for the go command · Internal Directories (go help internal / go.dev/doc/go1.4#internalpackages)"
---

# The ported package still imports the monolith's `internal/`

## The signature

```
services/orders/internal/orders/service.go:12:2: use of internal package github.com/acme/shop/internal/platform not allowed
```

## What it means

Go enforces `internal/` by **path**: a package under `a/b/internal/x` may be
imported only by code rooted at `a/b`. Inside the monolith every package lived
under `github.com/acme/shop`, so `orders` could import `internal/platform`
freely. After the port the orders code lives in its own module
(`github.com/acme/shop/services/orders`), whose root is no longer an ancestor of
`shop/internal/`, and the import that compiled for years is now illegal.

It usually names one of two things, and they need different fixes:

- **A cross-cutting helper** (`platform.Open`, `platform.Bus`, the auth
  middleware). It belongs in each service's own `internal/`, or — only when it is
  genuinely identical everywhere — in the shared kernel `pkg/`.
- **Another context's package** (`internal/catalog`). That is a cross-context
  call the port did not convert: it must become a typed client, never a copy.

## What to do

- Never "fix" it by moving the monolith's `internal/` to `pkg/` wholesale — that
  re-creates the shared-everything monolith as a library every service links.
- Helpers: copy into `services/<ctx>/internal/platform/` (it is a small, owned
  copy per service) or into `pkg/` when it is pure and stable.
- Another context: the seam is `// CW-SEAM[kind=cross-context-call]` — the call
  goes through `internal/clients/<ctx>_client.go`.
