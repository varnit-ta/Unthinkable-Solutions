# Backend Documentation

## Overview

The Smart Recipe Generator backend is a RESTful API service built with Go that provides recipe management, AI-powered ingredient detection, user authentication, and personalized recommendations. It uses PostgreSQL for data persistence and integrates with Hugging Face AI for vision-based ingredient detection.

## Technology Stack

- **Go 1.21+** - Programming language
- **Chi Router** - HTTP router and middleware
- **PostgreSQL** - Relational database
- **SQLC** - Type-safe SQL code generator
- **JWT** - JSON Web Tokens for authentication
- **Bcrypt** - Password hashing
- **Hugging Face API** - AI vision/image captioning
- **Docker** - Containerization

## Architecture

### Layered Architecture

```
┌─────────────────────────────────────────────┐
│           HTTP Handlers Layer               │
│  (Request/Response, Validation, Errors)     │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│          Service/Business Logic             │
│  (Recipe Matching, Recommendations, etc.)   │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│         Data Access Layer (SQLC)            │
│       (Type-safe Database Queries)          │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│           PostgreSQL Database               │
└─────────────────────────────────────────────┘
```

### External Integrations

- **Hugging Face AI** - Image captioning and ingredient detection
- **JWT** - Stateless authentication

## Core Components

### 1. Main Application (`cmd/server/main.go`)

**Purpose**: Application initialization and server setup

**Key Features**:
- Environment configuration loading
- Database connection with retry logic and connection pooling
- Vision service initialization
- HTTP router setup with middleware
- Route registration
- Server startup

**Configuration Options**:
- Database pooling (max connections, idle connections, timeouts)
- JWT settings (secret, expiry)
- Vision API settings
- Port configuration
- Retry logic for database connections

**Database Connection Pooling**:
```go
db.SetMaxOpenConns(cfg.DBMaxOpenConns)      // Max concurrent connections
db.SetMaxIdleConns(cfg.DBMaxIdleConns)      // Max idle connections
db.SetConnMaxIdleTime(cfg.DBConnMaxIdle)    // Max idle time
db.SetConnMaxLifetime(cfg.DBConnMaxLife)    // Max connection lifetime
```

### 2. Configuration (`internal/config/config.go`)

**Purpose**: Centralized configuration management from environment variables

**Configuration Structure**:
```go
type Config struct {
    DatabaseURL       string        // PostgreSQL connection string
    Port              string        // HTTP server port
    JWTSecret         string        // JWT signing secret
    JWTExpiryHours    int          // Token expiration time
    DBMaxOpenConns    int          // Max open DB connections
    DBMaxIdleConns    int          // Max idle DB connections
    DBConnMaxIdle     time.Duration // Max idle connection time
    DBConnMaxLife     time.Duration // Max connection lifetime
    DBRetryMax        int          // Max DB connection retries
    DBRetryBackoff    time.Duration // Retry backoff duration
    HuggingFaceAPIKey string        // HuggingFace API key
    MaxImageSizeMB    int          // Max image upload size
}
```

**Environment Variables**:
- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - Server port (default: 8081)
- `JWT_SECRET` - JWT signing secret
- `DB_MAX_OPEN_CONNS` - Max open connections (default: 20)
- `DB_MAX_IDLE_CONNS` - Max idle connections (default: 10)
- `DB_CONN_MAX_IDLE` - Max idle time (default: 15m)
- `DB_CONN_MAX_LIFE` - Max connection lifetime (default: 1h)
- `DB_RETRY_MAX` - Max retry attempts (default: 8)
- `DB_RETRY_BACKOFF` - Initial retry delay (default: 500ms)
- `HUGGINGFACE_API_KEY` - Hugging Face API key
- `MAX_IMAGE_SIZE_MB` - Max image size (default: 10MB)

### 3. Authentication (`internal/auth/auth.go`)

**Purpose**: JWT token management and password security

**Key Functions**:

#### `HashPassword(password string) (string, error)`
Hashes a password using bcrypt with default cost (10).

#### `VerifyPassword(hash, password string) error`
Verifies a password against its bcrypt hash.

#### `GenerateJWT(secret string, userID int, expiryHours int) (string, error)`
Generates a JWT token with user ID claim.

**Token Structure**:
```go
type Claims struct {
    UserID int `json:"userId"`
    jwt.RegisteredClaims
}
```

#### `ParseJWT(secret, tokenStr string) (*Claims, error)`
Parses and validates a JWT token, returns claims.

#### `RandomSecret() (string, error)`
Generates a cryptographically secure random secret (32 bytes, base64 encoded).

**Security Features**:
- Bcrypt password hashing (cost factor 10)
- HMAC-SHA256 JWT signing
- Token expiration validation
- Cryptographically secure random secrets

### 4. Database Layer (`internal/db/`)

**Purpose**: Type-safe database access using SQLC

**SQLC Generated Files**:
- `db.go` - Base queries interface
- `models.go` - Database model structs
- `recipes.sql.go` - Recipe-related queries
- `users.sql.go` - User management queries
- `favorites.sql.go` - Favorites management queries

**Key Database Models**:

```go
type Recipe struct {
    ID               int32
    Title            string
    Description      sql.NullString
    Ingredients      []string
    Steps            []string
    Tags             []string
    CookTimeMinutes  sql.NullInt32
    TotalTimeMinutes sql.NullInt32
    Servings         sql.NullInt32
    Difficulty       sql.NullString
    Cuisine          sql.NullString
    DietType         sql.NullString
    AverageRating    sql.NullString
    CreatedAt        time.Time
    UpdatedAt        sql.NullTime
}

type User struct {
    ID           int32
    Username     sql.NullString
    Email        sql.NullString
    PasswordHash sql.NullString
    CreatedAt    time.Time
}

type Favorite struct {
    ID        int32
    UserID    sql.NullInt32
    RecipeID  sql.NullInt32
    CreatedAt time.Time
}

type Rating struct {
    ID        int32
    UserID    sql.NullInt32
    RecipeID  sql.NullInt32
    Rating    sql.NullInt32
    CreatedAt time.Time
}
```

**Benefits of SQLC**:
- Compile-time SQL validation
- Type-safe query execution
- Automatic NULL handling
- Performance (no reflection)
- Maintainable SQL queries

### 5. Service Layer (`internal/service/service.go`)

**Purpose**: Business logic and data orchestration

**Key Methods**:

#### Recipe Management

**`ListRecipes(ctx, limit, offset int) ([]db.ListRecipesRow, error)`**
- Fetches paginated list of recipes
- Returns basic recipe information

**`GetRecipe(ctx, id int) (db.GetRecipeByIDRow, error)`**
- Fetches complete recipe details by ID
- Includes all ingredients, steps, and metadata

**`SearchAndFilterRecipes(ctx, query, diet, difficulty string, maxTime *int, cuisine string, limit, offset int) ([]db.SearchRecipesRow, error)`**
- Full-text search by title/tags
- Multiple filter criteria:
  - Diet type (matches against tags)
  - Difficulty level (exact match)
  - Cuisine type (exact match)
  - Maximum cooking time
- Pagination support

#### Recipe Matching

**`MatchRecipes(ctx, detected []string, limit, offset int) ([]RecipeSummary, error)`**
- Scores recipes based on ingredient overlap
- Uses tag matching and title search
- Returns sorted results (highest score first)

**`MatchWithFilters(ctx, ingredients []string, filters MatchFilters) ([]RecipeWithScore, error)`**
- Combines filtering and ingredient matching
- Applies diet, difficulty, cuisine, and time filters
- Scores filtered candidates by ingredient overlap

**Scoring Algorithm**:
1. Tag matching: +1 point per matched tag
2. Title matching: +1 point per ingredient in title
3. Sort by descending score

#### User Management

**`CreateUser(ctx, username, email, password string) (db.CreateUserRow, error)`**
- Creates new user account
- Hashes password before storage
- Returns created user

**`Authenticate(ctx, email, password string) (db.GetUserByEmailRow, error)`**
- Validates user credentials
- Returns user record on success

#### Ratings

**`AddRating(ctx, userID sql.NullInt32, recipeID, rating int) (db.Rating, error)`**
- Adds or updates recipe rating
- Validates rating value (1-5)

#### Favorites

**`AddFavorite(ctx, userID, recipeID int) (db.Favorite, error)`**
- Adds recipe to user's favorites
- Handles duplicate prevention

**`RemoveFavorite(ctx, userID, recipeID int) error`**
- Removes recipe from favorites

**`ListFavorites(ctx, userID int) ([]db.ListFavoritesByUserRow, error)`**
- Returns all user's favorite recipes with details

**`IsFavorite(ctx, userID, recipeID int) (bool, error)`**
- Checks if a recipe is in user's favorites

#### Recommendations

**`GetSuggestions(ctx, userID int, limit int) ([]RecipeWithScore, error)`**
- Content-based filtering algorithm
- Analyzes user's favorite recipes
- Extracts common tags/preferences
- Scores candidate recipes by tag overlap
- Returns top N recommendations

**Recommendation Algorithm**:
1. Fetch user's favorite recipes
2. Build tag frequency map from favorites
3. Fetch broad set of candidate recipes
4. Score each candidate by tag overlap with favorites
5. Sort by score (descending)
6. Return top N results

### 6. HTTP Handlers (`internal/handlers/`)

**Purpose**: HTTP request/response handling

#### Public Endpoints

**`GET /health`**
- Health check endpoint
- Returns 200 OK

**`GET /recipes`**
- List recipes with optional filters
- Query parameters:
  - `q` - Search query
  - `diet` - Diet type filter
  - `difficulty` - Difficulty filter
  - `cuisine` - Cuisine filter
  - `maxTime` - Maximum cooking time
  - `limit` - Results per page (max 200, default 50)
  - `offset` - Pagination offset

**`GET /recipes/{id}`**
- Get recipe details by ID
- Returns 404 if not found

**`POST /match`**
- Find recipes matching ingredients
- Request body:
  ```json
  {
    "detectedIngredients": ["tomato", "onion", "garlic"]
  }
  ```
- Query parameters: same filters as `/recipes`
- Returns recipes with match scores

**`POST /detect-ingredients`**
- AI-powered ingredient detection from image
- Multipart form data with `image` field
- Supported formats: JPEG, PNG, GIF, WebP
- Max size: configurable (default 10MB)
- Returns:
  ```json
  {
    "detectedIngredients": ["tomato", "onion"],
    "confidence": 0.85,
    "provider": "huggingface",
    "caption": "a plate of tomatoes and onions"
  }
  ```

#### Authentication Endpoints

**`POST /auth/register`**
- Create new user account
- Request body:
  ```json
  {
    "username": "johndoe",
    "email": "john@example.com",
    "password": "secure123"
  }
  ```
- Returns JWT token

**`POST /auth/login`**
- Authenticate user
- Request body:
  ```json
  {
    "email": "john@example.com",
    "password": "secure123"
  }
  ```
- Returns JWT token

#### Protected Endpoints (Require JWT)

**`POST /ratings`**
- Submit recipe rating
- Request body:
  ```json
  {
    "recipeId": 123,
    "rating": 5
  }
  ```

**`POST /favorites/{id}`**
- Add recipe to favorites
- URL parameter: recipe ID

**`DELETE /favorites/{id}`**
- Remove recipe from favorites
- URL parameter: recipe ID

**`GET /favorites`**
- List user's favorite recipes

**`GET /favorites/{id}`**
- Check if recipe is favorited
- Returns:
  ```json
  {
    "isFavorite": true
  }
  ```

**`GET /suggestions`**
- Get personalized recipe recommendations
- Query parameters:
  - `limit` - Number of suggestions (max 100, default 10)

### 7. Middleware (`internal/middleware/`)

#### JWT Authentication (`auth.go`)

**`JWTAuth(secret string) func(http.Handler) http.Handler`**
- Validates JWT tokens from Authorization header
- Format: `Bearer <token>`
- Extracts user ID from token claims
- Stores user ID in request context
- Returns 401 Unauthorized on failure

**Context Key**: `UserIDKey` - Access user ID in handlers via `r.Context().Value(middleware.UserIDKey)`

#### Request Logging (`logging.go`)

**`Logging(next http.Handler) http.Handler`**
- Logs all HTTP requests
- Format: `{METHOD} {PATH} {DURATION}`
- Example: `GET /recipes/123 15.2ms`

### 8. Vision Service (`internal/vision/`)

**Purpose**: AI-powered ingredient detection from images

#### Interface (`vision.go`)

```go
type VisionService interface {
    DetectIngredients(ctx, imageData []byte, filename string) (*DetectionResult, error)
}

type DetectionResult struct {
    Ingredients []string               // Detected ingredient names
    RawResponse string                 // AI model's raw caption
    Confidence  float64                // Detection confidence (0-1)
    Provider    string                 // Service provider name
    Metadata    map[string]interface{} // Additional metadata
}
```

#### Hugging Face Implementation (`huggingface.go`)

**Model**: Salesforce BLIP Image Captioning Large
- Excellent for food/ingredient recognition
- Generates descriptive captions
- API endpoint: `https://api-inference.huggingface.co/models/Salesforce/blip-image-captioning-large`

**Features**:
- Automatic retry on model loading
- Timeout handling (30 seconds)
- Error classification
- Confidence scoring

**Confidence Calculation**:
- Base confidence: 0.7
- +0.1 for 3+ ingredients
- +0.1 for 5+ ingredients
- +0.05 for food-related keywords
- Capped at 0.95

#### Ingredient Parser (`parser.go`)

**`ParseIngredientsFromText(text string) []string`**
- Extracts ingredient names from AI caption
- Uses comprehensive ingredient database (200+ items)
- Supports multi-word ingredients ("bell pepper", "soy sauce")
- Normalizes plural forms
- Removes noise words (descriptors, measurements)

**Ingredient Database Categories**:
- Vegetables (40+ items)
- Proteins (20+ items)
- Dairy (10+ items)
- Grains & Pasta (10+ items)
- Fruits (25+ items)
- Legumes & Nuts (10+ items)
- Condiments & Seasonings (20+ items)

**Text Processing**:
1. Convert to lowercase
2. Remove noise words (adjectives, measurements)
3. Split into words
4. Check 1-word, 2-word, and 3-word combinations
5. Match against ingredient database
6. Normalize and deduplicate

## Database Schema

### Tables

#### `recipes`
```sql
id                SERIAL PRIMARY KEY
title             VARCHAR(255) NOT NULL
description       TEXT
ingredients       TEXT[] -- Array of ingredient strings
steps             TEXT[] -- Array of instruction steps
tags              TEXT[] -- Tags for categorization
cook_time_minutes INTEGER
total_time_minutes INTEGER
servings          INTEGER
difficulty        VARCHAR(50) -- easy, medium, hard
cuisine           VARCHAR(100) -- italian, mexican, etc.
diet_type         VARCHAR(100) -- vegetarian, vegan, etc.
average_rating    NUMERIC(3,2)
created_at        TIMESTAMP DEFAULT NOW()
updated_at        TIMESTAMP
```

#### `users`
```sql
id            SERIAL PRIMARY KEY
username      VARCHAR(100) UNIQUE
email         VARCHAR(255) UNIQUE NOT NULL
password_hash VARCHAR(255) NOT NULL
created_at    TIMESTAMP DEFAULT NOW()
```

#### `favorites`
```sql
id         SERIAL PRIMARY KEY
user_id    INTEGER REFERENCES users(id)
recipe_id  INTEGER REFERENCES recipes(id)
created_at TIMESTAMP DEFAULT NOW()
UNIQUE(user_id, recipe_id)
```

#### `ratings`
```sql
id         SERIAL PRIMARY KEY
user_id    INTEGER REFERENCES users(id)
recipe_id  INTEGER REFERENCES recipes(id)
rating     INTEGER CHECK (rating >= 1 AND rating <= 5)
created_at TIMESTAMP DEFAULT NOW()
UNIQUE(user_id, recipe_id)
```

### Indexes

```sql
CREATE INDEX idx_recipes_tags ON recipes USING GIN(tags);
CREATE INDEX idx_recipes_difficulty ON recipes(difficulty);
CREATE INDEX idx_recipes_cuisine ON recipes(cuisine);
CREATE INDEX idx_favorites_user_id ON favorites(user_id);
CREATE INDEX idx_ratings_recipe_id ON ratings(recipe_id);
```

## API Response Formats

### Recipe Detail Response
```json
{
  "id": 1,
  "title": "Spaghetti Carbonara",
  "description": "Classic Italian pasta dish",
  "ingredients": ["pasta", "eggs", "bacon", "parmesan"],
  "steps": ["Boil pasta", "Cook bacon", "Mix with eggs"],
  "tags": ["italian", "pasta", "quick"],
  "cookTimeMinutes": 20,
  "totalTimeMinutes": 30,
  "servings": 4,
  "difficulty": "easy",
  "cuisine": "italian",
  "dietType": null,
  "averageRating": "4.5"
}
```

### Recipe with Score
```json
{
  "id": 1,
  "title": "Tomato Soup",
  "score": 3,
  ...other recipe fields
}
```

### Favorite Recipe Response
```json
{
  "id": 1,
  "recipeId": 123,
  "title": "Pasta Primavera",
  "description": "Fresh vegetable pasta",
  "cuisine": "italian",
  "difficulty": "medium",
  "cookTimeMinutes": 25,
  "averageRating": "4.2"
}
```

## Error Handling

### HTTP Status Codes

- **200 OK** - Successful request
- **201 Created** - Resource created
- **204 No Content** - Successful deletion
- **400 Bad Request** - Invalid input
- **401 Unauthorized** - Authentication required or failed
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server error
- **503 Service Unavailable** - Vision service not configured

### Error Response Format
```json
{
  "message": "error description"
}
```

## CORS Configuration

Allowed origins:
- `http://localhost:5173` (Vite dev server)
- `http://localhost:3000` (Alternative dev)
- `http://localhost:4173` (Vite preview)

Allowed methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
Allowed headers: Accept, Authorization, Content-Type, X-CSRF-Token

## Development

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- SQLC (for regenerating queries)
- Migrate (for database migrations)

### Setup

1. **Install dependencies**:
```bash
go mod download
```

2. **Set environment variables**:
```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/recipes?sslmode=disable"
export JWT_SECRET="your-secure-secret"
export HUGGINGFACE_API_KEY="your-api-key"
```

3. **Run migrations**:
```bash
migrate -path migrations -database $DATABASE_URL up
```

4. **Generate SQLC code** (if queries changed):
```bash
sqlc generate
```

### Running the Server

```bash
go run cmd/server/main.go
```

Server starts on `http://localhost:8081`

### Building

```bash
go build -o server cmd/server/main.go
```

### Running with Docker

```bash
docker build -t recipe-backend .
docker run -p 8081:8081 \
  -e DATABASE_URL="..." \
  -e JWT_SECRET="..." \
  -e HUGGINGFACE_API_KEY="..." \
  recipe-backend
```

## Testing

### Unit Tests

Test individual components:

```bash
go test ./internal/auth/...
go test ./internal/service/...
go test ./internal/vision/...
```

### Integration Tests

Test with database:

```bash
go test -tags=integration ./...
```

### API Testing with curl

**Register user**:
```bash
curl -X POST http://localhost:8081/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"test123"}'
```

**Login**:
```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}'
```

**Get recipes**:
```bash
curl http://localhost:8081/recipes?limit=10
```

**Search recipes**:
```bash
curl "http://localhost:8081/recipes?q=pasta&cuisine=italian&difficulty=easy"
```

**Detect ingredients**:
```bash
curl -X POST http://localhost:8081/detect-ingredients \
  -F "image=@/path/to/food.jpg"
```

**Add favorite** (requires auth):
```bash
curl -X POST http://localhost:8081/favorites/123 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Performance Considerations

### Database Optimization

1. **Connection Pooling**:
   - Max open connections: 20
   - Max idle connections: 10
   - Connection reuse for efficiency

2. **Indexes**:
   - GIN index on tags for fast array searches
   - B-tree indexes on frequently filtered columns

3. **Query Optimization**:
   - LIMIT/OFFSET for pagination
   - Selective column fetching
   - Efficient JOIN operations

### Caching Strategies

Consider implementing:
- Recipe list caching (Redis)
- User session caching
- Ingredient database in-memory caching

### Rate Limiting

Recommended for production:
- Per-IP rate limiting
- Per-user rate limiting for authenticated endpoints
- Vision API rate limiting (external API costs)

## Security Best Practices

### Implemented

✅ Password hashing with bcrypt
✅ JWT-based stateless authentication
✅ SQL injection prevention (parameterized queries via SQLC)
✅ CORS configuration
✅ Input validation
✅ File upload size limits
✅ Image type validation

### Recommended Additions

- Rate limiting
- Request size limits
- HTTPS/TLS in production
- Secrets management (HashiCorp Vault, AWS Secrets Manager)
- API key rotation
- SQL query logging
- Security headers (helmet middleware)

## Monitoring & Observability

### Logging

Current: Basic request logging with duration

Recommendations:
- Structured logging (JSON format)
- Log levels (DEBUG, INFO, WARN, ERROR)
- Correlation IDs for request tracking
- Error tracking (Sentry, Rollbar)

### Metrics

Consider adding:
- Request count by endpoint
- Response time percentiles
- Error rate
- Database connection pool metrics
- Vision API usage/costs

### Health Checks

Current: `/health` endpoint

Enhancements:
- Database connectivity check
- Vision API availability check
- Disk space check
- Memory usage check

## Deployment

### Environment Variables for Production

```bash
# Database
DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_IDLE=10m
DB_CONN_MAX_LIFE=30m

# Server
PORT=8081

# Authentication
JWT_SECRET=<secure-random-secret>
JWT_EXPIRY_HOURS=48

# Vision AI
HUGGINGFACE_API_KEY=<your-api-key>
MAX_IMAGE_SIZE_MB=10

# Retry
DB_RETRY_MAX=10
DB_RETRY_BACKOFF=1s
```

### Docker Compose Example

```yaml
version: '3.8'
services:
  db:
    image: postgres:14
    environment:
      POSTGRES_DB: recipes
      POSTGRES_USER: recipeuser
      POSTGRES_PASSWORD: securepass
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    ports:
      - "8081:8081"
    environment:
      DATABASE_URL: postgres://recipeuser:securepass@db:5432/recipes?sslmode=disable
      JWT_SECRET: ${JWT_SECRET}
      HUGGINGFACE_API_KEY: ${HUGGINGFACE_API_KEY}
    depends_on:
      - db

volumes:
  pgdata:
```

### Kubernetes Deployment

Consider:
- Horizontal pod autoscaling
- Database replica sets
- Secret management with Kubernetes Secrets
- Ingress with TLS termination
- Health probes (liveness, readiness)

## Troubleshooting

### Common Issues

**Database Connection Failures**:
- Check `DATABASE_URL` format
- Verify PostgreSQL is running
- Check network connectivity
- Review connection pool settings
- Check logs for retry attempts

**JWT Authentication Errors**:
- Verify `JWT_SECRET` is set
- Check token expiration
- Ensure correct Authorization header format
- Validate token signature

**Vision API Errors**:
- Confirm `HUGGINGFACE_API_KEY` is set
- Check API quota/limits
- Verify image format and size
- Review model availability status

**High Database Load**:
- Review slow query logs
- Check connection pool metrics
- Consider read replicas
- Optimize frequently-used queries

## Future Enhancements

### Features

1. **Advanced Search**:
   - Full-text search with PostgreSQL FTS
   - Elasticsearch integration
   - Fuzzy matching for typos

2. **Recipe Management**:
   - User-submitted recipes
   - Recipe editing and versioning
   - Recipe collections/cookbooks

3. **Social Features**:
   - Recipe sharing
   - Comments and reviews
   - Follow users
   - Activity feed

4. **Analytics**:
   - Popular recipes tracking
   - User behavior analysis
   - Recommendation improvement via ML

5. **Notifications**:
   - Email notifications
   - Push notifications
   - Weekly recipe suggestions

### Technical Improvements

1. **Caching Layer**: Redis for frequently accessed data
2. **Message Queue**: RabbitMQ/Kafka for async processing
3. **GraphQL API**: Alternative to REST
4. **WebSocket**: Real-time updates
5. **Microservices**: Split into smaller services
6. **API Versioning**: Support multiple API versions
7. **OpenAPI/Swagger**: Auto-generated API documentation
8. **End-to-End Encryption**: For sensitive data

## API Documentation Tools

Consider adding:
- **Swagger/OpenAPI** - Interactive API documentation
- **Postman Collection** - Shareable API collection
- **API Blueprint** - Human-readable API docs

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Chi Router](https://github.com/go-chi/chi)
- [SQLC](https://sqlc.dev/)
- [PostgreSQL](https://www.postgresql.org/docs/)
- [JWT.io](https://jwt.io/)
- [Hugging Face API](https://huggingface.co/docs/api-inference)

## License

Part of the Smart Recipe Generator project.
