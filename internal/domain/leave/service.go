package leave

import (
	"context"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// DomainService defines domain-specific business logic that doesn't belong to a single entity
type DomainService interface {
	// ValidateEmployeeExists checks if the employee exists and is valid
	ValidateEmployeeExists(ctx context.Context, employeeID EmployeeID) error

	// ValidateApproverExists checks if the approver exists and is valid
	ValidateApproverExists(ctx context.Context, approverID EmployeeID) error

	// CanDeleteLeave checks if a leave request can be safely deleted
	CanDeleteLeave(ctx context.Context, leaveID LeaveID) error

	// ValidateNoOverlappingLeaves checks if there are no overlapping leave requests
	ValidateNoOverlappingLeaves(ctx context.Context, employeeID EmployeeID, startDate, endDate time.Time, excludeLeaveID *LeaveID) error
}

// domainService implements the DomainService interface
type domainService struct {
	leaveRepo Repository
}

// NewDomainService creates a new domain service instance
func NewDomainService(leaveRepo Repository) DomainService {
	return &domainService{
		leaveRepo: leaveRepo,
	}
}

// ValidateEmployeeExists checks if the employee exists and is valid
func (s *domainService) ValidateEmployeeExists(ctx context.Context, employeeID EmployeeID) error {
	// In a real implementation, you would check against the employee repository
	// For now, we'll just validate the format
	if employeeID.String() == "" {
		return errors.ErrInvalidEmployeeID
	}

	return nil
}

// ValidateApproverExists checks if the approver exists and is valid
func (s *domainService) ValidateApproverExists(ctx context.Context, approverID EmployeeID) error {
	// In a real implementation, you would check against the employee repository
	// and verify the approver has appropriate permissions
	// For now, we'll just validate the format
	if approverID.String() == "" {
		return errors.ErrInvalidEmployeeID
	}

	return nil
}

// CanDeleteLeave checks if a leave request can be safely deleted
func (s *domainService) CanDeleteLeave(ctx context.Context, leaveID LeaveID) error {
	// Check if leave exists
	leave, err := s.leaveRepo.FindByID(ctx, leaveID)
	if err != nil {
		return err
	}

	// Only allow deletion of pending or cancelled leave requests
	if !leave.Status().IsPending() && !leave.Status().IsCancelled() {
		return errors.DomainError{
			Code:    "INVALID_DELETE_OPERATION",
			Message: "Cannot delete approved or rejected leave requests",
			Field:   "status",
		}
	}

	return nil
}

// ValidateNoOverlappingLeaves checks if there are no overlapping leave requests
func (s *domainService) ValidateNoOverlappingLeaves(ctx context.Context, employeeID EmployeeID, startDate, endDate time.Time, excludeLeaveID *LeaveID) error {
	// Get all leave requests for the employee in the date range
	leaves, _, err := s.leaveRepo.FindByDateRange(ctx, startDate, endDate, 100, "")
	if err != nil {
		return err
	}

	// Check for overlapping leaves
	for _, leave := range leaves {
		// Skip the leave we're updating (if any)
		if excludeLeaveID != nil && leave.ID().Equals(*excludeLeaveID) {
			continue
		}

		// Skip leaves for different employees
		if !leave.EmployeeID().Equals(employeeID) {
			continue
		}

		// Skip cancelled leaves
		if leave.Status().IsCancelled() {
			continue
		}

		// Check for overlap
		if (startDate.Before(leave.EndDate()) || startDate.Equal(leave.EndDate())) &&
			(endDate.After(leave.StartDate()) || endDate.Equal(leave.StartDate())) {
			return errors.DomainError{
				Code:    "OVERLAPPING_LEAVE",
				Message: "Leave request overlaps with existing approved or pending leave",
				Field:   "dateRange",
			}
		}
	}

	return nil
}