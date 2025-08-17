package employee

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/captain-corgi/go-graphql-example/internal/domain/employee"
	"github.com/captain-corgi/go-graphql-example/internal/domain/user"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks

// Service defines the interface for employee application operations
type Service interface {
	GetEmployee(ctx context.Context, req GetEmployeeRequest) (*GetEmployeeResponse, error)
	ListEmployees(ctx context.Context, req ListEmployeesRequest) (*ListEmployeesResponse, error)
	ListEmployeesByDepartment(ctx context.Context, req ListEmployeesByDepartmentRequest) (*ListEmployeesByDepartmentResponse, error)
	ListEmployeesByStatus(ctx context.Context, req ListEmployeesByStatusRequest) (*ListEmployeesByStatusResponse, error)
	CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*CreateEmployeeResponse, error)
	UpdateEmployee(ctx context.Context, req UpdateEmployeeRequest) (*UpdateEmployeeResponse, error)
	DeleteEmployee(ctx context.Context, req DeleteEmployeeRequest) (*DeleteEmployeeResponse, error)
}

// service implements the Service interface
type service struct {
	employeeRepo employee.Repository
	userRepo     user.Repository
	logger       *slog.Logger
}

// NewService creates a new employee service
func NewService(employeeRepo employee.Repository, userRepo user.Repository, logger *slog.Logger) Service {
	return &service{
		employeeRepo: employeeRepo,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// GetEmployee retrieves an employee by ID
func (s *service) GetEmployee(ctx context.Context, req GetEmployeeRequest) (*GetEmployeeResponse, error) {
	s.logger.Info("Getting employee", "id", req.ID)

	// Validate employee ID
	employeeID, err := employee.NewEmployeeID(req.ID)
	if err != nil {
		s.logger.Error("Invalid employee ID", "id", req.ID, "error", err)
		return &GetEmployeeResponse{}, err
	}

	// Get employee from repository
	emp, err := s.employeeRepo.FindByID(ctx, employeeID)
	if err != nil {
		s.logger.Error("Failed to get employee", "id", req.ID, "error", err)
		return &GetEmployeeResponse{}, err
	}

	// Get associated user
	user, err := s.userRepo.FindByID(ctx, emp.UserID())
	if err != nil {
		s.logger.Error("Failed to get user for employee", "user_id", emp.UserID().String(), "error", err)
		// Continue without user data rather than failing completely
	}

	// Map to DTO
	employeeDTO := MapEmployeeToDTO(emp, user)

	s.logger.Info("Successfully retrieved employee", "id", req.ID)
	return &GetEmployeeResponse{
		Employee: employeeDTO,
	}, nil
}

// ListEmployees retrieves a list of employees with pagination
func (s *service) ListEmployees(ctx context.Context, req ListEmployeesRequest) (*ListEmployeesResponse, error) {
	s.logger.Info("Listing employees", "limit", req.Limit, "cursor", req.Cursor)

	// Set default limit if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Get employees from repository
	employees, nextCursor, err := s.employeeRepo.FindAll(ctx, limit, req.Cursor)
	if err != nil {
		s.logger.Error("Failed to list employees", "error", err)
		return &ListEmployeesResponse{}, err
	}

	// Get associated users
	userIDs := make([]user.UserID, len(employees))
	for i, emp := range employees {
		userIDs[i] = emp.UserID()
	}

	users := make(map[string]*user.User)
	for _, emp := range employees {
		user, err := s.userRepo.FindByID(ctx, emp.UserID())
		if err != nil {
			s.logger.Warn("Failed to get user for employee", "user_id", emp.UserID().String(), "error", err)
			continue
		}
		users[emp.UserID().String()] = user
	}

	// Map to DTOs
	employeeDTOs := MapEmployeesToDTOs(employees, users)

	s.logger.Info("Successfully listed employees", "count", len(employeeDTOs))
	return &ListEmployeesResponse{
		Employees:  employeeDTOs,
		NextCursor: nextCursor,
	}, nil
}

// ListEmployeesByDepartment retrieves employees by department
func (s *service) ListEmployeesByDepartment(ctx context.Context, req ListEmployeesByDepartmentRequest) (*ListEmployeesByDepartmentResponse, error) {
	s.logger.Info("Listing employees by department", "department", req.Department, "limit", req.Limit, "cursor", req.Cursor)

	// Validate department
	dept, err := employee.NewDepartment(req.Department)
	if err != nil {
		s.logger.Error("Invalid department", "department", req.Department, "error", err)
		return &ListEmployeesByDepartmentResponse{}, err
	}

	// Set default limit if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Get employees from repository
	employees, nextCursor, err := s.employeeRepo.FindByDepartment(ctx, dept, limit, req.Cursor)
	if err != nil {
		s.logger.Error("Failed to list employees by department", "department", req.Department, "error", err)
		return &ListEmployeesByDepartmentResponse{}, err
	}

	// Get associated users
	users := make(map[string]*user.User)
	for _, emp := range employees {
		user, err := s.userRepo.FindByID(ctx, emp.UserID())
		if err != nil {
			s.logger.Warn("Failed to get user for employee", "user_id", emp.UserID().String(), "error", err)
			continue
		}
		users[emp.UserID().String()] = user
	}

	// Map to DTOs
	employeeDTOs := MapEmployeesToDTOs(employees, users)

	s.logger.Info("Successfully listed employees by department", "department", req.Department, "count", len(employeeDTOs))
	return &ListEmployeesByDepartmentResponse{
		Employees:  employeeDTOs,
		NextCursor: nextCursor,
	}, nil
}

// ListEmployeesByStatus retrieves employees by status
func (s *service) ListEmployeesByStatus(ctx context.Context, req ListEmployeesByStatusRequest) (*ListEmployeesByStatusResponse, error) {
	s.logger.Info("Listing employees by status", "status", req.Status, "limit", req.Limit, "cursor", req.Cursor)

	// Validate status
	status, err := employee.NewStatus(req.Status)
	if err != nil {
		s.logger.Error("Invalid status", "status", req.Status, "error", err)
		return &ListEmployeesByStatusResponse{}, err
	}

	// Set default limit if not provided
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Get employees from repository
	employees, nextCursor, err := s.employeeRepo.FindByStatus(ctx, status, limit, req.Cursor)
	if err != nil {
		s.logger.Error("Failed to list employees by status", "status", req.Status, "error", err)
		return &ListEmployeesByStatusResponse{}, err
	}

	// Get associated users
	users := make(map[string]*user.User)
	for _, emp := range employees {
		user, err := s.userRepo.FindByID(ctx, emp.UserID())
		if err != nil {
			s.logger.Warn("Failed to get user for employee", "user_id", emp.UserID().String(), "error", err)
			continue
		}
		users[emp.UserID().String()] = user
	}

	// Map to DTOs
	employeeDTOs := MapEmployeesToDTOs(employees, users)

	s.logger.Info("Successfully listed employees by status", "status", req.Status, "count", len(employeeDTOs))
	return &ListEmployeesByStatusResponse{
		Employees:  employeeDTOs,
		NextCursor: nextCursor,
	}, nil
}

// CreateEmployee creates a new employee
func (s *service) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*CreateEmployeeResponse, error) {
	s.logger.Info("Creating employee", "user_id", req.UserID, "employee_code", req.EmployeeCode)

	var errors []error

	// Validate user ID
	userID, err := user.NewUserID(req.UserID)
	if err != nil {
		s.logger.Error("Invalid user ID", "user_id", req.UserID, "error", err)
		errors = append(errors, err)
	}

	// Check if user exists
	if err == nil {
		_, err = s.userRepo.FindByID(ctx, userID)
		if err != nil {
			s.logger.Error("User not found", "user_id", req.UserID, "error", err)
			errors = append(errors, err)
		}
	}

	// Check if employee already exists for this user
	if err == nil {
		exists, err := s.employeeRepo.ExistsByUserID(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to check if employee exists", "user_id", req.UserID, "error", err)
			errors = append(errors, err)
		} else if exists {
			s.logger.Error("Employee already exists for user", "user_id", req.UserID)
			errors = append(errors, fmt.Errorf("employee already exists for user"))
		}
	}

	// Check if employee code already exists
	employeeCode, err := employee.NewEmployeeCode(req.EmployeeCode)
	if err != nil {
		s.logger.Error("Invalid employee code", "employee_code", req.EmployeeCode, "error", err)
		errors = append(errors, err)
	} else {
		exists, err := s.employeeRepo.ExistsByEmployeeCode(ctx, employeeCode)
		if err != nil {
			s.logger.Error("Failed to check if employee code exists", "employee_code", req.EmployeeCode, "error", err)
			errors = append(errors, err)
		} else if exists {
			s.logger.Error("Employee code already exists", "employee_code", req.EmployeeCode)
			errors = append(errors, fmt.Errorf("employee code already exists"))
		}
	}

	// If there are validation errors, return them
	if len(errors) > 0 {
		return &CreateEmployeeResponse{
			Errors: MapErrorsToDTOs(errors),
		}, nil
	}

	// Create employee domain entity
	emp, err := employee.NewEmployee(
		userID,
		req.EmployeeCode,
		req.Department,
		req.Position,
		req.HireDate,
		req.Salary,
		req.Status,
	)
	if err != nil {
		s.logger.Error("Failed to create employee entity", "error", err)
		return &CreateEmployeeResponse{
			Errors: []ErrorDTO{MapErrorToDTO(err)},
		}, nil
	}

	// Validate employee
	if err := emp.Validate(); err != nil {
		s.logger.Error("Employee validation failed", "error", err)
		return &CreateEmployeeResponse{
			Errors: []ErrorDTO{MapErrorToDTO(err)},
		}, nil
	}

	// Save to repository
	if err := s.employeeRepo.Create(ctx, emp); err != nil {
		s.logger.Error("Failed to save employee", "error", err)
		return &CreateEmployeeResponse{
			Errors: []ErrorDTO{MapErrorToDTO(err)},
		}, nil
	}

	// Get associated user for response
	user, err := s.userRepo.FindByID(ctx, emp.UserID())
	if err != nil {
		s.logger.Warn("Failed to get user for created employee", "user_id", emp.UserID().String(), "error", err)
	}

	// Map to DTO
	employeeDTO := MapEmployeeToDTO(emp, user)

	s.logger.Info("Successfully created employee", "id", emp.ID().String())
	return &CreateEmployeeResponse{
		Employee: employeeDTO,
	}, nil
}

// UpdateEmployee updates an existing employee
func (s *service) UpdateEmployee(ctx context.Context, req UpdateEmployeeRequest) (*UpdateEmployeeResponse, error) {
	s.logger.Info("Updating employee", "id", req.ID)

	var errors []error

	// Validate employee ID
	employeeID, err := employee.NewEmployeeID(req.ID)
	if err != nil {
		s.logger.Error("Invalid employee ID", "id", req.ID, "error", err)
		errors = append(errors, err)
	}

	// Get existing employee
	var emp *employee.Employee
	if err == nil {
		emp, err = s.employeeRepo.FindByID(ctx, employeeID)
		if err != nil {
			s.logger.Error("Employee not found", "id", req.ID, "error", err)
			errors = append(errors, err)
		}
	}

	// If there are validation errors, return them
	if len(errors) > 0 {
		return &UpdateEmployeeResponse{
			Errors: MapErrorsToDTOs(errors),
		}, nil
	}

	// Update fields if provided
	if req.EmployeeCode != nil {
		if err := emp.UpdateEmployeeCode(*req.EmployeeCode); err != nil {
			s.logger.Error("Failed to update employee code", "employee_code", *req.EmployeeCode, "error", err)
			errors = append(errors, err)
		}
	}

	if req.Department != nil {
		if err := emp.UpdateDepartment(*req.Department); err != nil {
			s.logger.Error("Failed to update department", "department", *req.Department, "error", err)
			errors = append(errors, err)
		}
	}

	if req.Position != nil {
		if err := emp.UpdatePosition(*req.Position); err != nil {
			s.logger.Error("Failed to update position", "position", *req.Position, "error", err)
			errors = append(errors, err)
		}
	}

	if req.HireDate != nil {
		if err := emp.UpdateHireDate(*req.HireDate); err != nil {
			s.logger.Error("Failed to update hire date", "hire_date", *req.HireDate, "error", err)
			errors = append(errors, err)
		}
	}

	if req.Salary != nil {
		if err := emp.UpdateSalary(*req.Salary); err != nil {
			s.logger.Error("Failed to update salary", "salary", *req.Salary, "error", err)
			errors = append(errors, err)
		}
	}

	if req.Status != nil {
		if err := emp.UpdateStatus(*req.Status); err != nil {
			s.logger.Error("Failed to update status", "status", *req.Status, "error", err)
			errors = append(errors, err)
		}
	}

	// If there are update errors, return them
	if len(errors) > 0 {
		return &UpdateEmployeeResponse{
			Errors: MapErrorsToDTOs(errors),
		}, nil
	}

	// Validate employee
	if err := emp.Validate(); err != nil {
		s.logger.Error("Employee validation failed", "error", err)
		return &UpdateEmployeeResponse{
			Errors: []ErrorDTO{MapErrorToDTO(err)},
		}, nil
	}

	// Save to repository
	if err := s.employeeRepo.Update(ctx, emp); err != nil {
		s.logger.Error("Failed to save employee", "error", err)
		return &UpdateEmployeeResponse{
			Errors: []ErrorDTO{MapErrorToDTO(err)},
		}, nil
	}

	// Get associated user for response
	user, err := s.userRepo.FindByID(ctx, emp.UserID())
	if err != nil {
		s.logger.Warn("Failed to get user for updated employee", "user_id", emp.UserID().String(), "error", err)
	}

	// Map to DTO
	employeeDTO := MapEmployeeToDTO(emp, user)

	s.logger.Info("Successfully updated employee", "id", req.ID)
	return &UpdateEmployeeResponse{
		Employee: employeeDTO,
	}, nil
}

// DeleteEmployee deletes an employee
func (s *service) DeleteEmployee(ctx context.Context, req DeleteEmployeeRequest) (*DeleteEmployeeResponse, error) {
	s.logger.Info("Deleting employee", "id", req.ID)

	// Validate employee ID
	employeeID, err := employee.NewEmployeeID(req.ID)
	if err != nil {
		s.logger.Error("Invalid employee ID", "id", req.ID, "error", err)
		return &DeleteEmployeeResponse{
			Success: false,
			Errors:  []ErrorDTO{MapErrorToDTO(err)},
		}, nil
	}

	// Check if employee exists
	_, err = s.employeeRepo.FindByID(ctx, employeeID)
	if err != nil {
		s.logger.Error("Employee not found", "id", req.ID, "error", err)
		return &DeleteEmployeeResponse{
			Success: false,
			Errors:  []ErrorDTO{MapErrorToDTO(err)},
		}, nil
	}

	// Delete from repository
	if err := s.employeeRepo.Delete(ctx, employeeID); err != nil {
		s.logger.Error("Failed to delete employee", "error", err)
		return &DeleteEmployeeResponse{
			Success: false,
			Errors:  []ErrorDTO{MapErrorToDTO(err)},
		}, nil
	}

	s.logger.Info("Successfully deleted employee", "id", req.ID)
	return &DeleteEmployeeResponse{
		Success: true,
	}, nil
}