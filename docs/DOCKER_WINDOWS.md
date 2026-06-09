# Запуск через Docker (Windows, без WSL на хосте)

## Что в Docker, а что нет

| Сервис | Контейнер | Порт на ПК |
|--------|-----------|------------|
| PostgreSQL | да | 5432 |
| Redis | да | 6379 |
| Go API | да | 8080 |
| Rust engine | да | 50052 |
| UI (React) | да | 1420 (браузер) |
| Tauri (.exe окно) | **нет** | только `npm run tauri:dev` на Windows |

Первый `docker-compose.yml` — только БД. Файл `docker-compose.app.yml` добавляет API, engine и веб-UI.

## Docker Desktop без WSL

1. **Docker Desktop** → Settings → General  
2. Снимите галочку **Use the WSL 2 based engine** (если WSL не установлен).  
3. Включите **Use Hyper-V** / backend **Hyper-V** (нужны Windows Pro/Enterprise и виртуализация в BIOS).  
4. Перезапустите Docker Desktop.

На **Windows Home** без WSL Linux-контейнеры часто недоступны — тогда см. раздел «Без Docker» в README.

## Запуск всего стека

PowerShell в корне проекта:

```powershell
cd C:\Users\Михаил\Desktop\fin-helper
docker compose -f docker-compose.yml -f docker-compose.app.yml up -d --build
```

Первый раз сборка **10–20 минут** (Rust + npm).

Проверка:

```powershell
docker compose -f docker-compose.yml -f docker-compose.app.yml ps
curl http://127.0.0.1:8080/health
```

Откройте в браузере: **http://127.0.0.1:1420**

Вход: автоматически **dev-token** (demo-пользователь), как в `.env.example`.

## Остановка

```powershell
docker compose -f docker-compose.yml -f docker-compose.app.yml down
```

С удалением данных БД:

```powershell
docker compose -f docker-compose.yml -f docker-compose.app.yml down -v
```

## Логи при ошибках

```powershell
docker compose -f docker-compose.yml -f docker-compose.app.yml logs -f api
docker compose -f docker-compose.yml -f docker-compose.app.yml logs -f engine
docker compose -f docker-compose.yml -f docker-compose.app.yml logs -f web
```

## Только БД в Docker (API на Windows)

Если полная сборка в Docker тяжёлая:

```powershell
docker compose up -d
```

Дальше в отдельных терминалах на Windows: `cargo run`, `go run`, `npm run tauri:dev` — см. README.
