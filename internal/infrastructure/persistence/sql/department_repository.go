package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/department"
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/infrastructure/database"
	"github.com/lib/pq"
)

// departmentRepository implements the department.Repository interface using SQL
type departmentRepository struct {
	db     *database.DB
	logger *slog.Logger
}

// NewDepartmentRepository creates a new SQL-based department repository
func NewDepartmentRepository(db *database.DB, logger *slog.Logger) department.Repository {
	return &departmentRepository{
		db:     db,
		logger: logger,
	}
}

// FindByID retrieves a department by its ID
func (r *departmentRepository) FindByID(ctx context.Context, id department.DepartmentID) (*department.Department, error) {
	r.logger.DebugContext(ctx, "Finding department by ID", "department_id", id.String())

	query := `
		SELECT id, name, description, manager_id, created_at, updated_at 
		FROM departments 
		WHERE id = $1`

	var deptID, name, description string
	var managerID sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&deptID, &name, &description, &managerID, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.DebugContext(ctx, "Department not found", "department_id", id.String())
			return nil, errors.ErrDepartmentNotFound
		}
		r.logger.ErrorContext(ctx, "Failed to find department by ID", "error", err, "department_id", id.String())
		return nil, fmt.Errorf("failed to find department by ID: %w", err)
	}

	var managerIDPtr *string
	if managerID.Valid {
		managerIDPtr = &managerID.String
	}

	domainDept, err := department.NewDepartmentWithID(deptID, name, description, managerIDPtr, createdAt, updatedAt)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create domain department from database record", "error", err)
		return nil, fmt.Errorf("failed to create domain department: %w", err)
	}

	r.logger.DebugContext(ctx, "Successfully found department by ID", "department_id", id.String())
	return domainDept, nil
}

// FindByName retrieves a department by its name
func (r *departmentRepository) FindByName(ctx context.Context, name department.Name) (*department.Department, error) {
	r.logger.DebugContext(ctx, "Finding department by name", "name", name.String())

	query := `
		SELECT id, name, description, manager_id, created_at, updated_at 
		FROM departments 
		WHERE name = $1`

	var deptID, deptName, description string
	var managerID sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, name.String()).Scan(
		&deptID, &deptName, &description, &managerID, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.DebugContext(ctx, "Department not found by name", "name", name.String())
			return nil, errors.ErrDepartmentNotFound
		}
		r.logger.ErrorContext(ctx, "Failed to find department by name", "error", err, "name", name.String())
		return nil, fmt.Errorf("failed to find department by name: %w", err)
	}

	var managerIDPtr *string
	if managerID.Valid {
		managerIDPtr = &managerID.String
	}

	domainDept, err := department.NewDepartmentWithID(deptID, deptName, description, managerIDPtr, createdAt, updatedAt)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to create domain department from database record", "error", err)
		return nil, fmt.Errorf("failed to create domain department: %w", err)
	}

	r.logger.DebugContext(ctx, "Successfully found department by name", "name", name.String())
	return domainDept, nil
}

// FindAll retrieves departments with pagination support
func (r *departmentRepository) FindAll(ctx context.Context, limit int, cursor string) ([]*department.Department, string, error) {
	r.logger.DebugContext(ctx, "Finding all departments", "limit", limit, "cursor", cursor)

	var query string
	var args []interface{}

	if cursor != "" {
		query = `
			SELECT id, name, description, manager_id, created_at, updated_at 
			FROM departments 
			WHERE id > $1
			ORDER BY id 
			LIMIT $2`
		args = []interface{}{cursor, limit}
	} else {
		query = `
			SELECT id, name, description, manager_id, created_at, updated_at 
			FROM departments 
			ORDER BY id 
			LIMIT $1`
		args = []interface{}{limit}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query departments", "error", err)
		return nil, "", fmt.Errorf("failed to query departments: %w", err)
	}
	defer rows.Close()

	var departments []*department.Department
	var lastID string

	for rows.Next() {
		var deptID, name, description string
		var managerID sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(&deptID, &name, &description, &managerID, &createdAt, &updatedAt)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan department row", "error", err)
			return nil, "", fmt.Errorf("failed to scan department row: %w", err)
		}

		var managerIDPtr *string
		if managerID.Valid {
			managerIDPtr = &managerID.String
		}

		domainDept, err := department.NewDepartmentWithID(deptID, name, description, managerIDPtr, createdAt, updatedAt)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to create domain department from database record", "error", err)
			return nil, "", fmt.Errorf("failed to create domain department: %w", err)
		}

		departments = append(departments, domainDept)
		lastID = deptID
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating department rows", "error", err)
		return nil, "", fmt.Errorf("error iterating department rows: %w", err)
	}

	var nextCursor string
	if len(departments) == limit {
		nextCursor = lastID
	}

	r.logger.DebugContext(ctx, "Successfully found departments", "count", len(departments), "nextCursor", nextCursor)
	return departments, nextCursor, nil
}

// FindByManager retrieves departments by manager ID
func (r *departmentRepository) FindByManager(ctx context.Context, managerID department.EmployeeID, limit int, cursor string) ([]*department.Department, string, error) {
	r.logger.DebugContext(ctx, "Finding departments by manager", "manager_id", managerID.String(), "limit", limit, "cursor", cursor)

	var query string
	var args []interface{}

	if cursor != "" {
		query = `
			SELECT id, name, description, manager_id, created_at, updated_at 
			FROM departments 
			WHERE manager_id = $1 AND id > $2
			ORDER BY id 
			LIMIT $3`
		args = []interface{}{managerID.String(), cursor, limit}
	} else {
		query = `
			SELECT id, name, description, manager_id, created_at, updated_at 
			FROM departments 
			WHERE manager_id = $1
			ORDER BY id 
			LIMIT $2`
		args = []interface{}{managerID.String(), limit}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to query departments by manager", "error", err)
		return nil, "", fmt.Errorf("failed to query departments by manager: %w", err)
	}
	defer rows.Close()

	var departments []*department.Department
	var lastID string

	for rows.Next() {
		var deptID, name, description string
		var managerID sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(&deptID, &name, &description, &managerID, &createdAt, &updatedAt)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan department row", "error", err)
			return nil, "", fmt.Errorf("failed to scan department row: %w", err)
		}

		var managerIDPtr *string
		if managerID.Valid {
			managerIDPtr = &managerID.String
		}

		domainDept, err := department.NewDepartmentWithID(deptID, name, description, managerIDPtr, createdAt, updatedAt)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to create domain department from database record", "error", err)
			return nil, "", fmt.Errorf("failed to create domain department: %w", err)
		}

		departments = append(departments, domainDept)
		lastID = deptID
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "Error iterating department rows", "error", err)
		return nil, "", fmt.Errorf("error iterating department rows: %w", err)
	}

	var nextCursor string
	if len(departments) == limit {
		nextCursor = lastID
	}

	r.logger.DebugContext(ctx, "Successfully found departments by manager", "count", len(departments), "nextCursor", nextCursor)
	return departments, nextCursor, nil
}

// Save persists a department (create or update)
func (r *departmentRepository) Save(ctx context.Context, dept *department.Department) error {
	r.logger.DebugContext(ctx, "Saving department", "department_id", dept.ID().String())

	// Check if department exists
	exists, err := r.Exists(ctx, dept.ID())
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to check if department exists", "error", err)
		return fmt.Errorf("failed to check if department exists: %w", err)
	}

	var managerID sql.NullString
	if dept.ManagerID() != nil {
		managerID.String = dept.ManagerID().String()
		managerID.Valid = true
	}

	if exists {
		// Update existing department
		query := `
			UPDATE departments 
			SET name = $1, description = $2, manager_id = $3, updated_at = $4
			WHERE id = $5`

		result, err := r.db.ExecContext(ctx, query,
			dept.Name().String(),
			dept.Description().String(),
			managerID,
			dept.UpdatedAt(),
			dept.ID().String(),
		)

		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to update department", "error", err)
			return fmt.Errorf("failed to update department: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to get rows affected", "error", err)
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			r.logger.WarnContext(ctx, "No rows affected when updating department", "department_id", dept.ID().String())
			return errors.ErrDepartmentNotFound
		}
	} else {
		// Insert new department
		query := `
			INSERT INTO departments (id, name, description, manager_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err := r.db.ExecContext(ctx, query,
			dept.ID().String(),
			dept.Name().String(),
			dept.Description().String(),
			managerID,
			dept.CreatedAt(),
			dept.UpdatedAt(),
		)

		if err != nil {
			// Check for unique constraint violation
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				r.logger.WarnContext(ctx, "Department name already exists", "error", err)
				return errors.ErrDuplicateDepartmentName
			}
			r.logger.ErrorContext(ctx, "Failed to insert department", "error", err)
			return fmt.Errorf("failed to insert department: %w", err)
		}
	}

	r.logger.DebugContext(ctx, "Successfully saved department", "department_id", dept.ID().String())
	return nil
}

// Delete removes a department by ID
func (r *departmentRepository) Delete(ctx context.Context, id department.DepartmentID) error {
	r.logger.DebugContext(ctx, "Deleting department", "department_id", id.String())

	query := `DELETE FROM departments WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to delete department", "error", err)
		return fmt.Errorf("failed to delete department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to get rows affected", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		r.logger.WarnContext(ctx, "No rows affected when deleting department", "department_id", id.String())
		return errors.ErrDepartmentNotFound
	}

	r.logger.DebugContext(ctx, "Successfully deleted department", "department_id", id.String())
	return nil
}

// Exists checks if a department exists by ID
func (r *departmentRepository) Exists(ctx context.Context, id department.DepartmentID) (bool, error) {
	r.logger.DebugContext(ctx, "Checking if department exists", "department_id", id.String())

	query := `SELECT EXISTS(SELECT 1 FROM departments WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&exists)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to check if department exists", "error", err)
		return false, fmt.Errorf("failed to check if department exists: %w", err)
	}

	r.logger.DebugContext(ctx, "Department exists check result", "department_id", id.String(), "exists", exists)
	return exists, nil
}

// ExistsByName checks if a department exists by name
func (r *departmentRepository) ExistsByName(ctx context.Context, name department.Name, excludeID *department.DepartmentID) (bool, error) {
	r.logger.DebugContext(ctx, "Checking if department exists by name", "name", name.String(), "exclude_id", excludeID)

	var query string
	var args []interface{}

	if excludeID != nil {
		query = `SELECT EXISTS(SELECT 1 FROM departments WHERE name = $1 AND id != $2)`
		args = []interface{}{name.String(), excludeID.String()}
	} else {
		query = `SELECT EXISTS(SELECT 1 FROM departments WHERE name = $1)`
		args = []interface{}{name.String()}
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		r.logger.ErrorContext(ctx, "Failed to check if department exists by name", "error", err)
		return false, fmt.Errorf("failed to check if department exists by name: %w", err)
	}

	r.logger.DebugContext(ctx, "Department exists by name check result", "name", name.String(), "exists", exists)
	return exists, nil
}