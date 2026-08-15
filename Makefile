.PHONY: dev dev-web dev-backend test build-web build up down logs smoke

dev: dev-backend dev-web

dev-backend:
	go run ./cmd/server

dev-web:
	cd web && npm run dev

test:
	go test ./...

build-web:
	cd web && npm run build

build:
	go build -o server ./cmd/server

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

# Quick smoke test against a running container (uses values from config.yml).
smoke:
	@echo "== login =="
	@curl -s -X POST localhost:1325/api/login -H 'Content-Type: application/json' \
		-d '{"password":"CHANGE_ME_PASSWORD"}'
