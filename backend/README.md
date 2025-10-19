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

Environment Variables
- `DATABASE_URL` (optional) — Postgres connection string. Default: `postgres://unthinkable:unthinkable@localhost:5432/unthinkable_recipes?sslmode=disable`.
- `PORT` (optional) — HTTP server port. Default: `8081`.
- `JWT_SECRET` (optional) — Secret key for JWT signing. Default: `change-me-to-a-secure-secret`.
- `AI_SERVICE_URL` (required) — URL for local Python AI service. Default: `http://localhost:8000`. Use `http://ai-service:8000` in Docker.
- `MAX_IMAGE_SIZE_MB` (optional) — Maximum image upload size in MB. Default: `10`.
- `ALLOWED_ORIGINS` (optional) — Comma-separated list of allowed CORS origins. Default includes localhost ports.

## AI Service Configuration

The backend connects to a local Python AI service for ingredient detection from images.

### Configuration

Set the AI service URL via environment variable:

```env
# For Docker Compose (service name)
AI_SERVICE_URL=http://ai-service:8000

# For local development
AI_SERVICE_URL=http://localhost:8000
```

### Starting the AI Service

**With Docker Compose (Recommended):**
```cmd
docker-compose up ai-service
```

**Standalone:**
```cmd
cd ai-service
pip install -r requirements.txt
python main.py
```

The backend will automatically connect to the AI service at the configured URL.
