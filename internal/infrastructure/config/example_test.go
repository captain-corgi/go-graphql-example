package config_test

import (
	"fmt"
	"log"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
)

func ExampleLoad() {
	// Load configuration from files and environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Server will run on port: %s\n", cfg.Server.Port)
	fmt.Printf("Database URL: %s\n", cfg.Database.URL)
	fmt.Printf("Log level: %s\n", cfg.Logging.Level)

	// Output:
	// Server will run on port: 8080
	// Database URL: postgres://user:password@localhost:5432/graphql_service?sslmode=disable
	// Log level: info
}

func ExampleMustLoad() {
	// Load configuration and panic if it fails
	cfg := config.MustLoad()

	fmt.Printf("Server will run on port: %s\n", cfg.Server.Port)
	// Output:
	// Server will run on port: 8080
}
