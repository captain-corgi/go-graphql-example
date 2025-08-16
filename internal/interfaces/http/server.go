package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/generated"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/resolver"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/http/middleware"
)

// Server represents the HTTP server
type Server struct {
	server   *http.Server
	router   *gin.Engine
	config   *config.ServerConfig
	logger   *slog.Logger
	resolver *resolver.Resolver
}

// NewServer creates a new HTTP server with the given configuration and dependencies
func NewServer(cfg *config.ServerConfig, resolver *resolver.Resolver, logger *slog.Logger) *Server {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Recovery middleware to handle panics
	router.Use(gin.Recovery())

	// Add custom middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())

	server := &Server{
		router:   router,
		config:   cfg,
		logger:   logger,
		resolver: resolver,
	}

	server.setupRoutes()
	server.setupHTTPServer()

	return server
}

// setupRoutes configures all the HTTP routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthHandler)

	// GraphQL handler
	graphqlHandler := s.createGraphQLHandler()
	s.router.POST("/query", gin.WrapH(graphqlHandler))

	// GraphQL Playground (only in development)
	playgroundHandler := playground.Handler("GraphQL Playground", "/query")
	s.router.GET("/playground", gin.WrapH(playgroundHandler))
	s.router.GET("/", gin.WrapH(playgroundHandler)) // Redirect root to playground
}

// createGraphQLHandler creates the GraphQL handler with the schema and resolvers
func (s *Server) createGraphQLHandler() *handler.Server {
	// Create the executable schema
	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: s.resolver,
	})

	// Create the GraphQL handler
	srv := handler.NewDefaultServer(schema)

	// Add custom error handling if needed
	// srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
	//     // Custom error handling logic
	//     return graphql.DefaultErrorPresenter(ctx, e)
	// })

	return srv
}

// setupHTTPServer configures the underlying HTTP server
func (s *Server) setupHTTPServer() {
	s.server = &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server",
		slog.String("port", s.config.Port),
		slog.Duration("read_timeout", s.config.ReadTimeout),
		slog.Duration("write_timeout", s.config.WriteTimeout),
		slog.Duration("idle_timeout", s.config.IdleTimeout),
	)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	s.logger.Info("HTTP server stopped")
	return nil
}

// healthHandler handles health check requests
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "graphql-service",
	})
}
