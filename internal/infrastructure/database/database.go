package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/config"
)

// Manager manages database connections, migrations, and transactions
type Manager struct {
	DB               *DB
	TxManager        *TxManager
	MigrationManager *MigrationManager
	logger           *slog.Logger
}

// NewManager creates a new database manager with all components
func NewManager(cfg config.DatabaseConfig, logger *slog.Logger) (*Manager, error) {
	// Create database connection
	db, err := NewConnection(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create transaction manager
	txManager := NewTxManager(db, logger)

	// Create migration manager
	migrationManager := NewMigrationManager(db, logger)

	return &Manager{
		DB:               db,
		TxManager:        txManager,
		MigrationManager: migrationManager,
		logger:           logger,
	}, nil
}

// Initialize sets up the database by running migrations
func (m *Manager) Initialize(ctx context.Context, migrationsPath string) error {
	m.logger.InfoContext(ctx, "Initializing database", "migrations_path", migrationsPath)

	// Run migrations
	if err := m.MigrationManager.RunMigrations(ctx, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.InfoContext(ctx, "Database initialization completed successfully")
	return nil
}

// Close closes all database connections
func (m *Manager) Close() error {
	m.logger.Info("Closing database manager")
	return m.DB.Close()
}

// Health checks the health of all database components
func (m *Manager) Health(ctx context.Context) error {
	return m.DB.Health(ctx)
}
