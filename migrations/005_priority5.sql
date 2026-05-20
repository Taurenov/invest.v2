-- Приоритет 5: теги, повторяющиеся транзакции, бюджеты по категориям

CREATE TABLE IF NOT EXISTS tags (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    color       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, name)
);

CREATE TABLE IF NOT EXISTS transaction_tags (
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    tag_id         UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (transaction_id, tag_id)
);

-- repeating rules
CREATE TABLE IF NOT EXISTS recurring_transactions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind          TEXT NOT NULL CHECK (kind IN ('income', 'expense')),
    category_id   UUID REFERENCES categories(id) ON DELETE SET NULL,
    amount        NUMERIC(18, 2) NOT NULL CHECK (amount > 0),
    currency      CHAR(3) NOT NULL DEFAULT 'RUB',
    description   TEXT,
    schedule      TEXT NOT NULL CHECK (schedule IN ('daily', 'weekly', 'monthly')),
    day_of_month  INT,
    day_of_week   INT,
    next_run_at   DATE NOT NULL,
    is_active     BOOLEAN NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_recurring_user_next ON recurring_transactions (user_id, next_run_at);

-- budgets per category/month
CREATE TABLE IF NOT EXISTS budgets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    period      TEXT NOT NULL DEFAULT 'monthly' CHECK (period IN ('monthly')),
    amount      NUMERIC(18, 2) NOT NULL CHECK (amount > 0),
    currency    CHAR(3) NOT NULL DEFAULT 'RUB',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, category_id, period)
);

-- demo: one tag and one budget (optional)
