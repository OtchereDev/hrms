package hr

import (
	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// Company represents an organization/company
type Company struct {
	base.BaseModel
	CompanyName      string  `gorm:"size:255;not null" json:"company_name"`
	Abbr             string  `gorm:"size:10;uniqueIndex" json:"abbr"` // Abbreviation
	Domain           string  `gorm:"size:100" json:"domain"`
	DefaultCurrency  string  `gorm:"size:10;default:'GHS'" json:"default_currency"`
	Country          string  `gorm:"size:100;not null" json:"country"`
	DateOfEstablishment *gorm.DeletedAt `json:"date_of_establishment,omitempty"`

	// Contact Information
	Phone            string  `gorm:"size:20" json:"phone"`
	Email            string  `gorm:"size:255" json:"email"`
	Website          string  `gorm:"size:255" json:"website"`

	// Address
	Address          string  `gorm:"type:text" json:"address"`
	City             string  `gorm:"size:100" json:"city"`
	State            string  `gorm:"size:100" json:"state"`
	PostalCode       string  `gorm:"size:20" json:"postal_code"`

	// Tax Information
	TaxID            string  `gorm:"size:50" json:"tax_id"`            // TIN/Tax ID
	RegistrationNo   string  `gorm:"size:50" json:"registration_no"`   // Company registration number

	// Settings
	DefaultHolidayList string `gorm:"size:255" json:"default_holiday_list"`
	IsGroup          bool    `gorm:"default:false" json:"is_group"`
	ParentCompany    string  `gorm:"size:255" json:"parent_company,omitempty"`

	// Relationships
	Departments      []Department `gorm:"foreignKey:CompanyName;references:CompanyName" json:"departments,omitempty"`
	Branches         []Branch     `gorm:"foreignKey:CompanyName;references:CompanyName" json:"branches,omitempty"`
}

// TableName specifies the table name
func (Company) TableName() string {
	return "companies"
}

// Department represents an organizational department
type Department struct {
	base.BaseModel
	DepartmentName   string  `gorm:"size:255;not null;index" json:"department_name"`
	CompanyName      string  `gorm:"size:255;index" json:"company"`
	ParentDepartment string  `gorm:"size:255" json:"parent_department,omitempty"`
	IsGroup          bool    `gorm:"default:false" json:"is_group"`
	Disabled         bool    `gorm:"default:false" json:"disabled"`

	// Department Head
	DepartmentHead   string  `gorm:"size:255" json:"department_head,omitempty"` // Link to Employee

	// Contact
	Email            string  `gorm:"size:255" json:"email"`
	Phone            string  `gorm:"size:20" json:"phone"`

	// Leave Settings
	LeaveBlockList   string  `gorm:"size:255" json:"leave_block_list,omitempty"`

	// Approvers
	LeaveApprovers   []DepartmentApprover  `gorm:"foreignKey:ParentID" json:"leave_approvers,omitempty"`
	ExpenseApprovers []DepartmentApprover  `gorm:"foreignKey:ParentID" json:"expense_approvers,omitempty"`

	// Relationships
	Company          *Company    `gorm:"foreignKey:CompanyName;references:CompanyName" json:"-"`
	Employees        []Employee  `gorm:"foreignKey:Department;references:DepartmentName" json:"employees,omitempty"`
}

// TableName specifies the table name
func (Department) TableName() string {
	return "departments"
}

// DepartmentApprover represents an approver for a department
type DepartmentApprover struct {
	gorm.Model
	ParentID     uint   `gorm:"index" json:"parent_id"`
	ApproverType string `gorm:"size:50" json:"approver_type"` // "Leave" or "Expense"
	Approver     string `gorm:"size:255" json:"approver"`     // Link to User
	Required     bool   `gorm:"default:false" json:"required"`
}

// TableName specifies the table name
func (DepartmentApprover) TableName() string {
	return "department_approvers"
}

// Branch represents a company branch/location
type Branch struct {
	base.BaseModel
	BranchName       string  `gorm:"size:255;not null;index" json:"branch"`
	CompanyName      string  `gorm:"size:255;index" json:"company"`

	// Address
	Address          string  `gorm:"type:text" json:"address"`
	City             string  `gorm:"size:100" json:"city"`
	State            string  `gorm:"size:100" json:"state"`
	Country          string  `gorm:"size:100" json:"country"`
	PostalCode       string  `gorm:"size:20" json:"postal_code"`

	// Contact
	Phone            string  `gorm:"size:20" json:"phone"`
	Email            string  `gorm:"size:255" json:"email"`

	// Manager
	BranchManager    string  `gorm:"size:255" json:"branch_manager,omitempty"`

	// Relationships
	Company          *Company   `gorm:"foreignKey:CompanyName;references:CompanyName" json:"-"`
	Employees        []Employee `gorm:"foreignKey:Branch;references:BranchName" json:"employees,omitempty"`
}

// TableName specifies the table name
func (Branch) TableName() string {
	return "branches"
}

// Designation represents a job position/title
type Designation struct {
	base.BaseModel
	DesignationName  string  `gorm:"size:255;not null;uniqueIndex" json:"designation_name"`
	Description      string  `gorm:"type:text" json:"description"`

	// Skills Required (JSON array or separate table)
	RequiredSkills   string  `gorm:"type:text" json:"required_skills"` // JSON array

	// Relationships
	Employees        []Employee `gorm:"foreignKey:Designation;references:DesignationName" json:"employees,omitempty"`
}

// TableName specifies the table name
func (Designation) TableName() string {
	return "designations"
}

// EmploymentType represents types of employment (Full-time, Part-time, etc.)
type EmploymentType struct {
	base.BaseModel
	EmploymentTypeName string `gorm:"size:255;not null;uniqueIndex" json:"employment_type"`

	// Relationships
	Employees          []Employee `gorm:"foreignKey:EmploymentType;references:EmploymentTypeName" json:"employees,omitempty"`
}

// TableName specifies the table name
func (EmploymentType) TableName() string {
	return "employment_types"
}

// EmployeeGrade represents employee grade/level
type EmployeeGrade struct {
	base.BaseModel
	GradeName        string  `gorm:"size:255;not null;uniqueIndex" json:"name"`
	DefaultBasePay   float64 `gorm:"type:decimal(18,2)" json:"default_base_pay"`
	DefaultLeaves    int     `gorm:"default:0" json:"default_leaves_per_year"`

	// Relationships
	Employees        []Employee `gorm:"foreignKey:Grade;references:GradeName" json:"employees,omitempty"`
}

// TableName specifies the table name
func (EmployeeGrade) TableName() string {
	return "employee_grades"
}

// HolidayList represents a list of holidays
type HolidayList struct {
	base.BaseModel
	HolidayListName  string    `gorm:"size:255;not null;uniqueIndex" json:"holiday_list_name"`
	FromDate         string    `gorm:"type:date" json:"from_date"`
	ToDate           string    `gorm:"type:date" json:"to_date"`
	TotalHolidays    int       `gorm:"default:0" json:"total_holidays"`
	Country          string    `gorm:"size:100" json:"country"`

	// Weekly offs
	WeeklySunday     bool      `gorm:"default:false" json:"weekly_off_sunday"`
	WeeklyMonday     bool      `gorm:"default:false" json:"weekly_off_monday"`
	WeeklyTuesday    bool      `gorm:"default:false" json:"weekly_off_tuesday"`
	WeeklyWednesday  bool      `gorm:"default:false" json:"weekly_off_wednesday"`
	WeeklyThursday   bool      `gorm:"default:false" json:"weekly_off_thursday"`
	WeeklyFriday     bool      `gorm:"default:false" json:"weekly_off_friday"`
	WeeklySaturday   bool      `gorm:"default:false" json:"weekly_off_saturday"`

	// Relationships
	Holidays         []Holiday `gorm:"foreignKey:ParentID" json:"holidays,omitempty"`
}

// TableName specifies the table name
func (HolidayList) TableName() string {
	return "holiday_lists"
}

// Holiday represents a single holiday
type Holiday struct {
	gorm.Model
	ParentID         uint   `gorm:"index" json:"parent_id"`
	HolidayDate      string `gorm:"type:date;not null" json:"holiday_date"`
	Description      string `gorm:"size:255;not null" json:"description"`
	WeeklyOff        bool   `gorm:"default:false" json:"weekly_off"`
}

// TableName specifies the table name
func (Holiday) TableName() string {
	return "holidays"
}
