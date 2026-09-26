# State flow

Which service owns which table, and which writes stopped being one transaction.

| Table | Owner |
|---|---|
| Category | catalog-service |
| Product | catalog-service |
| Customer | customers-service |
| StockItem | inventory-service |
| Payment | payments-service |
| Order | orders-service |
| OrderLine | orders-service |

## Writes that are no longer atomic

- `CreateOrder` spanned orders and inventory, payments in one database transaction. It is now a saga; each step commits alone and is compensated on failure.
