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

**Приоритет 1** — [docs/PRIORITY_1.md](docs/PRIORITY_1.md) · **Приоритет 2** — [docs/PRIORITY_2.md](docs/PRIORITY_2.md) (портфель, watchlist, аналитика, настройки, уведомления).

## Структура

| Путь | Описание |
|------|----------|
| `apps/desktop/` | Tauri 2 + React, i18n ru/en, экраны |
| `backend/` | REST API, WebSocket котировок |
| `engine/fin-grpc/` | Rust: gRPC + HTTP JSON |
| `engine/fin-math/` | ROI, линейный прогноз |
| `migrations/` | SQL-схема и seed |
| `docker-compose.yml` | Postgres + Redis |
| `docs/ARCHITECTURE.md` | Архитектура |

## API (основное)

| Метод | Путь |
|-------|------|
| GET | `/api/v1/users/{id}/transactions` |
| GET | `/api/v1/users/{id}/goals` |
| GET | `/api/v1/markets/{symbol}/quote?exchange=MOEX` |
| GET | `/api/v1/markets/{symbol}/forecast?horizon_days=7&locale=ru` |
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
