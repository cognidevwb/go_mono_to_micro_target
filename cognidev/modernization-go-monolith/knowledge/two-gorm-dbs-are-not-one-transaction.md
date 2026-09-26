---
id: two-gorm-dbs-are-not-one-transaction
runtime: go
since: 0
category: correctness
tags: [gorm, pgx, transactions, saga, outbox, database-per-service, consistency]
symptoms:
  - "sql: transaction has already been committed or rolled back"
  - "tx *gorm.DB"
  - "db.Transaction(func(tx *gorm.DB) error"
  - "pgx: tx is closed"
remedy: { kind: seam }
reference: "database/sql Tx semantics (pkg.go.dev/database/sql#Tx) · a Tx is bound to one connection of one database"
---

# The `db.Transaction` the monolith had does not survive the cut

## The signature

The loud version, where a `*gorm.DB` transaction handle was carried across a
boundary that now runs on another connection:

```
sql: transaction has already been committed or rolled back
```

The dangerous version prints nothing at all. In the monolith:

```go
err = s.db.Transaction(func(tx *gorm.DB) error {
    ok, err := s.inventory.Reserve(tx, item.ProductID, item.Quantity)
    ...
    return s.payments.Charge(tx, order.ID, order.Total)
})
```

`Reserve` and `Charge` took `tx *gorm.DB` as a parameter — the transaction is
part of the cross-context signature.

## What it means

In the monolith one `db.Transaction` committed a stock reservation, an order and
a payment together, and Postgres made that atomic for free. Database per service
removes the shared connection: a `*gorm.DB`, `*sql.Tx` or `pgx.Tx` cannot be
sent over HTTP, and Go has no distributed transaction coordinator to put them
back together. A signature carrying `tx` is therefore not portable — the call
cannot become a network hop as written.

Two service calls in sequence, both succeeding, are not one transaction: the
process can stop between them and the second write is lost. A test suite that
never kills the process mid-flow will not notice.

## What to do

- One local transaction per service per step. The cross-context writes become a
  saga (`internal/saga/`) with a compensating action per step — reserve stock,
  charge payment, confirm order; on failure release the reservation.
- Publish each step's event in the SAME local transaction via the outbox table;
  a publish after `Commit()` is lost on a crash between them.
- Drop `tx *gorm.DB` from every signature that crosses a boundary. Check `sagas`
  in `port-report.json` against the multi-context transactions the map found.
