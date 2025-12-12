package handlers

import (
	"strconv"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"github.com/OtchereDev/hrms-go/internal/core/services/payroll"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PayrollHandler handles payroll and salary-related requests
type PayrollHandler struct {
	payrollService *payroll.PayrollService
}

// NewPayrollHandler creates a new payroll handler
func NewPayrollHandler(db *gorm.DB) *PayrollHandler {
	return &PayrollHandler{
		payrollService: payroll.NewPayrollService(db),
	}
}

// ========== Salary Slip Operations ==========

// GenerateSalarySlip generates a salary slip for an employee
// POST /api/method/hrms.payroll.doctype.salary_slip.salary_slip.generate_salary_slip
func (h *PayrollHandler) GenerateSalarySlip(c *fiber.Ctx) error {
	var req payroll.GenerateSalarySlipRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	slip, err := h.payrollService.GenerateSalarySlip(c.Context(), &req)
	if err != nil {
		switch err {
		case payroll.ErrEmployeeRequired:
			return response.BadRequest(c, "employee is required")
		case payroll.ErrNoActiveAssignment:
			return response.BadRequest(c, "no active salary assignment for employee")
		case payroll.ErrSalaryStructureNotFound:
			return response.NotFound(c, "salary structure not found")
		default:
			return response.InternalServerError(c, "failed to generate salary slip")
		}
	}

	return response.Created(c, slip, "Salary slip generated successfully")
}

// GetSalarySlip retrieves a salary slip by ID
// GET /api/method/hrms.payroll.doctype.salary_slip.salary_slip.get_salary_slip
func (h *PayrollHandler) GetSalarySlip(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	slip, err := h.payrollService.GetSalarySlipByID(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrSalarySlipNotFound {
			return response.NotFound(c, "salary slip not found")
		}
		return response.InternalServerError(c, "failed to get salary slip")
	}

	return response.Success(c, slip, "success")
}

// SubmitSalarySlip submits a salary slip
// POST /api/method/hrms.payroll.doctype.salary_slip.salary_slip.submit_salary_slip
func (h *PayrollHandler) SubmitSalarySlip(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	slip, err := h.payrollService.SubmitSalarySlip(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrSalarySlipNotFound {
			return response.NotFound(c, "salary slip not found")
		}
		return response.InternalServerError(c, "failed to submit salary slip")
	}

	return response.Success(c, slip, "Salary slip submitted successfully")
}

// ListSalarySlips retrieves salary slips with filters
// GET /api/method/hrms.payroll.doctype.salary_slip.salary_slip.list_salary_slips
func (h *PayrollHandler) ListSalarySlips(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	filters := repositories.SalarySlipFilters{
		Employee:   c.Query("employee"),
		Company:    c.Query("company"),
		Department: c.Query("department"),
		Status:     c.Query("status"),
	}

	// Parse dates if provided
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = endDate
		}
	}

	slips, total, err := h.payrollService.ListSalarySlips(c.Context(), filters, page, pageSize)
	if err != nil {
		return response.InternalServerError(c, "failed to list salary slips")
	}

	return response.SuccessWithPagination(c, slips, total, page, pageSize, "success")
}

// ========== Salary Structure Operations ==========

// AssignSalaryStructure assigns a salary structure to an employee
// POST /api/method/hrms.payroll.doctype.salary_structure_assignment.salary_structure_assignment.assign_structure
func (h *PayrollHandler) AssignSalaryStructure(c *fiber.Ctx) error {
	var req payroll.AssignSalaryStructureRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	assignment, err := h.payrollService.AssignSalaryStructure(c.Context(), &req)
	if err != nil {
		switch err {
		case payroll.ErrEmployeeRequired:
			return response.BadRequest(c, "employee is required")
		case payroll.ErrSalaryStructureNotFound:
			return response.NotFound(c, "salary structure not found")
		default:
			return response.InternalServerError(c, "failed to assign salary structure")
		}
	}

	return response.Created(c, assignment, "Salary structure assigned successfully")
}

// GetActiveSalaryAssignment retrieves the active salary assignment for an employee
// GET /api/method/hrms.payroll.doctype.salary_structure_assignment.salary_structure_assignment.get_active_assignment
func (h *PayrollHandler) GetActiveSalaryAssignment(c *fiber.Ctx) error {
	employee := c.Query("employee")
	if employee == "" {
		return response.BadRequest(c, "employee is required")
	}

	assignment, err := h.payrollService.GetActiveSalaryAssignment(c.Context(), employee)
	if err != nil {
		if err == payroll.ErrNoActiveAssignment {
			return response.NotFound(c, "no active salary assignment found")
		}
		return response.InternalServerError(c, "failed to get active assignment")
	}

	return response.Success(c, assignment, "success")
}

// ========== Loan Operations ==========

// CreateLoan creates a new loan application
// POST /api/method/hrms.payroll.doctype.loan.loan.create_loan
func (h *PayrollHandler) CreateLoan(c *fiber.Ctx) error {
	var req payroll.CreateLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	loan, err := h.payrollService.CreateLoan(c.Context(), &req)
	if err != nil {
		if err == payroll.ErrEmployeeRequired {
			return response.BadRequest(c, "applicant is required")
		}
		return response.InternalServerError(c, "failed to create loan")
	}

	return response.Created(c, loan, "Loan created successfully")
}

// ApproveLoan approves a loan application
// POST /api/method/hrms.payroll.doctype.loan.loan.approve_loan
func (h *PayrollHandler) ApproveLoan(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	loan, err := h.payrollService.ApproveLoan(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrLoanNotFound {
			return response.NotFound(c, "loan not found")
		}
		return response.InternalServerError(c, "failed to approve loan")
	}

	return response.Success(c, loan, "Loan approved successfully")
}

// ListLoans retrieves loans with filters
// GET /api/method/hrms.payroll.doctype.loan.loan.list_loans
func (h *PayrollHandler) ListLoans(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	filters := repositories.LoanFilters{
		Applicant: c.Query("applicant"),
		Company:   c.Query("company"),
		LoanType:  c.Query("loan_type"),
		Status:    c.Query("status"),
	}

	loans, total, err := h.payrollService.ListLoans(c.Context(), filters, page, pageSize)
	if err != nil {
		return response.InternalServerError(c, "failed to list loans")
	}

	return response.SuccessWithPagination(c, loans, total, page, pageSize, "success")
}

// ========== Employee Advance Operations ==========

// CreateEmployeeAdvance creates a new employee advance
// POST /api/method/hrms.payroll.doctype.employee_advance.employee_advance.create_advance
func (h *PayrollHandler) CreateEmployeeAdvance(c *fiber.Ctx) error {
	var req payroll.CreateAdvanceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	advance, err := h.payrollService.CreateEmployeeAdvance(c.Context(), &req)
	if err != nil {
		if err == payroll.ErrEmployeeRequired {
			return response.BadRequest(c, "employee is required")
		}
		return response.InternalServerError(c, "failed to create advance")
	}

	return response.Created(c, advance, "Employee advance created successfully")
}

// ApproveAdvance approves an employee advance
// POST /api/method/hrms.payroll.doctype.employee_advance.employee_advance.approve_advance
func (h *PayrollHandler) ApproveAdvance(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	advance, err := h.payrollService.ApproveAdvance(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrEmployeeAdvanceNotFound {
			return response.NotFound(c, "employee advance not found")
		}
		return response.InternalServerError(c, "failed to approve advance")
	}

	return response.Success(c, advance, "Employee advance approved successfully")
}

// ========== Expense Claim Operations ==========

// CreateExpenseClaim creates a new expense claim
// POST /api/method/hrms.payroll.doctype.expense_claim.expense_claim.create_expense_claim
func (h *PayrollHandler) CreateExpenseClaim(c *fiber.Ctx) error {
	var req payroll.CreateExpenseClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	claim, err := h.payrollService.CreateExpenseClaim(c.Context(), &req)
	if err != nil {
		if err == payroll.ErrEmployeeRequired {
			return response.BadRequest(c, "employee is required")
		}
		return response.InternalServerError(c, "failed to create expense claim")
	}

	return response.Created(c, claim, "Expense claim created successfully")
}

// ApproveExpenseClaim approves an expense claim
// POST /api/method/hrms.payroll.doctype.expense_claim.expense_claim.approve_claim
func (h *PayrollHandler) ApproveExpenseClaim(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	claim, err := h.payrollService.ApproveExpenseClaim(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrExpenseClaimNotFound {
			return response.NotFound(c, "expense claim not found")
		}
		return response.InternalServerError(c, "failed to approve expense claim")
	}

	return response.Success(c, claim, "Expense claim approved successfully")
}

// ListExpenseClaims retrieves expense claims with filters
// GET /api/method/hrms.payroll.doctype.expense_claim.expense_claim.list_expense_claims
func (h *PayrollHandler) ListExpenseClaims(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	filters := repositories.ExpenseClaimFilters{
		Employee:        c.Query("employee"),
		Company:         c.Query("company"),
		Status:          c.Query("status"),
		ExpenseApprover: c.Query("expense_approver"),
	}

	claims, total, err := h.payrollService.ListExpenseClaims(c.Context(), filters, page, pageSize)
	if err != nil {
		return response.InternalServerError(c, "failed to list expense claims")
	}

	return response.SuccessWithPagination(c, claims, total, page, pageSize, "success")
}
