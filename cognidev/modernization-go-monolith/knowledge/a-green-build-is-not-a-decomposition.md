---
id: a-green-build-is-not-a-decomposition
runtime: go
since: 0
category: decomposition-failure
tags: [hollow, port-report, cross-context, sagas, verification, go-work]
symptoms:
  - "cross_context_client_calls: 0"
  - "sagas: 0"
  - "clients_wired: 0"
  - "go build ./... exit 0"
remedy: { kind: investigate }
reference: "Carried over from the .NET twin (nopCommerce, five chained defects, every one exiting 0) — the failure is language-independent"
---

# `go build ./...` is green, and nothing was decomposed

## The signature

`go build ./...` and `go vet ./...` exit 0 in every module listed in `go.work`,
and in `.cognidev/feature/port-report.json`:

```json
"summary": { "cross_context_client_calls": 0, "clients_wired": 0, "sagas": 0 }
```

Sometimes with service names that are suspiciously structural — `internal`,
`cmd`, `pkg`, `platform` — rather than `catalog`, `orders`, `payments`.

## What it means

A green build proves the generated tree compiles. It proves nothing about the
cut, because the emptiest decomposition — every package copied into a module,
no boundary crossed, no call converted — also compiles, and Go compiles it
fast. Three shapes produce it, each exiting 0:

- **Services named after layout, not domain.** The clustering fell back to
  directories (`internal/`, `cmd/`), so the plan and the file map disagree and
  code lands in modules `go.work` never lists.
- **Every cross-context call was dropped.** The port removed the `catalog`
  import and the call with it, rather than turning it into a typed client.
- **No commit point was recognised.** A monolith that writes through a
  repository interface or raw `pgx` never literally says `db.Transaction`, and
  a saga detector keyed on that string finds nothing.

## What to do

Read these before the build result:

| Number | In | Zero on a real monolith means |
| --- | --- | --- |
| `cross_context_client_calls` | `port-report.json` | boundaries not crossed, or calls dropped |
| `clients_wired` | `port-report.json` | clients written, never constructed in `main.go` |
| `sagas` | `port-report.json` | no multi-context transaction recognised |

Check the service names against `cognidev/domain.md`: if they read like a
directory listing, the cut is structural and not a decomposition.
