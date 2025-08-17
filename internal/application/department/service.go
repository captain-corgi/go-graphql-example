package department

import (
	"context"
	"log/slog"

	"github.com/captain-corgi/go-graphql-example/internal/domain/department"
	"github.com/captain-corgi/go-graphql-example/internal/domain/errors"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Service defines the interface for department application services
type Service interface {
	// GetDepartment retrieves a department by ID
	GetDepartment(ctx context.Context, req GetDepartmentRequest) (*GetDepartmentResponse, error)

	// ListDepartments retrieves a paginated list of departments
	ListDepartments(ctx context.Context, req ListDepartmentsRequest) (*ListDepartmentsResponse, error)

	// ListDepartmentsByManager retrieves departments by manager ID
	ListDepartmentsByManager(ctx context.Context, req ListDepartmentsByManagerRequest) (*ListDepartmentsByManagerResponse, error)

	// CreateDepartment creates a new department
	CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*CreateDepartmentResponse, error)

	// UpdateDepartment updates an existing department
	UpdateDepartment(ctx context.Context, req UpdateDepartmentRequest) (*UpdateDepartmentResponse, error)

	// DeleteDepartment deletes a department by ID
	DeleteDepartment(ctx context.Context, req DeleteDepartmentRequest) (*DeleteDepartmentResponse, error)
}

// service implements the Service interface
type service struct {
	deptRepo     department.Repository
	deptService  department.DomainService
	logger       *slog.Logger
}

// NewService creates a new department service
func NewService(deptRepo department.Repository, deptService department.DomainService, logger *slog.Logger) Service {
	return &service{
		deptRepo:    deptRepo,
		deptService: deptService,
		logger:      logger,
	}
}

// GetDepartment retrieves a department by ID
func (s *service) GetDepartment(ctx context.Context, req GetDepartmentRequest) (*GetDepartmentResponse, error) {
	s.logger.InfoContext(ctx, "Getting department", "departmentID", req.ID)

	// Validate request
	if err := s.validateGetDepartmentRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid get department request", "error", err, "departmentID", req.ID)
		return &GetDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Convert string ID to domain DepartmentID
	deptID, err := department.NewDepartmentID(req.ID)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid department ID format", "error", err, "departmentID", req.ID)
		return &GetDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Retrieve department from repository
	domainDept, err := s.deptRepo.FindByID(ctx, deptID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get department from repository", "error", err, "departmentID", req.ID)
		return &GetDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully retrieved department", "departmentID", req.ID)
	return &GetDepartmentResponse{
		Department: mapDomainDepartmentToDTO(domainDept),
	}, nil
}

// ListDepartments retrieves a paginated list of departments
func (s *service) ListDepartments(ctx context.Context, req ListDepartmentsRequest) (*ListDepartmentsResponse, error) {
	s.logger.InfoContext(ctx, "Listing departments", "limit", req.Limit, "cursor", req.Cursor)

	// Validate request
	if err := s.validateListDepartmentsRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid list departments request", "error", err)
		return &ListDepartmentsResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Set default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 10 // Default page size
	}

	// Retrieve departments from repository
	domainDepts, nextCursor, err := s.deptRepo.FindAll(ctx, limit+1, req.Cursor) // +1 to check if there's a next page
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list departments from repository", "error", err)
		return &ListDepartmentsResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Check if there are more pages
	hasNextPage := len(domainDepts) > limit
	if hasNextPage {
		domainDepts = domainDepts[:limit] // Remove the extra item
	}

	s.logger.InfoContext(ctx, "Successfully listed departments", "count", len(domainDepts), "hasNextPage", hasNextPage)
	return &ListDepartmentsResponse{
		Departments: &DepartmentConnectionDTO{
			Departments: mapDomainDepartmentsToDTOs(domainDepts),
			NextCursor:  nextCursor,
		},
	}, nil
}

// ListDepartmentsByManager retrieves departments by manager ID
func (s *service) ListDepartmentsByManager(ctx context.Context, req ListDepartmentsByManagerRequest) (*ListDepartmentsByManagerResponse, error) {
	s.logger.InfoContext(ctx, "Listing departments by manager", "managerID", req.ManagerID, "limit", req.Limit, "cursor", req.Cursor)

	// Validate request
	if err := s.validateListDepartmentsByManagerRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid list departments by manager request", "error", err)
		return &ListDepartmentsByManagerResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Convert string manager ID to domain EmployeeID
	managerID, err := department.NewEmployeeID(req.ManagerID)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid manager ID format", "error", err, "managerID", req.ManagerID)
		return &ListDepartmentsByManagerResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Set default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 10 // Default page size
	}

	// Retrieve departments from repository
	domainDepts, nextCursor, err := s.deptRepo.FindByManager(ctx, managerID, limit+1, req.Cursor) // +1 to check if there's a next page
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list departments by manager from repository", "error", err)
		return &ListDepartmentsByManagerResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Check if there are more pages
	hasNextPage := len(domainDepts) > limit
	if hasNextPage {
		domainDepts = domainDepts[:limit] // Remove the extra item
	}

	s.logger.InfoContext(ctx, "Successfully listed departments by manager", "count", len(domainDepts), "hasNextPage", hasNextPage)
	return &ListDepartmentsByManagerResponse{
		Departments: &DepartmentConnectionDTO{
			Departments: mapDomainDepartmentsToDTOs(domainDepts),
			NextCursor:  nextCursor,
		},
	}, nil
}

// CreateDepartment creates a new department
func (s *service) CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*CreateDepartmentResponse, error) {
	s.logger.InfoContext(ctx, "Creating department", "name", req.Name)

	// Validate request
	if err := s.validateCreateDepartmentRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid create department request", "error", err)
		return &CreateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Validate unique name
	name, err := department.NewName(req.Name)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid department name", "error", err)
		return &CreateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	if err := s.deptService.ValidateUniqueName(ctx, name, nil); err != nil {
		s.logger.WarnContext(ctx, "Department name already exists", "error", err)
		return &CreateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Create domain department
	domainDept, err := department.NewDepartment(req.Name, req.Description, req.ManagerID)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to create domain department", "error", err)
		return &CreateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Validate domain entity
	if err := domainDept.Validate(); err != nil {
		s.logger.WarnContext(ctx, "Invalid domain department", "error", err)
		return &CreateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Save to repository
	if err := s.deptRepo.Save(ctx, domainDept); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save department to repository", "error", err)
		return &CreateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully created department", "departmentID", domainDept.ID().String())
	return &CreateDepartmentResponse{
		Department: mapDomainDepartmentToDTO(domainDept),
	}, nil
}

// UpdateDepartment updates an existing department
func (s *service) UpdateDepartment(ctx context.Context, req UpdateDepartmentRequest) (*UpdateDepartmentResponse, error) {
	s.logger.InfoContext(ctx, "Updating department", "departmentID", req.ID)

	// Validate request
	if err := s.validateUpdateDepartmentRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid update department request", "error", err)
		return &UpdateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Convert string ID to domain DepartmentID
	deptID, err := department.NewDepartmentID(req.ID)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid department ID format", "error", err, "departmentID", req.ID)
		return &UpdateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Retrieve existing department
	domainDept, err := s.deptRepo.FindByID(ctx, deptID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get department from repository", "error", err, "departmentID", req.ID)
		return &UpdateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Update fields if provided
	if req.Name != nil {
		if err := domainDept.UpdateName(*req.Name); err != nil {
			s.logger.WarnContext(ctx, "Failed to update department name", "error", err)
			return &UpdateDepartmentResponse{
				Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
			}, nil
		}
	}

	if req.Description != nil {
		if err := domainDept.UpdateDescription(*req.Description); err != nil {
			s.logger.WarnContext(ctx, "Failed to update department description", "error", err)
			return &UpdateDepartmentResponse{
				Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
			}, nil
		}
	}

	if req.ManagerID != nil {
		if err := domainDept.UpdateManager(req.ManagerID); err != nil {
			s.logger.WarnContext(ctx, "Failed to update department manager", "error", err)
			return &UpdateDepartmentResponse{
				Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
			}, nil
		}
	}

	// Validate unique name if name was updated
	if req.Name != nil {
		name, err := department.NewName(*req.Name)
		if err != nil {
			s.logger.WarnContext(ctx, "Invalid department name", "error", err)
			return &UpdateDepartmentResponse{
				Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
			}, nil
		}

		if err := s.deptService.ValidateUniqueName(ctx, name, &deptID); err != nil {
			s.logger.WarnContext(ctx, "Department name already exists", "error", err)
			return &UpdateDepartmentResponse{
				Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
			}, nil
		}
	}

	// Validate domain entity
	if err := domainDept.Validate(); err != nil {
		s.logger.WarnContext(ctx, "Invalid domain department", "error", err)
		return &UpdateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Save to repository
	if err := s.deptRepo.Save(ctx, domainDept); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save department to repository", "error", err)
		return &UpdateDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully updated department", "departmentID", req.ID)
	return &UpdateDepartmentResponse{
		Department: mapDomainDepartmentToDTO(domainDept),
	}, nil
}

// DeleteDepartment deletes a department by ID
func (s *service) DeleteDepartment(ctx context.Context, req DeleteDepartmentRequest) (*DeleteDepartmentResponse, error) {
	s.logger.InfoContext(ctx, "Deleting department", "departmentID", req.ID)

	// Validate request
	if err := s.validateDeleteDepartmentRequest(req); err != nil {
		s.logger.WarnContext(ctx, "Invalid delete department request", "error", err)
		return &DeleteDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Convert string ID to domain DepartmentID
	deptID, err := department.NewDepartmentID(req.ID)
	if err != nil {
		s.logger.WarnContext(ctx, "Invalid department ID format", "error", err, "departmentID", req.ID)
		return &DeleteDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Check if department can be deleted
	if err := s.deptService.CanDeleteDepartment(ctx, deptID); err != nil {
		s.logger.WarnContext(ctx, "Department cannot be deleted", "error", err)
		return &DeleteDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	// Delete from repository
	if err := s.deptRepo.Delete(ctx, deptID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete department from repository", "error", err)
		return &DeleteDepartmentResponse{
			Errors: []ErrorDTO{mapDomainErrorToDTO(err)},
		}, nil
	}

	s.logger.InfoContext(ctx, "Successfully deleted department", "departmentID", req.ID)
	return &DeleteDepartmentResponse{
		Success: true,
	}, nil
}

// Validation methods
func (s *service) validateGetDepartmentRequest(req GetDepartmentRequest) error {
	if req.ID == "" {
		return department.ValidateDepartmentID(req.ID)
	}
	return nil
}

func (s *service) validateListDepartmentsRequest(req ListDepartmentsRequest) error {
	if req.Limit < 0 {
		return errors.DomainError{
			Message: "Limit cannot be negative",
			Field:   "limit",
			Code:    "INVALID_LIMIT",
		}
	}
	return nil
}

func (s *service) validateListDepartmentsByManagerRequest(req ListDepartmentsByManagerRequest) error {
	if req.ManagerID == "" {
		return errors.DomainError{
			Message: "Manager ID cannot be empty",
			Field:   "managerId",
			Code:    "INVALID_MANAGER_ID",
		}
	}
	if req.Limit < 0 {
		return errors.DomainError{
			Message: "Limit cannot be negative",
			Field:   "limit",
			Code:    "INVALID_LIMIT",
		}
	}
	return nil
}

func (s *service) validateCreateDepartmentRequest(req CreateDepartmentRequest) error {
	if err := department.ValidateDepartmentName(req.Name); err != nil {
		return err
	}
	if err := department.ValidateDepartmentDescription(req.Description); err != nil {
		return err
	}
	return nil
}

func (s *service) validateUpdateDepartmentRequest(req UpdateDepartmentRequest) error {
	if req.ID == "" {
		return department.ValidateDepartmentID(req.ID)
	}
	if req.Name != nil {
		if err := department.ValidateDepartmentName(*req.Name); err != nil {
			return err
		}
	}
	if req.Description != nil {
		if err := department.ValidateDepartmentDescription(*req.Description); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) validateDeleteDepartmentRequest(req DeleteDepartmentRequest) error {
	if req.ID == "" {
		return department.ValidateDepartmentID(req.ID)
	}
	return nil
}