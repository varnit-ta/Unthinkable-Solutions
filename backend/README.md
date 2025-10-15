Go backend scaffold (sqlc)

This folder contains a minimal Go backend scaffold that uses `sqlc` to generate type-safe DB code from SQL.

Files of interest:
- `cmd/server/main.go` — minimal HTTP server with a `/health` and `/recipes` endpoint (reads directly using database/sql before sqlc generation).
- `migrations/001_create_schema.sql` — DB schema used by sqlc and migrations.
- `queries/recipes.sql` — SQL queries for sqlc to generate typed code.
- `sqlc.yaml` — sqlc configuration.
- `go.mod` — module file with basic deps.

Getting started
1. Install Go (1.20+), sqlc (https://sqlc.dev), and `go` tools.
2. Start Postgres using the top-level `docker-compose.yml` in the project root.
3. Generate sqlc code:

```cmd
cd backend
sqlc generate
```

4. Build and run the server:

```cmd
go run ./cmd/server
```

Environment
- `DATABASE_URL` (optional) — Postgres connection string. Default used in `main.go` is `postgres://unthinkable:unthinkable@localhost:5432/unthinkable_recipes?sslmode=disable`.
