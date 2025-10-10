CREATE TABLE IF NOT EXISTS customers (
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  email VARCHAR(50) UNIQUE NOT NULL,
  password VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,
  status VARCHAR(50) NOT NULL,
  name VARCHAR(50) NOT NULL,
  description TEXT,
  price BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS inventories (
  id UUID PRIMARY KEY,
  product_id UUID NOT NULL,
  stock_quantity INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (product_id) REFERENCES products(id)
);