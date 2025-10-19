// Package main is the entry point for the Smart Recipe Generator backend server.
// It initializes the database connection, configures services, sets up HTTP routes,
// and starts the web server.
package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/varnit-ta/smart-recipe-generator/backend/internal/config"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/handlers"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/middleware"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/service"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/vision"
)

// main is the application entry point.
// It performs the following initialization steps:
// 1. Load configuration from environment variables
// 2. Establish database connection with retry logic
// 3. Configure connection pooling
// 4. Initialize vision service (Hugging Face AI)
// 5. Set up service and handler layers
// 6. Configure HTTP router with middleware and routes
// 7. Start HTTP server
func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	configureConnectionPool(db, cfg)
	connectToDatabase(db, cfg)
	initializeServices(db, cfg)
}

// configureConnectionPool sets up database connection pooling parameters
// to optimize connection reuse and prevent resource exhaustion.
func configureConnectionPool(db *sql.DB, cfg config.Config) {
	if cfg.DBMaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	}
	if cfg.DBMaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	}
	if cfg.DBConnMaxIdle > 0 {
		db.SetConnMaxIdleTime(cfg.DBConnMaxIdle)
	}
	if cfg.DBConnMaxLife > 0 {
		db.SetConnMaxLifetime(cfg.DBConnMaxLife)
	}
}

// connectToDatabase attempts to establish a database connection with
// exponential backoff retry logic to handle temporary connection failures.
func connectToDatabase(db *sql.DB, cfg config.Config) {
	var pingErr error
	backoff := cfg.DBRetryBackoff
	if backoff <= 0 {
		backoff = 500 * time.Millisecond
	}
	attempts := cfg.DBRetryMax
	if attempts <= 0 {
		attempts = 8
	}
	for i := 0; i < attempts; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			log.Printf("connected to db")
			break
		}
		wait := backoff << i
		if wait > 5*time.Second {
			wait = 5 * time.Second
		}
		log.Printf("db ping failed (attempt %d/%d): %v; retrying in %s", i+1, attempts, pingErr, wait)
		time.Sleep(wait)
	}
	if pingErr != nil {
		log.Fatalf("could not connect to db after %d attempts: %v", attempts, pingErr)
	}
}

// initializeServices sets up all application services including
// vision AI, business logic, handlers, and HTTP routing.
func initializeServices(db *sql.DB, cfg config.Config) {
	visionService := setupVisionService(cfg)
	svc := service.NewService(db)
	h := handlers.New(svc, visionService, cfg.MaxImageSizeMB)
	authH := &handlers.AuthHandler{
		Service:   svc,
		JWTSecret: cfg.JWTSecret,
		JWTExpiry: cfg.JWTExpiryHours,
	}

	r := setupRouter(cfg, h, authH)

	addr := ":" + cfg.Port
	log.Printf("starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// setupVisionService initializes the AI vision service for ingredient detection.
// Priority order:
// 1. Local AI service (if AI_SERVICE_URL is set) - preferred for better ingredient extraction
// 2. Hugging Face API (if HUGGINGFACE_TOKEN is set)
// Returns nil if no service is configured.
func setupVisionService(cfg config.Config) vision.VisionService {
	// Prefer local AI service (has better ingredient extraction)
	if cfg.AIServiceURL != "" {
		log.Printf("✅ Local AI service configured at: %s", cfg.AIServiceURL)
		return vision.NewLocalAIService(cfg.AIServiceURL)
	}

	// Fall back to Hugging Face if token is provided
	if cfg.HuggingFaceToken != "" {
		log.Printf("✅ Hugging Face AI service configured with model: %s", cfg.HuggingFaceModel)
		return vision.NewHuggingFaceService(cfg.HuggingFaceToken, cfg.HuggingFaceModel)
	}

	log.Printf("❌ WARNING: No AI service configured - ingredient detection disabled")
	log.Printf("💡 Set AI_SERVICE_URL env var to use local AI service (recommended)")
	log.Printf("💡 Or set HUGGINGFACE_TOKEN env var to use Hugging Face API")
	log.Printf("💡 Start local AI service with: docker-compose up ai-service")
	return nil
}

// setupRouter configures the HTTP router with middleware, CORS, and all application routes.
func setupRouter(cfg config.Config, h *handlers.Handler, authH *handlers.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logging)

	setupRoutes(r, cfg, h, authH)

	return r
}

// setupRoutes registers all HTTP endpoints for the application.
// Routes are organized into public and protected (JWT-authenticated) endpoints.
func setupRoutes(r *chi.Mux, cfg config.Config, h *handlers.Handler, authH *handlers.AuthHandler) {
	r.Get("/health", healthCheck)

	r.Get("/recipes", h.ListRecipes)
	r.Get("/recipes/{id}", h.GetRecipe)
	r.Post("/match", h.Match)
	r.Post("/detect-ingredients", h.DetectIngredients)

	r.Post("/auth/register", authH.Register)
	r.Post("/auth/login", authH.Login)

	jwtAuth := middleware.JWTAuth(cfg.JWTSecret)
	r.With(jwtAuth).Post("/ratings", h.PostRating)
	r.With(jwtAuth).Post("/favorites/{id}", h.AddFavorite)
	r.With(jwtAuth).Delete("/favorites/{id}", h.RemoveFavorite)
	r.With(jwtAuth).Get("/favorites", h.ListFavorites)
	r.With(jwtAuth).Get("/favorites/{id}", h.IsFavorite)
	r.With(jwtAuth).Get("/suggestions", h.GetSuggestions)
}

// healthCheck is a simple endpoint that returns 200 OK to indicate server health.
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
