-- Tables catalog owns and nothing else: categories (Category), products
-- (Product). Replaces the monolith's shared AutoMigrate; a column another
-- service needs is served by catalog's API, never by a cross-schema join.

CREATE TABLE IF NOT EXISTS categories (
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL PRIMARY KEY,
    sku         VARCHAR(40) NOT NULL,
    name        VARCHAR(200) NOT NULL,
    price       DOUBLE PRECISION NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_products_sku ON products (sku);
