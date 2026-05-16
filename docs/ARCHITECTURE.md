# Fin Helper — архитектура десктопного финансового помощника

## 1. Обзор

Многослойное приложение: **десктопный клиент** (Tauri + React — рекомендуемый вариант для современного UI), **Go API** (пользователи, транзакции, цели, портфель, прокси к биржам), **Rust engine** (ROI, временные ряды, ML-прогнозы), **PostgreSQL/TimescaleDB** + **Redis**.

```
┌─────────────────────────────────────────────────────────────────┐
│  Desktop (Tauri 2 + React/TS)                                   │
│  • UI, i18n (ru/en), локальный secure storage (OS keychain)     │
│  • HTTP → Go API  |  gRPC/HTTP → Rust engine (тяжёлые расчёты)  │
└───────────────┬─────────────────────────────┬───────────────────┘
                │ REST + WebSocket             │ gRPC (internal)
                ▼                              ▼
┌───────────────────────────┐    ┌──────────────────────────────┐
│  Go API (api-service)      │    │  Rust engine (fin-engine)     │
│  • Auth, CRUD, аналитика   │◄──►│  • ROI, CAGR, прогнозы TS     │
│  • Кэш Redis, rate limits  │    │  • ONNX / linreg (опционально)│
│  • Jobs: синк котировок    │    └──────────────────────────────┘
└───────────────┬───────────┘
                │
        ┌───────┴───────┐
        ▼               ▼
   PostgreSQL      Redis
   (+ Timescale)   (котировки, сессии)
        │
        ▼
   External APIs (Alpha Vantage, MOEX ISS, Finnhub…)
```

### Почему Tauri, а не чистый Fyne

| Критерий | Tauri + React | Fyne / Wails |
|----------|---------------|--------------|
| Современный UI | Charts (Recharts/ECharts), анимации, design tokens | Ограниченнее |
| i18n, темы | react-i18next, CSS variables | Возможно, но труднее |
| Размер бинарника | Малый (WebView OS) | Fyne — CGO; Wails — WebView |
| Единый стек с вебом | Да | Wails — частично |

**Wails** — запасной вариант, если команда хочет максимум Go на клиенте. В репозитории — минимальный **Wails/HTML** прототип в `desktop/`.

---

## 2. Сервисы и модули

### 2.1 `desktop` (Tauri + React)

| Модуль | Назначение |
|--------|------------|
| `src/views/` | Home, Transactions, Analytics, Markets, Forecast, Goals, Settings |
| `src/components/` | BalanceCard, TransactionList, Charts, DisclaimerBanner |
| `src/api/` | Клиент к Go API (axios/fetch), WebSocket котировок |
| `src/i18n/` | `ru.json`, `en.json` |
| `src-tauri/` | Rust: окно, deep links, автообновление, вызов локального API |

**Связь:** `https://127.0.0.1:18443` (локальный Go) или Unix-socket / named pipe для IPC без TLS в dev.

### 2.2 `api-service` (Go)

| Пакет | Назначение |
|-------|------------|
| `cmd/api` | Точка входа, graceful shutdown |
| `internal/http` | REST handlers, middleware (auth, CORS, request ID) |
| `internal/ws` | WebSocket: live quotes, goal progress |
| `internal/domain` | Сущности, валидация |
| `internal/repo` | PostgreSQL (sqlc или pgx) |
| `internal/cache` | Redis |
| `internal/market` | Адаптеры внешних API + нормализация тикеров |
| `internal/engine` | gRPC-клиент к Rust |
| `internal/jobs` | Cron: обновление цен, агрегаты в Timescale |

**Протоколы:**

- **REST (JSON)** — CRUD транзакций, целей, настройки, сводки компаний.
- **WebSocket** — поток цен и уведомления.
- **gRPC** — только Go ↔ Rust (прогноз, пакетные расчёты ROI).

### 2.3 `fin-engine` (Rust)

| Crate | Назначение |
|-------|------------|
| `fin-math` | ROI, CAGR, NPV, простые прогнозы |
| `fin-ts` | Скользящие средние, линейная регрессия на рядах |
| `fin-ml` | (фаза 2) ONNX-модели |
| `fin-grpc` | tonic server |

Экспорт: **gRPC** (prod) или **cdylib + FFI** (если engine встроен в Tauri sidecar).

### 2.4 Инфраструктура

- **PostgreSQL 16** — OLTP.
- **TimescaleDB** — hypertable `market_prices`, `portfolio_snapshots`.
- **Redis** — TTL-кэш котировок (30–60 с), сессии, rate limit по API-ключам бирж.

---

## 3. Коммуникация

| Канал | Кто | Что |
|-------|-----|-----|
| REST `/api/v1/*` | Desktop → Go | Транзакции, цели, портфель, сводки |
| WS `/ws/quotes` | Desktop → Go | Live цены |
| gRPC `Predict`, `CalculateRoi` | Go → Rust | Тяжёлая математика |
| Локальный IPC (опционально) | Tauri ↔ Go | Unix socket `fin-helper.sock` — без открытого порта |
| Внешние HTTPS | Go → биржи | Котировки, fundamentals |

**Безопасность:** JWT (access + refresh), refresh в httpOnly при веб-версии; в десктопе — OS keychain через Tauri plugin. Все секреты API бирж — только на сервере Go, не в клиенте.

---

## 4. Схема БД

См. `migrations/001_initial.sql`.

Кратко:

- `users`, `user_settings` (locale, currency, theme)
- `categories`, `transactions`
- `goals`, `goal_contributions`
- `instruments`, `watchlists`, `watchlist_items`
- `market_prices` (Timescale hypertable)
- `portfolios`, `portfolio_holdings`
- `company_summaries` (кэш текстовых обзоров)
- `ai_forecasts` (с disclaimer hash / version модели)

---

## 5. UX — десктоп (не мобильный)

### 5.1 Макет главного окна

```
┌──────────────────────────────────────────────────────────────────┐
│ [≡] Fin Helper          🔍 Поиск...          RU | EN    [👤]     │
├──────────┬───────────────────────────────────────────────────────┤
│ Dashboard│  Добрый день, Иван                    Март 2026      │
│ Транзакции│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐    │
│ Аналитика │ │ Баланс      │ │ Доход ▲     │ │ Расход ▼    │    │
│ Рынки     │ │ 124 500 ₽   │ │ +85 000     │ │ −42 300     │    │
│ Прогноз AI│ └─────────────┘ └─────────────┘ └─────────────┘    │
│ Цели      │ ┌──────────────────────────────┐ ┌──────────────┐  │
│ Настройки │ │ Доход / расход (12 мес.)     │ │ Цели (3)     │  │
│           │ │ [линейный график]            │ │ ████░░ 67%   │  │
│           │ └──────────────────────────────┘ └──────────────┘  │
│           │ Последние транзакции              [+ Добавить]      │
│           │ • Зарплата    +80 000    01.03                          │
│           │ • Продукты    −3 200     28.02                          │
└──────────┴───────────────────────────────────────────────────────┘
```

- **Sidebar** 240px (сворачивается до иконок 64px на ноутбуках &lt; 1280px).
- **Content** max-width ~1400px, карточки в grid 12 колонок.

### 5.2 Добавление транзакции (модал / drawer справа 480px)

Поля: тип (доход/расход), сумма, валюта, категория, дата, заметка, повтор (опционально). Primary CTA «Сохранить», secondary «Отмена».

### 5.3 Аналитика

Вкладки: «По категориям» (donut), «По времени» (stacked bar / line), фильтр периода. Таблица drill-down под графиком.

### 5.4 Рынки

Таблица: тикер, название, цена, Δ%, мини-sparkline. Клик → деталь: большой график, кнопки «Сводка», «В прогноз AI», «В watchlist».

### 5.5 AI-прогноз

График история + пунктир прогноз. Блок disclaimer (жёлтый/нейтральный фон): *«Не является инвестиционной рекомендацией…»*. Текст модели + confidence.

### 5.6 Design tokens

```css
--color-bg: #0f1419;
--color-surface: #1a2332;
--color-primary: #3b82f6;
--color-text: #e8edf4;
--color-muted: #8b9cb3;
--color-success: #22c55e;  /* только рост */
--color-danger: #ef4444;   /* только падение */
--font-sans: "Inter", "Segoe UI", system-ui;
--radius: 12px;
--shadow: 0 4px 24px rgba(0,0,0,.25);
```

Светлая тема: фон `#f4f6f9`, surface `#ffffff`.

---

## 6. i18n

Ключи в `desktop/src/i18n/ru.json` (default) и `en.json`. Формат денег: `Intl.NumberFormat`. Даты: `date-fns` с локалью.

---

## 7. Roadmap реализации

1. Миграции БД + Go CRUD транзакций/категорий  
2. Rust `fin-math` + gRPC  
3. Tauri shell + Home + Transactions  
4. Redis + market adapter (MOEX / Finnhub)  
5. WebSocket котировок  
6. AI forecast + disclaimer UI  
7. Упаковка: MSI (Windows), DMG, AppImage  

---

## 8. Запуск примеров в репозитории

```bash
# API
cd backend && go run ./cmd/api

# Rust (тесты)
cd engine/fin-math && cargo test

# Desktop прототип (откройте в браузере или через Wails)
cd desktop && go run ./cmd/app
```
