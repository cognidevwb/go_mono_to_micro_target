---
id: import-cycle-not-allowed
runtime: go
since: 0
category: build-failure
tags: [packages, cycles, bounded-context, shared-kernel]
symptoms:
  - "import cycle not allowed"
  - "package github.com/"
  - "imports github.com/"
remedy: { kind: seam }
reference: "The Go Programming Language Specification · Import declarations (an import cycle is a compile error)"
---

# The cut split a package pair that imported each other through a third

## The signature

```
package github.com/acme/shop/services/orders/internal/orders
	imports github.com/acme/shop/pkg/contracts/inventory
	imports github.com/acme/shop/services/orders/internal/orders: import cycle not allowed
```

## What it means

Go forbids import cycles outright, so a monolith that compiles has none — and a
decomposition can still create one. The usual shape is a "shared" contract or
helper moved into `pkg/` that references a type which stayed with a service,
while that service imports the contract. The package graph that was a clean DAG
inside one module becomes a cycle across two.

The second shape is an event handler: `orders` publishes `OrderPlaced`,
`inventory` subscribes, and the event struct was declared in `orders`. When
`inventory` imports `orders` for the struct and `orders` imports `inventory` for
`Reserve`, the monolith was already cyclic in intent and only compiled because
`platform.Bus` took `any`.

## What to do

- The shared kernel (`pkg/`) imports NOTHING from `services/`. A contract type is
  a plain struct with JSON tags, declared once in `pkg/contracts/<provider>/`.
- An event's payload belongs to the publisher's contract package, never to its
  domain package.
- A cycle between two services is a boundary decision, not a build error: the
  pair is either one service, or one direction becomes an event.
