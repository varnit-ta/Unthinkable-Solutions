SHELL := /bin/sh

.PHONY: help frontend backend sqlc migrate-up migrate-down docker-build docker-up docker-restart

help:
	@echo "Makefile targets:"
	@echo "  make frontend        # Run frontend dev server (from frontend folder)"
	@echo "  make backend         # Run Go backend server"
	@echo "  make sqlc            # Run sqlc generate (requires sqlc installed)"
	@echo "  make migrate-up      # Apply DB migrations (uses psql against DATABASE_URL)"
	@echo "  make migrate-down    # Drop/rollback migrations (manual drop)"
	@echo "  make docker-build    # Build docker images (frontend/backend)"
	@echo "  make docker-up       # docker compose up -d"
	@echo "  make docker-restart  # docker compose restart"

DATABASE_URL ?= postgres://unthinkable:unthinkable@localhost:5432/unthinkable_recipes?sslmode=disable

frontend:
	@echo "Starting frontend dev server..."
	@cd frontend && npm install && npm run dev

backend:
	@echo "Running Go backend..."
	@cd backend && go run ./cmd/server

sqlc:
	@echo "Generating sqlc code..."
	@cd backend && sqlc generate

migrate-up:
	@echo "Applying migrations..."
	@for f in backend/migrations/*.sql; do \
		echo "Applying $$f"; \
		psql "$(DATABASE_URL)" -f $$f; \
	done

migrate-seed: migrate-up sqlc
	@echo "migrations applied and sqlc generated"

migrate-down:
	@echo "migrate-down is destructive. Dropping tables..."
	@psql "$(DATABASE_URL)" -c "DROP TABLE IF EXISTS ratings, users, recipes;"

docker-build:
	@echo "Building docker images (if Dockerfile present)..."
	@docker build -t unthinkable-frontend ./frontend || true
	@docker build -t unthinkable-backend ./backend || true

docker-up:
	@echo "Starting docker compose..."
	@docker compose up -d

docker-restart:
	@echo "Restarting docker compose services..."
	@docker compose restart
