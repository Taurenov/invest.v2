# Приоритет 5 — выполнено

## Что добавлено

### 1) Поиск/фильтры транзакций

- UI: строка поиска + фильтр `income/expense` в `Транзакциях`.
- API: `GET /api/v1/me/transactions?q=...&kind=...&from=YYYY-MM-DD&to=YYYY-MM-DD`.

### 2) Теги

- Таблицы: `tags`, `transaction_tags`.
- API:
  - `GET/POST/DELETE /api/v1/me/tags`
  - `PUT /api/v1/me/transactions/{txID}/tags`
- UI: экран `Дополнительно → Теги` (создание/удаление).

### 3) Бюджеты по категориям

- Таблица: `budgets` (пока период `monthly`).
- API:
  - `GET/POST/DELETE /api/v1/me/budgets`
- UI: `Дополнительно → Бюджеты` (создать/удалить).

### 4) Повторяющиеся операции

- Таблица: `recurring_transactions`.
- API:
  - `GET/POST/DELETE /api/v1/me/recurring`
  - `POST /api/v1/me/recurring/{id}/toggle?active=true|false`
- UI: `Дополнительно → Повторяющиеся операции` (создать/вкл/выкл/удалить).

## Миграция

Добавлено: `migrations/005_priority5.sql` (подключено в `docker-compose.yml`).

Если база уже поднята и вы хотите пересоздать её с нуля:

```powershell
docker compose down -v
docker compose up -d
```

## Проверка сборки

- Backend: `go build ./cmd/api`
- Desktop: `npm run build`, `cargo build` (src-tauri)

## Что логично улучшить следующим шагом

- Привязка тегов к транзакциям в UI (мультиселект в форме).
- Реальный расчёт «сколько потрачено из бюджета за месяц» и предупреждения.
- Планировщик (cron/job) который создаёт транзакции из `recurring_transactions`.
