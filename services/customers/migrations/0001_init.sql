-- Table customers owns and nothing else: customers (Customer). Replaces the
-- monolith's shared AutoMigrate; a column another service needs is served by
-- customers's API, never by a cross-schema join.

CREATE TABLE IF NOT EXISTS customers (
    id         BIGSERIAL PRIMARY KEY,
    email      VARCHAR(200) NOT NULL,
    name       VARCHAR(200) NOT NULL,
    active     BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_email ON customers (email);
