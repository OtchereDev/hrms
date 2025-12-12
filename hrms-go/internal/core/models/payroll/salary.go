package payroll

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// SalaryComponent represents an earnings or deduction component
type SalaryComponent struct {
	base.BaseModel
	SalaryComponentName string  `gorm:"size:255;not null;uniqueIndex" json:"salary_component"`
	Type                string  `gorm:"size:50;not null" json:"type"` // Earning, Deduction

	// Component Settings
	Description         string  `gorm:"type:text" json:"description"`
	Company             string  `gorm:"size:255" json:"company"`
	IsActive            bool    `gorm:"default:true" json:"is_active"`

	// Tax & Statutory
	IsTaxable           bool    `gorm:"default:false" json:"is_tax_applicable"`
	ExemptFromTotalInTax bool   `gorm:"default:false" json:"exempted_from_income_tax"`
	DoNotIncludeInTotal bool    `gorm:"default:false" json:"do_not_include_in_total"`

	// Flexible Benefit
	IsFlexibleBenefit   bool    `gorm:"default:false" json:"is_flexible_benefit"`
	MaxBenefitAmount    float64 `gorm:"type:decimal(18,2)" json:"max_benefit_amount"`

	// Statutory Component (SSNIT, etc.)
	IsStatutory         bool    `gorm:"default:false" json:"statistical_component"`
	ComponentType       string  `gorm:"size:100" json:"component_type"` // Provident Fund, Additional Provident Fund, etc.

	// Deduction Settings (if type is Deduction)
	RoundToNearestNumber float64 `gorm:"type:decimal(10,2)" json:"round_to_nearest_number"`

	// Formula
	AmountBasedOnFormula bool   `gorm:"default:false" json:"amount_based_on_formula"`
	Formula              string `gorm:"type:text" json:"formula"`
	FormulaHelp          string `gorm:"type:text" json:"formula_help"`

	// Accounting
	Accounts             []SalaryComponentAccount `gorm:"foreignKey:ParentID" json:"accounts"`
}

// TableName specifies the table name
func (SalaryComponent) TableName() string {
	return "salary_components"
}

// SalaryComponentAccount represents accounting entries for salary components
type SalaryComponentAccount struct {
	gorm.Model
	ParentID       uint   `gorm:"index" json:"parent_id"`
	Company        string `gorm:"size:255;not null" json:"company"`
	Account        string `gorm:"size:255" json:"account"`
	DefaultAccount string `gorm:"size:255" json:"default_account"`
}

// TableName specifies the table name
func (SalaryComponentAccount) TableName() string {
	return "salary_component_accounts"
}

// SalaryStructure represents a salary structure template
type SalaryStructure struct {
	base.BaseModel
	SalaryStructureName string `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Company             string `gorm:"size:255;not null" json:"company"`
	PayrollFrequency    string `gorm:"size:50;not null;default:'Monthly'" json:"payroll_frequency"` // Monthly, Fortnightly, Weekly, Daily

	// Pay Period
	IsActive            bool   `gorm:"default:true" json:"is_active"`
	PayrollPeriod       string `gorm:"size:255" json:"payroll_period"`

	// Currency
	Currency            string `gorm:"size:10;default:'GHS'" json:"currency"`

	// Settings
	HourRate            float64 `gorm:"type:decimal(18,2)" json:"hour_rate"`
	LeavesEncashment    bool    `gorm:"default:false" json:"include_in_leave_encashment_calculation"`
	MaxBenefits         float64 `gorm:"type:decimal(18,2)" json:"max_benefits"`

	// Mode of Payment
	ModeOfPayment       string  `gorm:"size:50;default:'Bank'" json:"mode_of_payment"` // Bank, Cash, Cheque

	// Components
	Earnings            []SalaryDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"earnings"`
	Deductions          []SalaryDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"deductions"`

	// Other Benefits
	OtherBenefits       []SalaryDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"other_benefits"`
}

// TableName specifies the table name
func (SalaryStructure) TableName() string {
	return "salary_structures"
}

// SalaryDetail represents earnings/deductions in a salary structure
type SalaryDetail struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	ParentType          string  `gorm:"size:50" json:"parent_type"` // "SalaryStructure" or "SalarySlip"
	SalaryComponent     string  `gorm:"size:255;not null" json:"salary_component"`
	ComponentType       string  `gorm:"size:50" json:"component_type"` // Earning, Deduction

	// Amount Calculation
	Amount              float64 `gorm:"type:decimal(18,2)" json:"amount"`
	AmountBasedOnFormula bool   `gorm:"default:false" json:"amount_based_on_formula"`
	Formula             string  `gorm:"type:text" json:"formula"`

	// Condition
	DoNotIncludeInTotal bool    `gorm:"default:false" json:"do_not_include_in_total"`
	StatisticalComponent bool   `gorm:"default:false" json:"statistical_component"`

	// Tax
	IsTaxable           bool    `gorm:"default:false" json:"is_tax_applicable"`
	IsFlexibleBenefit   bool    `gorm:"default:false" json:"is_flexible_benefit"`

	// Account
	DependOnPaymentDays bool    `gorm:"default:false" json:"depend_on_payment_days"`
}

// TableName specifies the table name
func (SalaryDetail) TableName() string {
	return "salary_details"
}

// SalaryStructureAssignment assigns a salary structure to an employee
type SalaryStructureAssignment struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Salary Structure
	SalaryStructure     string     `gorm:"size:255;not null" json:"salary_structure"`

	// Effective Date
	FromDate            time.Time  `gorm:"type:date;not null;index" json:"from_date"`
	ToDate              *time.Time `gorm:"type:date" json:"to_date"`

	// Base Amount
	Base                float64    `gorm:"type:decimal(18,2)" json:"base"`
	Variable            float64    `gorm:"type:decimal(18,2)" json:"variable"`

	// Income Tax Slab
	IncomeTaxSlab       string     `gorm:"size:255" json:"income_tax_slab"`

	// Payroll Cost Center
	PayrollCostCenter   string     `gorm:"size:255" json:"payroll_cost_center"`
}

// TableName specifies the table name
func (SalaryStructureAssignment) TableName() string {
	return "salary_structure_assignments"
}

// SalarySlip represents an employee's salary slip for a period
type SalarySlip struct {
	base.BaseModel

	// Naming
	NamingSeries        string     `gorm:"size:50" json:"naming_series"`

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`
	Department          string     `gorm:"size:255" json:"department"`
	Designation         string     `gorm:"size:255" json:"designation"`
	Branch              string     `gorm:"size:255" json:"branch"`
	Grade               string     `gorm:"size:255" json:"grade"`

	// Period
	StartDate           time.Time  `gorm:"type:date;not null;index" json:"start_date"`
	EndDate             time.Time  `gorm:"type:date;not null;index" json:"end_date"`
	PostingDate         time.Time  `gorm:"type:date;not null" json:"posting_date"`

	// Salary Structure
	SalaryStructure     string     `gorm:"size:255" json:"salary_structure"`
	PayrollFrequency    string     `gorm:"size:50" json:"payroll_frequency"`

	// Bank Details
	BankName            string     `gorm:"size:255" json:"bank_name"`
	BankAccountNo       string     `gorm:"size:50" json:"bank_account_no"`

	// Attendance & Leave
	TotalWorkingDays    float64    `gorm:"type:decimal(10,2)" json:"total_working_days"`
	PaymentDays         float64    `gorm:"type:decimal(10,2)" json:"payment_days"`
	LeavesWithoutPay    float64    `gorm:"type:decimal(10,2)" json:"leave_without_pay"`
	AbsentDays          float64    `gorm:"type:decimal(10,2)" json:"absent_days"`

	// Earnings & Deductions
	Earnings            []SalaryDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"earnings"`
	Deductions          []SalaryDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"deductions"`

	// Totals
	GrossPay            float64    `gorm:"type:decimal(18,2)" json:"gross_pay"`
	TotalDeduction      float64    `gorm:"type:decimal(18,2)" json:"total_deduction"`
	NetPay              float64    `gorm:"type:decimal(18,2)" json:"net_pay"`
	RoundedTotal        float64    `gorm:"type:decimal(18,2)" json:"rounded_total"`

	// Tax
	TotalTaxableAmount  float64    `gorm:"type:decimal(18,2)" json:"total_taxable_amount"`
	TaxDeducted         float64    `gorm:"type:decimal(18,2)" json:"total_tax_deducted"`

	// Mode of Payment
	ModeOfPayment       string     `gorm:"size:50" json:"mode_of_payment"`
	PaymentReference    string     `gorm:"size:255" json:"payment_reference"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft';index" json:"status"` // Draft, Submitted, Paid, Cancelled

	// Email
	EmailSent           bool       `gorm:"default:false" json:"email_sent"`

	// Payroll Entry Reference
	PayrollEntry        string     `gorm:"size:255" json:"payroll_entry"`

	// Letter Head & Print
	LetterHead          string     `gorm:"size:255" json:"letter_head"`

	// Journal Entry
	JournalEntry        string     `gorm:"size:255" json:"journal_entry"`
}

// TableName specifies the table name
func (SalarySlip) TableName() string {
	return "salary_slips"
}

// BeforeCreate hook for salary slip
func (s *SalarySlip) BeforeCreate(tx *gorm.DB) error {
	if err := s.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Set posting date if not set
	if s.PostingDate.IsZero() {
		s.PostingDate = time.Now()
	}

	return nil
}

// AdditionalSalary represents one-time additions to salary
type AdditionalSalary struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Salary Component
	SalaryComponent     string     `gorm:"size:255;not null" json:"salary_component"`
	Type                string     `gorm:"size:50" json:"type"` // Earning, Deduction

	// Amount
	Amount              float64    `gorm:"type:decimal(18,2);not null" json:"amount"`
	IsRecurring         bool       `gorm:"default:false" json:"is_recurring"`

	// Period
	PayrollDate         time.Time  `gorm:"type:date;not null" json:"payroll_date"`
	FromDate            *time.Time `gorm:"type:date" json:"from_date"`
	ToDate              *time.Time `gorm:"type:date" json:"to_date"`

	// References
	SalarySlip          string     `gorm:"size:255" json:"ref_docname"`
	Overwrite           bool       `gorm:"default:false" json:"overwrite_salary_structure_amount"`

	// Reason
	Reason              string     `gorm:"type:text" json:"reason"`
}

// TableName specifies the table name
func (AdditionalSalary) TableName() string {
	return "additional_salaries"
}

// RetentionBonus represents employee retention bonus
type RetentionBonus struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Bonus Details
	BonusPaymentDate    time.Time  `gorm:"type:date;not null" json:"bonus_payment_date"`
	BonusAmount         float64    `gorm:"type:decimal(18,2);not null" json:"bonus_amount"`

	// Additional Salary Reference
	AdditionalSalary    string     `gorm:"size:255" json:"additional_salary"`
}

// TableName specifies the table name
func (RetentionBonus) TableName() string {
	return "retention_bonuses"
}

// EmployeeIncentive represents performance-based incentives
type EmployeeIncentive struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`
	Department          string     `gorm:"size:255" json:"department"`

	// Incentive Details
	IncentiveAmount     float64    `gorm:"type:decimal(18,2);not null" json:"incentive_amount"`
	PayrollDate         time.Time  `gorm:"type:date;not null" json:"payroll_date"`

	// Currency
	Currency            string     `gorm:"size:10;default:'GHS'" json:"currency"`

	// Additional Salary Reference
	AdditionalSalary    string     `gorm:"size:255" json:"additional_salary"`
}

// TableName specifies the table name
func (EmployeeIncentive) TableName() string {
	return "employee_incentives"
}
