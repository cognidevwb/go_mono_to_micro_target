---
id: the-cut-kept-the-code-and-dropped-the-pipelines
runtime: go
since: 0
category: cutover-failure
tags: [cutover, ci, cd, keep-set, github-actions, deploy, go-work]
symptoms:
  - "no such file or directory: .github/workflows"
  - "CI is missing after cutover"
  - "go: cannot load module pkg listed in go.work file"
remedy: { kind: build-file }
reference: "cutover-go keep_set — an allow-list, not a deny-list (same rule as cutover-dotnet)"
---

# The scaffold wrote it, the cutover threw it away

## The signature

Everything is green, the services are on the new branch, and a top-level entry
the scaffold definitely produced is simply not there — most often `.github/`,
`deploy/`, `pkg/`, `docs/` or `go.work` itself. When it is `pkg/` the first
build on the new branch prints:

```
go: cannot load module pkg listed in go.work file: open pkg/go.mod: no such file or directory
```

## What it means

The cutover builds the new branch from an **allow-list** of top-level entries.
Every entry it does not name is dropped — which is what makes the branch clean,
and what makes any new top-level scaffold output invisible unless someone adds
it. It is never a scaffold bug, and looking for it there costs the most time.

## What to do

Any new top-level output (a directory or `go.work`) is added to the cutover's
`keep_set` in the same change that introduces it. Verify on the cut branch with
`go work sync && go build ./...` from the root, not on the working tree — the
working tree still holds both worlds and looks correct until the cutover runs.
