CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- for gen_random_uuid()

-- users
CREATE TABLE users (
  id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
  full_name       TEXT         NOT NULL,
  email           TEXT         NOT NULL UNIQUE,
  password_hash   TEXT         NOT NULL,
  role            TEXT         NOT NULL DEFAULT 'customer'
                                 CHECK (role IN ('customer', 'admin')),
  created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- products (soft-delete via deleted_at)
CREATE TABLE products (
  id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
  name            TEXT          NOT NULL,
  description     TEXT,
  category        TEXT          NOT NULL,
  price           NUMERIC(10,2) NOT NULL CHECK (price >= 0),
  stock_quantity  INT           NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
  deleted_at      TIMESTAMPTZ,
  created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- one cart per user
CREATE TABLE carts (
  id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         UUID         NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- cart items — UNIQUE constraint enables upsert on quantity
CREATE TABLE cart_items (
  id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
  cart_id         UUID         NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
  product_id      UUID         NOT NULL REFERENCES products(id),
  quantity        INT          NOT NULL DEFAULT 1 CHECK (quantity > 0),
  created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
  UNIQUE (cart_id, product_id)
);

-- orders
CREATE TABLE orders (
  id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         UUID          NOT NULL REFERENCES users(id),
  total_amount    NUMERIC(10,2) NOT NULL CHECK (total_amount >= 0),
  status          TEXT          NOT NULL DEFAULT 'pending'
                                  CHECK (status IN ('pending', 'paid', 'shipped', 'cancelled')),
  created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- order items — price snapshot at checkout time
CREATE TABLE order_items (
  id              UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id        UUID          NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id      UUID          NOT NULL REFERENCES products(id),
  quantity        INT           NOT NULL CHECK (quantity > 0),
  unit_price      NUMERIC(10,2) NOT NULL, -- locked at checkout
  created_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- payments — one per order, Stripe intent tracked here
CREATE TABLE payments (
  id                       UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id                 UUID          NOT NULL UNIQUE REFERENCES orders(id),
  stripe_payment_intent_id TEXT          NOT NULL UNIQUE,
  amount                   NUMERIC(10,2) NOT NULL,
  status                   TEXT          NOT NULL DEFAULT 'pending'
                                          CHECK (status IN ('pending', 'succeeded', 'failed')),
  created_at               TIMESTAMPTZ   NOT NULL DEFAULT now(),
  updated_at               TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- indexes
CREATE INDEX ON cart_items (cart_id);
CREATE INDEX ON order_items (order_id);
CREATE INDEX ON orders (user_id);
CREATE INDEX ON products (category)
  WHERE deleted_at IS NULL; -- partial index: only active products