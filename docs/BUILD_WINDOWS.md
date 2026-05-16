# Сборка Fin Helper как .exe (Windows)

## Это не «сайт в браузере»

| Режим | Что получаете |
|-------|----------------|
| `npm run dev` | Только для разработки UI в браузере (localhost:1420) |
| **`npm run tauri:dev`** | **Окно приложения** на ПК (как будет у пользователя) |
| **`npm run tauri:build`** | **Установщик .exe / .msi** для распространения |

Tauri упаковывает интерфейс (React) в **нативное окно** через WebView2 (компонент Windows). Пользователь **не открывает Chrome** — запускает `Fin Helper.exe` с рабочего стола или из меню «Пуск».

Аналогия: Discord, Spotify, VS Code — тот же принцип (нативная оболочка + веб-UI внутри).

## Что нужно установить (один раз)

1. [Node.js](https://nodejs.org/) LTS  
2. [Rust](https://rustup.rs/)  
3. [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/) (часто уже есть в Windows 11)  
4. [Visual Studio Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/) — workload «Desktop development with C++»

## Сборка установщика

```powershell
# 1. API и engine (в отдельных терминалах или как службы — см. README)
docker compose up -d
cd engine\fin-grpc
cargo run --bin fin-engine

cd backend
$env:API_TOKEN = "dev-token"
$env:DATABASE_URL = "postgres://fin:fin@127.0.0.1:5432/finhelper?sslmode=disable"
go run ./cmd/api

# 2. Десктоп .exe
cd apps\desktop
copy .env.example .env
npm install
npm run tauri:build
```

## Где лежит .exe

После успешной сборки:

```
apps\desktop\src-tauri\target\release\
  fin-helper-desktop.exe          ← можно запускать напрямую

apps\desktop\src-tauri\target\release\bundle\nsis\
  Fin Helper_0.1.0_x64-setup.exe  ← установщик для пользователей
```

## Альтернатива: чистый Go .exe (Fyne)

В `desktop/cmd/app` — простое окно без Tauri:

```powershell
cd desktop
go build -o FinHelper.exe ./cmd/app
```

Меньше возможностей по UI, но один файл без Node.js.

## Продакшен: один установщик «всё в одном»

Сейчас: **клиент .exe** + **локальный API** (Go) + **engine** (Rust). Для одного инсталлятора позже можно:

- вложить `fin-api.exe` и `fin-engine.exe` как sidecar в bundle Tauri;
- поднимать их при старте приложения автоматически.

Текущая схема удобна для разработки и отладки.
