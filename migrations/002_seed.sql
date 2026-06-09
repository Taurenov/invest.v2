-- Demo user and sample data
INSERT INTO users (id, email, display_name)
VALUES (
    '00000000-0000-4000-8000-000000000001',
    'demo@fin-helper.local',
    'Demo User'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO user_settings (user_id, locale, base_currency)
VALUES ('00000000-0000-4000-8000-000000000001', 'ru', 'RUB')
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO categories (user_id, name, kind, icon, color, is_system)
SELECT '00000000-0000-4000-8000-000000000001', v.name, v.kind, v.icon, v.color, true
FROM (VALUES
    ('Зарплата', 'income', '💼', '#22c55e'),
    ('Подработка', 'income', '📈', '#3b82f6'),
    ('Продукты', 'expense', '🛒', '#f59e0b'),
    ('Транспорт', 'expense', '🚌', '#8b5cf6'),
    ('Жильё', 'expense', '🏠', '#ef4444'),
    ('Развлечения', 'expense', '🎬', '#ec4899')
) AS v(name, kind, icon, color)
WHERE NOT EXISTS (
    SELECT 1 FROM categories c
    WHERE c.user_id = '00000000-0000-4000-8000-000000000001' AND c.name = v.name AND c.kind = v.kind
);

INSERT INTO transactions (user_id, kind, amount, currency, description, occurred_at)
VALUES
    ('00000000-0000-4000-8000-000000000001', 'income', 85000, 'RUB', 'Зарплата', now() - interval '2 days'),
    ('00000000-0000-4000-8000-000000000001', 'expense', 3200, 'RUB', 'Продукты', now() - interval '5 days'),
    ('00000000-0000-4000-8000-000000000001', 'expense', 890, 'RUB', 'Транспорт', now() - interval '6 days')
ON CONFLICT DO NOTHING;

INSERT INTO goals (user_id, title, goal_type, target_amount, current_amount, currency)
VALUES (
    '00000000-0000-4000-8000-000000000001',
    'Подушка безопасности',
    'savings',
    300000,
    200000,
    'RUB'
) ON CONFLICT DO NOTHING;

INSERT INTO instruments (id, symbol, exchange, name, asset_type, currency)
VALUES (
    '10000000-0000-4000-8000-000000000001',
    'SBER',
    'MOEX',
    'Сбербанк',
    'stock',
    'RUB'
) ON CONFLICT (symbol, exchange) DO NOTHING;

INSERT INTO market_prices (instrument_id, time, close)
SELECT '10000000-0000-4000-8000-000000000001', now() - (n || ' hours')::interval, 280 + n * 0.15
FROM generate_series(0, 120) AS n
ON CONFLICT DO NOTHING;

INSERT INTO instruments (id, symbol, exchange, name, asset_type, currency)
VALUES
    ('10000000-0000-4000-8000-000000000002', 'GAZP', 'MOEX', 'Газпром', 'stock', 'RUB'),
    ('10000000-0000-4000-8000-000000000003', 'LKOH', 'MOEX', 'Лукойл', 'stock', 'RUB')
ON CONFLICT (symbol, exchange) DO NOTHING;

INSERT INTO market_prices (instrument_id, time, close)
SELECT '10000000-0000-4000-8000-000000000002', now() - (n || ' hours')::interval, 165 + n * 0.08
FROM generate_series(0, 120) AS n
ON CONFLICT DO NOTHING;

INSERT INTO market_prices (instrument_id, time, close)
SELECT '10000000-0000-4000-8000-000000000003', now() - (n || ' hours')::interval, 7200 + n * 2.5
FROM generate_series(0, 120) AS n
ON CONFLICT DO NOTHING;
