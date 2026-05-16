-- Fin Helper — начальная схема PostgreSQL (+ TimescaleDB для цен)

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "timescaledb" CASCADE;

-- Пользователи
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    display_name    TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_settings (
    user_id         UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    locale          TEXT NOT NULL DEFAULT 'ru' CHECK (locale IN ('ru', 'en')),
    base_currency   CHAR(3) NOT NULL DEFAULT 'RUB',
    theme           TEXT NOT NULL DEFAULT 'dark' CHECK (theme IN ('dark', 'light', 'system')),
    timezone        TEXT NOT NULL DEFAULT 'Europe/Moscow'
);

-- Категории (системные + пользовательские)
CREATE TABLE categories (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    kind            TEXT NOT NULL CHECK (kind IN ('income', 'expense')),
    icon            TEXT,
    color           TEXT,
    is_system       BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, name, kind)
);

-- Транзакции
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
CREATE INDEX idx_transactions_user_kind ON transactions (user_id, kind);

-- Цели
CREATE TABLE goals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           TEXT NOT NULL,
    goal_type       TEXT NOT NULL CHECK (goal_type IN ('savings', 'investment', 'purchase')),
    target_amount   NUMERIC(18, 2) NOT NULL CHECK (target_amount > 0),
    current_amount  NUMERIC(18, 2) NOT NULL DEFAULT 0 CHECK (current_amount >= 0),
    currency        CHAR(3) NOT NULL DEFAULT 'RUB',
    deadline        DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE goal_contributions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goal_id         UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    amount          NUMERIC(18, 2) NOT NULL CHECK (amount > 0),
    contributed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    note            TEXT
);

-- Инструменты (акции, ETF, облигации)
CREATE TABLE instruments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol          TEXT NOT NULL,
    exchange        TEXT NOT NULL,
    name            TEXT NOT NULL,
    asset_type      TEXT NOT NULL CHECK (asset_type IN ('stock', 'etf', 'bond', 'fund')),
    currency        CHAR(3) NOT NULL,
    UNIQUE (symbol, exchange)
);

CREATE TABLE watchlists (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL DEFAULT 'Основной',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE watchlist_items (
    watchlist_id    UUID NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
    instrument_id   UUID NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    sort_order      INT NOT NULL DEFAULT 0,
    PRIMARY KEY (watchlist_id, instrument_id)
);

-- Котировки (Timescale hypertable)
CREATE TABLE market_prices (
    instrument_id   UUID NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    time            TIMESTAMPTZ NOT NULL,
    open            NUMERIC(18, 6),
    high            NUMERIC(18, 6),
    low             NUMERIC(18, 6),
    close           NUMERIC(18, 6) NOT NULL,
    volume          BIGINT,
    PRIMARY KEY (instrument_id, time)
);

SELECT create_hypertable('market_prices', 'time', if_not_exists => TRUE);

-- Портфель
CREATE TABLE portfolios (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL DEFAULT 'Основной',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE portfolio_holdings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    portfolio_id    UUID NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    instrument_id   UUID NOT NULL REFERENCES instruments(id),
    quantity        NUMERIC(18, 8) NOT NULL CHECK (quantity > 0),
    avg_cost        NUMERIC(18, 6) NOT NULL,
    currency        CHAR(3) NOT NULL,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (portfolio_id, instrument_id)
);

CREATE TABLE portfolio_snapshots (
    portfolio_id    UUID NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    time            TIMESTAMPTZ NOT NULL,
    total_value     NUMERIC(18, 2) NOT NULL,
    PRIMARY KEY (portfolio_id, time)
);

SELECT create_hypertable('portfolio_snapshots', 'time', if_not_exists => TRUE);

-- Текстовые сводки по компаниям (кэш)
CREATE TABLE company_summaries (
    instrument_id   UUID PRIMARY KEY REFERENCES instruments(id) ON DELETE CASCADE,
    summary_text    TEXT NOT NULL,
    key_metrics     JSONB NOT NULL DEFAULT '{}',
    source          TEXT,
    fetched_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL
);

-- AI-прогнозы (с версией модели и disclaimer)
CREATE TABLE ai_forecasts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    instrument_id   UUID NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    horizon_days    INT NOT NULL CHECK (horizon_days > 0),
    predicted_change_pct NUMERIC(8, 4),
    confidence      NUMERIC(5, 4) CHECK (confidence >= 0 AND confidence <= 1),
    narrative       TEXT NOT NULL,
    model_version   TEXT NOT NULL,
    disclaimer_id   TEXT NOT NULL DEFAULT 'v1_not_advice',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ai_forecasts_instrument ON ai_forecasts (instrument_id, created_at DESC);
