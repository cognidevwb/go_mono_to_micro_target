# What this folder is

The other kind of fact.

Everything else this playbook knows is derived from your repository — which
packages exist, which structs each `*gorm.DB` migrates, which calls cross a
package boundary, which `db.Transaction` spans more than one context. This
folder holds what reading a repository cannot reveal: what goes wrong when a Go
monolith is cut into modules and services, and how each failure announces
itself — in `go build`, `go vet`, `go test`, `go work`, `govulncheck`, or the
running binary's log.

## Every document is keyed on what the toolchain prints

Not on a topic. A reader arrives here holding a message from the compiler, the
module loader or a crashed service, and the fastest possible answer is the
document whose `## The signature` block contains that exact text.

So each file opens with the literal output, and the front matter repeats the
searchable fragments in `symptoms:` for the indexer:

```yaml
---
id: the-document-id
runtime: go
since: 0              # Go 1.x MINOR this landed in; 0 = applies to every target
category: build-failure   # build-failure · decomposition-failure · correctness · runtime-failure · cutover-failure
tags: [modules, go-work]
symptoms: ["a fragment the toolchain prints"]
remedy: { kind: build-file }   # build-file · seam · environment · investigate
reference: "the primary source this is checked against"
---
```

**This page must not quote any signature.** An index that reproduces the
searchable text outranks the documents it points at, and every query lands here
instead of on the answer.

## `since`

`since: 0` means the fact is a property of the decomposition rather than of a
Go release, and applies whatever toolchain the services land on. A non-zero
value is the MINOR of a Go 1.x release (`22` = Go 1.22), and the fact is emitted
only when the run targets that release or later. The target is the
`go_version` answer, else the highest `go`/`toolchain` directive in the
monolith's `go.mod`, else Go 1.26.

## Adding one

Write the markdown file, give it a signature block and a citation, then rebuild
the index:

```
index-go-knowledge          # rewrites index.json from the .md files
```

A card with no `reference:` is refused — an uncited fact is folklore. A document
with no `## The signature` block loses retrieval to a document full of generic
terms, however much better it is.
