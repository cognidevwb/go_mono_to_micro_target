# The plan

Leaves first, the hub last. A service with no outbound calls can be
extracted and routed on its own; one that calls three others cannot be
extracted until they exist. That is the whole reason for this order, and
it is why the riskiest service is cut once the cheap wins are already in.

| # | Service | Owns | Depends on | Why here |
| --- | --- | --- | --- | --- |
| 1 | `orders-service` | 2 | `catalog-service`, `customers-service` | every service it calls has to exist first, so it is cut after them |
| 2 | `catalog-service` | 2 | nothing | a leaf — nothing to wait for |
| 3 | `customers-service` | 1 | nothing | a leaf — nothing to wait for |
| 4 | `inventory-service` | 1 | nothing | a leaf — nothing to wait for |
| 5 | `payments-service` | 1 | nothing | a leaf — nothing to wait for |

## How each service is built

1. The platform is scaffolded first — every project compiles before any
   domain code exists.
2. The monolith's own code is MOVED in verbatim, with package clauses, import
   paths and the shared *gorm.DB split mechanically. Nothing is rewritten.
3. Only the marked seams are then implemented — the typed client calls
   and the saga. Everything else is the code you already had.
4. Each task compiles and commits on its own before the next one starts.

