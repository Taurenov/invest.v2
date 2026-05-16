.PHONY: infra api engine desktop web

infra:
	docker compose up -d

api:
	cd backend && set API_TOKEN=dev-token&& set DATABASE_URL=postgres://fin:fin@127.0.0.1:5432/finhelper?sslmode=disable&& set REDIS_URL=redis://127.0.0.1:6379&& go run ./cmd/api

engine:
	cd engine/fin-grpc && cargo run --bin fin-engine

desktop:
	cd apps/desktop && npm install && npm run tauri dev

web:
	cd apps/desktop && npm install && npm run dev
