package payroll

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// PayrollEntry represents bulk salary slip generation and processing
type PayrollEntry struct {
	base.BaseModel

	// Naming
	NamingSeries        string     `gorm:"size:50" json:"naming_series"`

	// Payroll Details
	PayrollFrequency    string     `gorm:"size:50;not null" json:"payroll_frequency"` // Monthly, Fortnightly, Weekly, Daily
	Company             string     `gorm:"size:255;not null" json:"company"`
	PayrollPeriod       string     `gorm:"size:255" json:"payroll_period"`
	PostingDate         time.Time  `gorm:"type:date;not null" json:"posting_date"`

	// Period
	StartDate           time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate             time.Time  `gorm:"type:date;not null" json:"end_date"`

	// Filters
	Department          string     `gorm:"size:255" json:"department"`
	Designation         string     `gorm:"size:255" json:"designation"`
	Branch              string     `gorm:"size:255" json:"branch"`
	Project             string     `gorm:"size:255" json:"project"`

	// Employees
	Employees           []PayrollEmployeeDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"employees"`

	// Process Details
	ValidateAttendance  bool       `gorm:"default:false" json:"validate_attendance"`
	SalarySlipBasedOnTimesheet bool `gorm:"default:false" json:"salary_slip_based_on_timesheet"`

	// Totals
	TotalEmployees      int        `gorm:"default:0" json:"number_of_employees"`
	GrossPayTotal       float64    `gorm:"type:decimal(18,2)" json:"total_gross_pay"`
	NetPayTotal         float64    `gorm:"type:decimal(18,2)" json:"total_net_pay"`

	// Deductions Summary
	TotalDeductions     []PayrollDeductionDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"deductions"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft'" json:"status"` // Draft, Submitted, Cancelled

	// Payment
	PaymentAccount      string     `gorm:"size:255" json:"payment_account"`
}

// TableName specifies the table name
func (PayrollEntry) TableName() string {
	return "payroll_entries"
}

// PayrollEmployeeDetail represents employees in a payroll run
type PayrollEmployeeDetail struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	Employee            string  `gorm:"size:255;not null" json:"employee"`
	EmployeeName        string  `gorm:"size:255" json:"employee_name"`

	// Salary Slip
	SalarySlip          string  `gorm:"size:255" json:"salary_slip"`

	// Amounts
	GrossPay            float64 `gorm:"type:decimal(18,2)" json:"gross_pay"`
	NetPay              float64 `gorm:"type:decimal(18,2)" json:"net_pay"`
}

// TableName specifies the table name
func (PayrollEmployeeDetail) TableName() string {
	return "payroll_employee_details"
}

// PayrollDeductionDetail represents summary of deductions in payroll
type PayrollDeductionDetail struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	Component           string  `gorm:"size:255;not null" json:"component"`
	Amount              float64 `gorm:"type:decimal(18,2);not null" json:"amount"`
}

// TableName specifies the table name
func (PayrollDeductionDetail) TableName() string {
	return "payroll_deduction_details"
}

// PayrollPeriod represents a payroll period (fiscal year)
type PayrollPeriod struct {
	base.BaseModel
	PeriodName          string    `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Company             string    `gorm:"size:255;not null" json:"company"`
	StartDate           time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate             time.Time `gorm:"type:date;not null" json:"end_date"`
	CostCenter          string    `gorm:"size:255" json:"cost_center"`
}

// TableName specifies the table name
func (PayrollPeriod) TableName() string {
	return "payroll_periods"
}

// Gratuity represents gratuity calculation and payment
type Gratuity struct {
	base.BaseModel

	// Naming
	NamingSeries        string     `gorm:"size:50" json:"naming_series"`

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`
	Department          string     `gorm:"size:255" json:"department"`
	Designation         string     `gorm:"size:255" json:"designation"`

	// Service Period
	DateOfJoining       time.Time  `gorm:"type:date" json:"date_of_joining"`
	RelievingDate       time.Time  `gorm:"type:date" json:"relieving_date"`
	TotalWorkingDays    int        `json:"total_working_days"`
	YearsOfService      int        `json:"number_of_years_of_service"`

	// Calculation
	GratuitySlab        string     `gorm:"size:255" json:"gratuity_rule"`
	CurrentSalary       float64    `gorm:"type:decimal(18,2)" json:"current_work_experience"`
	Amount              float64    `gorm:"type:decimal(18,2)" json:"gratuity_amount"`

	// Payment
	PayViaAdditionalSalary bool    `gorm:"default:false" json:"pay_via_salary_slip"`
	AdditionalSalary    string     `gorm:"size:255" json:"additional_salary_component"`
	PayrollDate         *time.Time `gorm:"type:date" json:"payroll_date"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft'" json:"status"` // Draft, Submitted, Paid
}

// TableName specifies the table name
func (Gratuity) TableName() string {
	return "gratuities"
}

// GratuitySlab represents gratuity calculation rules
type GratuitySlab struct {
	base.BaseModel
	SlabName            string  `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Company             string  `gorm:"size:255" json:"company"`

	// Calculation Method
	CalculationBasis    string  `gorm:"size:100;default:'15/26 * last drawn salary * number of service years'" json:"calculate_gratuity_amount_based_on"`
	WorkExperienceCalculationMethod string `gorm:"size:100" json:"work_experience_calculation_method"`
	FractionOfYears     float64 `gorm:"type:decimal(5,2)" json:"fraction_of_applicable_earnings"`
	MinimumYearsOfService int   `gorm:"default:5" json:"minimum_year_for_gratuity"`

	// Slabs
	Slabs               []GratuitySlabDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"gratuity_rule_slabs"`
}

// TableName specifies the table name
func (GratuitySlab) TableName() string {
	return "gratuity_slabs"
}

// GratuitySlabDetail represents gratuity calculation brackets
type GratuitySlabDetail struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	FromYears           int     `gorm:"not null" json:"from_year_of_service"`
	ToYears             int     `json:"to_year_of_service"`
	SlabPercent         float64 `gorm:"type:decimal(5,2);not null" json:"fraction_of_applicable_earnings"`
}

// TableName specifies the table name
func (GratuitySlabDetail) TableName() string {
	return "gratuity_slab_details"
}

// Loan represents employee loan
type Loan struct {
	base.BaseModel

	// Naming
	NamingSeries        string     `gorm:"size:50" json:"naming_series"`

	// Applicant
	Applicant           string     `gorm:"size:255;not null;index" json:"applicant"`
	ApplicantType       string     `gorm:"size:50;default:'Employee'" json:"applicant_type"` // Employee, Customer, Member
	ApplicantName       string     `gorm:"size:255" json:"applicant_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Loan Details
	LoanType            string     `gorm:"size:255;not null" json:"loan_type"`
	LoanApplication     string     `gorm:"size:255" json:"loan_application"`
	LoanAmount          float64    `gorm:"type:decimal(18,2);not null" json:"loan_amount"`
	RateOfInterest      float64    `gorm:"type:decimal(5,2)" json:"rate_of_interest"`
	IsTermLoan          bool       `gorm:"default:true" json:"is_term_loan"`

	// Repayment
	RepaymentMethod     string     `gorm:"size:100;default:'Repay Over Number of Periods'" json:"repayment_method"`
	MonthlyRepaymentAmount float64 `gorm:"type:decimal(18,2)" json:"monthly_repayment_amount"`
	RepaymentPeriods    int        `gorm:"default:0" json:"repayment_periods"`
	RepaymentStartDate  time.Time  `gorm:"type:date" json:"repayment_start_date"`

	// Dates
	PostingDate         time.Time  `gorm:"type:date;not null" json:"posting_date"`
	DisbursementDate    time.Time  `gorm:"type:date" json:"disbursement_date"`

	// Repayment Schedule
	RepaymentSchedule   []LoanRepaymentSchedule `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"repayment_schedule"`

	// Amounts
	TotalPayment        float64    `gorm:"type:decimal(18,2)" json:"total_payment"`
	TotalInterestPayable float64   `gorm:"type:decimal(18,2)" json:"total_interest_payable"`
	TotalAmountPaid     float64    `gorm:"type:decimal(18,2)" json:"total_amount_paid"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft'" json:"status"` // Draft, Sanctioned, Disbursed, Closed, Cancelled

	// Accounting
	Mode                string     `gorm:"size:50" json:"mode_of_payment"`
	PaymentAccount      string     `gorm:"size:255" json:"payment_account"`
	LoanAccount         string     `gorm:"size:255" json:"loan_account"`
	InterestIncomeAccount string   `gorm:"size:255" json:"interest_income_account"`
}

// TableName specifies the table name
func (Loan) TableName() string {
	return "loans"
}

// LoanRepaymentSchedule represents loan repayment installments
type LoanRepaymentSchedule struct {
	gorm.Model
	ParentID            uint      `gorm:"index" json:"parent_id"`
	PaymentDate         time.Time `gorm:"type:date;not null" json:"payment_date"`
	PrincipalAmount     float64   `gorm:"type:decimal(18,2);not null" json:"principal_amount"`
	InterestAmount      float64   `gorm:"type:decimal(18,2)" json:"interest_amount"`
	TotalPayment        float64   `gorm:"type:decimal(18,2);not null" json:"total_payment"`
	Balance             float64   `gorm:"type:decimal(18,2)" json:"balance_loan_amount"`
	IsPaid              bool      `gorm:"default:false" json:"is_accrued"`
}

// TableName specifies the table name
func (LoanRepaymentSchedule) TableName() string {
	return "loan_repayment_schedules"
}

// LoanType represents types of loans
type LoanType struct {
	base.BaseModel
	LoanTypeName        string  `gorm:"size:255;not null;uniqueIndex" json:"loan_name"`
	Company             string  `gorm:"size:255" json:"company"`
	MaximumLoanAmount   float64 `gorm:"type:decimal(18,2)" json:"maximum_loan_amount"`
	RateOfInterest      float64 `gorm:"type:decimal(5,2)" json:"rate_of_interest"`
	Description         string  `gorm:"type:text" json:"description"`
	Disabled            bool    `gorm:"default:false" json:"disabled"`
}

// TableName specifies the table name
func (LoanType) TableName() string {
	return "loan_types"
}

// EmployeeAdvance represents employee cash advance
type EmployeeAdvance struct {
	base.BaseModel

	// Naming
	NamingSeries        string     `gorm:"size:50" json:"naming_series"`

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`
	Department          string     `gorm:"size:255" json:"department"`

	// Advance Details
	PostingDate         time.Time  `gorm:"type:date;not null" json:"posting_date"`
	AdvanceAmount       float64    `gorm:"type:decimal(18,2);not null" json:"advance_amount"`
	Purpose             string     `gorm:"type:text;not null" json:"purpose"`

	// Repayment
	RepayFromSalary     bool       `gorm:"default:false" json:"repay_unclaimed_amount_from_salary"`
	RepaymentDate       *time.Time `gorm:"type:date" json:"return_date"`

	// Amounts
	PaidAmount          float64    `gorm:"type:decimal(18,2)" json:"paid_amount"`
	ClaimedAmount       float64    `gorm:"type:decimal(18,2)" json:"claimed_amount"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft'" json:"status"` // Draft, Paid, Claimed, Returned, Unpaid

	// Accounting
	Mode                string     `gorm:"size:50" json:"mode_of_payment"`
	AdvanceAccount      string     `gorm:"size:255" json:"advance_account"`
}

// TableName specifies the table name
func (EmployeeAdvance) TableName() string {
	return "employee_advances"
}

// ExpenseClaim represents employee expense reimbursement
type ExpenseClaim struct {
	base.BaseModel

	// Naming
	NamingSeries        string     `gorm:"size:50" json:"naming_series"`

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`
	Department          string     `gorm:"size:255" json:"department"`

	// Expense Details
	PostingDate         time.Time  `gorm:"type:date;not null" json:"posting_date"`
	ExpenseApprover     string     `gorm:"size:255" json:"expense_approver"`

	// Advance
	EmployeeAdvance     string     `gorm:"size:255" json:"advance_paid"`
	AdvanceAmount       float64    `gorm:"type:decimal(18,2)" json:"advance_amount"`

	// Expenses
	Expenses            []ExpenseClaimDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"expenses"`

	// Amounts
	TotalClaimedAmount  float64    `gorm:"type:decimal(18,2)" json:"total_claimed_amount"`
	TotalSanctionedAmount float64  `gorm:"type:decimal(18,2)" json:"total_sanctioned_amount"`
	TotalTaxes          float64    `gorm:"type:decimal(18,2)" json:"total_taxes_and_charges"`
	GrandTotal          float64    `gorm:"type:decimal(18,2)" json:"grand_total"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft'" json:"approval_status"` // Draft, Submitted, Approved, Rejected, Paid, Unpaid

	// Accounting
	Mode                string     `gorm:"size:50" json:"mode_of_payment"`
	PayableAccount      string     `gorm:"size:255" json:"payable_account"`
}

// TableName specifies the table name
func (ExpenseClaim) TableName() string {
	return "expense_claims"
}

// ExpenseClaimDetail represents individual expense items
type ExpenseClaimDetail struct {
	gorm.Model
	ParentID            uint      `gorm:"index" json:"parent_id"`
	ExpenseDate         time.Time `gorm:"type:date;not null" json:"expense_date"`
	ExpenseType         string    `gorm:"size:255;not null" json:"expense_type"`
	Description         string    `gorm:"type:text" json:"description"`
	Amount              float64   `gorm:"type:decimal(18,2);not null" json:"amount"`
	SanctionedAmount    float64   `gorm:"type:decimal(18,2)" json:"sanctioned_amount"`

	// Accounting
	DefaultAccount      string    `gorm:"size:255" json:"default_account"`
	CostCenter          string    `gorm:"size:255" json:"cost_center"`
}

// TableName specifies the table name
func (ExpenseClaimDetail) TableName() string {
	return "expense_claim_details"
}

// ExpenseClaimType represents types of expense claims
type ExpenseClaimType struct {
	base.BaseModel
	ExpenseType         string `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Description         string `gorm:"type:text" json:"description"`
	DefaultAccount      string `gorm:"size:255" json:"deferred_expense_account"`
	Disabled            bool   `gorm:"default:false" json:"disabled"`
}

// TableName specifies the table name
func (ExpenseClaimType) TableName() string {
	return "expense_claim_types"
}
