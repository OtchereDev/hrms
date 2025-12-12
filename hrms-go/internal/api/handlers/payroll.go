package handlers

import (
	frappeResponse "github.com/OtchereDev/hrms-go/internal/core/frappe"
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
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	slip, err := h.payrollService.GenerateSalarySlip(c.Context(), &req)
	if err != nil {
		switch err {
		case payroll.ErrEmployeeRequired:
			return frappeResponse.SendBadRequest(c, "employee is required")
		case payroll.ErrNoActiveAssignment:
			return frappeResponse.SendBadRequest(c, "no active salary assignment for employee")
		case payroll.ErrSalaryStructureNotFound:
			return frappeResponse.SendNotFound(c, "salary structure not found")
		default:
			return frappeResponse.SendInternalError(c, "failed to generate salary slip")
		}
	}

	return response.Created(c, slip, "Salary slip generated successfully")
}

// GetSalarySlip retrieves a salary slip by ID
// GET /api/method/hrms.payroll.doctype.salary_slip.salary_slip.get_salary_slip
func (h *PayrollHandler) GetSalarySlip(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	slip, err := h.payrollService.GetSalarySlipByID(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrSalarySlipNotFound {
			return frappeResponse.SendNotFound(c, "salary slip not found")
		}
		return frappeResponse.SendInternalError(c, "failed to get salary slip")
	}

	return frappeResponse.SendSuccess(c, slip)
}

// SubmitSalarySlip submits a salary slip
// POST /api/method/hrms.payroll.doctype.salary_slip.salary_slip.submit_salary_slip
func (h *PayrollHandler) SubmitSalarySlip(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	slip, err := h.payrollService.SubmitSalarySlip(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrSalarySlipNotFound {
			return frappeResponse.SendNotFound(c, "salary slip not found")
		}
		return frappeResponse.SendInternalError(c, "failed to submit salary slip")
	}

	return frappeResponse.SendSuccess(c, slip)
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
		return frappeResponse.SendInternalError(c, "failed to list salary slips")
	}

	return response.Paginated(c, slips, page, pageSize, total)
}

// ========== Salary Structure Operations ==========

// AssignSalaryStructure assigns a salary structure to an employee
// POST /api/method/hrms.payroll.doctype.salary_structure_assignment.salary_structure_assignment.assign_structure
func (h *PayrollHandler) AssignSalaryStructure(c *fiber.Ctx) error {
	var req payroll.AssignSalaryStructureRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	assignment, err := h.payrollService.AssignSalaryStructure(c.Context(), &req)
	if err != nil {
		switch err {
		case payroll.ErrEmployeeRequired:
			return frappeResponse.SendBadRequest(c, "employee is required")
		case payroll.ErrSalaryStructureNotFound:
			return frappeResponse.SendNotFound(c, "salary structure not found")
		default:
			return frappeResponse.SendInternalError(c, "failed to assign salary structure")
		}
	}

	return response.Created(c, assignment, "Salary structure assigned successfully")
}

// GetActiveSalaryAssignment retrieves the active salary assignment for an employee
// GET /api/method/hrms.payroll.doctype.salary_structure_assignment.salary_structure_assignment.get_active_assignment
func (h *PayrollHandler) GetActiveSalaryAssignment(c *fiber.Ctx) error {
	employee := c.Query("employee")
	if employee == "" {
		return frappeResponse.SendBadRequest(c, "employee is required")
	}

	assignment, err := h.payrollService.GetActiveSalaryAssignment(c.Context(), employee)
	if err != nil {
		if err == payroll.ErrNoActiveAssignment {
			return frappeResponse.SendNotFound(c, "no active salary assignment found")
		}
		return frappeResponse.SendInternalError(c, "failed to get active assignment")
	}

	return frappeResponse.SendSuccess(c, assignment)
}

// ========== Loan Operations ==========

// CreateLoan creates a new loan application
// POST /api/method/hrms.payroll.doctype.loan.loan.create_loan
func (h *PayrollHandler) CreateLoan(c *fiber.Ctx) error {
	var req payroll.CreateLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	loan, err := h.payrollService.CreateLoan(c.Context(), &req)
	if err != nil {
		if err == payroll.ErrEmployeeRequired {
			return frappeResponse.SendBadRequest(c, "applicant is required")
		}
		return frappeResponse.SendInternalError(c, "failed to create loan")
	}

	return response.Created(c, loan, "Loan created successfully")
}

// ApproveLoan approves a loan application
// POST /api/method/hrms.payroll.doctype.loan.loan.approve_loan
func (h *PayrollHandler) ApproveLoan(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	loan, err := h.payrollService.ApproveLoan(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrLoanNotFound {
			return frappeResponse.SendNotFound(c, "loan not found")
		}
		return frappeResponse.SendInternalError(c, "failed to approve loan")
	}

	return frappeResponse.SendSuccess(c, loan)
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
		return frappeResponse.SendInternalError(c, "failed to list loans")
	}

	return response.Paginated(c, loans, page, pageSize, total)
}

// ========== Employee Advance Operations ==========

// CreateEmployeeAdvance creates a new employee advance
// POST /api/method/hrms.payroll.doctype.employee_advance.employee_advance.create_advance
func (h *PayrollHandler) CreateEmployeeAdvance(c *fiber.Ctx) error {
	var req payroll.CreateAdvanceRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	advance, err := h.payrollService.CreateEmployeeAdvance(c.Context(), &req)
	if err != nil {
		if err == payroll.ErrEmployeeRequired {
			return frappeResponse.SendBadRequest(c, "employee is required")
		}
		return frappeResponse.SendInternalError(c, "failed to create advance")
	}

	return response.Created(c, advance, "Employee advance created successfully")
}

// ApproveAdvance approves an employee advance
// POST /api/method/hrms.payroll.doctype.employee_advance.employee_advance.approve_advance
func (h *PayrollHandler) ApproveAdvance(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	advance, err := h.payrollService.ApproveAdvance(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrEmployeeAdvanceNotFound {
			return frappeResponse.SendNotFound(c, "employee advance not found")
		}
		return frappeResponse.SendInternalError(c, "failed to approve advance")
	}

	return frappeResponse.SendSuccess(c, advance)
}

// ========== Expense Claim Operations ==========

// CreateExpenseClaim creates a new expense claim
// POST /api/method/hrms.payroll.doctype.expense_claim.expense_claim.create_expense_claim
func (h *PayrollHandler) CreateExpenseClaim(c *fiber.Ctx) error {
	var req payroll.CreateExpenseClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	claim, err := h.payrollService.CreateExpenseClaim(c.Context(), &req)
	if err != nil {
		if err == payroll.ErrEmployeeRequired {
			return frappeResponse.SendBadRequest(c, "employee is required")
		}
		return frappeResponse.SendInternalError(c, "failed to create expense claim")
	}

	return response.Created(c, claim, "Expense claim created successfully")
}

// ApproveExpenseClaim approves an expense claim
// POST /api/method/hrms.payroll.doctype.expense_claim.expense_claim.approve_claim
func (h *PayrollHandler) ApproveExpenseClaim(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	claim, err := h.payrollService.ApproveExpenseClaim(c.Context(), uint(id))
	if err != nil {
		if err == payroll.ErrExpenseClaimNotFound {
			return frappeResponse.SendNotFound(c, "expense claim not found")
		}
		return frappeResponse.SendInternalError(c, "failed to approve expense claim")
	}

	return frappeResponse.SendSuccess(c, claim)
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
		return frappeResponse.SendInternalError(c, "failed to list expense claims")
	}

	return response.Paginated(c, claims, page, pageSize, total)
}

// GetSalarySlipDetails returns detailed breakdown of a salary slip
// GET /api/method/hrms.payroll.doctype.salary_slip.salary_slip.get_salary_slip_details
func (h *PayrollHandler) GetSalarySlipDetails(c *fiber.Ctx) error {
	employee := c.Query("employee")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if employee == "" || startDateStr == "" || endDateStr == "" {
		return frappeResponse.SendBadRequest(c, "employee, start_date, and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid start_date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid end_date format, expected YYYY-MM-DD")
	}

	details, err := h.payrollService.GetSalarySlipDetails(c.Context(), employee, startDate, endDate)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get salary slip details")
	}

	return frappeResponse.SendSuccess(c, details)
}

// CalculateLoanAmounts calculates loan repayment amounts
// GET /api/method/hrms.payroll.doctype.loan.loan.calculate_amounts
func (h *PayrollHandler) CalculateLoanAmounts(c *fiber.Ctx) error {
	loanType := c.Query("loan_type")
	loanAmountStr := c.Query("loan_amount")
	rateOfInterestStr := c.Query("rate_of_interest")
	repaymentPeriodsStr := c.Query("repayment_periods")

	if loanAmountStr == "" || repaymentPeriodsStr == "" {
		return frappeResponse.SendBadRequest(c, "loan_amount and repayment_periods are required")
	}

	loanAmount, err := strconv.ParseFloat(loanAmountStr, 64)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid loan_amount")
	}

	rateOfInterest := 0.0
	if rateOfInterestStr != "" {
		rateOfInterest, err = strconv.ParseFloat(rateOfInterestStr, 64)
		if err != nil {
			return frappeResponse.SendBadRequest(c, "invalid rate_of_interest")
		}
	}

	repaymentPeriods, err := strconv.Atoi(repaymentPeriodsStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid repayment_periods")
	}

	amounts, err := h.payrollService.CalculateLoanAmounts(c.Context(), loanType, loanAmount, rateOfInterest, repaymentPeriods)
	if err != nil {
		return frappeResponse.SendBadRequest(c, err.Error())
	}

	return frappeResponse.SendSuccess(c, amounts)
}

// CalculateNetPay calculates net pay from a salary slip
// GET /api/method/hrms.payroll.doctype.salary_slip.salary_slip.calculate_net_pay
func (h *PayrollHandler) CalculateNetPay(c *fiber.Ctx) error {
	slipIDStr := c.Query("slip_id")

	if slipIDStr == "" {
		return frappeResponse.SendBadRequest(c, "slip_id is required")
	}

	slipID, err := strconv.ParseUint(slipIDStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid slip_id")
	}

	netPay, err := h.payrollService.CalculateNetPay(c.Context(), uint(slipID))
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to calculate net pay")
	}

	return frappeResponse.SendSuccess(c, map[string]interface{}{"net_pay": netPay})
}

// GetPayrollSummary returns payroll summary for a period
// GET /api/method/hrms.payroll.doctype.payroll_entry.payroll_entry.get_payroll_summary
func (h *PayrollHandler) GetPayrollSummary(c *fiber.Ctx) error {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	company := c.Query("company")

	if startDateStr == "" || endDateStr == "" {
		return frappeResponse.SendBadRequest(c, "start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid start_date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid end_date format, expected YYYY-MM-DD")
	}

	summary, err := h.payrollService.GetPayrollSummary(c.Context(), startDate, endDate, company)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get payroll summary")
	}

	return frappeResponse.SendSuccess(c, summary)
}
