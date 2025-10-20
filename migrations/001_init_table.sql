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

CREATE TABLE IF NOT EXISTS carts (
  id UUID PRIMARY KEY,
  customer_id UUID UNIQUE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE TABLE IF NOT EXISTS cart_items (
  id UUID PRIMARY KEY,
  cart_id UUID NOT NULL,
  product_id UUID NOT NULL,
  quantity INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (cart_id) REFERENCES carts(id),
  FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY,
  customer_id UUID NOT NULL,
  total_price BIGINT NOT NULL,
  total_quantity INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE TABLE IF NOT EXISTS order_items (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL,
  product_id UUID NOT NULL,
  quantity INT NOT NULL,
  price BIGINT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id),
  FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS payments (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL,
  payment_gateway VARCHAR(50) NOT NULL,
  payment_gateway_transaction_id VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS addresses (
  id UUID PRIMARY KEY,
  customer_id UUID NOT NULL,
  is_default BOOLEAN NOT NULL,
  street VARCHAR(100) NOT NULL,
  city VARCHAR(50) NOT NULL,
  state VARCHAR(20) NOT NULL,
  number VARCHAR(20) NOT NULL,
  zip_code VARCHAR(20) NOT NULL,
  address_line VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (customer_id) REFERENCES customers(id)
);