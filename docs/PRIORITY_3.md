# Приоритет 3 — выполнено

## Что добавлено

- **История прогнозов**: каждый `GET /api/v1/markets/{symbol}/forecast` теперь сохраняется в `ai_forecasts`.
- **API истории прогнозов**: `GET /api/v1/markets/{symbol}/forecast/history?limit=20`.
- **API истории цен для графика**: `GET /api/v1/markets/{symbol}/history?points=120`.
- **Экран прогноза (UI)**:
  - график на базе реальных исторических цен + точка прогноза;
  - таблица последних прогнозов (дата, горизонт, изменение, confidence).
- **Пакетирование .exe**:
  - sidecar-автопоиск `fin-api` и `fin-engine` рядом с exe или в `binaries`;
  - ресурсы `src-tauri/binaries/*` включаются в bundle;
  - добавлен рабочий `icon.ico`, Tauri сборка проходит.

## Новые/обновленные эндпоинты

- `GET /api/v1/markets/{symbol}/forecast`
- `GET /api/v1/markets/{symbol}/forecast/history`
- `GET /api/v1/markets/{symbol}/history`

Все под JWT/dev-token через `Authorization: Bearer ...`.

## Проверка

```powershell
# backend
cd backend
go build -o fin-api.exe ./cmd/api

# tauri rust side
cd ../apps/desktop/src-tauri
cargo build
```

## Для «одного установщика»

1. Соберите sidecar бинарники:
   - `scripts/build-sidecars.ps1`
2. Убедитесь, что в `apps/desktop/src-tauri/binaries` есть:
   - `fin-api.exe`
   - `fin-engine.exe`
3. Сборка инсталлятора:
   - `cd apps/desktop`
   - `npm run tauri:build`

При запуске десктоп-клиент попробует поднять sidecar автоматически.
