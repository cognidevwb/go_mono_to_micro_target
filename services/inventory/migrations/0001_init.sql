-- Table inventory owns and nothing else: stock_items (StockItem). Replaces
-- the monolith's shared AutoMigrate; a column another service needs is
-- served by inventory's API, never by a cross-schema join.

CREATE TABLE IF NOT EXISTS stock_items (
    id         BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    on_hand    INTEGER NOT NULL DEFAULT 0,
    reserved   INTEGER NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_stock_items_product_id ON stock_items (product_id);
