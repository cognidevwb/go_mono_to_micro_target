CREATE TABLE IF NOT EXISTS orders (
    id          BIGSERIAL PRIMARY KEY,
    customer_id BIGINT,
    status      TEXT,
    total       DOUBLE PRECISION
);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders (customer_id);

CREATE TABLE IF NOT EXISTS order_lines (
    id         BIGSERIAL PRIMARY KEY,
    order_id   BIGINT,
    product_id BIGINT,
    quantity   BIGINT,
    unit_price DOUBLE PRECISION
);

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
