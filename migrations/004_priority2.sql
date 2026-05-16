-- Приоритет 2: портфель, watchlist

CREATE TABLE IF NOT EXISTS portfolios (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL DEFAULT 'Основной',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS portfolio_holdings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    portfolio_id    UUID NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    instrument_id   UUID NOT NULL REFERENCES instruments(id),
    quantity        NUMERIC(18, 8) NOT NULL CHECK (quantity > 0),
    avg_cost        NUMERIC(18, 6) NOT NULL,
    currency        CHAR(3) NOT NULL DEFAULT 'RUB',
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (portfolio_id, instrument_id)
);

CREATE TABLE IF NOT EXISTS watchlists (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL DEFAULT 'Основной',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS watchlist_items (
    watchlist_id    UUID NOT NULL REFERENCES watchlists(id) ON DELETE CASCADE,
    instrument_id   UUID NOT NULL REFERENCES instruments(id) ON DELETE CASCADE,
    sort_order      INT NOT NULL DEFAULT 0,
    PRIMARY KEY (watchlist_id, instrument_id)
);

-- Demo watchlist
INSERT INTO watchlists (id, user_id, name)
VALUES ('20000000-0000-4000-8000-000000000001', '00000000-0000-4000-8000-000000000001', 'Основной')
ON CONFLICT (id) DO NOTHING;

INSERT INTO watchlist_items (watchlist_id, instrument_id, sort_order)
SELECT '20000000-0000-4000-8000-000000000001', id, row_number() OVER ()::int
FROM instruments WHERE symbol IN ('SBER', 'GAZP', 'LKOH')
ON CONFLICT DO NOTHING;

INSERT INTO portfolios (id, user_id, name)
VALUES ('30000000-0000-4000-8000-000000000001', '00000000-0000-4000-8000-000000000001', 'Основной')
ON CONFLICT (id) DO NOTHING;

INSERT INTO portfolio_holdings (portfolio_id, instrument_id, quantity, avg_cost, currency)
SELECT '30000000-0000-4000-8000-000000000001', id, 10, 250, 'RUB'
FROM instruments WHERE symbol = 'SBER'
ON CONFLICT DO NOTHING;
