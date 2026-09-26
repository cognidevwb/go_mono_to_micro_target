# Contributing

1. `make sync build vet test` must pass.
2. `make hooks` installs the pre-commit hook (gofmt, no `.env`, no unfilled `CW-SEAM`).
3. A contract change is a provider change: run the contracts drift check.
