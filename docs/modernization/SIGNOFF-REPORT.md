# Migration sign-off

**SIGNED OFF** — every blocking clause passed.

## The predicate

| clause | required | actual | |
|---|---|---|---|
| dropped_unexplained == 0 | 0 | 0 | pass |
| regressed == 0 | 0 | 0 | pass |
| unclaimed_target_behavior == 0 | 0 | 0 | pass |
| proven_fraction >= policy.min_proven | 50.0% | 74.5% | pass |
| every waiver has {author, reason, scope, timestamp} | 0 incomplete | 0 | pass |
| ledger.source_merkle == current source merkle | not stale | current | pass |

## The ledger

59 source behaviour-bearing unit(s), each with exactly one disposition.

| disposition | units | blocks |
|---|---:|---|
| `PORTED_PROVEN` | 39 | no |
| `PORTED_EQUIVALENT` | 5 | no |
| `PORTED_UNPROVEN` | 15 | no |
| `TRANSFORMED` | 0 | no |
| `DROPPED_INTENTIONAL` | 0 | no |
| `DROPPED_UNEXPLAINED` | 0 | yes |
| `REGRESSED` | 0 | yes |
| **total** | **59** | |

## Proof-method coverage

How each row was decided — the coverage of the *proof method*, not of the code.

| method | units | share |
|---|---:|---:|
| fingerprint | 49 | 83.0% |
| transform rule | 5 | 8.4% |
| cluster adjudication | 5 | 8.4% |

> 59 source unit(s) · 44 carry positive evidence (74.5%) · 83.0% by fingerprint · 8.4% by transform rule · 8.4% by cluster adjudication.

**No row carries behavioural evidence.** Every proven unit below is proven *structurally*.

`PORTED_PROVEN` by fingerprint means **structural** confidence: fingerprint equality is a strong necessary condition for equivalence, cheap and deterministic, but it is not a proof of semantic equivalence. `pin replay` is different in kind — a characterization test of the source was run against the target and passed, which is observed behaviour. A full runtime differential (booting the frozen source as an oracle) is stronger still and this cycle does not run one, so no row claims `runtime differential`.

### Unproven (15) — listed, not hidden

- **36 field(s) conserved exactly but 4 unknown on one side — partial evidence is not proof** — 4 unit(s)
  - `internal/platform/bus.go#Bus`
  - `internal/platform/bus.go#Bus.Publish`
  - `internal/platform/bus.go#Bus.Subscribe`
  - `internal/platform/bus.go#NewBus`
- **7 field(s) conserved exactly but 1 unknown on one side — partial evidence is not proof** — 6 unit(s)
  - `internal/payments/gateway.go#Gateway`
  - `internal/payments/gateway.go#Gateway.Charge`
  - `internal/payments/gateway.go#Gateway.NewGateway`
  - `internal/payments/service.go#Service.Charge`
  - `internal/payments/service.go#Service.NewService`
  - `internal/payments/service.go#payments.Service`
- **cluster C2078bcd8f7d27f80 adjudicated as an extractor blind spot — not comparable** — 5 unit(s)
  - `cmd/shop/main.go#main`
  - `internal/orders/service.go#Service.ForCustomer`
  - `internal/orders/service.go#Service.NewService`
  - `internal/orders/service.go#Service.PlaceOrder`
  - `internal/orders/service.go#orders.Service`

## Unclaimed target code (70 total, 0 carrying behaviour)

Target code with no source ancestor. Some is legitimate — scaffolding, framework glue, the new saga. Some is behaviour nobody asked for. Both are named; only the ones carrying behaviour (an endpoint, a write, a branch, an external call) block.


## Waivers (0)

None.

## Rank-pruned clusters (0)

None — every cluster was carried through.

## Money / decimal invariant (0)

No arithmetic drift in type, scale or rounding was detected on a money path.

## Correspondence quality

46 pair(s), 0 of them (0.0%) anchored only by the name canonicalizer (tier 3, heuristic).

## Cluster adjudications (15)

Verdicts a model reached on a cluster the deterministic layers had already isolated. Each one is recorded and replayable; none of them is the sole basis for a blocking verdict without a waiver.

| cluster | verdict | propagated to | rationale |
|---|---|---:|---|
| `Cf059feb3af56b4df` | extractor-blind-spot | 4 finding(s) | All four members are unanchored target-only scaffolding files (pre-commit hook, Bicep landing zone, .env.example) that scaffold-report.json declared generated,  |
| `Cf93ff891a5ca551d` | extractor-blind-spot | 4 finding(s) | All four members are target-only Go functions (NewHealth, NewMemorySeen, Routes) in files scaffold-report.json declared generated, so they are unanchored by con |
| `C0c0fb0ed9c0f75c0` | extractor-blind-spot | 13 finding(s) | All three exemplars are bare Go type declarations (`type Memory struct {`, `type Health struct {`, `type SQLSeen struct{ DB *sql.DB }`) in files scaffold-report |
| `Cf47228bf0162fe16` | extractor-blind-spot | 49 finding(s) | All three exemplars are target-only units (internalToken, a proxy header-stripping test, MemorySeen.MarkSeen) with no source unit, and every source span is unre |
| `Ce6f76633e78e285d` | extractor-blind-spot | 1 finding(s) | The exemplar pairs source `cmd/shop/main.go#main` (monolith wiring) with a target span that shows `orders/service.go#Service.PlaceOrder`, not the target `catalo |
| `Cab27fa1b441f9686` | extractor-blind-spot | 1 finding(s) | The target span shown is only the 6-line main() wrapper that calls run(); the flagged catch of gorm.ErrRecordNotFound (e:catch:gormerrrecordnotfound) is not in  |
| `C979ef9ce3175cb0d` | extractor-blind-spot | 1 finding(s) | The finding pairs the source `main` (monolith wiring, no transaction) with a target `main (+56 more)` that aggregates many units, yet the target span shown is ` |
| `C7fb6324ecbb6ef01` | extractor-blind-spot | 1 finding(s) | The finding claims a target-only `ar:mul:int:?:unknown` arithmetic op, but the shown target span (a routes_test.go test) contains no multiplication at all, and  |
| `C6cb7e277f2e478ba` | extractor-blind-spot | 1 finding(s) | The source unit is a monolith main() and the target span shown is an unrelated init() in the inventory service's compensation file, so the pairing compares a mo |
| `C5e67bbc580028999` | extractor-blind-spot | 1 finding(s) | The source unit is the monolith's `main` (wiring all contexts: catalog, customers, inventory, payments, orders), paired with only the catalog service's `run` (s |
| `C568df80e5038d8aa` | extractor-blind-spot | 1 finding(s) | The source `main` is a monolith wiring function, but the target span shown is `remote_operations.go#init` (a route mount), not the named target `main.go#main (+ |
| `C49b2110f8845602f` | extractor-blind-spot | 1 finding(s) | The exemplar's source and target ForCustomer spans are byte-for-byte identical (Preload/Where/Find, return out, err), so the flagged e:catch:error and e:catch:s |
| `C2ccf4b7c847a46d5` | extractor-blind-spot | 1 finding(s) | The spans do not line up: the source span is `Gateway.Charge` in payments/gateway.go, while the target span is `Service.PlaceOrder` in orders/service.go, and th |
| `C25f4debceb5bf93d` | extractor-blind-spot | 1 finding(s) | The finding pairs source unit ForCustomer (+3 more) with target ForCustomer (+4 more), but the spans shown are Gateway.Charge (source) and Service.releaseStock  |
| `C2078bcd8f7d27f80` | extractor-blind-spot | 2 finding(s) | The cluster is heterogeneous and the spans do not line up: Exemplar 1's SOURCE span is payments Gateway.Charge while the TARGET span is orders PlaceOrder, so no |

## Attestation

Unattested. A computed verdict is not a sign-off until a named human accepts it — set `signoff_attestation` in the intake, or add the name to this file in review.

Source Merkle: `191424716611929f67cc1862dda5ca22ac6dd8a1bc5675ac08334d4c0b03fce6`  ·  shards `../go_mono_to_micro_source/.cognidev/understand`
Target Merkle: `5ca5c3120db67b5dbc2fcb4ec113403ff3d51db72f73c0e29895c941712bbae3`  ·  shards `./.cognidev/understand`

Regenerate with `signoff-verdict`. If the source moves, this verdict invalidates itself rather than quietly aging.
