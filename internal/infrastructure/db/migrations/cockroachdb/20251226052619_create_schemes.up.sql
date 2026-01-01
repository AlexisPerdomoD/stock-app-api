CREATE TABLE markets (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        STRING NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE companies (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    market_id   INT8 NOT NULL REFERENCES markets(id),
    name        STRING NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(market_id, name)
);

CREATE TABLE brokerages (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        STRING NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username    STRING NOT NULL UNIQUE,
    firstname   STRING NOT NULL,
    lastname    STRING NOT NULL,
    password    BYTES NOT NULL,
    active      BOOL NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE stocks (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_id  INT8 NOT NULL REFERENCES companies(id),
    market_id   INT8 NOT NULL REFERENCES markets(id),
    ticker      STRING NOT NULL,
    isin        STRING DEFAULT NULL UNIQUE,
    name        STRING DEFAULT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(market_id, ticker)
);

CREATE TABLE stock_registers(
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    stock_id    INT8 NOT NULL REFERENCES stocks(id),
    price       FLOAT8 NOT NULL,
    tendency    INT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON stock_registers (stock_id, created_at DESC);

CREATE TABLE stock_recommendations (
    id                  INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    brokerage_id        INT8 NOT NULL REFERENCES brokerages(id),
    stock_register_id   INT8 NOT NULL REFERENCES stock_registers(id),
    rating_from         INT NOT NULL,
    rating_to           INT NOT NULL,
    target_from         FLOAT8 NOT NULL,
    target_to           FLOAT8 NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE stock_users (
    user_id INT8 NOT NULL REFERENCES users(id), 
    stock_id INT8 NOT NULL REFERENCES stocks(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, stock_id)
);

CREATE TABLE stock_tendency_stats(
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    stock_id    INT8 NOT NULL REFERENCES stocks(id) UNIQUE,
    up_count    INT8 NOT NULL DEFAULT 0,
    side_count  INT8 NOT NULL DEFAULT 0,
    down_count  INT8 NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(stock_id)
);
