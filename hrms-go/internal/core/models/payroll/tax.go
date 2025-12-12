package payroll

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// IncomeTaxSlab represents tax slabs for income tax calculation
type IncomeTaxSlab struct {
	base.BaseModel
	SlabName            string `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Company             string `gorm:"size:255" json:"company"`
	EffectiveFrom       time.Time `gorm:"type:date;not null" json:"effective_from"`
	Currency            string `gorm:"size:10;default:'GHS'" json:"currency"`
	Disabled            bool   `gorm:"default:false" json:"disabled"`

	// Allow Tax Exemption
	AllowTaxExemption   bool   `gorm:"default:false" json:"allow_tax_exemption"`
	StandardTaxExemptionAmount float64 `gorm:"type:decimal(18,2)" json:"standard_tax_exemption_amount"`

	// Slabs (Ghana Tax Rates)
	Slabs               []IncomeTaxSlabDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"slabs"`

	// Other Tax Settings
	OtherTaxes          []IncomeTaxSlabOtherTax `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"other_taxes_and_charges"`
}

// TableName specifies the table name
func (IncomeTaxSlab) TableName() string {
	return "income_tax_slabs"
}

// IncomeTaxSlabDetail represents individual tax brackets
type IncomeTaxSlabDetail struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	FromAmount          float64 `gorm:"type:decimal(18,2);not null" json:"from_amount"`
	ToAmount            float64 `gorm:"type:decimal(18,2)" json:"to_amount"`
	PercentDeduction    float64 `gorm:"type:decimal(5,2);not null" json:"percent_deduction"`
	Condition           string  `gorm:"type:text" json:"condition"`
}

// TableName specifies the table name
func (IncomeTaxSlabDetail) TableName() string {
	return "income_tax_slab_details"
}

// IncomeTaxSlabOtherTax represents other statutory taxes (SSNIT, etc.)
type IncomeTaxSlabOtherTax struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	Description         string  `gorm:"size:255;not null" json:"description"`
	Category            string  `gorm:"size:100" json:"category"` // "Employee Tax", "Employer Tax", "Other"
	PercentageOrAmount  string  `gorm:"size:50;default:'Percentage'" json:"percentage_or_amount"`
	Percentage          float64 `gorm:"type:decimal(5,2)" json:"percent"`
	Amount              float64 `gorm:"type:decimal(18,2)" json:"amount"`
	MinTaxableIncome    float64 `gorm:"type:decimal(18,2)" json:"min_taxable_income"`
	MaxTaxableIncome    float64 `gorm:"type:decimal(18,2)" json:"max_taxable_income"`
}

// TableName specifies the table name
func (IncomeTaxSlabOtherTax) TableName() string {
	return "income_tax_slab_other_taxes"
}

// EmployeeTaxExemptionDeclaration represents employee's tax exemption claims
type EmployeeTaxExemptionDeclaration struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Period
	PayrollPeriod       string     `gorm:"size:255;not null" json:"payroll_period"`
	Currency            string     `gorm:"size:10;default:'GHS'" json:"currency"`

	// Declarations
	Declarations        []EmployeeTaxExemptionDeclarationCategory `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"declarations"`

	// Totals
	TotalExemptionAmount float64   `gorm:"type:decimal(18,2)" json:"total_exemption_amount"`

	// Submission
	SubmissionDate      *time.Time `gorm:"type:date" json:"submission_date"`
}

// TableName specifies the table name
func (EmployeeTaxExemptionDeclaration) TableName() string {
	return "employee_tax_exemption_declarations"
}

// EmployeeTaxExemptionDeclarationCategory represents exemption categories
type EmployeeTaxExemptionDeclarationCategory struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	ExemptionSubCategory string `gorm:"size:255;not null" json:"exemption_sub_category"`
	ExemptionCategory   string  `gorm:"size:255" json:"exemption_category"`
	MaxExemptionAmount  float64 `gorm:"type:decimal(18,2)" json:"maximum_exempted_amount"`
	DeclaredAmount      float64 `gorm:"type:decimal(18,2);not null" json:"amount"`
}

// TableName specifies the table name
func (EmployeeTaxExemptionDeclarationCategory) TableName() string {
	return "employee_tax_exemption_declaration_categories"
}

// EmployeeTaxExemptionProofSubmission represents proof submission for exemptions
type EmployeeTaxExemptionProofSubmission struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Period
	PayrollPeriod       string     `gorm:"size:255;not null" json:"payroll_period"`
	Currency            string     `gorm:"size:10;default:'GHS'" json:"currency"`

	// Submission Details
	TaxExemptionProofs  []EmployeeTaxExemptionProofDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"tax_exemption_proofs"`

	// Total
	TotalActualAmount   float64    `gorm:"type:decimal(18,2)" json:"total_actual_amount"`

	// Submission
	SubmissionDate      *time.Time `gorm:"type:date" json:"submission_date"`
}

// TableName specifies the table name
func (EmployeeTaxExemptionProofSubmission) TableName() string {
	return "employee_tax_exemption_proof_submissions"
}

// EmployeeTaxExemptionProofDetail represents individual proof details
type EmployeeTaxExemptionProofDetail struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	ExemptionSubCategory string `gorm:"size:255;not null" json:"exemption_sub_category"`
	ExemptionCategory   string  `gorm:"size:255" json:"exemption_category"`
	Type                string  `gorm:"size:255" json:"type_of_proof"`
	MaxExemptionAmount  float64 `gorm:"type:decimal(18,2)" json:"maximum_exempted_amount"`
	Amount              float64 `gorm:"type:decimal(18,2);not null" json:"amount"`
}

// TableName specifies the table name
func (EmployeeTaxExemptionProofDetail) TableName() string {
	return "employee_tax_exemption_proof_details"
}

// EmployeeOtherIncome represents other sources of income for tax calculation
type EmployeeOtherIncome struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Income Details
	Source              string     `gorm:"size:255;not null" json:"source"`
	Amount              float64    `gorm:"type:decimal(18,2);not null" json:"amount"`

	// Period
	PayrollPeriod       string     `gorm:"size:255;not null" json:"payroll_period"`
	Currency            string     `gorm:"size:10;default:'GHS'" json:"currency"`

	// Description
	Description         string     `gorm:"type:text" json:"description"`
}

// TableName specifies the table name
func (EmployeeOtherIncome) TableName() string {
	return "employee_other_incomes"
}

// TaxWithholdingCategory represents tax categories for vendors/customers
type TaxWithholdingCategory struct {
	base.BaseModel
	CategoryName        string  `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Category            string  `gorm:"size:100" json:"category_name"`
	TaxWithholdingRate  float64 `gorm:"type:decimal(5,2);not null" json:"rate"`
}

// TableName specifies the table name
func (TaxWithholdingCategory) TableName() string {
	return "tax_withholding_categories"
}

// EmployeeBenefitApplication represents benefit applications (HRA, LTA, etc.)
type EmployeeBenefitApplication struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Period
	PayrollPeriod       string     `gorm:"size:255;not null" json:"payroll_period"`
	Currency            string     `gorm:"size:10;default:'GHS'" json:"currency"`

	// Benefits
	Benefits            []EmployeeBenefitApplicationDetail `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"employee_benefits"`

	// Total
	TotalAmount         float64    `gorm:"type:decimal(18,2)" json:"total_amount"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft'" json:"status"` // Draft, Submitted, Approved, Rejected
}

// TableName specifies the table name
func (EmployeeBenefitApplication) TableName() string {
	return "employee_benefit_applications"
}

// EmployeeBenefitApplicationDetail represents benefit details
type EmployeeBenefitApplicationDetail struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	EarningComponent    string  `gorm:"size:255;not null" json:"earning_component"`
	Amount              float64 `gorm:"type:decimal(18,2);not null" json:"amount"`
	Date                time.Time `gorm:"type:date" json:"date"`
}

// TableName specifies the table name
func (EmployeeBenefitApplicationDetail) TableName() string {
	return "employee_benefit_application_details"
}

// EmployeeBenefitClaim represents claims against applied benefits
type EmployeeBenefitClaim struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`

	// Claim Details
	ClaimDate           time.Time  `gorm:"type:date;not null" json:"claim_date"`
	ClaimBenefit        string     `gorm:"size:255;not null" json:"claim_benefit"`
	ClaimedAmount       float64    `gorm:"type:decimal(18,2);not null" json:"claimed_amount"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft'" json:"status"` // Draft, Submitted, Approved, Rejected
}

// TableName specifies the table name
func (EmployeeBenefitClaim) TableName() string {
	return "employee_benefit_claims"
}
