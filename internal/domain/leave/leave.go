package leave

import (
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

// Leave represents the leave request domain entity
type Leave struct {
	id          LeaveID
	employeeID  EmployeeID
	leaveType   LeaveType
	startDate   time.Time
	endDate     time.Time
	reason      Reason
	status      Status
	approvedBy  *EmployeeID
	approvedAt  *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

// NewLeave creates a new Leave entity with validation
func NewLeave(employeeID string, leaveType string, startDate, endDate time.Time, reason string) (*Leave, error) {
	empID, err := NewEmployeeID(employeeID)
	if err != nil {
		return nil, err
	}

	leaveTypeVO, err := NewLeaveType(leaveType)
	if err != nil {
		return nil, err
	}

	reasonVO, err := NewReason(reason)
	if err != nil {
		return nil, err
	}

	// Validate dates
	if startDate.IsZero() {
		return nil, errors.DomainError{
			Code:    "INVALID_START_DATE",
			Message: "Start date cannot be zero",
			Field:   "startDate",
		}
	}

	if endDate.IsZero() {
		return nil, errors.DomainError{
			Code:    "INVALID_END_DATE",
			Message: "End date cannot be zero",
			Field:   "endDate",
		}
	}

	if startDate.After(endDate) {
		return nil, errors.DomainError{
			Code:    "INVALID_DATE_RANGE",
			Message: "Start date cannot be after end date",
			Field:   "startDate",
		}
	}

	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.DomainError{
			Code:    "INVALID_START_DATE",
			Message: "Start date cannot be in the past",
			Field:   "startDate",
		}
	}

	now := time.Now()

	return &Leave{
		id:         GenerateLeaveID(),
		employeeID: empID,
		leaveType:  leaveTypeVO,
		startDate:  startDate,
		endDate:    endDate,
		reason:     reasonVO,
		status:     Status{value: StatusPending},
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// NewLeaveWithID creates a Leave entity with a specific ID (for reconstruction from persistence)
func NewLeaveWithID(id, employeeID, leaveType, reason string, startDate, endDate time.Time, status string, approvedBy *string, approvedAt *time.Time, createdAt, updatedAt time.Time) (*Leave, error) {
	leaveID, err := NewLeaveID(id)
	if err != nil {
		return nil, err
	}

	empID, err := NewEmployeeID(employeeID)
	if err != nil {
		return nil, err
	}

	leaveTypeVO, err := NewLeaveType(leaveType)
	if err != nil {
		return nil, err
	}

	reasonVO, err := NewReason(reason)
	if err != nil {
		return nil, err
	}

	statusVO, err := NewStatus(status)
	if err != nil {
		return nil, err
	}

	var approvedByVO *EmployeeID
	if approvedBy != nil {
		empID, err := NewEmployeeID(*approvedBy)
		if err != nil {
			return nil, err
		}
		approvedByVO = &empID
	}

	return &Leave{
		id:         leaveID,
		employeeID: empID,
		leaveType:  leaveTypeVO,
		startDate:  startDate,
		endDate:    endDate,
		reason:     reasonVO,
		status:     statusVO,
		approvedBy: approvedByVO,
		approvedAt: approvedAt,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}, nil
}

// ID returns the leave's ID
func (l *Leave) ID() LeaveID {
	return l.id
}

// EmployeeID returns the leave's employee ID
func (l *Leave) EmployeeID() EmployeeID {
	return l.employeeID
}

// LeaveType returns the leave's type
func (l *Leave) LeaveType() LeaveType {
	return l.leaveType
}

// StartDate returns the leave's start date
func (l *Leave) StartDate() time.Time {
	return l.startDate
}

// EndDate returns the leave's end date
func (l *Leave) EndDate() time.Time {
	return l.endDate
}

// Reason returns the leave's reason
func (l *Leave) Reason() Reason {
	return l.reason
}

// Status returns the leave's status
func (l *Leave) Status() Status {
	return l.status
}

// ApprovedBy returns the leave's approver ID
func (l *Leave) ApprovedBy() *EmployeeID {
	return l.approvedBy
}

// ApprovedAt returns when the leave was approved
func (l *Leave) ApprovedAt() *time.Time {
	return l.approvedAt
}

// CreatedAt returns when the leave was created
func (l *Leave) CreatedAt() time.Time {
	return l.createdAt
}

// UpdatedAt returns when the leave was last updated
func (l *Leave) UpdatedAt() time.Time {
	return l.updatedAt
}

// Approve approves the leave request
func (l *Leave) Approve(approvedBy string) error {
	approvedByVO, err := NewEmployeeID(approvedBy)
	if err != nil {
		return err
	}

	if !l.status.IsPending() {
		return errors.DomainError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "Only pending leave requests can be approved",
			Field:   "status",
		}
	}

	approvedStatus, _ := NewStatus(StatusApproved)
	l.status = approvedStatus
	l.approvedBy = &approvedByVO
	now := time.Now()
	l.approvedAt = &now
	l.updatedAt = now
	return nil
}

// Reject rejects the leave request
func (l *Leave) Reject(approvedBy string) error {
	approvedByVO, err := NewEmployeeID(approvedBy)
	if err != nil {
		return err
	}

	if !l.status.IsPending() {
		return errors.DomainError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "Only pending leave requests can be rejected",
			Field:   "status",
		}
	}

	rejectedStatus, _ := NewStatus(StatusRejected)
	l.status = rejectedStatus
	l.approvedBy = &approvedByVO
	now := time.Now()
	l.approvedAt = &now
	l.updatedAt = now
	return nil
}

// Cancel cancels the leave request
func (l *Leave) Cancel() error {
	if !l.status.IsPending() && !l.status.IsApproved() {
		return errors.DomainError{
			Code:    "INVALID_STATUS_TRANSITION",
			Message: "Only pending or approved leave requests can be cancelled",
			Field:   "status",
		}
	}

	cancelledStatus, _ := NewStatus(StatusCancelled)
	l.status = cancelledStatus
	l.updatedAt = time.Now()
	return nil
}

// UpdateReason updates the leave's reason with validation
func (l *Leave) UpdateReason(reason string) error {
	if !l.status.IsPending() {
		return errors.DomainError{
			Code:    "INVALID_UPDATE",
			Message: "Cannot update reason for non-pending leave requests",
			Field:   "reason",
		}
	}

	reasonVO, err := NewReason(reason)
	if err != nil {
		return err
	}

	l.reason = reasonVO
	l.updatedAt = time.Now()
	return nil
}

// UpdateDates updates the leave's dates with validation
func (l *Leave) UpdateDates(startDate, endDate time.Time) error {
	if !l.status.IsPending() {
		return errors.DomainError{
			Code:    "INVALID_UPDATE",
			Message: "Cannot update dates for non-pending leave requests",
			Field:   "dates",
		}
	}

	if startDate.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_START_DATE",
			Message: "Start date cannot be zero",
			Field:   "startDate",
		}
	}

	if endDate.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_END_DATE",
			Message: "End date cannot be zero",
			Field:   "endDate",
		}
	}

	if startDate.After(endDate) {
		return errors.DomainError{
			Code:    "INVALID_DATE_RANGE",
			Message: "Start date cannot be after end date",
			Field:   "startDate",
		}
	}

	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.DomainError{
			Code:    "INVALID_START_DATE",
			Message: "Start date cannot be in the past",
			Field:   "startDate",
		}
	}

	l.startDate = startDate
	l.endDate = endDate
	l.updatedAt = time.Now()
	return nil
}

// Validate performs comprehensive validation of the leave entity
func (l *Leave) Validate() error {
	if l.id.String() == "" {
		return errors.ErrInvalidLeaveID
	}

	if l.employeeID.String() == "" {
		return errors.ErrInvalidEmployeeID
	}

	if l.leaveType.String() == "" {
		return errors.ErrInvalidLeaveType
	}

	if l.reason.String() == "" {
		return errors.ErrInvalidLeaveReason
	}

	if l.status.String() == "" {
		return errors.ErrInvalidLeaveStatus
	}

	if l.startDate.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_START_DATE",
			Message: "Start date cannot be zero",
			Field:   "startDate",
		}
	}

	if l.endDate.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_END_DATE",
			Message: "End date cannot be zero",
			Field:   "endDate",
		}
	}

	if l.startDate.After(l.endDate) {
		return errors.DomainError{
			Code:    "INVALID_DATE_RANGE",
			Message: "Start date cannot be after end date",
			Field:   "startDate",
		}
	}

	if l.createdAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_CREATED_AT",
			Message: "Created at timestamp cannot be zero",
			Field:   "createdAt",
		}
	}

	if l.updatedAt.IsZero() {
		return errors.DomainError{
			Code:    "INVALID_UPDATED_AT",
			Message: "Updated at timestamp cannot be zero",
			Field:   "updatedAt",
		}
	}

	if l.updatedAt.Before(l.createdAt) {
		return errors.DomainError{
			Code:    "INVALID_TIMESTAMPS",
			Message: "Updated at cannot be before created at",
			Field:   "updatedAt",
		}
	}

	return nil
}

// Equals checks if two leaves are equal based on their ID
func (l *Leave) Equals(other *Leave) bool {
	if other == nil {
		return false
	}
	return l.id.Equals(other.id)
}