CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(200),
    name TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_customers_email ON customers (email);
