CREATE TABLE IF NOT EXISTS goal_contributions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goal_id         UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    amount          NUMERIC(18, 2) NOT NULL CHECK (amount > 0),
    contributed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    note            TEXT
);

CREATE INDEX IF NOT EXISTS idx_goal_contributions_goal ON goal_contributions (goal_id, contributed_at DESC);
