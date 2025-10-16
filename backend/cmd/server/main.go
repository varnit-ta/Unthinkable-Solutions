package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"github.com/varnit-ta/smart-recipe-generator/backend/internal/config"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/handlers"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/middleware"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// Configure DB pooling
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

	// Retry DB connection with backoff
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
		wait := backoff << i // exponential backoff
		if wait > 5*time.Second {
			wait = 5 * time.Second
		}
		log.Printf("db ping failed (attempt %d/%d): %v; retrying in %s", i+1, attempts, pingErr, wait)
		time.Sleep(wait)
	}
	if pingErr != nil {
		log.Fatalf("could not connect to db after %d attempts: %v", attempts, pingErr)
	}

	svc := service.NewService(db)
	h := handlers.New(svc)
	authH := &handlers.AuthHandler{Service: svc, JWTSecret: cfg.JWTSecret, JWTExpiry: cfg.JWTExpiryHours}

	r := chi.NewRouter()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logging)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("ok")) })
	r.Get("/recipes", h.ListRecipes)
	r.Get("/recipes/{id}", h.GetRecipe)
	r.Post("/match", h.Match)
	r.Post("/detect-ingredients", h.DetectIngredients)
	r.With(middleware.JWTAuth(cfg.JWTSecret)).Post("/ratings", h.PostRating)
	// auth
	r.Post("/auth/register", authH.Register)
	r.Post("/auth/login", authH.Login)
	// favorites (protected)
	r.With(middleware.JWTAuth(cfg.JWTSecret)).Post("/favorites/{id}", h.AddFavorite)
	r.With(middleware.JWTAuth(cfg.JWTSecret)).Delete("/favorites/{id}", h.RemoveFavorite)
	r.With(middleware.JWTAuth(cfg.JWTSecret)).Get("/favorites", h.ListFavorites)
	r.With(middleware.JWTAuth(cfg.JWTSecret)).Get("/suggestions", h.GetSuggestions)

	addr := ":" + cfg.Port
	log.Printf("starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
