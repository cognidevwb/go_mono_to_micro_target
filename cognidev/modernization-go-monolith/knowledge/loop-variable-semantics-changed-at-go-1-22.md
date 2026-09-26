---
id: loop-variable-semantics-changed-at-go-1-22
runtime: go
since: 22
category: correctness
tags: [loopvar, go-directive, semantics, goroutines, upgrade]
symptoms:
  - "loop variable"
  - "captured by func literal"
  - "loopclosure"
remedy: { kind: investigate }
reference: "Go 1.22 release notes · per-iteration loop variables (go.dev/blog/loopvar-preview), gated on the module's go line"
---

# The service module's `go` line changed what a `for` loop means

## The signature

`go vet` in the monolith (go line < 1.22):

```
internal/orders/service.go:48:15: loop variable item captured by func literal
```

and after the port the warning disappears — because behaviour changed, not
because the code was fixed.

## What it means

Loop variable semantics are decided by the `go` line of the module that
declares the loop. The monolith at `go 1.21` shares one variable across every
iteration; the same file ported into a service module at `go 1.26` gets a fresh
variable per iteration. Code that captured the loop variable in a goroutine or
closure now behaves differently — usually more correctly, occasionally in a way
a test pinned.

A decomposition that also raises the `go` line is two changes at once, and a
behaviour difference afterwards has two possible causes.

## What to do

- Run the baseline `go vet` and keep its `loopclosure` findings; each one is a
  place where the ported service will behave differently.
- `bisect -godebug loopvar=...` / `GOEXPERIMENT` is not needed — compare the
  monolith's test run to the service's for the packages those findings name.
