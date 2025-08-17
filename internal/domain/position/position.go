package position

import (
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

// Position represents the position/job domain entity
type Position struct {
	id          PositionID
	title       Title
	description Description
	departmentID *DepartmentID
	requirements Requirements
	salaryRange  SalaryRange
	createdAt    time.Time
	updatedAt    time.Time
}

// NewPosition creates a new Position entity with validation
func NewPosition(title, description, requirements string, departmentID *string, minSalary, maxSalary float64) (*Position, error) {
	titleVO, err := NewTitle(title)
	if err != nil {
		return nil, err
	}

	descriptionVO, err := NewDescription(description)
	if err != nil {
		return nil, err
	}

	requirementsVO, err := NewRequirements(requirements)
	if err != nil {
		return nil, err
	}

	var deptIDVO *DepartmentID
	if departmentID != nil {
		deptID, err := NewDepartmentID(*departmentID)
		if err != nil {
			return nil, err
		}
		deptIDVO = &deptID
	}

	salaryRangeVO, err := NewSalaryRange(minSalary, maxSalary)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Position{
		id:           GeneratePositionID(),
		title:        titleVO,
		description:  descriptionVO,
		departmentID: deptIDVO,
		requirements: requirementsVO,
		salaryRange:  salaryRangeVO,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// NewPositionWithID creates a Position entity with a specific ID (for reconstruction from persistence)
func NewPositionWithID(id, title, description, requirements string, departmentID *string, minSalary, maxSalary float64, createdAt, updatedAt time.Time) (*Position, error) {
	positionID, err := NewPositionID(id)
	if err != nil {
		return nil, err
	}

	titleVO, err := NewTitle(title)
	if err != nil {
		return nil, err
	}

	descriptionVO, err := NewDescription(description)
	if err != nil {
		return nil, err
	}

	requirementsVO, err := NewRequirements(requirements)
	if err != nil {
		return nil, err
	}

	var deptIDVO *DepartmentID
	if departmentID != nil {
		deptID, err := NewDepartmentID(*departmentID)
		if err != nil {
			return nil, err
		}
		deptIDVO = &deptID
	}

	salaryRangeVO, err := NewSalaryRange(minSalary, maxSalary)
	if err != nil {
		return nil, err
	}

	return &Position{
		id:           positionID,
		title:        titleVO,
		description:  descriptionVO,
		departmentID: deptIDVO,
		requirements: requirementsVO,
		salaryRange:  salaryRangeVO,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}

// ID returns the position's ID
func (p *Position) ID() PositionID {
	return p.id
}

// Title returns the position's title
func (p *Position) Title() Title {
	return p.title
}

// Description returns the position's description
func (p *Position) Description() Description {
	return p.description
}

// DepartmentID returns the position's department ID
func (p *Position) DepartmentID() *DepartmentID {
	return p.departmentID
}

// Requirements returns the position's requirements
func (p *Position) Requirements() Requirements {
	return p.requirements
}

// SalaryRange returns the position's salary range
func (p *Position) SalaryRange() SalaryRange {
	return p.salaryRange
}

// CreatedAt returns when the position was created
func (p *Position) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns when the position was last updated
func (p *Position) UpdatedAt() time.Time {
	return p.updatedAt
}

// UpdateTitle updates the position's title with validation
func (p *Position) UpdateTitle(title string) error {
	titleVO, err := NewTitle(title)
	if err != nil {
		return err
	}

	p.title = titleVO
	p.updatedAt = time.Now()
	return nil
}

// UpdateDescription updates the position's description with validation
func (p *Position) UpdateDescription(description string) error {
	descriptionVO, err := NewDescription(description)
	if err != nil {
		return err
	}

	p.description = descriptionVO
	p.updatedAt = time.Now()
	return nil
}

// UpdateRequirements updates the position's requirements with validation
func (p *Position) UpdateRequirements(requirements string) error {
	requirementsVO, err := NewRequirements(requirements)
	if err != nil {
		return err
	}

	p.requirements = requirementsVO
	p.updatedAt = time.Now()
	return nil
}

// UpdateDepartment updates the position's department with validation
func (p *Position) UpdateDepartment(departmentID *string) error {
	var deptIDVO *DepartmentID
	if departmentID != nil {
		deptID, err := NewDepartmentID(*departmentID)
		if err != nil {
			return err
		}
		deptIDVO = &deptID
	}

	p.departmentID = deptIDVO
	p.updatedAt = time.Now()
	return nil
}

// UpdateSalaryRange updates the position's salary range with validation
func (p *Position) UpdateSalaryRange(minSalary, maxSalary float64) error {
	salaryRangeVO, err := NewSalaryRange(minSalary, maxSalary)
	if err != nil {
		return err
	}

	p.salaryRange = salaryRangeVO
	p.updatedAt = time.Now()
	return nil
}

// Validate performs comprehensive validation of the position entity
func (p *Position) Validate() error {
	if p.id.String() == "" {
		return errors.ErrInvalidPositionID
	}

	if p.title.String() == "" {
		return errors.ErrInvalidPositionTitle
	}

	if p.description.String() == "" {
		return errors.ErrInvalidPositionDescription
	}

	if p.requirements.String() == "" {
		return errors.ErrInvalidPositionRequirements
	}

	if p.createdAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_CREATED_AT",
			Message: "Created at timestamp cannot be zero",
			Field:   "createdAt",
		}
	}

	if p.updatedAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_UPDATED_AT",
			Message: "Updated at timestamp cannot be zero",
			Field:   "updatedAt",
		}
	}

	if p.updatedAt.Before(p.createdAt) {
		return errors.DomainError{
			Code:    "INVALID_TIMESTAMPS",
			Message: "Updated at cannot be before created at",
			Field:   "updatedAt",
		}
	}

	return nil
}

// Equals checks if two positions are equal based on their ID
func (p *Position) Equals(other *Position) bool {
	if other == nil {
		return false
	}
	return p.id.Equals(other.id)
}