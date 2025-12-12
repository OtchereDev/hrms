package payroll

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/payroll"
	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"gorm.io/gorm"
)

var (
	ErrSalarySlipNotFound       = errors.New("salary slip not found")
	ErrSalaryStructureNotFound  = errors.New("salary structure not found")
	ErrSalaryComponentNotFound  = errors.New("salary component not found")
	ErrNoActiveAssignment       = errors.New("no active salary assignment for employee")
	ErrPayrollEntryNotFound     = errors.New("payroll entry not found")
	ErrLoanNotFound             = errors.New("loan not found")
	ErrEmployeeAdvanceNotFound  = errors.New("employee advance not found")
	ErrExpenseClaimNotFound     = errors.New("expense claim not found")
	ErrInsufficientLoanBalance  = errors.New("insufficient loan balance")
	ErrInvalidStatus            = errors.New("invalid status")
	ErrEmployeeRequired         = errors.New("employee is required")
)

// PayrollService handles payroll and salary business logic
type PayrollService struct {
	// Salary repositories
	componentRepo    *repositories.SalaryComponentRepository
	structureRepo    *repositories.SalaryStructureRepository
	assignmentRepo   *repositories.SalaryStructureAssignmentRepository
	salarySlipRepo   *repositories.SalarySlipRepository
	additionalRepo   *repositories.AdditionalSalaryRepository

	// Payroll repositories
	payrollEntryRepo *repositories.PayrollEntryRepository
	loanRepo         *repositories.LoanRepository
	advanceRepo      *repositories.EmployeeAdvanceRepository
	expenseRepo      *repositories.ExpenseClaimRepository
	gratuityRepo     *repositories.GratuityRepository

	db *gorm.DB
}

// NewPayrollService creates a new payroll service
func NewPayrollService(db *gorm.DB) *PayrollService {
	return &PayrollService{
		componentRepo:    repositories.NewSalaryComponentRepository(db),
		structureRepo:    repositories.NewSalaryStructureRepository(db),
		assignmentRepo:   repositories.NewSalaryStructureAssignmentRepository(db),
		salarySlipRepo:   repositories.NewSalarySlipRepository(db),
		additionalRepo:   repositories.NewAdditionalSalaryRepository(db),
		payrollEntryRepo: repositories.NewPayrollEntryRepository(db),
		loanRepo:         repositories.NewLoanRepository(db),
		advanceRepo:      repositories.NewEmployeeAdvanceRepository(db),
		expenseRepo:      repositories.NewExpenseClaimRepository(db),
		gratuityRepo:     repositories.NewGratuityRepository(db),
		db:               db,
	}
}

// ========== Salary Slip Operations ==========

// GenerateSalarySlipRequest represents a request to generate salary slip
type GenerateSalarySlipRequest struct {
	Employee       string    `json:"employee" validate:"required"`
	StartDate      time.Time `json:"start_date" validate:"required"`
	EndDate        time.Time `json:"end_date" validate:"required"`
	PostingDate    time.Time `json:"posting_date"`
	PaymentDays    float64   `json:"payment_days"`
}

// GenerateSalarySlip generates a salary slip for an employee
func (s *PayrollService) GenerateSalarySlip(ctx context.Context, req *GenerateSalarySlipRequest) (*payroll.SalarySlip, error) {
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}

	// Check if slip already exists
	existing, err := s.salarySlipRepo.GetByEmployeeAndPeriod(ctx, req.Employee, req.StartDate, req.EndDate)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing slip: %w", err)
	}
	if existing != nil {
		return existing, nil // Return existing slip
	}

	// Get active salary assignment
	assignment, err := s.assignmentRepo.GetActiveForEmployee(ctx, req.Employee)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoActiveAssignment
		}
		return nil, fmt.Errorf("failed to get salary assignment: %w", err)
	}

	// Get salary structure
	structure, err := s.structureRepo.GetByName(ctx, assignment.SalaryStructure)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSalaryStructureNotFound
		}
		return nil, fmt.Errorf("failed to get salary structure: %w", err)
	}

	// Create salary slip
	slip := &payroll.SalarySlip{
		Employee:        req.Employee,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		PostingDate:     req.PostingDate,
		SalaryStructure: structure.SalaryStructureName,
		PayrollFrequency: structure.PayrollFrequency,
		PaymentDays:     req.PaymentDays,
		Status:          "Draft",
	}

	// Calculate working days (simplified - actual implementation would be more complex)
	slip.TotalWorkingDays = req.EndDate.Sub(req.StartDate).Hours() / 24

	// Copy earnings from structure
	slip.Earnings = make([]payroll.SalaryDetail, len(structure.Earnings))
	var grossPay float64
	for i, earning := range structure.Earnings {
		slip.Earnings[i] = payroll.SalaryDetail{
			ParentType:      "SalarySlip",
			SalaryComponent: earning.SalaryComponent,
			ComponentType:   "Earning",
			Amount:          earning.Amount,
		}
		grossPay += earning.Amount
	}

	// Copy deductions from structure
	slip.Deductions = make([]payroll.SalaryDetail, len(structure.Deductions))
	var totalDeduction float64
	for i, deduction := range structure.Deductions {
		slip.Deductions[i] = payroll.SalaryDetail{
			ParentType:      "SalarySlip",
			SalaryComponent: deduction.SalaryComponent,
			ComponentType:   "Deduction",
			Amount:          deduction.Amount,
		}
		totalDeduction += deduction.Amount
	}

	// Get additional salaries for this period
	additionalSalaries, err := s.additionalRepo.GetForEmployeeAndDate(ctx, req.Employee, req.PostingDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get additional salaries: %w", err)
	}

	// Apply additional salaries
	for _, additional := range additionalSalaries {
		if additional.Type == "Earning" {
			grossPay += additional.Amount
		} else if additional.Type == "Deduction" {
			totalDeduction += additional.Amount
		}
	}

	slip.GrossPay = grossPay
	slip.TotalDeduction = totalDeduction
	slip.NetPay = grossPay - totalDeduction
	slip.RoundedTotal = slip.NetPay

	if err := s.salarySlipRepo.Create(ctx, slip); err != nil {
		return nil, fmt.Errorf("failed to create salary slip: %w", err)
	}

	return slip, nil
}

// GetSalarySlipByID retrieves a salary slip by ID
func (s *PayrollService) GetSalarySlipByID(ctx context.Context, id uint) (*payroll.SalarySlip, error) {
	slip, err := s.salarySlipRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSalarySlipNotFound
		}
		return nil, fmt.Errorf("failed to get salary slip: %w", err)
	}
	return slip, nil
}

// SubmitSalarySlip submits a salary slip for approval
func (s *PayrollService) SubmitSalarySlip(ctx context.Context, id uint) (*payroll.SalarySlip, error) {
	slip, err := s.salarySlipRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSalarySlipNotFound
		}
		return nil, fmt.Errorf("failed to get salary slip: %w", err)
	}

	if slip.Status != "Draft" {
		return nil, fmt.Errorf("cannot submit salary slip with status: %s", slip.Status)
	}

	slip.Status = "Submitted"

	if err := s.salarySlipRepo.Update(ctx, slip); err != nil {
		return nil, fmt.Errorf("failed to submit salary slip: %w", err)
	}

	return slip, nil
}

// ListSalarySlips retrieves salary slips with filters
func (s *PayrollService) ListSalarySlips(ctx context.Context, filters repositories.SalarySlipFilters, page, pageSize int) ([]payroll.SalarySlip, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.salarySlipRepo.List(ctx, filters, page, pageSize)
}

// ========== Salary Structure Operations ==========

// CreateSalaryStructureRequest represents a request to create salary structure
type CreateSalaryStructureRequest struct {
	Name             string  `json:"name" validate:"required"`
	Company          string  `json:"company" validate:"required"`
	PayrollFrequency string  `json:"payroll_frequency" validate:"required"`
	Currency         string  `json:"currency"`
}

// AssignSalaryStructureRequest represents a request to assign salary structure
type AssignSalaryStructureRequest struct {
	Employee        string     `json:"employee" validate:"required"`
	SalaryStructure string     `json:"salary_structure" validate:"required"`
	FromDate        time.Time  `json:"from_date" validate:"required"`
	ToDate          *time.Time `json:"to_date"`
	Base            float64    `json:"base" validate:"required"`
	Company         string     `json:"company" validate:"required"`
}

// AssignSalaryStructure assigns a salary structure to an employee
func (s *PayrollService) AssignSalaryStructure(ctx context.Context, req *AssignSalaryStructureRequest) (*payroll.SalaryStructureAssignment, error) {
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}

	// Verify salary structure exists
	_, err := s.structureRepo.GetByName(ctx, req.SalaryStructure)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSalaryStructureNotFound
		}
		return nil, fmt.Errorf("failed to get salary structure: %w", err)
	}

	assignment := &payroll.SalaryStructureAssignment{
		Employee:        req.Employee,
		SalaryStructure: req.SalaryStructure,
		FromDate:        req.FromDate,
		ToDate:          req.ToDate,
		Base:            req.Base,
		Company:         req.Company,
	}

	if err := s.assignmentRepo.Create(ctx, assignment); err != nil {
		return nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	return assignment, nil
}

// GetActiveLeaveTypes retrieves the active salary structure for an employee
func (s *PayrollService) GetActiveSalaryAssignment(ctx context.Context, employee string) (*payroll.SalaryStructureAssignment, error) {
	if employee == "" {
		return nil, ErrEmployeeRequired
	}

	assignment, err := s.assignmentRepo.GetActiveForEmployee(ctx, employee)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoActiveAssignment
		}
		return nil, fmt.Errorf("failed to get active assignment: %w", err)
	}

	return assignment, nil
}

// ========== Loan Operations ==========

// CreateLoanRequest represents a request to create a loan
type CreateLoanRequest struct {
	Applicant              string    `json:"applicant" validate:"required"`
	Company                string    `json:"company" validate:"required"`
	LoanType               string    `json:"loan_type" validate:"required"`
	LoanAmount             float64   `json:"loan_amount" validate:"required"`
	RateOfInterest         float64   `json:"rate_of_interest"`
	RepaymentPeriods       int       `json:"repayment_periods" validate:"required"`
	RepaymentStartDate     time.Time `json:"repayment_start_date" validate:"required"`
	PostingDate            time.Time `json:"posting_date" validate:"required"`
}

// CreateLoan creates a new loan application
func (s *PayrollService) CreateLoan(ctx context.Context, req *CreateLoanRequest) (*payroll.Loan, error) {
	if req.Applicant == "" {
		return nil, ErrEmployeeRequired
	}

	// Calculate monthly repayment
	monthlyRepayment := req.LoanAmount / float64(req.RepaymentPeriods)
	totalInterest := (req.LoanAmount * req.RateOfInterest * float64(req.RepaymentPeriods)) / 1200
	totalPayment := req.LoanAmount + totalInterest

	loan := &payroll.Loan{
		Applicant:              req.Applicant,
		ApplicantType:          "Employee",
		Company:                req.Company,
		LoanType:               req.LoanType,
		LoanAmount:             req.LoanAmount,
		RateOfInterest:         req.RateOfInterest,
		IsTermLoan:             true,
		RepaymentMethod:        "Repay Over Number of Periods",
		MonthlyRepaymentAmount: monthlyRepayment,
		RepaymentPeriods:       req.RepaymentPeriods,
		RepaymentStartDate:     req.RepaymentStartDate,
		PostingDate:            req.PostingDate,
		TotalPayment:           totalPayment,
		TotalInterestPayable:   totalInterest,
		Status:                 "Draft",
	}

	if err := s.loanRepo.Create(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	return loan, nil
}

// ApproveLoan approves a loan
func (s *PayrollService) ApproveLoan(ctx context.Context, id uint) (*payroll.Loan, error) {
	loan, err := s.loanRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLoanNotFound
		}
		return nil, fmt.Errorf("failed to get loan: %w", err)
	}

	if loan.Status != "Draft" {
		return nil, fmt.Errorf("cannot approve loan with status: %s", loan.Status)
	}

	loan.Status = "Sanctioned"

	if err := s.loanRepo.Update(ctx, loan); err != nil {
		return nil, fmt.Errorf("failed to approve loan: %w", err)
	}

	return loan, nil
}

// ListLoans retrieves loans with filters
func (s *PayrollService) ListLoans(ctx context.Context, filters repositories.LoanFilters, page, pageSize int) ([]payroll.Loan, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.loanRepo.List(ctx, filters, page, pageSize)
}

// ========== Employee Advance Operations ==========

// CreateAdvanceRequest represents a request to create an advance
type CreateAdvanceRequest struct {
	Employee      string    `json:"employee" validate:"required"`
	Company       string    `json:"company" validate:"required"`
	AdvanceAmount float64   `json:"advance_amount" validate:"required"`
	Purpose       string    `json:"purpose" validate:"required"`
	PostingDate   time.Time `json:"posting_date" validate:"required"`
}

// CreateEmployeeAdvance creates a new employee advance
func (s *PayrollService) CreateEmployeeAdvance(ctx context.Context, req *CreateAdvanceRequest) (*payroll.EmployeeAdvance, error) {
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}

	advance := &payroll.EmployeeAdvance{
		Employee:      req.Employee,
		Company:       req.Company,
		AdvanceAmount: req.AdvanceAmount,
		Purpose:       req.Purpose,
		PostingDate:   req.PostingDate,
		Status:        "Draft",
	}

	if err := s.advanceRepo.Create(ctx, advance); err != nil {
		return nil, fmt.Errorf("failed to create advance: %w", err)
	}

	return advance, nil
}

// ApproveAdvance approves an employee advance
func (s *PayrollService) ApproveAdvance(ctx context.Context, id uint) (*payroll.EmployeeAdvance, error) {
	advance, err := s.advanceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmployeeAdvanceNotFound
		}
		return nil, fmt.Errorf("failed to get advance: %w", err)
	}

	if advance.Status != "Draft" {
		return nil, fmt.Errorf("cannot approve advance with status: %s", advance.Status)
	}

	advance.Status = "Paid"
	advance.PaidAmount = advance.AdvanceAmount

	if err := s.advanceRepo.Update(ctx, advance); err != nil {
		return nil, fmt.Errorf("failed to approve advance: %w", err)
	}

	return advance, nil
}

// ========== Expense Claim Operations ==========

// CreateExpenseClaimRequest represents a request to create expense claim
type CreateExpenseClaimRequest struct {
	Employee        string    `json:"employee" validate:"required"`
	Company         string    `json:"company" validate:"required"`
	PostingDate     time.Time `json:"posting_date" validate:"required"`
	ExpenseApprover string    `json:"expense_approver"`
}

// CreateExpenseClaim creates a new expense claim
func (s *PayrollService) CreateExpenseClaim(ctx context.Context, req *CreateExpenseClaimRequest) (*payroll.ExpenseClaim, error) {
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}

	claim := &payroll.ExpenseClaim{
		Employee:        req.Employee,
		Company:         req.Company,
		PostingDate:     req.PostingDate,
		ExpenseApprover: req.ExpenseApprover,
		Status:          "Draft",
	}

	if err := s.expenseRepo.Create(ctx, claim); err != nil {
		return nil, fmt.Errorf("failed to create expense claim: %w", err)
	}

	return claim, nil
}

// ApproveExpenseClaim approves an expense claim
func (s *PayrollService) ApproveExpenseClaim(ctx context.Context, id uint) (*payroll.ExpenseClaim, error) {
	claim, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrExpenseClaimNotFound
		}
		return nil, fmt.Errorf("failed to get expense claim: %w", err)
	}

	if claim.Status != "Submitted" && claim.Status != "Draft" {
		return nil, fmt.Errorf("cannot approve expense claim with status: %s", claim.Status)
	}

	claim.Status = "Approved"

	if err := s.expenseRepo.Update(ctx, claim); err != nil {
		return nil, fmt.Errorf("failed to approve expense claim: %w", err)
	}

	return claim, nil
}

// ListExpenseClaims retrieves expense claims with filters
func (s *PayrollService) ListExpenseClaims(ctx context.Context, filters repositories.ExpenseClaimFilters, page, pageSize int) ([]payroll.ExpenseClaim, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.expenseRepo.List(ctx, filters, page, pageSize)
}
