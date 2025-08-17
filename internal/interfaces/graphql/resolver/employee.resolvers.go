package resolver

import (
	"context"
	"time"

	"github.com/captain-corgi/go-graphql-example/internal/application/employee"
	"github.com/captain-corgi/go-graphql-example/internal/interfaces/graphql/model"
)

// Employee is the resolver for the employee field.
func (r *queryResolver) Employee(ctx context.Context, id string) (*model.Employee, error) {
	// Log operation start
	r.logOperation(ctx, "Employee", map[string]interface{}{
		"id": id,
	})

	// Validate and sanitize input
	sanitizedID := sanitizeString(id)
	if err := r.validateInput(ctx, "Employee", func() error {
		return validateEmployeeID(sanitizedID)
	}); err != nil {
		return nil, err
	}

	// Call application service
	req := employee.GetEmployeeRequest{ID: sanitizedID}
	resp, err := r.employeeService.GetEmployee(ctx, req)
	if err != nil {
		return nil, r.handleGraphQLError(ctx, err, "Employee")
	}

	// No application-level errors for GetEmployeeResponse

	// Map result to GraphQL model
	result := mapEmployeeDTOToGraphQL(resp.Employee)
	r.logOperationSuccess(ctx, "Employee", result)

	return result, nil
}

// Employees is the resolver for the employees field.
func (r *queryResolver) Employees(ctx context.Context, first *int, after *string) (*model.EmployeeConnection, error) {
	// Log operation start
	r.logOperation(ctx, "Employees", map[string]interface{}{
		"first": first,
		"after": after,
	})

	// Validate pagination parameters
	if err := r.validateInput(ctx, "Employees", func() error {
		return validatePaginationParams(first, after)
	}); err != nil {
		return nil, err
	}

	// Sanitize after parameter
	var sanitizedAfter string
	if after != nil {
		sanitizedAfter = sanitizeString(*after)
	}

	// Set default first value
	firstValue := 10 // Default page size
	if first != nil {
		firstValue = *first
	}

	// Call application service
	req := employee.ListEmployeesRequest{
		Limit:  firstValue,
		Cursor: sanitizedAfter,
	}
	resp, err := r.employeeService.ListEmployees(ctx, req)
	if err != nil {
		return nil, r.handleGraphQLError(ctx, err, "Employees")
	}

	// No application-level errors for ListEmployeesResponse

	// Map result to GraphQL model
	result := mapEmployeeConnectionDTOToGraphQL(resp.Employees, resp.NextCursor)
	r.logOperationSuccess(ctx, "Employees", result)

	return result, nil
}

// EmployeesByDepartment is the resolver for the employeesByDepartment field.
func (r *queryResolver) EmployeesByDepartment(ctx context.Context, department string, first *int, after *string) (*model.EmployeeConnection, error) {
	// Log operation start
	r.logOperation(ctx, "EmployeesByDepartment", map[string]interface{}{
		"department": department,
		"first":      first,
		"after":      after,
	})

	// Validate pagination parameters
	if err := r.validateInput(ctx, "EmployeesByDepartment", func() error {
		return validatePaginationParams(first, after)
	}); err != nil {
		return nil, err
	}

	// Sanitize inputs
	sanitizedDepartment := sanitizeString(department)
	var sanitizedAfter string
	if after != nil {
		sanitizedAfter = sanitizeString(*after)
	}

	// Set default first value
	firstValue := 10 // Default page size
	if first != nil {
		firstValue = *first
	}

	// Call application service
	req := employee.ListEmployeesByDepartmentRequest{
		Department: sanitizedDepartment,
		Limit:      firstValue,
		Cursor:     sanitizedAfter,
	}
	resp, err := r.employeeService.ListEmployeesByDepartment(ctx, req)
	if err != nil {
		return nil, r.handleGraphQLError(ctx, err, "EmployeesByDepartment")
	}

	// No application-level errors for ListEmployeesByDepartmentResponse

	// Map result to GraphQL model
	result := mapEmployeeConnectionDTOToGraphQL(resp.Employees, resp.NextCursor)
	r.logOperationSuccess(ctx, "EmployeesByDepartment", result)

	return result, nil
}

// EmployeesByStatus is the resolver for the employeesByStatus field.
func (r *queryResolver) EmployeesByStatus(ctx context.Context, status string, first *int, after *string) (*model.EmployeeConnection, error) {
	// Log operation start
	r.logOperation(ctx, "EmployeesByStatus", map[string]interface{}{
		"status": status,
		"first":  first,
		"after":  after,
	})

	// Validate pagination parameters
	if err := r.validateInput(ctx, "EmployeesByStatus", func() error {
		return validatePaginationParams(first, after)
	}); err != nil {
		return nil, err
	}

	// Sanitize inputs
	sanitizedStatus := sanitizeString(status)
	var sanitizedAfter string
	if after != nil {
		sanitizedAfter = sanitizeString(*after)
	}

	// Set default first value
	firstValue := 10 // Default page size
	if first != nil {
		firstValue = *first
	}

	// Call application service
	req := employee.ListEmployeesByStatusRequest{
		Status: sanitizedStatus,
		Limit:  firstValue,
		Cursor: sanitizedAfter,
	}
	resp, err := r.employeeService.ListEmployeesByStatus(ctx, req)
	if err != nil {
		return nil, r.handleGraphQLError(ctx, err, "EmployeesByStatus")
	}

	// No application-level errors for ListEmployeesByStatusResponse

	// Map result to GraphQL model
	result := mapEmployeeConnectionDTOToGraphQL(resp.Employees, resp.NextCursor)
	r.logOperationSuccess(ctx, "EmployeesByStatus", result)

	return result, nil
}

// CreateEmployee is the resolver for the createEmployee field.
func (r *mutationResolver) CreateEmployee(ctx context.Context, input model.CreateEmployeeInput) (*model.CreateEmployeePayload, error) {
	// Log operation start
	r.logOperation(ctx, "CreateEmployee", map[string]interface{}{
		"userId":       input.UserID,
		"employeeCode": input.EmployeeCode,
		"department":   input.Department,
		"position":     input.Position,
		"hireDate":     input.HireDate,
		"salary":       input.Salary,
		"status":       input.Status,
	})

	// Validate and sanitize input
	sanitizedInput := model.CreateEmployeeInput{
		UserID:       sanitizeString(input.UserID),
		EmployeeCode: sanitizeString(input.EmployeeCode),
		Department:   sanitizeString(input.Department),
		Position:     sanitizeString(input.Position),
		HireDate:     sanitizeString(input.HireDate),
		Salary:       input.Salary,
		Status:       sanitizeString(input.Status),
	}

	if err := r.validateInput(ctx, "CreateEmployee", func() error {
		return validateCreateEmployeeInput(sanitizedInput)
	}); err != nil {
		return &model.CreateEmployeePayload{
			Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
				Message: err.Error(),
				Code:    "VALIDATION_ERROR",
			})},
		}, nil
	}

	// Parse hire date
	hireDate, err := time.Parse("2006-01-02", sanitizedInput.HireDate)
	if err != nil {
		return &model.CreateEmployeePayload{
			Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
				Message: "Invalid hire date format. Use YYYY-MM-DD",
				Code:    "INVALID_HIRE_DATE",
				Field:   "hireDate",
			})},
		}, nil
	}

	// Call application service
	req := employee.CreateEmployeeRequest{
		UserID:       sanitizedInput.UserID,
		EmployeeCode: sanitizedInput.EmployeeCode,
		Department:   sanitizedInput.Department,
		Position:     sanitizedInput.Position,
		HireDate:     hireDate,
		Salary:       sanitizedInput.Salary,
		Status:       sanitizedInput.Status,
	}
	resp, err := r.employeeService.CreateEmployee(ctx, req)
	if err != nil {
		return &model.CreateEmployeePayload{
			Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
				Message: "Failed to create employee",
				Code:    "INTERNAL_ERROR",
			})},
		}, nil
	}

	// Handle application-level errors
	if len(resp.Errors) > 0 {
		return &model.CreateEmployeePayload{
			Errors: mapEmployeeErrorDTOsToGraphQL(resp.Errors),
		}, nil
	}

	// Map successful result
	result := &model.CreateEmployeePayload{
		Employee: mapEmployeeDTOToGraphQL(resp.Employee),
	}

	r.logOperationSuccess(ctx, "CreateEmployee", result)
	return result, nil
}

// UpdateEmployee is the resolver for the updateEmployee field.
func (r *mutationResolver) UpdateEmployee(ctx context.Context, id string, input model.UpdateEmployeeInput) (*model.UpdateEmployeePayload, error) {
	// Log operation start
	r.logOperation(ctx, "UpdateEmployee", map[string]interface{}{
		"id":           id,
		"employeeCode": input.EmployeeCode,
		"department":   input.Department,
		"position":     input.Position,
		"hireDate":     input.HireDate,
		"salary":       input.Salary,
		"status":       input.Status,
	})

	// Validate and sanitize input
	sanitizedID := sanitizeString(id)
	sanitizedInput := model.UpdateEmployeeInput{
		EmployeeCode: sanitizeStringPointer(input.EmployeeCode),
		Department:   sanitizeStringPointer(input.Department),
		Position:     sanitizeStringPointer(input.Position),
		HireDate:     sanitizeStringPointer(input.HireDate),
		Salary:       input.Salary,
		Status:       sanitizeStringPointer(input.Status),
	}

	// Validate employee ID
	if err := r.validateInput(ctx, "UpdateEmployee", func() error {
		return validateEmployeeID(sanitizedID)
	}); err != nil {
		return &model.UpdateEmployeePayload{
			Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
				Message: err.Error(),
				Code:    "VALIDATION_ERROR",
			})},
		}, nil
	}

	// Validate update input
	if err := r.validateInput(ctx, "UpdateEmployee", func() error {
		return validateUpdateEmployeeInput(sanitizedInput)
	}); err != nil {
		return &model.UpdateEmployeePayload{
			Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
				Message: err.Error(),
				Code:    "VALIDATION_ERROR",
			})},
		}, nil
	}

	// Parse hire date if provided
	var hireDate *time.Time
	if sanitizedInput.HireDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *sanitizedInput.HireDate)
		if err != nil {
			return &model.UpdateEmployeePayload{
				Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
					Message: "Invalid hire date format. Use YYYY-MM-DD",
					Code:    "INVALID_HIRE_DATE",
					Field:   "hireDate",
				})},
			}, nil
		}
		hireDate = &parsedDate
	}

	// Call application service
	req := employee.UpdateEmployeeRequest{
		ID:           sanitizedID,
		EmployeeCode: sanitizedInput.EmployeeCode,
		Department:   sanitizedInput.Department,
		Position:     sanitizedInput.Position,
		HireDate:     hireDate,
		Salary:       sanitizedInput.Salary,
		Status:       sanitizedInput.Status,
	}
	resp, err := r.employeeService.UpdateEmployee(ctx, req)
	if err != nil {
		return &model.UpdateEmployeePayload{
			Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
				Message: "Failed to update employee",
				Code:    "INTERNAL_ERROR",
			})},
		}, nil
	}

	// Handle application-level errors
	if len(resp.Errors) > 0 {
		return &model.UpdateEmployeePayload{
			Errors: mapEmployeeErrorDTOsToGraphQL(resp.Errors),
		}, nil
	}

	// Map successful result
	result := &model.UpdateEmployeePayload{
		Employee: mapEmployeeDTOToGraphQL(resp.Employee),
	}

	r.logOperationSuccess(ctx, "UpdateEmployee", result)
	return result, nil
}

// DeleteEmployee is the resolver for the deleteEmployee field.
func (r *mutationResolver) DeleteEmployee(ctx context.Context, id string) (*model.DeleteEmployeePayload, error) {
	// Log operation start
	r.logOperation(ctx, "DeleteEmployee", map[string]interface{}{
		"id": id,
	})

	// Validate and sanitize input
	sanitizedID := sanitizeString(id)
	if err := r.validateInput(ctx, "DeleteEmployee", func() error {
		return validateEmployeeID(sanitizedID)
	}); err != nil {
		return &model.DeleteEmployeePayload{
			Success: false,
			Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
				Message: err.Error(),
				Code:    "VALIDATION_ERROR",
			})},
		}, nil
	}

	// Call application service
	req := employee.DeleteEmployeeRequest{ID: sanitizedID}
	resp, err := r.employeeService.DeleteEmployee(ctx, req)
	if err != nil {
		return &model.DeleteEmployeePayload{
			Success: false,
			Errors: []*model.Error{mapEmployeeErrorDTOToGraphQL(employee.ErrorDTO{
				Message: "Failed to delete employee",
				Code:    "INTERNAL_ERROR",
			})},
		}, nil
	}

	// Handle application-level errors
	if len(resp.Errors) > 0 {
		return &model.DeleteEmployeePayload{
			Success: false,
			Errors:  mapEmployeeErrorDTOsToGraphQL(resp.Errors),
		}, nil
	}

	// Map successful result
	result := &model.DeleteEmployeePayload{
		Success: resp.Success,
	}

	r.logOperationSuccess(ctx, "DeleteEmployee", result)
	return result, nil
}