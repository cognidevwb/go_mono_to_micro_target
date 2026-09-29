# Integration report

## Met
- ✓ [01-catalog-core] No CW-SEAM or TODO(cognidev) marker in the shown catalog files; the shared-state price cache is now TTL-bound (30s)
- ✓ [02-catalog-docs] README.md has a real business-capability paragraph, no marker
- ✓ [03-catalog-cicd] deploy.targets.yml names staging/production and smoke_check_path /readyz
- ✓ [04-catalog-deployment] strangler.values.yaml has weight 50 and legacyPaths, no marker
- ✓ [05-customers-core] No CW-SEAM or TODO(cognidev) marker in the shown customers files
- ✓ [06-customers-docs] README.md has a real business-capability paragraph, no marker
- ✓ [07-customers-cicd] deploy.targets.yml names environments and /readyz
- ✓ [08-customers-deployment] strangler.values.yaml has weight 50 and legacyPaths, no marker

## Missing / broken

## Not shown — no verdict either way
- ? [01-catalog-core] go build/vet green (not run here; the compile note says the project compiles)
- ? [05-customers-core] go build/vet green
- ? [09-inventory-core] The shown inventory files have no markers, but scheduler.go, inventory_test.go and 0001_init.sql are not shown
- ? [09-inventory-core] go build/vet green
- ? [10-inventory-docs] README.md not shown
- ? [11-inventory-cicd] deploy.targets.yml not shown
- ? [12-inventory-deployment] strangler.values.yaml not shown
- ? [13-orders-clients] all orders client files not shown
- ? [14-orders-core] orders core files not shown
- ? [15-orders-saga] saga/create_order.go and its tests not shown, so compensation binding and go test cannot be judged
- ? [16-orders-docs] orders README not shown
- ? [17-orders-cicd] orders deploy.targets.yml not shown
- ? [18-orders-deployment] orders strangler.values.yaml not shown
- ? [19-payments-core] payments core files not shown apart from none
- ? [20-payments-docs] payments README not shown
- ? [21-payments-cicd] payments deploy.targets.yml not shown
- ? [22-payments-deployment] payments strangler.values.yaml not shown
- ? [23-docs-architecture] ARCHITECTURE.md not shown

## Notes
- The feature request text is empty, so the check is against the sign-off criteria only.
- catalog and customers remote_operations.go return err.Error() in a plain {"error":...} body. That leaks database error strings and is not RFC 9457 problem JSON, which breaks the guardrail. inventory's handlers hide the error text, though its remote_operations.go and compensation handler also use plain gin.H errors rather than problem+json.
- inventory's reserveHandler locks the row with SELECT ... FOR UPDATE, then calls the ported svc.Reserve. It cannot be confirmed here that Reserve uses the single guarded UPDATE with a RowsAffected check that the guardrail requires.
- inventory creates its compensation ledger table at runtime through `schemas`. The migration file is not shown, so it cannot be checked that the table is also in 0001_init.sql.
- The price cache has a TTL but no visible invalidation event, which the shared-state rule asks for on a cache.
- inventory wire.go guards Bus and SQL for a zero-value Deps{}, but the compensate and reserve mounts guard d.DB inside the handlers, which is fine. The signal.NotifyContext cancel is discarded, and there is no join on the jobs goroutine at shutdown.

## Coverage
26 of 63 delivered file(s) are shown below. These 37 were NOT shown and cannot be judged either way:
- services/inventory/internal/inventory/inventory_test.go
- services/inventory/internal/jobs/scheduler.go
- services/inventory/migrations/0001_init.sql
- services/inventory/test/equivalence/inventory_equivalence_test.go
- services/inventory/README.md
- services/inventory/deploy/deploy.targets.yml
- services/inventory/deploy/strangler.values.yaml
- services/orders/internal/clients/catalog/client.go
- services/orders/internal/clients/catalog_contract_test.go
- services/orders/internal/clients/customers/client.go
- services/orders/internal/clients/customers_contract_test.go
- services/orders/internal/clients/inventory/client.go
- services/orders/internal/clients/payments/client.go
- services/orders/internal/acl/legacy.go
- services/orders/internal/events/order_placed.go
- services/orders/internal/orders/orders_test.go
- services/orders/internal/orders/service.go
- services/orders/internal/outbox/outbox.go
- services/orders/internal/store/db.go
- services/orders/migrations/0001_init.sql
- services/orders/test/equivalence/orders_equivalence_test.go
- services/orders/internal/saga/create_order.go
- services/orders/README.md
- services/orders/deploy/deploy.targets.yml
- services/orders/deploy/strangler.values.yaml
- services/payments/internal/acl/legacy.go
- services/payments/internal/app/create_order_compensation.go
- services/payments/internal/app/remote_operations.go
- services/payments/internal/httpapi/payment_handlers.go
- services/payments/internal/payments/payments_test.go
- services/payments/internal/store/db.go
- services/payments/migrations/0001_init.sql
- services/payments/test/equivalence/payments_equivalence_test.go
- services/payments/README.md
- services/payments/deploy/deploy.targets.yml
- services/payments/deploy/strangler.values.yaml
- ARCHITECTURE.md
