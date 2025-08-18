package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	authApp "github.com/captain-corgi/go-graphql-example/internal/application/auth"
	"github.com/captain-corgi/go-graphql-example/internal/application/user"
	authInfra "github.com/captain-corgi/go-graphql-example/internal/infrastructure/auth"
	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/database"
	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/persistence/sql"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/resolver"
	httpserver "github.com/captain-corgi/go-graphql-example/internal/interfaces/http"
)

// Application represents the main application with all its dependencies
type Application struct {
	config    *config.Config
	logger    *slog.Logger
	dbManager *database.Manager
	server    *httpserver.Server
}

// NewApplication creates a new application instance with all dependencies wired
func NewApplication() (*Application, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Initialize logger
	logger := initLogger(cfg.Logging)

	// Initialize database manager
	dbManager, err := database.NewManager(cfg.Database, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database manager: %w", err)
	}

	// Initialize database (run migrations)
	ctx := context.Background()
	if err := dbManager.Initialize(ctx, "migrations"); err != nil {
		dbManager.Close() // Clean up on error
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Perform startup health checks
	if err := performStartupHealthChecks(ctx, dbManager, logger); err != nil {
		dbManager.Close() // Clean up on error
		return nil, fmt.Errorf("startup health checks failed: %w", err)
	}

	// Initialize repositories
	userRepo := sql.NewUserRepository(dbManager.DB, logger)
	sessionRepo := sql.NewSessionRepository(dbManager.DB.DB) // Access the embedded sql.DB

	// Initialize infrastructure services
	jwtService := authInfra.NewJWTService(
		cfg.Auth.JWTSecret,
		cfg.Auth.AccessTokenTTL,
		cfg.Auth.RefreshTokenTTL,
		cfg.Server.Name,
	)
	passwordService := authInfra.NewPasswordService()

	// Initialize application services
	userService := user.NewService(userRepo, logger)
	authService := authApp.NewService(userRepo, sessionRepo, jwtService, passwordService, logger)

	// Initialize resolver with all dependencies
	resolver := resolver.NewResolver(userService, authService, logger)

	// Create HTTP server
	server := httpserver.NewServer(&cfg.Server, resolver, logger)

	return &Application{
		config:    cfg,
		logger:    logger,
		dbManager: dbManager,
		server:    server,
	}, nil
}

// Start starts the application
func (app *Application) Start() error {
	app.logger.Info("Starting GraphQL service",
		slog.String("version", "1.0.0"),
		slog.String("port", app.config.Server.Port))

	// Channel to receive server errors
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		app.logger.Info("HTTP server starting", slog.String("address", ":"+app.config.Server.Port))
		if err := app.server.Start(); err != nil {
			serverErrors <- err
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server failed to start: %w", err)
	case sig := <-quit:
		app.logger.Info("Received shutdown signal", slog.String("signal", sig.String()))
	}

	return nil
}

// Stop gracefully stops the application
func (app *Application) Stop() error {
	app.logger.Info("Initiating graceful shutdown...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop HTTP server
	if err := app.server.Stop(ctx); err != nil {
		app.logger.Error("Failed to stop HTTP server gracefully", slog.String("error", err.Error()))
		return fmt.Errorf("failed to stop server: %w", err)
	}

	// Close database connections
	if err := app.dbManager.Close(); err != nil {
		app.logger.Error("Failed to close database manager", slog.String("error", err.Error()))
		return fmt.Errorf("failed to close database: %w", err)
	}

	app.logger.Info("Application shutdown completed successfully")
	return nil
}

func main() {
	// Create application with all dependencies
	app, err := NewApplication()
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Ensure cleanup on exit
	defer func() {
		if err := app.Stop(); err != nil {
			app.logger.Error("Failed to stop application cleanly", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Start the application
	if err := app.Start(); err != nil {
		app.logger.Error("Application failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

// performStartupHealthChecks validates that all critical components are working
func performStartupHealthChecks(ctx context.Context, dbManager *database.Manager, logger *slog.Logger) error {
	logger.Info("Performing startup health checks...")

	// Check database connectivity
	if err := dbManager.Health(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	logger.Info("Database health check passed")

	// Add more health checks here as needed
	// For example: external service connectivity, cache availability, etc.

	logger.Info("All startup health checks passed")
	return nil
}

// initLogger initializes the structured logger based on configuration
func initLogger(cfg config.LoggingConfig) *slog.Logger {
	var handler slog.Handler

	// Parse log level
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Choose handler based on format
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
