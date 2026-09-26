# Migration log

What moved out of the monolith, in the order it moved. One entry per task the run executed.

This is the chronology; `MIGRATION-SCORECARD.md` is the grade and `VALIDATION-REPORT.html` is the verdict. A task that never ran has no entry here — the gap between the 23 planned and the 23 below is the record that it did not.

| # | Task | Outcome | Touched |
|---|---|---|---|
| 1 | Implement catalog-service — domain + wiring | COMPLETED | 6 file(s) |
| 2 | Docs: catalog | COMPLETED | 1 file(s) |
| 3 | CI/CD: catalog | COMPLETED | 1 file(s) |
| 4 | Deployment: catalog | COMPLETED | 1 file(s) |
| 5 | Implement customers-service — domain + wiring | COMPLETED | 5 file(s) |
| 6 | Docs: customers | COMPLETED | 1 file(s) |
| 7 | CI/CD: customers | COMPLETED | 1 file(s) |
| 8 | Deployment: customers | COMPLETED | 1 file(s) |
| 9 | Implement inventory-service — domain + wiring | COMPLETED | 9 file(s) |
| 10 | Docs: inventory | COMPLETED | 1 file(s) |
| 11 | CI/CD: inventory | COMPLETED | 1 file(s) |
| 12 | Deployment: inventory | COMPLETED | 1 file(s) |
| 13 | Implement orders-service — cross-service calls | COMPLETED | 6 file(s) |
| 14 | Implement orders-service — domain + wiring | COMPLETED | 8 file(s) |
| 15 | Implement orders-service — orchestration saga | COMPLETED | 1 file(s) |
| 16 | Docs: orders | COMPLETED | 1 file(s) |
| 17 | CI/CD: orders | COMPLETED | 1 file(s) |
| 18 | Deployment: orders | COMPLETED | 1 file(s) |
| 19 | Implement payments-service — domain + wiring | COMPLETED | 6 file(s) |
| 20 | Docs: payments | COMPLETED | 1 file(s) |
| 21 | CI/CD: payments | COMPLETED | 1 file(s) |
| 22 | Deployment: payments | COMPLETED | 1 file(s) |
| 23 | Docs: architecture | COMPLETED | 1 file(s) |

## Services this run produced

- **catalog-service** — `services/catalog-service`
- **customers-service** — `services/customers-service`
- **inventory-service** — `services/inventory-service`
- **payments-service** — `services/payments-service`
- **orders-service** — `services/orders-service`

Each entry links to its own record under `cognidev/docs/status/`, beside the plan it was written against in `cognidev/docs/tasks/`. Holding the two side by side is how intent and outcome are compared; this file only says what order they happened in.

