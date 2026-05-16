-- Core schema (PostgreSQL без TimescaleDB — для локального docker-compose)

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL DEFAULT '',
    display_name    TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_settings (
    user_id         UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    locale          TEXT NOT NULL DEFAULT 'ru' CHECK (locale IN ('ru', 'en')),
    base_currency   CHAR(3) NOT NULL DEFAULT 'RUB',
    theme           TEXT NOT NULL DEFAULT 'dark',
    timezone        TEXT NOT NULL DEFAULT 'Europe/Moscow'
);

CREATE TABLE categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    kind            TEXT NOT NULL CHECK (kind IN ('income', 'expense')),
    icon            TEXT,
    color           TEXT,
    is_system       BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id     UUID REFERENCES categories(id) ON DELETE SET NULL,
    kind            TEXT NOT NULL CHECK (kind IN ('income', 'expense')),
    amount          NUMERIC(18, 2) NOT NULL CHECK (amount > 0),
    currency        CHAR(3) NOT NULL DEFAULT 'RUB',
    description     TEXT,
    occurred_at     TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_user_date ON transactions (user_id, occurred_at DESC);

CREATE TABLE goals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    goal_type       TEXT NOT NULL CHECK (goal_type IN ('savings', 'investment', 'purchase')),
    target_amount   NUMERIC(18, 2) NOT NULL,
    current_amount  NUMERIC(18, 2) NOT NULL DEFAULT 0,
    currency        CHAR(3) NOT NULL DEFAULT 'RUB',
    deadline        DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE goal_contributions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goal_id         UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    amount          NUMERIC(18, 2) NOT NULL CHECK (amount > 0),
    contributed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    note            TEXT
);

CREATE INDEX idx_goal_contributions_goal ON goal_contributions (goal_id, contributed_at DESC);

CREATE TABLE instruments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol          TEXT NOT NULL,
    exchange        TEXT NOT NULL,
    name            TEXT NOT NULL,
    asset_type      TEXT NOT NULL DEFAULT 'stock',
    currency        CHAR(3) NOT NULL DEFAULT 'RUB',
    UNIQUE (symbol, exchange)
);

CREATE TABLE market_prices (
    instrument_id   UUID NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    time            TIMESTAMPTZ NOT NULL,
    close           NUMERIC(18, 6) NOT NULL,
    volume          BIGINT,
    PRIMARY KEY (instrument_id, time)
);

CREATE TABLE company_summaries (
    instrument_id   UUID PRIMARY KEY REFERENCES instruments(id) ON DELETE CASCADE,
    summary_text    TEXT NOT NULL,
    key_metrics     JSONB NOT NULL DEFAULT '{}',
    fetched_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL
);

CREATE TABLE ai_forecasts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instrument_id   UUID NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    horizon_days    INT NOT NULL,
    predicted_change_pct NUMERIC(8, 4),
    confidence      NUMERIC(5, 4),
    narrative       TEXT NOT NULL,
    model_version   TEXT NOT NULL,
    disclaimer_id   TEXT NOT NULL DEFAULT 'v1_not_advice',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
