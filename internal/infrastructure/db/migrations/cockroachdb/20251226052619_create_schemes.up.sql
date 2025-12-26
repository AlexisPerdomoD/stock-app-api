
-- Markets
CREATE TABLE markets (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        STRING NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Companies
CREATE TABLE companies (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    market_id   INT8 NOT NULL,
    name        STRING NOT NULL,
    isin        STRING NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_companies_market
        FOREIGN KEY (market_id)
        REFERENCES markets(id)
);

-- Brokerages
CREATE TABLE brokerages (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        STRING NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Users
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

-- Stocks
CREATE TABLE stocks (
    id          INT8 PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_id  INT8 NOT NULL,
    name        STRING NULL,
    ticker      STRING NOT NULL,
    price       FLOAT8 NOT NULL,
    tendency    STRING NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_stocks_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id)
);
