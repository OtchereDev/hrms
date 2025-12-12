package handlers

import (
	"strconv"

	"github.com/OtchereDev/hrms-go/internal/api/middleware"
	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"github.com/OtchereDev/hrms-go/internal/core/services/employee"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// EmployeeHandler handles employee-related requests
type EmployeeHandler struct {
	employeeService *employee.EmployeeService
}

// NewEmployeeHandler creates a new employee handler
func NewEmployeeHandler(db *gorm.DB) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employee.NewEmployeeService(db),
	}
}

// CreateEmployee creates a new employee
// POST /api/method/hrms.hr.doctype.employee.employee.create_employee
func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req employee.CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	emp, err := h.employeeService.CreateEmployee(c.Context(), req)
	if err != nil {
		switch err {
		case employee.ErrEmployeeExists:
			return response.Conflict(c, "employee already exists")
		case employee.ErrEmployeeNumberRequired, employee.ErrEmployeeNameRequired,
			employee.ErrCompanyRequired, employee.ErrDateOfJoiningRequired:
			return response.BadRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to create employee")
		}
	}

	return response.Created(c, emp, "Employee created successfully")
}

// GetEmployee retrieves an employee by ID
// GET /api/method/hrms.hr.doctype.employee.employee.get_employee
func (h *EmployeeHandler) GetEmployee(c *fiber.Ctx) error {
	id := c.Query("id")
	employeeNumber := c.Query("employee_number")

	if id == "" && employeeNumber == "" {
		return response.BadRequest(c, "id or employee_number is required")
	}

	var emp interface{}
	var err error

	if employeeNumber != "" {
		emp, err = h.employeeService.GetEmployeeByNumber(c.Context(), employeeNumber)
	} else {
		idInt, parseErr := strconv.ParseUint(id, 10, 32)
		if parseErr != nil {
			return response.BadRequest(c, "invalid id")
		}
		emp, err = h.employeeService.GetEmployee(c.Context(), uint(idInt))
	}

	if err != nil {
		if err == employee.ErrEmployeeNotFound {
			return response.NotFound(c, "employee not found")
		}
		return response.InternalServerError(c, "failed to get employee")
	}

	return response.Success(c, emp, "success")
}

// GetEmployeeDetails retrieves employee with all details
// GET /api/method/hrms.hr.doctype.employee.employee.get_employee_details
func (h *EmployeeHandler) GetEmployeeDetails(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	details, err := h.employeeService.GetEmployeeWithDetails(c.Context(), uint(id))
	if err != nil {
		if err == employee.ErrEmployeeNotFound {
			return response.NotFound(c, "employee not found")
		}
		return response.InternalServerError(c, "failed to get employee details")
	}

	return response.Success(c, details, "success")
}

// UpdateEmployee updates an employee
// POST /api/method/hrms.hr.doctype.employee.employee.update_employee
func (h *EmployeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	var req employee.UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	emp, err := h.employeeService.UpdateEmployee(c.Context(), uint(id), req)
	if err != nil {
		if err == employee.ErrEmployeeNotFound {
			return response.NotFound(c, "employee not found")
		}
		return response.InternalServerError(c, "failed to update employee")
	}

	return response.Success(c, emp, "Employee updated successfully")
}

// DeleteEmployee deletes an employee
// POST /api/method/hrms.hr.doctype.employee.employee.delete_employee
func (h *EmployeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	if err := h.employeeService.DeleteEmployee(c.Context(), uint(id)); err != nil {
		if err == employee.ErrEmployeeNotFound {
			return response.NotFound(c, "employee not found")
		}
		return response.InternalServerError(c, "failed to delete employee")
	}

	return response.Success(c, nil, "Employee deleted successfully")
}

// ListEmployees retrieves employees with pagination and filters
// GET /api/method/hrms.api.get_employees
func (h *EmployeeHandler) ListEmployees(c *fiber.Ctx) error {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	// Get filters
	filters := repositories.EmployeeFilters{
		Company:        c.Query("company"),
		Department:     c.Query("department"),
		Branch:         c.Query("branch"),
		Designation:    c.Query("designation"),
		Status:         c.Query("status", "Active"),
		EmploymentType: c.Query("employment_type"),
		Search:         c.Query("search"),
	}

	employees, total, err := h.employeeService.ListEmployees(c.Context(), filters, page, pageSize)
	if err != nil {
		return response.InternalServerError(c, "failed to list employees")
	}

	return response.Paginated(c, employees, page, pageSize, total)
}

// GetActiveEmployees retrieves all active employees
// GET /api/method/hrms.api.get_active_employees
func (h *EmployeeHandler) GetActiveEmployees(c *fiber.Ctx) error {
	employees, err := h.employeeService.GetActiveEmployees(c.Context())
	if err != nil {
		return response.InternalServerError(c, "failed to get active employees")
	}

	return response.Success(c, employees, "success")
}

// GetEmployeesByDepartment retrieves employees in a department
// GET /api/method/hrms.api.get_employees_by_department
func (h *EmployeeHandler) GetEmployeesByDepartment(c *fiber.Ctx) error {
	department := c.Query("department")
	if department == "" {
		return response.BadRequest(c, "department is required")
	}

	employees, err := h.employeeService.GetEmployeesByDepartment(c.Context(), department)
	if err != nil {
		return response.InternalServerError(c, "failed to get employees")
	}

	return response.Success(c, employees, "success")
}

// GetCurrentEmployeeInfo returns current logged-in employee information
// GET /api/method/hrms.api.get_current_employee_info
func (h *EmployeeHandler) GetCurrentEmployeeInfo(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Unauthorized(c, "not authenticated")
	}

	// Convert user ID to string (assuming user_id is stored as string in employee table)
	userIDStr := strconv.FormatUint(uint64(userID), 10)

	emp, err := h.employeeService.GetEmployeeByUserID(c.Context(), userIDStr)
	if err != nil {
		if err == employee.ErrEmployeeNotFound {
			return response.NotFound(c, "employee record not found for current user")
		}
		return response.InternalServerError(c, "failed to get employee info")
	}

	return response.Success(c, emp, "success")
}

// UpdateEmployeeStatus updates employee status
// POST /api/method/hrms.hr.doctype.employee.employee.update_employee_status
func (h *EmployeeHandler) UpdateEmployeeStatus(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if req.Status == "" {
		return response.BadRequest(c, "status is required")
	}

	if err := h.employeeService.UpdateEmployeeStatus(c.Context(), uint(id), req.Status); err != nil {
		return response.InternalServerError(c, "failed to update employee status")
	}

	return response.Success(c, nil, "Employee status updated successfully")
}
