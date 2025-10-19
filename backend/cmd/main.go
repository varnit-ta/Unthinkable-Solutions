// Package main is the entry point for the Smart Recipe Generator backend server.
package main

import (
	"log"

	"github.com/joho/godotenv"

	app "github.com/varnit-ta/smart-recipe-generator/backend/cmd/dependencies"
	"github.com/varnit-ta/smart-recipe-generator/backend/internal/config"
)

// main is the application entry point.
// It loads configuration and starts the application.
func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	defer application.Close()

	if err := application.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
