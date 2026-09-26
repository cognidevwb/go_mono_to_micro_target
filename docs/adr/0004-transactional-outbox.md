# 4. Transactional outbox on nats

Status: accepted

Events are written to the `outbox` table in the same transaction as the state change (`pkg/outbox.Enqueue`) and published by a relay. Consumers are idempotent (`pkg/events.Idempotent` over `processed_messages`). Delivery is at-least-once.
