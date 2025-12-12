package repositories

import (
	"context"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/payroll"
	"gorm.io/gorm"
)

// PayrollEntryRepository handles payroll entry data access
type PayrollEntryRepository struct {
	db *gorm.DB
}

// NewPayrollEntryRepository creates a new payroll entry repository
func NewPayrollEntryRepository(db *gorm.DB) *PayrollEntryRepository {
	return &PayrollEntryRepository{db: db}
}

// Create creates a new payroll entry
func (r *PayrollEntryRepository) Create(ctx context.Context, entry *payroll.PayrollEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// GetByID retrieves a payroll entry by ID
func (r *PayrollEntryRepository) GetByID(ctx context.Context, id uint) (*payroll.PayrollEntry, error) {
	var entry payroll.PayrollEntry
	err := r.db.WithContext(ctx).
		Preload("Employees").
		Preload("Deductions").
		First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// Update updates a payroll entry
func (r *PayrollEntryRepository) Update(ctx context.Context, entry *payroll.PayrollEntry) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

// List retrieves payroll entries with filters
func (r *PayrollEntryRepository) List(ctx context.Context, filters PayrollEntryFilters, page, pageSize int) ([]payroll.PayrollEntry, int64, error) {
	var entries []payroll.PayrollEntry
	var total int64

	query := r.db.WithContext(ctx).Model(&payroll.PayrollEntry{})

	// Apply filters
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Department != "" {
		query = query.Where("department = ?", filters.Department)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if !filters.StartDate.IsZero() {
		query = query.Where("start_date >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("end_date <= ?", filters.EndDate)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("start_date DESC").Find(&entries).Error

	return entries, total, err
}

// PayrollEntryFilters represents filters for payroll entry queries
type PayrollEntryFilters struct {
	Company    string
	Department string
	Status     string
	StartDate  time.Time
	EndDate    time.Time
}

// LoanRepository handles loan data access
type LoanRepository struct {
	db *gorm.DB
}

// NewLoanRepository creates a new loan repository
func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

// Create creates a new loan
func (r *LoanRepository) Create(ctx context.Context, loan *payroll.Loan) error {
	return r.db.WithContext(ctx).Create(loan).Error
}

// GetByID retrieves a loan by ID
func (r *LoanRepository) GetByID(ctx context.Context, id uint) (*payroll.Loan, error) {
	var loan payroll.Loan
	err := r.db.WithContext(ctx).
		Preload("RepaymentSchedule").
		First(&loan, id).Error
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

// Update updates a loan
func (r *LoanRepository) Update(ctx context.Context, loan *payroll.Loan) error {
	return r.db.WithContext(ctx).Save(loan).Error
}

// List retrieves loans with filters
func (r *LoanRepository) List(ctx context.Context, filters LoanFilters, page, pageSize int) ([]payroll.Loan, int64, error) {
	var loans []payroll.Loan
	var total int64

	query := r.db.WithContext(ctx).Model(&payroll.Loan{})

	// Apply filters
	if filters.Applicant != "" {
		query = query.Where("applicant = ?", filters.Applicant)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.LoanType != "" {
		query = query.Where("loan_type = ?", filters.LoanType)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("posting_date DESC").Find(&loans).Error

	return loans, total, err
}

// GetActiveLoansForEmployee retrieves active loans for an employee
func (r *LoanRepository) GetActiveLoansForEmployee(ctx context.Context, employee string) ([]payroll.Loan, error) {
	var loans []payroll.Loan
	err := r.db.WithContext(ctx).
		Where("applicant = ? AND status IN ?", employee, []string{"Sanctioned", "Disbursed"}).
		Order("posting_date DESC").
		Find(&loans).Error
	return loans, err
}

// LoanFilters represents filters for loan queries
type LoanFilters struct {
	Applicant string
	Company   string
	LoanType  string
	Status    string
}

// EmployeeAdvanceRepository handles employee advance data access
type EmployeeAdvanceRepository struct {
	db *gorm.DB
}

// NewEmployeeAdvanceRepository creates a new advance repository
func NewEmployeeAdvanceRepository(db *gorm.DB) *EmployeeAdvanceRepository {
	return &EmployeeAdvanceRepository{db: db}
}

// Create creates a new employee advance
func (r *EmployeeAdvanceRepository) Create(ctx context.Context, advance *payroll.EmployeeAdvance) error {
	return r.db.WithContext(ctx).Create(advance).Error
}

// GetByID retrieves an advance by ID
func (r *EmployeeAdvanceRepository) GetByID(ctx context.Context, id uint) (*payroll.EmployeeAdvance, error) {
	var advance payroll.EmployeeAdvance
	err := r.db.WithContext(ctx).First(&advance, id).Error
	if err != nil {
		return nil, err
	}
	return &advance, nil
}

// Update updates an advance
func (r *EmployeeAdvanceRepository) Update(ctx context.Context, advance *payroll.EmployeeAdvance) error {
	return r.db.WithContext(ctx).Save(advance).Error
}

// List retrieves advances with filters
func (r *EmployeeAdvanceRepository) List(ctx context.Context, filters EmployeeAdvanceFilters, page, pageSize int) ([]payroll.EmployeeAdvance, int64, error) {
	var advances []payroll.EmployeeAdvance
	var total int64

	query := r.db.WithContext(ctx).Model(&payroll.EmployeeAdvance{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("posting_date DESC").Find(&advances).Error

	return advances, total, err
}

// GetUnpaidAdvancesForEmployee retrieves unpaid advances for an employee
func (r *EmployeeAdvanceRepository) GetUnpaidAdvancesForEmployee(ctx context.Context, employee string) ([]payroll.EmployeeAdvance, error) {
	var advances []payroll.EmployeeAdvance
	err := r.db.WithContext(ctx).
		Where("employee = ? AND status IN ?", employee, []string{"Paid", "Unpaid"}).
		Where("paid_amount > claimed_amount").
		Order("posting_date ASC").
		Find(&advances).Error
	return advances, err
}

// EmployeeAdvanceFilters represents filters for advance queries
type EmployeeAdvanceFilters struct {
	Employee string
	Company  string
	Status   string
}

// ExpenseClaimRepository handles expense claim data access
type ExpenseClaimRepository struct {
	db *gorm.DB
}

// NewExpenseClaimRepository creates a new expense claim repository
func NewExpenseClaimRepository(db *gorm.DB) *ExpenseClaimRepository {
	return &ExpenseClaimRepository{db: db}
}

// Create creates a new expense claim
func (r *ExpenseClaimRepository) Create(ctx context.Context, claim *payroll.ExpenseClaim) error {
	return r.db.WithContext(ctx).Create(claim).Error
}

// GetByID retrieves an expense claim by ID
func (r *ExpenseClaimRepository) GetByID(ctx context.Context, id uint) (*payroll.ExpenseClaim, error) {
	var claim payroll.ExpenseClaim
	err := r.db.WithContext(ctx).
		Preload("Expenses").
		First(&claim, id).Error
	if err != nil {
		return nil, err
	}
	return &claim, nil
}

// Update updates an expense claim
func (r *ExpenseClaimRepository) Update(ctx context.Context, claim *payroll.ExpenseClaim) error {
	return r.db.WithContext(ctx).Save(claim).Error
}

// List retrieves expense claims with filters
func (r *ExpenseClaimRepository) List(ctx context.Context, filters ExpenseClaimFilters, page, pageSize int) ([]payroll.ExpenseClaim, int64, error) {
	var claims []payroll.ExpenseClaim
	var total int64

	query := r.db.WithContext(ctx).Model(&payroll.ExpenseClaim{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.ExpenseApprover != "" {
		query = query.Where("expense_approver = ?", filters.ExpenseApprover)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("posting_date DESC").Find(&claims).Error

	return claims, total, err
}

// ExpenseClaimFilters represents filters for expense claim queries
type ExpenseClaimFilters struct {
	Employee        string
	Company         string
	Status          string
	ExpenseApprover string
}

// GratuityRepository handles gratuity data access
type GratuityRepository struct {
	db *gorm.DB
}

// NewGratuityRepository creates a new gratuity repository
func NewGratuityRepository(db *gorm.DB) *GratuityRepository {
	return &GratuityRepository{db: db}
}

// Create creates a new gratuity record
func (r *GratuityRepository) Create(ctx context.Context, gratuity *payroll.Gratuity) error {
	return r.db.WithContext(ctx).Create(gratuity).Error
}

// GetByID retrieves a gratuity record by ID
func (r *GratuityRepository) GetByID(ctx context.Context, id uint) (*payroll.Gratuity, error) {
	var gratuity payroll.Gratuity
	err := r.db.WithContext(ctx).First(&gratuity, id).Error
	if err != nil {
		return nil, err
	}
	return &gratuity, nil
}

// Update updates a gratuity record
func (r *GratuityRepository) Update(ctx context.Context, gratuity *payroll.Gratuity) error {
	return r.db.WithContext(ctx).Save(gratuity).Error
}

// List retrieves gratuity records with filters
func (r *GratuityRepository) List(ctx context.Context, filters GratuityFilters, page, pageSize int) ([]payroll.Gratuity, int64, error) {
	var gratuities []payroll.Gratuity
	var total int64

	query := r.db.WithContext(ctx).Model(&payroll.Gratuity{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&gratuities).Error

	return gratuities, total, err
}

// GratuityFilters represents filters for gratuity queries
type GratuityFilters struct {
	Employee string
	Company  string
	Status   string
}
