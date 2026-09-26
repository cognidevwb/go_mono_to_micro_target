---
id: go-work-requires-a-newer-go
runtime: go
since: 21
category: build-failure
tags: [go-work, toolchain, go-directive, gotoolchain]
symptoms:
  - "listed in go.work file requires go >="
  - "but go.work lists go"
  - "go: downloading go1."
  - "toolchain not available"
remedy: { kind: build-file }
reference: "Go 1.21 release notes · Go toolchains (go.dev/doc/toolchain) — the go line is a minimum requirement"
---

# One service moved to a newer Go and the workspace did not

## The signature

```
go: module ./services/orders listed in go.work file requires go >= 1.26, but go.work lists go 1.24; to update it:
	go work use
```

or, on a CI runner pinned with `GOTOOLCHAIN=local`:

```
go: go.mod requires go >= 1.26 (running go 1.24.6; GOTOOLCHAIN=local)
```

## What it means

Since Go 1.21 the `go` line is a hard minimum, not advice, and the workspace's
`go` line must be at least every member module's. A decomposition multiplies
the places that state a version — one `go.mod` per service, the gateway, the
shared kernel, `go.work`, every Dockerfile's `golang:` base image, every CI
`setup-go` step — and they drift the first time one service is bumped alone.

With the default `GOTOOLCHAIN=auto` the go command silently downloads the newer
toolchain instead of failing, which is why a laptop builds and the pinned CI
runner does not.

## What to do

- One version, stated once: the `go_version` answer. `go.work`, every `go.mod`,
  every `golang:1.26` image tag and every `go-version:` in CI carry the same
  value; `go work use -r .` then `go work sync` after any bump.
- Keep the monolith's own `go` line out of it — the before-state builds on what
  it declared; only the services move.
- If CI must stay offline, set `GOTOOLCHAIN=local` everywhere (not just in CI) so
  the mismatch is a failure on the laptop too.
