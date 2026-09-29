# What changes

## The services

Each service owns its tables outright. No other service may read them —
that constraint is what makes an independent deploy possible, and it is
also the whole cost of the cut.

| Service | Owns | Aggregate roots | Calls | Files planned |
| --- | --- | --- | --- | --- |
| `orders-service` | `Order`, `OrderLine` | `Order` | `catalog-service`, `customers-service` | 22 |
| `catalog-service` | `Category`, `Product` | `Category` | nothing | 13 |
| `customers-service` | `Customer` | `Customer` | nothing | 13 |
| `inventory-service` | `StockItem` | `StockItem` | nothing | 17 |
| `payments-service` | `Payment` | `Payment` | nothing | 16 |

## What becomes a network call

2 calls that cannot fail today become ones that can. Each one gets a typed,
resilient HttpClient with a timeout, a retry and a circuit breaker.

| From | To | Becomes |
| --- | --- | --- |
| `orders-service` | `catalog-service` | typed HttpClient + contract DTO |
| `orders-service` | `customers-service` | typed HttpClient + contract DTO |

## Cross-service writes

A transaction that spans services cannot be a database transaction.
Each one below becomes a saga with an explicit compensation per step.

| Saga | Orchestrator | Participants |
| --- | --- | --- |
| **CreateOrder** | `orders-service` | `inventory-service`, `orders-service`, `payments-service` |

## Files that belong to the fleet, not to one service

| File | Purpose |
| --- | --- |
| `ARCHITECTURE.md` | Docs: fill the <!-- CW-SEAM[kind=decomposition-rationale] --> in ARCHITECTURE.md — WHY these bounded contexts were drawn this way (which candidates merged/split, which coupling was accepted vs redesigned) |

