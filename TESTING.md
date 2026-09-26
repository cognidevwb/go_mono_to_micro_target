# Testing

`make test` runs the tests in every module. Each service ships a config test and a router test that builds with no live dependency; `pkg/` tests cover the saga's reverse compensation, idempotent consumption and readiness. The saga orchestrator's test pins the pivot step last.
