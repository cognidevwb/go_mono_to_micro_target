-- inventory-service owned tables: stock_items (StockItem), plus the outbox,
-- processed-message and compensation-ledger tables its own code writes.

CREATE TABLE IF NOT EXISTS stock_items (
    id         BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    on_hand    BIGINT,
    reserved   BIGINT
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_stock_items_product_id ON stock_items (product_id);

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

CREATE TABLE IF NOT EXISTS inventory_compensations (
    step_key   TEXT        PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
