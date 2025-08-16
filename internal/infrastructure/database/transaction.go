package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// TxFunc represents a function that executes within a transaction
type TxFunc func(tx *sql.Tx) error

// WithTransaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func (db *DB) WithTransaction(ctx context.Context, fn TxFunc) error {
	return db.WithTransactionOptions(ctx, nil, fn)
}

// WithTransactionOptions executes a function within a database transaction with custom options
func (db *DB) WithTransactionOptions(ctx context.Context, opts *sql.TxOptions, fn TxFunc) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is always handled
	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				db.logger.ErrorContext(ctx, "Failed to rollback transaction after panic",
					"panic", p,
					"rollback_error", rollbackErr,
				)
			}
			panic(p) // Re-panic after rollback
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		// Rollback on error
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			db.logger.ErrorContext(ctx, "Failed to rollback transaction",
				"original_error", err,
				"rollback_error", rollbackErr,
			)
			return fmt.Errorf("transaction failed and rollback failed: original error: %w, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TxManager provides transaction management utilities
type TxManager struct {
	db     *DB
	logger *slog.Logger
}

// NewTxManager creates a new transaction manager
func NewTxManager(db *DB, logger *slog.Logger) *TxManager {
	return &TxManager{
		db:     db,
		logger: logger,
	}
}

// Execute runs a function within a transaction with logging
func (tm *TxManager) Execute(ctx context.Context, operation string, fn TxFunc) error {
	tm.logger.DebugContext(ctx, "Starting transaction", "operation", operation)

	start := time.Now()
	err := tm.db.WithTransaction(ctx, fn)
	duration := time.Since(start)

	if err != nil {
		tm.logger.ErrorContext(ctx, "Transaction failed",
			"operation", operation,
			"duration", duration,
			"error", err,
		)
		return err
	}

	tm.logger.DebugContext(ctx, "Transaction completed successfully",
		"operation", operation,
		"duration", duration,
	)

	return nil
}

// ExecuteWithOptions runs a function within a transaction with custom options and logging
func (tm *TxManager) ExecuteWithOptions(ctx context.Context, operation string, opts *sql.TxOptions, fn TxFunc) error {
	tm.logger.DebugContext(ctx, "Starting transaction with options",
		"operation", operation,
		"isolation_level", opts.Isolation,
		"read_only", opts.ReadOnly,
	)

	start := time.Now()
	err := tm.db.WithTransactionOptions(ctx, opts, fn)
	duration := time.Since(start)

	if err != nil {
		tm.logger.ErrorContext(ctx, "Transaction with options failed",
			"operation", operation,
			"duration", duration,
			"error", err,
		)
		return err
	}

	tm.logger.DebugContext(ctx, "Transaction with options completed successfully",
		"operation", operation,
		"duration", duration,
	)

	return nil
}
