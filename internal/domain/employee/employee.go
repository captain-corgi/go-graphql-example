package employee

import (
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

// Employee represents the employee domain entity
type Employee struct {
	id           EmployeeID
	userID       user.UserID
	employeeCode EmployeeCode
	department   Department
	position     Position
	hireDate     time.Time
	salary       Salary
	status       Status
	createdAt    time.Time
	updatedAt    time.Time
}

// NewEmployee creates a new Employee entity with validation
func NewEmployee(userID user.UserID, employeeCode, department, position string, hireDate time.Time, salary float64, status string) (*Employee, error) {
	empCode, err := NewEmployeeCode(employeeCode)
	if err != nil {
		return nil, err
	}

	dept, err := NewDepartment(department)
	if err != nil {
		return nil, err
	}

	pos, err := NewPosition(position)
	if err != nil {
		return nil, err
	}

	sal, err := NewSalary(salary)
	if err != nil {
		return nil, err
	}

	empStatus, err := NewStatus(status)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Employee{
		id:           GenerateEmployeeID(),
		userID:       userID,
		employeeCode: empCode,
		department:   dept,
		position:     pos,
		hireDate:     hireDate,
		salary:       sal,
		status:       empStatus,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// NewEmployeeWithID creates an Employee entity with a specific ID (for reconstruction from persistence)
func NewEmployeeWithID(id string, userID user.UserID, employeeCode, department, position string, hireDate time.Time, salary float64, status string, createdAt, updatedAt time.Time) (*Employee, error) {
	empID, err := NewEmployeeID(id)
	if err != nil {
		return nil, err
	}

	empCode, err := NewEmployeeCode(employeeCode)
	if err != nil {
		return nil, err
	}

	dept, err := NewDepartment(department)
	if err != nil {
		return nil, err
	}

	pos, err := NewPosition(position)
	if err != nil {
		return nil, err
	}

	sal, err := NewSalary(salary)
	if err != nil {
		return nil, err
	}

	empStatus, err := NewStatus(status)
	if err != nil {
		return nil, err
	}

	return &Employee{
		id:           empID,
		userID:       userID,
		employeeCode: empCode,
		department:   dept,
		position:     pos,
		hireDate:     hireDate,
		salary:       sal,
		status:       empStatus,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}

// ID returns the employee's ID
func (e *Employee) ID() EmployeeID {
	return e.id
}

// UserID returns the employee's user ID
func (e *Employee) UserID() user.UserID {
	return e.userID
}

// EmployeeCode returns the employee's code
func (e *Employee) EmployeeCode() EmployeeCode {
	return e.employeeCode
}

// Department returns the employee's department
func (e *Employee) Department() Department {
	return e.department
}

// Position returns the employee's position
func (e *Employee) Position() Position {
	return e.position
}

// HireDate returns the employee's hire date
func (e *Employee) HireDate() time.Time {
	return e.hireDate
}

// Salary returns the employee's salary
func (e *Employee) Salary() Salary {
	return e.salary
}

// Status returns the employee's status
func (e *Employee) Status() Status {
	return e.status
}

// CreatedAt returns when the employee was created
func (e *Employee) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt returns when the employee was last updated
func (e *Employee) UpdatedAt() time.Time {
	return e.updatedAt
}

// UpdateEmployeeCode updates the employee's code with validation
func (e *Employee) UpdateEmployeeCode(employeeCode string) error {
	empCode, err := NewEmployeeCode(employeeCode)
	if err != nil {
		return err
	}

	e.employeeCode = empCode
	e.updatedAt = time.Now()
	return nil
}

// UpdateDepartment updates the employee's department with validation
func (e *Employee) UpdateDepartment(department string) error {
	dept, err := NewDepartment(department)
	if err != nil {
		return err
	}

	e.department = dept
	e.updatedAt = time.Now()
	return nil
}

// UpdatePosition updates the employee's position with validation
func (e *Employee) UpdatePosition(position string) error {
	pos, err := NewPosition(position)
	if err != nil {
		return err
	}

	e.position = pos
	e.updatedAt = time.Now()
	return nil
}

// UpdateHireDate updates the employee's hire date
func (e *Employee) UpdateHireDate(hireDate time.Time) error {
	if hireDate.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_HIRE_DATE",
			Message: "Hire date cannot be zero",
			Field:   "hireDate",
		}
	}

	e.hireDate = hireDate
	e.updatedAt = time.Now()
	return nil
}

// UpdateSalary updates the employee's salary with validation
func (e *Employee) UpdateSalary(salary float64) error {
	sal, err := NewSalary(salary)
	if err != nil {
		return err
	}

	e.salary = sal
	e.updatedAt = time.Now()
	return nil
}

// UpdateStatus updates the employee's status with validation
func (e *Employee) UpdateStatus(status string) error {
	empStatus, err := NewStatus(status)
	if err != nil {
		return err
	}

	e.status = empStatus
	e.updatedAt = time.Now()
	return nil
}

// Validate performs comprehensive validation of the employee entity
func (e *Employee) Validate() error {
	if e.id.String() == "" {
		return errors.ErrInvalidEmployeeID
	}

	if e.userID.String() == "" {
		return errors.DomainError{
			Code:    "INVALID_USER_ID",
			Message: "User ID cannot be empty",
			Field:   "userID",
		}
	}

	if e.employeeCode.String() == "" {
		return errors.ErrInvalidEmployeeCode
	}

	if e.department.String() == "" {
		return errors.ErrInvalidDepartment
	}

	if e.position.String() == "" {
		return errors.ErrInvalidPosition
	}

	if e.hireDate.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_HIRE_DATE",
			Message: "Hire date cannot be zero",
			Field:   "hireDate",
		}
	}

	if e.salary.Value() <= 0 {
		return errors.ErrInvalidSalary
	}

	if e.status.String() == "" {
		return errors.ErrInvalidStatus
	}

	if e.createdAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_CREATED_AT",
			Message: "Created at timestamp cannot be zero",
			Field:   "createdAt",
		}
	}

	if e.updatedAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_UPDATED_AT",
			Message: "Updated at timestamp cannot be zero",
			Field:   "updatedAt",
		}
	}

	if e.updatedAt.Before(e.createdAt) {
		return errors.DomainError{
			Code:    "INVALID_TIMESTAMPS",
			Message: "Updated at cannot be before created at",
			Field:   "updatedAt",
		}
	}

	return nil
}

// Equals checks if two employees are equal based on their ID
func (e *Employee) Equals(other *Employee) bool {
	if other == nil {
		return false
	}
	return e.id.Equals(other.id)
}