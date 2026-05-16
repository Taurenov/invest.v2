# Сборка бинарников для упаковки с Tauri (положите рядом с .exe или задайте FIN_API_BIN / FIN_ENGINE_BIN)
$root = Split-Path -Parent $PSScriptRoot
$out = Join-Path $root "apps\desktop\src-tauri\binaries"
New-Item -ItemType Directory -Force -Path $out | Out-Null

Push-Location (Join-Path $root "backend")
go build -o (Join-Path $out "fin-api.exe") ./cmd/api
Pop-Location

Push-Location (Join-Path $root "engine\fin-grpc")
cargo build --release --bin fin-engine
Copy-Item "target\release\fin-engine.exe" (Join-Path $out "fin-engine.exe") -Force
Pop-Location

Write-Host "Built:" (Get-ChildItem $out)
