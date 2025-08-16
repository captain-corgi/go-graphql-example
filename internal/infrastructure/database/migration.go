package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	db     *DB
	logger *slog.Logger
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *DB, logger *slog.Logger) *MigrationManager {
	return &MigrationManager{
		db:     db,
		logger: logger,
	}
}

// RunMigrations executes all pending database migrations
func (mm *MigrationManager) RunMigrations(ctx context.Context, migrationsPath string) error {
	mm.logger.InfoContext(ctx, "Starting database migrations", "path", migrationsPath)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(mm.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Get current version
	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if dirty {
		mm.logger.WarnContext(ctx, "Database is in dirty state", "version", currentVersion)
		return fmt.Errorf("database is in dirty state at version %d", currentVersion)
	}

	mm.logger.InfoContext(ctx, "Current migration version", "version", currentVersion)

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			mm.logger.InfoContext(ctx, "No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Get new version
	newVersion, _, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get new migration version: %w", err)
	}

	mm.logger.InfoContext(ctx, "Migrations completed successfully", "new_version", newVersion)
	return nil
}

// RollbackMigration rolls back the last migration
func (mm *MigrationManager) RollbackMigration(ctx context.Context, migrationsPath string) error {
	mm.logger.InfoContext(ctx, "Rolling back last migration", "path", migrationsPath)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(mm.db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Get current version
	currentVersion, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			mm.logger.InfoContext(ctx, "No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to get current migration version: %w", err)
	}

	if dirty {
		mm.logger.WarnContext(ctx, "Database is in dirty state", "version", currentVersion)
		return fmt.Errorf("database is in dirty state at version %d", currentVersion)
	}

	// Rollback one step
	if err := m.Steps(-1); err != nil {
		if err == migrate.ErrNoChange {
			mm.logger.InfoContext(ctx, "No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	// Get new version
	newVersion, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get new migration version: %w", err)
	}

	mm.logger.InfoContext(ctx, "Migration rollback completed successfully", "new_version", newVersion)
	return nil
}

// GetMigrationVersion returns the current migration version
func (mm *MigrationManager) GetMigrationVersion(ctx context.Context, migrationsPath string) (uint, bool, error) {
	// Create postgres driver instance
	driver, err := postgres.WithInstance(mm.db.DB, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}
