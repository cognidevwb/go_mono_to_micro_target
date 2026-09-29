# Runtime comparison — Go shop monolith vs the five services

Both systems were started on this machine and sent the same 36 requests, in the same order. The monolith is `go_mono_to_micro_source` (cmd/shop, one database). The target is this repository (gateway plus catalog, customers, inventory, payments and orders services, one database each, NATS for events), built and run exactly as generated.

## 1. Result

- **34 of 36 requests answer the same.** 1 differs only in error wording (same status), 1 is different: the target's new `/healthz`.
- **What each side leaves in its databases is the same: 4 of 4 tables** (orders, order lines, payments, stock).
- A declined payment releases the stock it reserved, and a failed order leaves no order behind, as the monolith's single transaction did (requests 19 and 25, and the state table).

## 2. Request by request

Bodies are shown after normalization (`id`, `orderId`, `createdAt` → `·`).

| # | Request | Monolith | Target | Result | Why |
|---|---|---|---|---|---|
| 1 | list products, empty<br>`GET /api/products` | `200 []` | `200 []` | same |  |
| 2 | create product A1 (price 10)<br>`POST /api/products` | `201 {"id":"\u00b7","sku":"A1","name":"Widget","price":10,"categoryId":1}` | `201 {"id":"\u00b7","sku":"A1","name":"Widget","price":10,"categoryId":1}` | same |  |
| 3 | create product B2 (price 250)<br>`POST /api/products` | `201 {"id":"\u00b7","sku":"B2","name":"Gadget","price":250,"categoryId":1}` | `201 {"id":"\u00b7","sku":"B2","name":"Gadget","price":250,"categoryId":1}` | same |  |
| 4 | create product, price 0<br>`POST /api/products` | `400 {"error":"Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'required' tag"}` | `400 {"error":"Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'required' tag"}` | same |  |
| 5 | create product, negative price<br>`POST /api/products` | `400 {"error":"Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'gt' tag"}` | `400 {"error":"Key: 'Product.Price' Error:Field validation for 'Price' failed on the 'gt' tag"}` | same |  |
| 6 | create product, name missing<br>`POST /api/products` | `400 {"error":"Key: 'Product.Name' Error:Field validation for 'Name' failed on the 'required' tag"}` | `400 {"error":"Key: 'Product.Name' Error:Field validation for 'Name' failed on the 'required' tag"}` | same |  |
| 7 | create product, duplicate SKU A1<br>`POST /api/products` | `422 {"error":"ERROR: duplicate key value violates unique constraint \"idx_products_sku\" (SQLSTATE 23505)"}` | `422 {"error":"ERROR: duplicate key value violates unique constraint \"idx_products_sku\" (SQLSTATE 23505)"}` | same |  |
| 8 | create product, malformed JSON<br>`POST /api/products` | `400 {"error":"unexpected EOF"}` | `400 {"error":"unexpected EOF"}` | same |  |
| 9 | list products, two<br>`GET /api/products` | `200 [{"id":"\u00b7","sku":"A1","name":"Widget","price":10,"categoryId":1},{"id":"\u00b7","sku":"B2","name":"Gadget","price":250,"categoryId":1}]` | `200 [{"id":"\u00b7","sku":"A1","name":"Widget","price":10,"categoryId":1},{"id":"\u00b7","sku":"B2","name":"Gadget","price":250,"categoryId":1}]` | same |  |
| 10 | register customer Ann<br>`POST /api/customers` | `201 {"id":"\u00b7","email":"ann@example.com","name":"Ann","active":true,"createdAt":"\u00b7"}` | `201 {"id":"\u00b7","email":"ann@example.com","name":"Ann","active":true,"createdAt":"\u00b7"}` | same |  |
| 11 | register customer Bob<br>`POST /api/customers` | `201 {"id":"\u00b7","email":"bob@example.com","name":"Bob","active":true,"createdAt":"\u00b7"}` | `201 {"id":"\u00b7","email":"bob@example.com","name":"Bob","active":true,"createdAt":"\u00b7"}` | same |  |
| 12 | register, bad email<br>`POST /api/customers` | `400 {"error":"Key: 'Customer.Email' Error:Field validation for 'Email' failed on the 'email' tag"}` | `400 {"error":"Key: 'Customer.Email' Error:Field validation for 'Email' failed on the 'email' tag"}` | same |  |
| 13 | register, duplicate email<br>`POST /api/customers` | `409 {"error":"ERROR: duplicate key value violates unique constraint \"idx_customers_email\" (SQLSTATE 23505)"}` | `409 {"error":"ERROR: duplicate key value violates unique constraint \"idx_customers_email\" (SQLSTATE 23505)"}` | same |  |
| 14 | get customer 1<br>`GET /api/customers/1` | `200 {"id":"\u00b7","email":"ann@example.com","name":"Ann","active":true,"createdAt":"\u00b7"}` | `200 {"id":"\u00b7","email":"ann@example.com","name":"Ann","active":true,"createdAt":"\u00b7"}` | same |  |
| 15 | get customer 999 (unknown)<br>`GET /api/customers/999` | `404 {"error":"not found"}` | `404 {"error":"not found"}` | same |  |
| 16 | get customer abc (bad id)<br>`GET /api/customers/abc` | `400 {"error":"bad id"}` | `400 {"error":"bad id"}` | same |  |
| 17 | place order: Ann, 2 x A1<br>`POST /api/orders` | `201 {"id":"\u00b7","customerId":1,"status":"placed","total":20,"lines":[{"id":"\u00b7","orderId":"\u00b7","productId":1,"quantity":2,"unitPrice":10}]}` | `201 {"id":"\u00b7","customerId":1,"status":"placed","total":20,"lines":[{"id":"\u00b7","orderId":"\u00b7","productId":1,"quantity":2,"unitPrice":10}]}` | same |  |
| 18 | my orders (Ann)<br>`GET /api/orders/mine` | `200 [{"id":"\u00b7","customerId":1,"status":"placed","total":20,"lines":[{"id":"\u00b7","orderId":"\u00b7","productId":1,"quantity":2,"unitPrice":10}]}]` | `200 [{"id":"\u00b7","customerId":1,"status":"placed","total":20,"lines":[{"id":"\u00b7","orderId":"\u00b7","productId":1,"quantity":2,"unitPrice":10}]}]` | same |  |
| 19 | place order: out of stock (50 x A1)<br>`POST /api/orders` | `409 {"error":"out of stock"}` | `409 {"error":"out of stock"}` | same |  |
| 20 | place order: inactive customer 2<br>`POST /api/orders` | `409 {"error":"customer inactive"}` | `409 {"error":"customer inactive"}` | same |  |
| 21 | place order: unknown customer 999<br>`POST /api/orders` | `409 {"error":"record not found"}` | `409 {"error":"record not found"}` | same |  |
| 22 | place order: unknown product 77<br>`POST /api/orders` | `409 {"error":"record not found"}` | `409 {"error":"stock item not found"}` | same status, different wording | Wording only: the target reserves stock first and inventory answers for the unknown product. Status is the same (409). |
| 23 | place order: no items<br>`POST /api/orders` | `400 {"error":"Key: 'PlaceOrderRequest.Items' Error:Field validation for 'Items' failed on the 'min' tag"}` | `400 {"error":"Key: 'PlaceOrderRequest.Items' Error:Field validation for 'Items' failed on the 'min' tag"}` | same |  |
| 24 | place order: quantity 0<br>`POST /api/orders` | `400 {"error":"Key: 'PlaceOrderRequest.Items[0].Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag"}` | `400 {"error":"Key: 'PlaceOrderRequest.Items[0].Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag"}` | same |  |
| 25 | place order: provider declines (4 x B2 = 1000)<br>`POST /api/orders` | `409 {"error":"provider declined: 402"}` | `409 {"error":"provider declined: 402"}` | same |  |
| 26 | place order: two lines (1 x A1, 1 x B2)<br>`POST /api/orders` | `201 {"id":"\u00b7","customerId":1,"status":"placed","total":260,"lines":[{"id":"\u00b7","orderId":"\u00b7","productId":1,"quantity":1,"unitPrice":10},{"id…` | `201 {"id":"\u00b7","customerId":1,"status":"placed","total":260,"lines":[{"id":"\u00b7","orderId":"\u00b7","productId":1,"quantity":1,"unitPrice":10},{"id…` | same |  |
| 27 | my orders (Ann), after the failures<br>`GET /api/orders/mine` | `200 [{"id":"\u00b7","customerId":1,"status":"placed","total":20,"lines":[{"id":"\u00b7","orderId":"\u00b7","productId":1,"quantity":2,"unitPrice":10}]},{"…` | `200 [{"id":"\u00b7","customerId":1,"status":"placed","total":20,"lines":[{"id":"\u00b7","orderId":"\u00b7","productId":1,"quantity":2,"unitPrice":10}]},{"…` | same |  |
| 28 | my orders (Bob), none<br>`GET /api/orders/mine` | `200 []` | `200 []` | same |  |
| 29 | list products, no X-User-ID<br>`GET /api/products` | `401 (empty)` | `401 (empty)` | same |  |
| 30 | place order, no X-User-ID<br>`POST /api/orders` | `401 (empty)` | `401 (empty)` | same |  |
| 31 | my orders, X-User-ID not a number<br>`GET /api/orders/mine` | `401 (empty)` | `401 (empty)` | same |  |
| 32 | stock of A1 (target route /v1/stockitem)<br>`GET /v1/stockitem/1` | `404 404 page not found` | `404 404 page not found` | same |  |
| 33 | charge (target route /v1/payment/charge)<br>`POST /v1/payment/charge` | `404 404 page not found` | `404 404 page not found` | same |  |
| 34 | gateway prefix /api/inventory<br>`GET /api/inventory/1` | `404 404 page not found` | `404 404 page not found` | same |  |
| 35 | gateway prefix /api/payments<br>`GET /api/payments/1` | `404 404 page not found` | `404 404 page not found` | same |  |
| 36 | health<br>`GET /healthz` | `401 (empty)` | `200 {"status":"ok"}` | different | Intended: the target adds /healthz on the gateway. The monolith has no health route, and its auth middleware answers 401 to anything without X-User-ID. |

## 3. What each side left in its databases

| Table | Monolith | Target | Result |
|---|---|---|---|
| orders | `[{"id": "\u00b7", "customer_id": 1, "status": "placed", "total": 20}, {"id": "\u00b7", "customer_id": 1, "status": "placed", "total": 260}]` | `[{"id": "\u00b7", "customer_id": 1, "status": "placed", "total": 20}, {"id": "\u00b7", "customer_id": 1, "status": "placed", "total": 260}]` | same |
| order_lines | `[{"id": "\u00b7", "order_id": 1, "product_id": 1, "quantity": 2, "unit_price": 10}, {"id": "\u00b7", "order_id": 3, "product_id": 1, "quantity": 1, "unit_price": 10}, {"id": "\u00b7", "order_id": 3, "product_id": 2, "quantity": 1, "unit_price": 250}]` | `[{"id": "\u00b7", "order_id": 1, "product_id": 1, "quantity": 2, "unit_price": 10}, {"id": "\u00b7", "order_id": 3, "product_id": 1, "quantity": 1, "unit_price": 10}, {"id": "\u00b7", "order_id": 3, "product_id": 2, "quantity": 1, "unit_price": 250}]` | same |
| payments | `[{"id": "\u00b7", "order_id": 1, "amount": 20, "status": "captured"}, {"id": "\u00b7", "order_id": 3, "amount": 260, "status": "captured"}]` | `[{"id": "\u00b7", "order_id": 1, "amount": 20, "status": "captured"}, {"id": "\u00b7", "order_id": 3, "amount": 260, "status": "captured"}]` | same |
| stock_items | `[{"product_id": 1, "on_hand": 5, "reserved": 3}, {"product_id": 2, "on_hand": 10, "reserved": 1}]` | `[{"product_id": 1, "on_hand": 5, "reserved": 3}, {"product_id": 2, "on_hand": 10, "reserved": 1}]` | same |

## 4. Defects found earlier, fixed in the playbook

An earlier target generated by this playbook matched 20 of 36 requests. Each defect the replay found was fixed where it was made, in the playbook, and this target was generated again from scratch.

| # | What happened | Fixed in |
|---|---|---|
| 1 | inventory-service exited at startup: the NATS consumer name contained a dot | scaffold-services-go events template |
| 2 | orders-service exited at startup: OrderLine was migrated before Order | port-monolith-go keeps the monolith's registration order |
| 3 | every call between services was refused (401): the monolith's X-User-ID check was copied onto internal routes | services accept internal calls through a shared token; the gateway strips caller headers from outside requests |
| 4 | payments-service served no route | port-monolith-go mounts every method another service calls |
| 5 | a declined payment left its stock reserved: the compensation only logged | the saga brief requires a compensation that calls the participant; go-gates fails a compensation that does not |
| 6 | payments-service had no provider URL in compose | scaffold-services-go carries the monolith's external settings |
| 7 | a failed order stayed visible as "cancelled" under my orders | the saga brief writes the order only in the final commit after payment succeeds |
| 8 | `go build` inside the workspace failed: go.work said go 1.26 while modules said 1.26.0 | scaffold-services-go writes the full release in go.work |

## 5. Setup and scope

- Seeding that neither front door offers was done in SQL on both sides, the same rows: category 1 (no category route), stock for A1 = 5 and B2 = 10 (no stock route on the monolith), and customer 2 set inactive (the register route cannot create an inactive customer).
- The external payment provider is a small fake, one per side, that accepts amounts under 1000 and declines the rest with 402.
- The gateway ran with `LEGACY_UPSTREAM` empty, so nothing fell back to the monolith, and with no `AUTH_ISSUER`, as in docker-compose.yml. Every service got the same `INTERNAL_TOKEN`, as in docker-compose.yml.
- Postgres and NATS ran on the host; host ports replace compose service names.
- Not covered: concurrency (the monolith's read-then-write reservation can oversell), the restock job (runs once a minute), NATS redelivery.
