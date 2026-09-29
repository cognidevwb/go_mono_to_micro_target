-- payments owns the payments table, plus the saga compensation ledger.
CREATE TABLE IF NOT EXISTS payments (
    id       BIGSERIAL PRIMARY KEY,
    order_id BIGINT,
    amount   DOUBLE PRECISION,
    status   TEXT
);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments (order_id);

CREATE TABLE IF NOT EXISTS payment_compensations (
    step_key   TEXT        PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
