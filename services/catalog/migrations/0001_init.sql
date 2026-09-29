CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    sku VARCHAR(40),
    name VARCHAR(200),
    price DOUBLE PRECISION,
    category_id BIGINT REFERENCES categories (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_products_sku ON products (sku);
