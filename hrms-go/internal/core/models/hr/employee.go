package hr

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// Employee represents an employee in the organization
type Employee struct {
	base.BaseModel

	// Basic Information
	EmployeeNumber   string     `gorm:"uniqueIndex;size:50" json:"employee_number"`
	EmployeeName     string     `gorm:"size:255;not null;index" json:"employee_name"`
	FirstName        string     `gorm:"size:150;not null" json:"first_name"`
	MiddleName       string     `gorm:"size:150" json:"middle_name"`
	LastName         string     `gorm:"size:150" json:"last_name"`
	NamingSeries     string     `gorm:"size:50" json:"naming_series"`

	// Personal Details
	Gender           string     `gorm:"size:20" json:"gender"` // Male, Female, Other
	DateOfBirth      *time.Time `json:"date_of_birth"`
	DateOfJoining    *time.Time `gorm:"not null" json:"date_of_joining"`
	Status           string     `gorm:"size:50;default:'Active';index" json:"status"` // Active, Inactive, Suspended, Left

	// Contact Information
	CellNumber       string     `gorm:"size:20" json:"cell_number"`
	PersonalEmail    string     `gorm:"size:255" json:"personal_email"`
	CompanyEmail     string     `gorm:"size:255;index" json:"company_email"`
	PreferedEmail    string     `gorm:"size:255" json:"prefered_email"`

	// Emergency Contact
	EmergencyContactName   string `gorm:"size:255" json:"emergency_contact_name"`
	EmergencyContactPhone  string `gorm:"size:20" json:"emergency_contact_phone"`
	EmergencyContactRelation string `gorm:"size:100" json:"emergency_contact_relation"`

	// Current Address
	CurrentAddress   string     `gorm:"type:text" json:"current_address"`
	CurrentCity      string     `gorm:"size:100" json:"current_city"`
	CurrentState     string     `gorm:"size:100" json:"current_state"`
	CurrentCountry   string     `gorm:"size:100" json:"current_country"`
	CurrentPostalCode string    `gorm:"size:20" json:"current_postal_code"`

	// Permanent Address
	PermanentAddress string     `gorm:"type:text" json:"permanent_address"`
	PermanentCity    string     `gorm:"size:100" json:"permanent_city"`
	PermanentState   string     `gorm:"size:100" json:"permanent_state"`
	PermanentCountry string     `gorm:"size:100" json:"permanent_country"`
	PermanentPostalCode string  `gorm:"size:20" json:"permanent_postal_code"`

	// Organization Links
	Company          string     `gorm:"size:255;not null;index" json:"company"`
	Department       string     `gorm:"size:255;index" json:"department"`
	Designation      string     `gorm:"size:255;index" json:"designation"`
	Branch           string     `gorm:"size:255" json:"branch"`
	Grade            string     `gorm:"size:255" json:"grade"`
	EmploymentType   string     `gorm:"size:255" json:"employment_type"` // Full-time, Part-time, Contract, Intern

	// Reporting
	ReportsTo        string     `gorm:"size:255" json:"reports_to"` // Manager employee number

	// Attendance & Shift
	DefaultShift     string     `gorm:"size:255" json:"default_shift"`
	HolidayList      string     `gorm:"size:255" json:"holiday_list"`
	AttendanceDeviceID string   `gorm:"size:50" json:"attendance_device_id"`

	// Leave
	LeaveApprover    string     `gorm:"size:255" json:"leave_approver"`
	LeaveEncashmentAmount float64 `gorm:"type:decimal(18,2)" json:"leave_encashment_amount_per_day"`

	// Expense
	ExpenseApprover  string     `gorm:"size:255" json:"expense_approver"`

	// Shift Request
	ShiftRequestApprover string `gorm:"size:255" json:"shift_request_approver"`

	// User Account
	UserID           string     `gorm:"size:255;index" json:"user_id"` // Link to User
	CreateUserPermission bool   `gorm:"default:false" json:"create_user_permission"`

	// Employment Details
	DateOfRetirement *time.Time `json:"date_of_retirement"`
	ContractEndDate  *time.Time `json:"contract_end_date"`
	NoticeDays       int        `gorm:"default:30" json:"notice_number_of_days"`

	// Onboarding/Exit
	JobApplicant     string     `gorm:"size:255" json:"job_applicant"`
	ScheduledConfirmationDate *time.Time `json:"scheduled_confirmation_date"`
	FinalConfirmationDate     *time.Time `json:"final_confirmation_date"`

	// Salary Information
	PayrollCostCenter string    `gorm:"size:255" json:"payroll_cost_center"`

	// Identification
	PassportNumber   string     `gorm:"size:50" json:"passport_number"`
	DateOfIssue      *time.Time `json:"date_of_issue"`
	ValidUpto        *time.Time `json:"valid_upto"`
	PlaceOfIssue     string     `gorm:"size:100" json:"place_of_issue"`

	// National ID (Ghana Voter ID, etc.)
	NationalIDNumber string     `gorm:"size:50" json:"national_id_number"`

	// Tax Information (Ghana specific)
	TIN              string     `gorm:"size:50" json:"tin"` // Tax Identification Number
	SSNITNumber      string     `gorm:"size:50" json:"ssnit_number"` // Social Security Number

	// Bank Information
	BankName         string     `gorm:"size:255" json:"bank_name"`
	BankACNumber     string     `gorm:"size:50" json:"bank_ac_no"`
	BankBranch       string     `gorm:"size:255" json:"bank_branch"`
	IBANNumber       string     `gorm:"size:50" json:"iban"`

	// Bio & Image
	Bio              string     `gorm:"type:text" json:"bio"`
	Image            string     `gorm:"type:text" json:"image"`

	// Health Insurance
	HealthInsuranceProvider string `gorm:"size:255" json:"health_insurance_provider"`
	HealthInsuranceNo       string `gorm:"size:50" json:"health_insurance_no"`

	// Marital Status
	MaritalStatus    string     `gorm:"size:50" json:"marital_status"` // Single, Married, Divorced, Widowed
	BloodGroup       string     `gorm:"size:10" json:"blood_group"`

	// Person to be contacted
	PersonToBeContacted string  `gorm:"size:255" json:"person_to_be_contacted"`
	RelationWithCompany string  `gorm:"size:100" json:"relation"`

	// Relieving Details
	RelievingDate    *time.Time `json:"relieving_date"`
	ReasonForLeaving string     `gorm:"type:text" json:"reason_for_leaving"`
	LeaveEncashed    bool       `gorm:"default:false" json:"leave_encashed"`
	Encashment_date  *time.Time `json:"encashment_date"`
	HeldOn           *time.Time `json:"held_on"`

	// Additional Details (can be JSON or separate tables)
	// Education, Work History, Skills, etc. would be in separate tables

	// Unsubscribed (from notifications)
	Unsubscribed     bool       `gorm:"default:false" json:"unsubscribed"`

	// Prefered Contact Email
	PreferedContactEmail string `gorm:"size:255" json:"prefered_contact_email"`
}

// TableName specifies the table name
func (Employee) TableName() string {
	return "employees"
}

// BeforeCreate hook
func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	// Call base model hook
	if err := e.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Set employee name if not set
	if e.EmployeeName == "" {
		e.EmployeeName = e.FirstName
		if e.MiddleName != "" {
			e.EmployeeName += " " + e.MiddleName
		}
		if e.LastName != "" {
			e.EmployeeName += " " + e.LastName
		}
	}

	// Set status to Active if not set
	if e.Status == "" {
		e.Status = "Active"
	}

	return nil
}

// IsActive checks if employee is active
func (e *Employee) IsActive() bool {
	return e.Status == "Active"
}

// GetFullName returns the full name of the employee
func (e *Employee) GetFullName() string {
	if e.EmployeeName != "" {
		return e.EmployeeName
	}

	name := e.FirstName
	if e.MiddleName != "" {
		name += " " + e.MiddleName
	}
	if e.LastName != "" {
		name += " " + e.LastName
	}
	return name
}

// EmployeeEducation represents educational qualifications
type EmployeeEducation struct {
	gorm.Model
	EmployeeID       uint      `gorm:"not null;index" json:"employee_id"`
	Employee         Employee  `gorm:"foreignKey:EmployeeID" json:"-"`

	Qualification    string    `gorm:"size:255;not null" json:"school_univ"`
	Level            string    `gorm:"size:100" json:"level"` // Graduate, Post Graduate, etc.
	YearOfPassing    int       `json:"year_of_passing"`
	Class            string    `gorm:"size:50" json:"class_per"` // Grade/Percentage
	MajorOptional    string    `gorm:"size:255" json:"maj_opt_subj"`
}

// TableName specifies the table name
func (EmployeeEducation) TableName() string {
	return "employee_education"
}

// EmployeeExternalWorkHistory represents previous work experience
type EmployeeExternalWorkHistory struct {
	gorm.Model
	EmployeeID       uint      `gorm:"not null;index" json:"employee_id"`
	Employee         Employee  `gorm:"foreignKey:EmployeeID" json:"-"`

	CompanyName      string    `gorm:"size:255;not null" json:"company_name"`
	Designation      string    `gorm:"size:255" json:"designation"`
	Salary           float64   `gorm:"type:decimal(18,2)" json:"salary"`
	Address          string    `gorm:"type:text" json:"address"`
	Contact          string    `gorm:"size:255" json:"contact"`
	FromDate         *time.Time `json:"from_date"`
	ToDate           *time.Time `json:"to_date"`
	TotalExperience  float64   `gorm:"type:decimal(5,2)" json:"total_experience"` // in years
}

// TableName specifies the table name
func (EmployeeExternalWorkHistory) TableName() string {
	return "employee_external_work_history"
}

// EmployeeInternalWorkHistory represents work history within the company
type EmployeeInternalWorkHistory struct {
	gorm.Model
	EmployeeID       uint      `gorm:"not null;index" json:"employee_id"`
	Employee         Employee  `gorm:"foreignKey:EmployeeID" json:"-"`

	Department       string    `gorm:"size:255" json:"department"`
	Designation      string    `gorm:"size:255" json:"designation"`
	Branch           string    `gorm:"size:255" json:"branch"`
	FromDate         *time.Time `gorm:"not null" json:"from_date"`
	ToDate           *time.Time `json:"to_date"`
}

// TableName specifies the table name
func (EmployeeInternalWorkHistory) TableName() string {
	return "employee_internal_work_history"
}

// EmployeeSkill represents skills possessed by employee
type EmployeeSkill struct {
	gorm.Model
	EmployeeID       uint      `gorm:"not null;index" json:"employee_id"`
	Employee         Employee  `gorm:"foreignKey:EmployeeID" json:"-"`

	SkillName        string    `gorm:"size:255;not null" json:"skill"`
	Proficiency      string    `gorm:"size:50" json:"proficiency"` // Beginner, Intermediate, Expert
	YearsOfExperience float64  `gorm:"type:decimal(5,2)" json:"years_of_experience"`
	EvaluationDate   *time.Time `json:"evaluation_date"`
}

// TableName specifies the table name
func (EmployeeSkill) TableName() string {
	return "employee_skills"
}
