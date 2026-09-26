---
id: module-found-but-does-not-contain-package
runtime: go
since: 0
category: build-failure
tags: [modules, go-work, replace, module-path, shared-kernel]
symptoms:
  - "found, but does not contain package"
  - "no required module provides package"
  - "go: updates to go.mod needed"
remedy: { kind: build-file }
reference: "Go Modules Reference · go.dev/ref/mod#resolve-pkg-mod and #workspaces"
---

# The module resolves and the package inside it does not

## The signature

```
go: module github.com/acme/shop/pkg found (v0.0.0-00010101000000-000000000000, replaced by ../../pkg), but does not contain package github.com/acme/shop/pkg/contracts/catalog
```

or, when nothing was wired at all:

```
services/orders/internal/clients/catalog_client.go:9:2: no required module provides package github.com/acme/shop/pkg/contracts/catalog; to add it:
	go get github.com/acme/shop/pkg/contracts/catalog
```

## What it means

A multi-module cut has three places that must agree on one import path: the
`module` line in the provider's `go.mod`, the directory the package sits in
under it, and the consumer's `require` + `go.work use` (or `replace`) entry.
Porting moves directories; it does not move module paths. The first form means
the module resolved — through `go.work` or a `replace` — but the path under it
does not exist (a contract that was planned and never generated, or a directory
named differently from its import path). The second means the consumer's module
never learned the other module exists.

`go get` is the wrong fix in a workspace: it reaches for a network version of a
module that only exists on disk.

## What to do

- Check the three in order: `head -1 pkg/go.mod`, `ls pkg/contracts/`, and
  `go work edit -json | grep -A2 Use`.
- Every module the services share is listed in `go.work` with `go work use`;
  inside each service's `go.mod` a `require … v0.0.0` plus `replace => ../../pkg`
  keeps a single-service build (CI, `docker build`) working without the workspace.
- If the contract was planned but never generated, it is a scaffold gap, not a
  module problem — look for it in `scaffold-report.json`.
