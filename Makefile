SHELL := /bin/sh

.PHONY: help frontend backend sqlc migrate-up migrate-down migrateup migratedown migrateall resetdb docker-build docker-up docker-restart

help:
	@echo "Makefile targets:"
	@echo "  make frontend        # Run frontend dev server (from frontend folder)"
	@echo "  make backend         # Run Go backend server"
	@echo "  make sqlc            # Run sqlc generate (requires sqlc installed)"
	@echo "  make migrate-up      # Apply DB migrations using Go migrate runner"
	@echo "  make migrate-down    # Rollback DB migrations using Go migrate runner"
	@echo "  make migrateup       # Apply migrations using migrate CLI"
	@echo "  make migratedown     # Rollback migrations using migrate CLI"
	@echo "  make migrateall      # Apply all migrations in docker container"
	@echo "  make resetdb         # Reset database and reapply all migrations"
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
	@cd backend/cmd/migrate && \
		DATABASE_URL="$(DATABASE_URL)" go run . up

migrate-seed: migrate-up sqlc
	@echo "migrations applied and sqlc generated"

migrate-down:
	@echo "migrate-down is destructive. Dropping tables..."
	@cd backend/cmd/migrate && \
		DATABASE_URL="$(DATABASE_URL)" go run . down

docker-build:
	@echo "Building docker images (if Dockerfile present)..."
	-@docker build -t unthinkable-frontend ./frontend
	-@docker build -t unthinkable-backend ./backend

docker-up:
	@echo "Starting docker compose..."
	@docker compose up -d

docker-restart:
	@echo "Restarting docker compose services..."
	@docker compose restart

migrateup:
	migrate -path backend/migrations/ -database "$(DATABASE_URL)" up

migratedown:
	migrate -path backend/migrations/ -database "$(DATABASE_URL)" down

migrateall:
	@echo "Running all database migrations..."
	@type backend\migrations\001_create_schema.up.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\002_seed_recipes.up.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\003_users_and_favorites.up.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\004_add_updated_at_recipes.up.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@echo "All migrations completed successfully!"

resetdb:
	@echo "Resetting database and running all migrations..."
	@type backend\migrations\004_add_updated_at_recipes.down.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\003_users_and_favorites.down.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\002_seed_recipes.down.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\001_create_schema.down.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@echo "Running all migrations from scratch..."
	@type backend\migrations\001_create_schema.up.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\002_seed_recipes.up.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\003_users_and_favorites.up.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@type backend\migrations\004_add_updated_at_recipes.up.sql | docker exec -i unthinkablesolutions-db-1 psql -U unthinkable -d unthinkable_recipes
	@echo "Database reset and migrations completed successfully!"
