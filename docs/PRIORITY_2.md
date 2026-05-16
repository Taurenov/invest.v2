# Приоритет 2 — готово

## Реализовано

| Фича | API | UI |
|------|-----|-----|
| Портфель | `GET/POST/DELETE /api/v1/me/portfolio/holdings` | `/portfolio` |
| Watchlist | `GET/POST/DELETE /api/v1/me/watchlist/items` | Рынки |
| Сводка компании | `GET /api/v1/markets/{symbol}/summary` | Кнопка «Обзор» |
| Аналитика | `GET /api/v1/me/analytics`, `export.csv` | Фильтр дат, stacked bar, CSV |
| Настройки | `GET/PATCH /api/v1/me/settings` | `/settings` |
| Уведомления | `GET /api/v1/me/alerts` | Tauri + Web Notifications |
| Тема | `theme: dark/light/system` | `data-theme` на `<html>` |

Миграция: `migrations/004_priority2.sql`

## Уведомления

- При достижении цели (проверка при загрузке целей)
- Опрос `/api/v1/me/alerts` каждые 15 с в приложении
- Windows: разрешите уведомления при первом запуске

## Запуск

Как в PRIORITY_1; для новых таблиц:

```powershell
docker compose down -v
docker compose up -d
```
