package leave

import (
	"strings"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
	"github.com/google/uuid"
)

// LeaveID represents a leave identifier
type LeaveID struct {
	value string
}

// NewLeaveID creates a new LeaveID from a string
func NewLeaveID(value string) (LeaveID, error) {
	if value == "" {
		return LeaveID{}, errors.ErrInvalidLeaveID
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return LeaveID{}, errors.ErrInvalidLeaveID
	}

	return LeaveID{value: value}, nil
}

// GenerateLeaveID generates a new LeaveID
func GenerateLeaveID() LeaveID {
	return LeaveID{value: uuid.New().String()}
}

// String returns the string representation of the LeaveID
func (l LeaveID) String() string {
	return l.value
}

// Equals checks if two LeaveIDs are equal
func (l LeaveID) Equals(other LeaveID) bool {
	return l.value == other.value
}

// LeaveType represents the type of leave
type LeaveType struct {
	value string
}

// Predefined leave types
const (
	LeaveTypeVacation    = "VACATION"
	LeaveTypeSick        = "SICK"
	LeaveTypePersonal    = "PERSONAL"
	LeaveTypeMaternity   = "MATERNITY"
	LeaveTypePaternity   = "PATERNITY"
	LeaveTypeBereavement = "BEREAVEMENT"
	LeaveTypeUnpaid      = "UNPAID"
	LeaveTypeOther       = "OTHER"
)

// Valid leave types
var validLeaveTypes = map[string]bool{
	LeaveTypeVacation:    true,
	LeaveTypeSick:        true,
	LeaveTypePersonal:    true,
	LeaveTypeMaternity:   true,
	LeaveTypePaternity:   true,
	LeaveTypeBereavement: true,
	LeaveTypeUnpaid:      true,
	LeaveTypeOther:       true,
}

// NewLeaveType creates a new LeaveType from a string
func NewLeaveType(value string) (LeaveType, error) {
	if value == "" {
		return LeaveType{}, errors.ErrInvalidLeaveType
	}

	upperValue := strings.ToUpper(strings.TrimSpace(value))
	if !validLeaveTypes[upperValue] {
		return LeaveType{}, errors.ErrInvalidLeaveType
	}

	return LeaveType{value: upperValue}, nil
}

// String returns the string representation of the LeaveType
func (l LeaveType) String() string {
	return l.value
}

// Equals checks if two LeaveTypes are equal
func (l LeaveType) Equals(other LeaveType) bool {
	return l.value == other.value
}

// Status represents the status of a leave request
type Status struct {
	value string
}

// Predefined status values
const (
	StatusPending   = "PENDING"
	StatusApproved  = "APPROVED"
	StatusRejected  = "REJECTED"
	StatusCancelled = "CANCELLED"
)

// Valid status values
var validStatuses = map[string]bool{
	StatusPending:   true,
	StatusApproved:  true,
	StatusRejected:  true,
	StatusCancelled: true,
}

// NewStatus creates a new Status from a string
func NewStatus(value string) (Status, error) {
	if value == "" {
		return Status{}, errors.ErrInvalidLeaveStatus
	}

	upperValue := strings.ToUpper(strings.TrimSpace(value))
	if !validStatuses[upperValue] {
		return Status{}, errors.ErrInvalidLeaveStatus
	}

	return Status{value: upperValue}, nil
}

// String returns the string representation of the Status
func (s Status) String() string {
	return s.value
}

// Equals checks if two Statuses are equal
func (s Status) Equals(other Status) bool {
	return s.value == other.value
}

// IsPending checks if the status is pending
func (s Status) IsPending() bool {
	return s.value == StatusPending
}

// IsApproved checks if the status is approved
func (s Status) IsApproved() bool {
	return s.value == StatusApproved
}

// IsRejected checks if the status is rejected
func (s Status) IsRejected() bool {
	return s.value == StatusRejected
}

// IsCancelled checks if the status is cancelled
func (s Status) IsCancelled() bool {
	return s.value == StatusCancelled
}

// Reason represents the reason for leave
type Reason struct {
	value string
}

// NewReason creates a new Reason from a string
func NewReason(value string) (Reason, error) {
	if value == "" {
		return Reason{}, errors.ErrInvalidLeaveReason
	}

	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Reason{}, errors.ErrInvalidLeaveReason
	}

	if len(trimmed) > 500 {
		return Reason{}, errors.DomainError{
			Code:    "LEAVE_REASON_TOO_LONG",
			Message: "Leave reason cannot exceed 500 characters",
			Field:   "reason",
		}
	}

	return Reason{value: trimmed}, nil
}

// String returns the string representation of the Reason
func (r Reason) String() string {
	return r.value
}

// Equals checks if two Reasons are equal
func (r Reason) Equals(other Reason) bool {
	return r.value == other.value
}

// EmployeeID represents an employee identifier (reused from employee domain)
type EmployeeID struct {
	value string
}

// NewEmployeeID creates a new EmployeeID from a string
func NewEmployeeID(value string) (EmployeeID, error) {
	if value == "" {
		return EmployeeID{}, errors.ErrInvalidEmployeeID
	}

	// Validate UUID format
	if _, err := uuid.Parse(value); err != nil {
		return EmployeeID{}, errors.ErrInvalidEmployeeID
	}

	return EmployeeID{value: value}, nil
}

// String returns the string representation of the EmployeeID
func (e EmployeeID) String() string {
	return e.value
}

// Equals checks if two EmployeeIDs are equal
func (e EmployeeID) Equals(other EmployeeID) bool {
	return e.value == other.value
}

// ValidateLeaveID validates a leave ID string
func ValidateLeaveID(id string) error {
	if id == "" {
		return errors.ErrInvalidLeaveID
	}

	if _, err := uuid.Parse(id); err != nil {
		return errors.ErrInvalidLeaveID
	}

	return nil
}

// ValidateLeaveType validates a leave type string
func ValidateLeaveType(leaveType string) error {
	if leaveType == "" {
		return errors.ErrInvalidLeaveType
	}

	upperValue := strings.ToUpper(strings.TrimSpace(leaveType))
	if !validLeaveTypes[upperValue] {
		return errors.ErrInvalidLeaveType
	}

	return nil
}

// ValidateLeaveStatus validates a leave status string
func ValidateLeaveStatus(status string) error {
	if status == "" {
		return errors.ErrInvalidLeaveStatus
	}

	upperValue := strings.ToUpper(strings.TrimSpace(status))
	if !validStatuses[upperValue] {
		return errors.ErrInvalidLeaveStatus
	}

	return nil
}

// ValidateLeaveReason validates a leave reason string
func ValidateLeaveReason(reason string) error {
	if reason == "" {
		return errors.ErrInvalidLeaveReason
	}

	trimmed := strings.TrimSpace(reason)
	if trimmed == "" {
		return errors.ErrInvalidLeaveReason
	}

	if len(trimmed) > 500 {
		return errors.DomainError{
			Code:    "LEAVE_REASON_TOO_LONG",
			Message: "Leave reason cannot exceed 500 characters",
			Field:   "reason",
		}
	}

	return nil
}