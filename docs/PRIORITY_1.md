# Приоритет 1 — готово

## Реализовано

### Backend (Go)
- **JWT**: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /api/v1/me`
- **Транзакции**: `GET/POST /api/v1/me/transactions`, `PATCH/DELETE .../{txID}`
- **Категории**: `GET/POST /api/v1/me/categories` (seed при регистрации в Postgres)
- **Цели**: `GET/POST /api/v1/me/goals`, `POST .../{goalID}/contribute`
- **Калькулятор**: `POST /api/v1/calculator/roi|cagr|savings` → Rust engine
- Dev-токен `API_TOKEN=dev-token` по-прежнему работает (demo user)

### Rust engine
- `POST /v1/cagr`, `POST /v1/savings` (+ ROI, predict)

### Desktop (Tauri + React)
- Вход / регистрация
- Форма добавления транзакции, удаление
- Страница **Цели** (создание + взнос)
- Страница **Калькулятор** (ROI, CAGR, накопления)
- Sidecar: `FIN_API_BIN`, `FIN_ENGINE_BIN` при старте приложения

## Запуск

```powershell
docker compose up -d
cd engine\fin-grpc && cargo run --bin fin-engine
cd backend
$env:DATABASE_URL="postgres://fin:fin@127.0.0.1:5432/finhelper?sslmode=disable"
$env:API_TOKEN="dev-token"
$env:JWT_SECRET="change-me"
go run ./cmd/api

cd apps\desktop
copy .env.example .env
npm install
npm run tauri:dev
```

Регистрация: `/register` или dev-токен в `.env` (`VITE_API_TOKEN=dev-token`).

## Sidecar для одного дистрибутива

```powershell
.\scripts\build-sidecars.ps1
$env:FIN_API_BIN="путь\fin-api.exe"
$env:FIN_ENGINE_BIN="путь\fin-engine.exe"
npm run tauri:build
```
