# Fin Helper

Десктопный финансовый помощник: **Tauri + React**, **Go API**, **Rust engine**, **PostgreSQL**, **Redis**.

## Быстрый старт (Windows)

### 1. Инфраструктура

```powershell
docker compose up -d
```

### 2. Rust engine (HTTP :50052)

```powershell
cd engine\fin-grpc
cargo run --bin fin-engine
```

### 3. Go API

```powershell
$env:API_TOKEN = "dev-token"
$env:DATABASE_URL = "postgres://fin:fin@127.0.0.1:5432/finhelper?sslmode=disable"
$env:REDIS_URL = "redis://127.0.0.1:6379"
$env:ENGINE_HTTP_URL = "http://127.0.0.1:50052"
cd backend
go run ./cmd/api
```

Проверка:

```powershell
curl -H "Authorization: Bearer dev-token" http://localhost:8080/api/v1/users/00000000-0000-4000-8000-000000000001/transactions
curl -H "Authorization: Bearer dev-token" "http://localhost:8080/api/v1/markets/SBER/forecast?horizon_days=7"
```

### 4. Desktop — **приложение .exe** (Tauri + React)

Это **не** «открыть сайт в браузере». Tauri собирает **нативное Windows-приложение** (`.exe` + установщик NSIS).

Требуется [Node.js](https://nodejs.org/), [Rust](https://rustup.rs/), WebView2.

**Окно приложения (разработка):**

```powershell
cd apps\desktop
copy .env.example .env
npm install
npm run tauri:dev
```

**Сборка установщика .exe для пользователей:**

```powershell
npm run tauri:build
# Результат: src-tauri\target\release\bundle\nsis\*-setup.exe
```

Подробно: [docs/BUILD_WINDOWS.md](docs/BUILD_WINDOWS.md)

`npm run dev` — только для правки UI в браузере, **не** то, что отдаёте пользователям.

**Приоритет 1** — [docs/PRIORITY_1.md](docs/PRIORITY_1.md) · **Приоритет 2** — [docs/PRIORITY_2.md](docs/PRIORITY_2.md) · **Приоритет 3** — [docs/PRIORITY_3.md](docs/PRIORITY_3.md) · **Приоритет 4** — [docs/PRIORITY_4.md](docs/PRIORITY_4.md) · **Приоритет 5** — [docs/PRIORITY_5.md](docs/PRIORITY_5.md).

## Структура

| Путь | Описание |
|------|----------|
| `apps/desktop/` | Tauri 2 + React, i18n ru/en, экраны |
| `backend/` | REST API, WebSocket котировок |
| `engine/fin-grpc/` | Rust: HTTP JSON engine |
| `engine/fin-math/` | ROI, линейный прогноз |
| `migrations/` | SQL-схема и seed |
| `docker-compose.yml` | Postgres + Redis |
| `docs/ARCHITECTURE.md` | Архитектура |

## API (основное)

| Метод | Путь |
|-------|------|
| POST | `/api/v1/auth/login` / `/api/v1/auth/register` |
| GET/POST/PATCH/DELETE | `/api/v1/me/transactions` |
| GET/POST | `/api/v1/me/goals`, `/api/v1/me/goals/{goalID}/contribute` |
| GET/POST/DELETE | `/api/v1/me/watchlist`, `/api/v1/me/watchlist/items` |
| GET/POST/DELETE | `/api/v1/me/portfolio`, `/api/v1/me/portfolio/holdings` |
| GET/PATCH | `/api/v1/me/settings` |
| GET | `/api/v1/me/analytics`, `/api/v1/me/analytics/export.csv` |
| GET | `/api/v1/markets/{symbol}/quote` |
| GET | `/api/v1/markets/{symbol}/summary` |
| GET | `/api/v1/markets/{symbol}/forecast` |
| GET | `/api/v1/markets/{symbol}/forecast/history` |
| GET | `/api/v1/markets/{symbol}/history` |
| WS | `/ws/quotes` |

Заголовок: `Authorization: Bearer dev-token`

## Переменные окружения

| Переменная | По умолчанию |
|------------|--------------|
| `API_TOKEN` | `dev-token` |
| `DATABASE_URL` | — (memory mode) |
| `REDIS_URL` | — |
| `ENGINE_HTTP_URL` | `http://127.0.0.1:50052` |
| `VITE_API_URL` | `http://127.0.0.1:8080` |
