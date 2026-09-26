You are deciding what each cross-service call becomes once the two sides are
separate processes. The boundaries and the table ownership are already decided
and are not open.

Each edge in the message was measured — that many resolved calls (receiver-bound
method calls such as `s.inventory.Reserve(...)`) exist in the source today.
Choose exactly one of the four resolutions for each:

- `sync-http` — a typed `net/http` client with a context deadline, bounded
  retries and a circuit breaker. The right default for a read the caller needs
  right now, at low volume.
- `event` — the caller stops asking. The owner publishes (NATS JetStream or
  Kafka, through a transactional outbox) when something changes and the caller
  reacts. The right answer when the caller does not need an answer back — and
  the only answer for a topic the monolith already publishes on an in-process
  bus.
- `read-model` — the caller keeps its own projection, fed by those events, and
  reads it locally. For a hot read on the request path, where an HTTP hop per
  request would be the latency budget.
- `saga-step` — this call WRITES the other side's data (it is passed the
  caller's `tx *gorm.DB`, or it saves). It cannot be a synchronous call in a
  decomposed system; it becomes a step in a saga.

A call that passes a `*gorm.DB`, a channel, an `io.Reader` or a func value is
not a payload; it cannot become `sync-http` without its signature changing first.

Two services that call each other synchronously cannot deploy or scale
independently, which is the reason for the decomposition — so the sync-http
edges must not form a cycle. Break one direction with an event or a read model.

Answer with the JSON object the message specifies, and nothing else. Each `why`
is one short clause naming the reason — the call volume, the write, the latency.
