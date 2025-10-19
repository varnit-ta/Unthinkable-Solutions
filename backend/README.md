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
- `HUGGINGFACE_TOKEN` (optional) — Hugging Face API token for ingredient detection. See [HUGGINGFACE_SETUP.md](./HUGGINGFACE_SETUP.md) for setup instructions.
- `HUGGINGFACE_MODEL` (optional) — Hugging Face model ID. Default: `Salesforce/blip-image-captioning-large`.
- `AI_SERVICE_URL` (optional) — URL for local Python AI service. Default: `http://localhost:8000`.

## AI Service Configuration

The backend supports two options for ingredient detection from images:

### Option 1: Hugging Face (Recommended) ✨

Use Hugging Face's free Inference API - no local setup needed!

1. Get a free API token from [huggingface.co](https://huggingface.co/settings/tokens)
2. Set environment variable: `HUGGINGFACE_TOKEN=hf_your_token_here`
3. Start the backend

See [HUGGINGFACE_SETUP.md](./HUGGINGFACE_SETUP.md) for detailed setup instructions.

### Option 2: Local AI Service

Use the local Python service (requires Docker):

```cmd
docker-compose up ai-service
```

The backend automatically uses Hugging Face if `HUGGINGFACE_TOKEN` is set, otherwise falls back to the local service.
