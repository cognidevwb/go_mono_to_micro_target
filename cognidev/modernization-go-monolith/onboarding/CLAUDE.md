# House rules (read by the generator on every seam it fills)

The tool generates the platform and moves your packages deterministically; when it fills a
`// CW-SEAM` (a cross-service call, a saga step, an event consumer), it obeys the rules below.
Write here what a strict senior Go reviewer on *your* team would enforce. Delete any prompt
you don't care about — what's left is treated as binding.

## Language & style
<!-- e.g. Go 1.26, gofmt + goimports, golangci-lint with our .golangci.yml; errors wrapped with
%w; no naked returns; context.Context is the first parameter of anything that does I/O. -->

## Naming
<!-- e.g. package names are short and singular (`order`, not `orders_pkg`); constructors are
`NewX`; interfaces are declared where they are CONSUMED; events are past-tense (`OrderPlaced`). -->

## Error handling
<!-- e.g. sentinel errors via errors.Is; HTTP handlers map domain errors to RFC 9457 problem
JSON; never panic across a request; no `_ = err`. -->

## Libraries — prefer / forbid
<!-- e.g. PREFER: stdlib net/http, slog, pgx, nats.go. FORBID: logrus (use slog), any
GPL/AGPL module, reflection-heavy DI containers. If a module needs a licence, flag it, don't add it. -->

## Testing
<!-- e.g. table-driven tests with testify/require; one testcontainers-go Postgres test per
service; assert the happy path and the one guard that matters (out-of-stock, 404, validation). -->

## Review checklist (the generator self-checks against this)
<!-- e.g. stock decrement is one guarded UPDATE (no read-check-then-write); status fields are
typed constants, not bare strings; no service imports another service's module; every
cross-context write goes through the outbox; every goroutine has a way to stop. -->
