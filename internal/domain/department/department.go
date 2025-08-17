package department

import (
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

// Department represents the department domain entity
type Department struct {
	id          DepartmentID
	name        Name
	description Description
	managerID   *EmployeeID
	createdAt   time.Time
	updatedAt   time.Time
}

// NewDepartment creates a new Department entity with validation
func NewDepartment(name, description string, managerID *string) (*Department, error) {
	nameVO, err := NewName(name)
	if err != nil {
		return nil, err
	}

	descriptionVO, err := NewDescription(description)
	if err != nil {
		return nil, err
	}

	var managerIDVO *EmployeeID
	if managerID != nil {
		empID, err := NewEmployeeID(*managerID)
		if err != nil {
			return nil, err
		}
		managerIDVO = &empID
	}

	now := time.Now()

	return &Department{
		id:          GenerateDepartmentID(),
		name:        nameVO,
		description: descriptionVO,
		managerID:   managerIDVO,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// NewDepartmentWithID creates a Department entity with a specific ID (for reconstruction from persistence)
func NewDepartmentWithID(id, name, description string, managerID *string, createdAt, updatedAt time.Time) (*Department, error) {
	deptID, err := NewDepartmentID(id)
	if err != nil {
		return nil, err
	}

	nameVO, err := NewName(name)
	if err != nil {
		return nil, err
	}

	descriptionVO, err := NewDescription(description)
	if err != nil {
		return nil, err
	}

	var managerIDVO *EmployeeID
	if managerID != nil {
		empID, err := NewEmployeeID(*managerID)
		if err != nil {
			return nil, err
		}
		managerIDVO = &empID
	}

	return &Department{
		id:          deptID,
		name:        nameVO,
		description: descriptionVO,
		managerID:   managerIDVO,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}, nil
}

// ID returns the department's ID
func (d *Department) ID() DepartmentID {
	return d.id
}

// Name returns the department's name
func (d *Department) Name() Name {
	return d.name
}

// Description returns the department's description
func (d *Department) Description() Description {
	return d.description
}

// ManagerID returns the department's manager ID
func (d *Department) ManagerID() *EmployeeID {
	return d.managerID
}

// CreatedAt returns when the department was created
func (d *Department) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt returns when the department was last updated
func (d *Department) UpdatedAt() time.Time {
	return d.updatedAt
}

// UpdateName updates the department's name with validation
func (d *Department) UpdateName(name string) error {
	nameVO, err := NewName(name)
	if err != nil {
		return err
	}

	d.name = nameVO
	d.updatedAt = time.Now()
	return nil
}

// UpdateDescription updates the department's description with validation
func (d *Department) UpdateDescription(description string) error {
	descriptionVO, err := NewDescription(description)
	if err != nil {
		return err
	}

	d.description = descriptionVO
	d.updatedAt = time.Now()
	return nil
}

// UpdateManager updates the department's manager with validation
func (d *Department) UpdateManager(managerID *string) error {
	var managerIDVO *EmployeeID
	if managerID != nil {
		empID, err := NewEmployeeID(*managerID)
		if err != nil {
			return err
		}
		managerIDVO = &empID
	}

	d.managerID = managerIDVO
	d.updatedAt = time.Now()
	return nil
}

// Validate performs comprehensive validation of the department entity
func (d *Department) Validate() error {
	if d.id.String() == "" {
		return errors.ErrInvalidDepartmentID
	}

	if d.name.String() == "" {
		return errors.ErrInvalidDepartmentName
	}

	if d.description.String() == "" {
		return errors.ErrInvalidDepartmentDescription
	}

	if d.createdAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_CREATED_AT",
			Message: "Created at timestamp cannot be zero",
			Field:   "createdAt",
		}
	}

	if d.updatedAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_UPDATED_AT",
			Message: "Updated at timestamp cannot be zero",
			Field:   "updatedAt",
		}
	}

	if d.updatedAt.Before(d.createdAt) {
		return errors.DomainError{
			Code:    "INVALID_TIMESTAMPS",
			Message: "Updated at cannot be before created at",
			Field:   "updatedAt",
		}
	}

	return nil
}

// Equals checks if two departments are equal based on their ID
func (d *Department) Equals(other *Department) bool {
	if other == nil {
		return false
	}
	return d.id.Equals(other.id)
}