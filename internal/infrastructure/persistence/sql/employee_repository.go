package sql

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/employee"
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

// employeeRepository implements the employee.Repository interface
type employeeRepository struct {
	db *sql.DB
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(db *sql.DB) employee.Repository {
	return &employeeRepository{
		db: db,
	}
}

// FindByID retrieves an employee by their ID
func (r *employeeRepository) FindByID(ctx context.Context, id employee.EmployeeID) (*employee.Employee, error) {
	query := `
		SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
		FROM employees
		WHERE id = $1
	`

	var (
		empID, userIDStr, empCode, dept, pos, status string
		hireDate                                      time.Time
		salary                                        float64
		createdAt, updatedAt                          time.Time
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&empID, &userIDStr, &empCode, &dept, &pos, &hireDate, &salary, &status, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to query employee: %w", err)
	}

	// Create user ID
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in database: %w", err)
	}

	// Create employee
	emp, err := employee.NewEmployeeWithID(empID, userID, empCode, dept, pos, hireDate, salary, status, createdAt, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee from database data: %w", err)
	}

	return emp, nil
}

// FindByUserID retrieves an employee by their user ID
func (r *employeeRepository) FindByUserID(ctx context.Context, userID user.UserID) (*employee.Employee, error) {
	query := `
		SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
		FROM employees
		WHERE user_id = $1
	`

	var (
		empID, userIDStr, empCode, dept, pos, status string
		hireDate                                      time.Time
		salary                                        float64
		createdAt, updatedAt                          time.Time
	)

	err := r.db.QueryRowContext(ctx, query, userID.String()).Scan(
		&empID, &userIDStr, &empCode, &dept, &pos, &hireDate, &salary, &status, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to query employee by user ID: %w", err)
	}

	// Create employee
	emp, err := employee.NewEmployeeWithID(empID, userID, empCode, dept, pos, hireDate, salary, status, createdAt, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee from database data: %w", err)
	}

	return emp, nil
}

// FindByEmployeeCode retrieves an employee by their employee code
func (r *employeeRepository) FindByEmployeeCode(ctx context.Context, employeeCode employee.EmployeeCode) (*employee.Employee, error) {
	query := `
		SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
		FROM employees
		WHERE employee_code = $1
	`

	var (
		empID, userIDStr, empCode, dept, pos, status string
		hireDate                                      time.Time
		salary                                        float64
		createdAt, updatedAt                          time.Time
	)

	err := r.db.QueryRowContext(ctx, query, employeeCode.String()).Scan(
		&empID, &userIDStr, &empCode, &dept, &pos, &hireDate, &salary, &status, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to query employee by employee code: %w", err)
	}

	// Create user ID
	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in database: %w", err)
	}

	// Create employee
	emp, err := employee.NewEmployeeWithID(empID, userID, empCode, dept, pos, hireDate, salary, status, createdAt, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee from database data: %w", err)
	}

	return emp, nil
}

// FindAll retrieves employees with pagination support
func (r *employeeRepository) FindAll(ctx context.Context, limit int, cursor string) ([]*employee.Employee, string, error) {
	var query string
	var args []interface{}

	if cursor == "" {
		query = `
			SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
			FROM employees
			ORDER BY created_at ASC, id ASC
			LIMIT $1
		`
		args = append(args, limit)
	} else {
		// Decode cursor
		cursorData, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor format: %w", err)
		}

		// Parse cursor (format: "created_at:employee_id")
		var createdAtStr, empID string
		_, err = fmt.Sscanf(string(cursorData), "%s:%s", &createdAtStr, &empID)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor content: %w", err)
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor timestamp: %w", err)
		}

		query = `
			SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
			FROM employees
			WHERE (created_at, id) > ($1, $2)
			ORDER BY created_at ASC, id ASC
			LIMIT $3
		`
		args = append(args, createdAt, empID, limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query employees: %w", err)
	}
	defer rows.Close()

	var employees []*employee.Employee
	var lastCreatedAt time.Time
	var lastID string

	for rows.Next() {
		var (
			empID, userIDStr, empCode, dept, pos, status string
			hireDate                                      time.Time
			salary                                        float64
			createdAt, updatedAt                          time.Time
		)

		err := rows.Scan(&empID, &userIDStr, &empCode, &dept, &pos, &hireDate, &salary, &status, &createdAt, &updatedAt)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan employee row: %w", err)
		}

		// Create user ID
		userID, err := user.NewUserID(userIDStr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid user ID in database: %w", err)
		}

		// Create employee
		emp, err := employee.NewEmployeeWithID(empID, userID, empCode, dept, pos, hireDate, salary, status, createdAt, updatedAt)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create employee from database data: %w", err)
		}

		employees = append(employees, emp)
		lastCreatedAt = createdAt
		lastID = empID
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating employee rows: %w", err)
	}

	// Generate next cursor if we have results and reached the limit
	var nextCursor string
	if len(employees) == limit && len(employees) > 0 {
		cursorData := fmt.Sprintf("%s:%s", lastCreatedAt.Format(time.RFC3339), lastID)
		nextCursor = base64.StdEncoding.EncodeToString([]byte(cursorData))
	}

	return employees, nextCursor, nil
}

// FindByDepartment retrieves employees by department
func (r *employeeRepository) FindByDepartment(ctx context.Context, department employee.Department, limit int, cursor string) ([]*employee.Employee, string, error) {
	var query string
	var args []interface{}

	if cursor == "" {
		query = `
			SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
			FROM employees
			WHERE department = $1
			ORDER BY created_at ASC, id ASC
			LIMIT $2
		`
		args = append(args, department.String(), limit)
	} else {
		// Decode cursor
		cursorData, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor format: %w", err)
		}

		// Parse cursor (format: "created_at:employee_id")
		var createdAtStr, empID string
		_, err = fmt.Sscanf(string(cursorData), "%s:%s", &createdAtStr, &empID)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor content: %w", err)
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor timestamp: %w", err)
		}

		query = `
			SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
			FROM employees
			WHERE department = $1 AND (created_at, id) > ($2, $3)
			ORDER BY created_at ASC, id ASC
			LIMIT $4
		`
		args = append(args, department.String(), createdAt, empID, limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query employees by department: %w", err)
	}
	defer rows.Close()

	var employees []*employee.Employee
	var lastCreatedAt time.Time
	var lastID string

	for rows.Next() {
		var (
			empID, userIDStr, empCode, dept, pos, status string
			hireDate                                      time.Time
			salary                                        float64
			createdAt, updatedAt                          time.Time
		)

		err := rows.Scan(&empID, &userIDStr, &empCode, &dept, &pos, &hireDate, &salary, &status, &createdAt, &updatedAt)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan employee row: %w", err)
		}

		// Create user ID
		userID, err := user.NewUserID(userIDStr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid user ID in database: %w", err)
		}

		// Create employee
		emp, err := employee.NewEmployeeWithID(empID, userID, empCode, dept, pos, hireDate, salary, status, createdAt, updatedAt)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create employee from database data: %w", err)
		}

		employees = append(employees, emp)
		lastCreatedAt = createdAt
		lastID = empID
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating employee rows: %w", err)
	}

	// Generate next cursor if we have results and reached the limit
	var nextCursor string
	if len(employees) == limit && len(employees) > 0 {
		cursorData := fmt.Sprintf("%s:%s", lastCreatedAt.Format(time.RFC3339), lastID)
		nextCursor = base64.StdEncoding.EncodeToString([]byte(cursorData))
	}

	return employees, nextCursor, nil
}

// FindByStatus retrieves employees by status
func (r *employeeRepository) FindByStatus(ctx context.Context, status employee.Status, limit int, cursor string) ([]*employee.Employee, string, error) {
	var query string
	var args []interface{}

	if cursor == "" {
		query = `
			SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
			FROM employees
			WHERE status = $1
			ORDER BY created_at ASC, id ASC
			LIMIT $2
		`
		args = append(args, status.String(), limit)
	} else {
		// Decode cursor
		cursorData, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor format: %w", err)
		}

		// Parse cursor (format: "created_at:employee_id")
		var createdAtStr, empID string
		_, err = fmt.Sscanf(string(cursorData), "%s:%s", &createdAtStr, &empID)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor content: %w", err)
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid cursor timestamp: %w", err)
		}

		query = `
			SELECT id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at
			FROM employees
			WHERE status = $1 AND (created_at, id) > ($2, $3)
			ORDER BY created_at ASC, id ASC
			LIMIT $4
		`
		args = append(args, status.String(), createdAt, empID, limit)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query employees by status: %w", err)
	}
	defer rows.Close()

	var employees []*employee.Employee
	var lastCreatedAt time.Time
	var lastID string

	for rows.Next() {
		var (
			empID, userIDStr, empCode, dept, pos, status string
			hireDate                                      time.Time
			salary                                        float64
			createdAt, updatedAt                          time.Time
		)

		err := rows.Scan(&empID, &userIDStr, &empCode, &dept, &pos, &hireDate, &salary, &status, &createdAt, &updatedAt)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan employee row: %w", err)
		}

		// Create user ID
		userID, err := user.NewUserID(userIDStr)
		if err != nil {
			return nil, "", fmt.Errorf("invalid user ID in database: %w", err)
		}

		// Create employee
		emp, err := employee.NewEmployeeWithID(empID, userID, empCode, dept, pos, hireDate, salary, status, createdAt, updatedAt)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create employee from database data: %w", err)
		}

		employees = append(employees, emp)
		lastCreatedAt = createdAt
		lastID = empID
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating employee rows: %w", err)
	}

	// Generate next cursor if we have results and reached the limit
	var nextCursor string
	if len(employees) == limit && len(employees) > 0 {
		cursorData := fmt.Sprintf("%s:%s", lastCreatedAt.Format(time.RFC3339), lastID)
		nextCursor = base64.StdEncoding.EncodeToString([]byte(cursorData))
	}

	return employees, nextCursor, nil
}

// Create persists a new employee
func (r *employeeRepository) Create(ctx context.Context, emp *employee.Employee) error {
	query := `
		INSERT INTO employees (id, user_id, employee_code, department, position, hire_date, salary, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		emp.ID().String(),
		emp.UserID().String(),
		emp.EmployeeCode().String(),
		emp.Department().String(),
		emp.Position().String(),
		emp.HireDate(),
		emp.Salary().Value(),
		emp.Status().String(),
		emp.CreatedAt(),
		emp.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert employee: %w", err)
	}

	return nil
}

// Update modifies an existing employee
func (r *employeeRepository) Update(ctx context.Context, emp *employee.Employee) error {
	query := `
		UPDATE employees
		SET user_id = $2, employee_code = $3, department = $4, position = $5, hire_date = $6, salary = $7, status = $8, updated_at = $9
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		emp.ID().String(),
		emp.UserID().String(),
		emp.EmployeeCode().String(),
		emp.Department().String(),
		emp.Position().String(),
		emp.HireDate(),
		emp.Salary().Value(),
		emp.Status().String(),
		emp.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrEmployeeNotFound
	}

	return nil
}

// Delete removes an employee by their ID
func (r *employeeRepository) Delete(ctx context.Context, id employee.EmployeeID) error {
	query := `DELETE FROM employees WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrEmployeeNotFound
	}

	return nil
}

// ExistsByEmployeeCode checks if an employee with the given employee code exists
func (r *employeeRepository) ExistsByEmployeeCode(ctx context.Context, employeeCode employee.EmployeeCode) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM employees WHERE employee_code = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, employeeCode.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if employee code exists: %w", err)
	}

	return exists, nil
}

// ExistsByUserID checks if an employee with the given user ID exists
func (r *employeeRepository) ExistsByUserID(ctx context.Context, userID user.UserID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM employees WHERE user_id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if employee exists for user: %w", err)
	}

	return exists, nil
}

// Count returns the total number of employees
func (r *employeeRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM employees`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count employees: %w", err)
	}

	return count, nil
}

// CountByDepartment returns the number of employees in a department
func (r *employeeRepository) CountByDepartment(ctx context.Context, department employee.Department) (int64, error) {
	query := `SELECT COUNT(*) FROM employees WHERE department = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, department.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count employees by department: %w", err)
	}

	return count, nil
}

// CountByStatus returns the number of employees with a specific status
func (r *employeeRepository) CountByStatus(ctx context.Context, status employee.Status) (int64, error) {
	query := `SELECT COUNT(*) FROM employees WHERE status = $1`

	var count int64
	err := r.db.QueryRowContext(ctx, query, status.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count employees by status: %w", err)
	}

	return count, nil
}