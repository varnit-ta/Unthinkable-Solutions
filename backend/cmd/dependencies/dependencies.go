// Package app provides application initialization and dependency injection.
// It handles database setup, service configuration, and HTTP server initialization.
package app

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"github.com/varnit-ta/smart-recipe-generator/backend/internal/config"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/handlers"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/middleware"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/service"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/vision"
)

// App encapsulates the application dependencies and configuration.
type App struct {
	Config config.Config
	DB     *sql.DB
	Router *chi.Mux
}

// New creates and initializes a new App instance with all dependencies.
// It performs database connection, service setup, and router configuration.
func New(cfg config.Config) (*App, error) {
	app := &App{
		Config: cfg,
	}

	if err := app.initDatabase(); err != nil {
		return nil, err
	}

	app.initRouter()

	return app, nil
}

// initDatabase establishes database connection with retry logic and configures connection pooling.
func (app *App) initDatabase() error {
	db, err := sql.Open("postgres", app.Config.DatabaseURL)
	if err != nil {
		return err
	}

	app.configureConnectionPool(db)

	if err := app.connectWithRetry(db); err != nil {
		db.Close()
		return err
	}

	app.DB = db
	return nil
}

// configureConnectionPool sets up database connection pooling parameters.
func (app *App) configureConnectionPool(db *sql.DB) {
	if app.Config.DBMaxOpenConns > 0 {
		db.SetMaxOpenConns(app.Config.DBMaxOpenConns)
	}
	if app.Config.DBMaxIdleConns > 0 {
		db.SetMaxIdleConns(app.Config.DBMaxIdleConns)
	}
	if app.Config.DBConnMaxIdle > 0 {
		db.SetConnMaxIdleTime(app.Config.DBConnMaxIdle)
	}
	if app.Config.DBConnMaxLife > 0 {
		db.SetConnMaxLifetime(app.Config.DBConnMaxLife)
	}
}

// connectWithRetry attempts to establish a database connection with exponential backoff.
func (app *App) connectWithRetry(db *sql.DB) error {
	backoff := app.Config.DBRetryBackoff
	if backoff <= 0 {
		backoff = 500 * time.Millisecond
	}
	attempts := app.Config.DBRetryMax
	if attempts <= 0 {
		attempts = 8
	}

	var pingErr error
	for i := 0; i < attempts; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			log.Printf("connected to db")
			return nil
		}
		wait := backoff << i
		if wait > 5*time.Second {
			wait = 5 * time.Second
		}
		log.Printf("db ping failed (attempt %d/%d): %v; retrying in %s", i+1, attempts, pingErr, wait)
		time.Sleep(wait)
	}

	return pingErr
}

// initRouter sets up the HTTP router with all middleware and routes.
func (app *App) initRouter() {
	visionService := app.setupVisionService()
	svc := service.NewService(app.DB)
	h := handlers.New(svc, visionService, app.Config.MaxImageSizeMB)
	authH := &handlers.AuthHandler{
		Service:   svc,
		JWTSecret: app.Config.JWTSecret,
		JWTExpiry: app.Config.JWTExpiryHours,
	}

	r := chi.NewRouter()

	r.Use(app.corsMiddleware())
	r.Use(middleware.Logging)

	app.setupRoutes(r, h, authH)

	app.Router = r
}

// setupVisionService initializes the AI vision service for ingredient detection.
func (app *App) setupVisionService() vision.VisionService {
	if app.Config.AIServiceURL != "" {
		log.Printf("Local AI service configured at: %s", app.Config.AIServiceURL)
		return vision.NewLocalAIService(app.Config.AIServiceURL)
	}

	log.Printf("WARNING: No AI service configured - ingredient detection disabled")
	log.Printf("Set AI_SERVICE_URL env var to use local AI service")
	log.Printf("Start local AI service with: docker-compose up ai-service")
	return nil
}

// corsMiddleware configures CORS settings for the application.
func (app *App) corsMiddleware() func(http.Handler) http.Handler {
	allowedOrigins := strings.Split(app.Config.AllowedOrigins, ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

// setupRoutes registers all HTTP endpoints for the application.
func (app *App) setupRoutes(r *chi.Mux, h *handlers.Handler, authH *handlers.AuthHandler) {
	r.Get("/health", healthCheck)

	r.Get("/recipes", h.ListRecipes)
	r.Get("/recipes/{id}", h.GetRecipe)
	r.Post("/match", h.Match)
	r.Post("/detect-ingredients", h.DetectIngredients)

	r.Post("/auth/register", authH.Register)
	r.Post("/auth/login", authH.Login)

	jwtAuth := middleware.JWTAuth(app.Config.JWTSecret)
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

// Run starts the HTTP server on the configured port.
func (app *App) Run() error {
	addr := ":" + app.Config.Port
	log.Printf("starting server on %s", addr)
	return http.ListenAndServe(addr, app.Router)
}

// Close cleans up application resources.
func (app *App) Close() error {
	if app.DB != nil {
		return app.DB.Close()
	}
	return nil
}
