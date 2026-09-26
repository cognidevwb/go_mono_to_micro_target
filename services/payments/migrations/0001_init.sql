-- Tables payments owns and nothing else: payments (Payment), plus the outbox
-- and processed_messages tables payments' own code writes. Replaces the
-- monolith's shared AutoMigrate; a column another service needs is served by
-- payments's API, never by a cross-schema join.

CREATE TABLE IF NOT EXISTS payments (
    id       BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    amount   DOUBLE PRECISION NOT NULL,
    status   VARCHAR(20) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments (order_id);

CREATE TABLE IF NOT EXISTS outbox (
    id           BIGSERIAL PRIMARY KEY,
    subject      TEXT        NOT NULL,
    envelope     JSONB       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS processed_messages (
    consumer   TEXT        NOT NULL,
    message_id TEXT        NOT NULL,
    seen_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (consumer, message_id)
);
