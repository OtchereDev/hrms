package hr

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// LeaveType represents a type of leave (Annual, Sick, Casual, etc.)
type LeaveType struct {
	base.BaseModel
	LeaveTypeName    string  `gorm:"size:255;not null;uniqueIndex" json:"leave_type_name"`

	// Allocation
	MaxLeaves        float64 `gorm:"type:decimal(10,2)" json:"max_leaves_allowed"`
	ApplicablAfter   int     `gorm:"default:0" json:"applicable_after_first_login_days"`
	MaxContinuousDays int    `gorm:"default:0" json:"max_continuous_days_allowed"`

	// Encashment
	IsEncashable     bool    `gorm:"default:false" json:"is_encashable"`
	EncashmentThreshold float64 `gorm:"type:decimal(10,2)" json:"encashment_threshold_days"`

	// Carry Forward
	IsCarryForward   bool    `gorm:"default:false" json:"is_carry_forward"`
	MaximumCarryForwardedLeaves float64 `gorm:"type:decimal(10,2)" json:"maximum_carry_forwarded_leaves"`
	ExpireCarryForwardedLeaves bool `gorm:"default:false" json:"expire_carry_forwarded_leaves_after_days"`
	CarryForwardExpiry int   `gorm:"default:0" json:"carry_forward_expiry_days"`

	// Application Rules
	IsOptional       bool    `gorm:"default:false" json:"is_optional_leave"`
	AllowNegative    bool    `gorm:"default:false" json:"allow_negative_balance"`
	IncludeHoliday   bool    `gorm:"default:true" json:"include_holidays_within_leaves_as_leaves"`

	// Leave Approval
	IsLWP            bool    `gorm:"default:false" json:"is_lwp"` // Leave Without Pay

	// Color for UI
	Color            string  `gorm:"size:20" json:"color"`
}

// TableName specifies the table name
func (LeaveType) TableName() string {
	return "leave_types"
}

// LeavePolicy represents a leave policy template
type LeavePolicy struct {
	base.BaseModel
	PolicyName       string  `gorm:"size:255;not null;uniqueIndex" json:"title"`
	Company          string  `gorm:"size:255" json:"company"`

	// Annual Allocation
	AnnualAllocation []LeavePolicyDetail `gorm:"foreignKey:ParentID" json:"annual_allocation"`
}

// TableName specifies the table name
func (LeavePolicy) TableName() string {
	return "leave_policies"
}

// LeavePolicyDetail represents leave allocation in a policy
type LeavePolicyDetail struct {
	gorm.Model
	ParentID         uint    `gorm:"index" json:"parent_id"`
	LeaveType        string  `gorm:"size:255;not null" json:"leave_type"`
	AnnualAllocation float64 `gorm:"type:decimal(10,2)" json:"annual_allocation"`
}

// TableName specifies the table name
func (LeavePolicyDetail) TableName() string {
	return "leave_policy_details"
}

// LeavePolicyAssignment assigns a leave policy to an employee
type LeavePolicyAssignment struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	LeavePolicy      string     `gorm:"size:255;not null" json:"leave_policy"`

	// Effective Date
	EffectiveFrom    time.Time  `gorm:"type:date;not null" json:"effective_from"`
	EffectiveTo      *time.Time `gorm:"type:date" json:"effective_to"`

	// Assignment Period
	AssignmentBasedOn string    `gorm:"size:50" json:"assignment_based_on"` // Joining Date, Leave Period

	// Carry Forward
	CarryForward     bool       `gorm:"default:false" json:"carry_forward"`
}

// TableName specifies the table name
func (LeavePolicyAssignment) TableName() string {
	return "leave_policy_assignments"
}

// LeaveApplication represents an employee's leave request
type LeaveApplication struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`
	Department       string     `gorm:"size:255" json:"department"`

	// Leave Details
	LeaveType        string     `gorm:"size:255;not null;index" json:"leave_type"`
	FromDate         time.Time  `gorm:"type:date;not null;index" json:"from_date"`
	ToDate           time.Time  `gorm:"type:date;not null;index" json:"to_date"`
	HalfDay          bool       `gorm:"default:false" json:"half_day"`
	HalfDayDate      *time.Time `gorm:"type:date" json:"half_day_date"`
	TotalLeaveDays   float64    `gorm:"type:decimal(10,2)" json:"total_leave_days"`

	// Reason
	Reason           string     `gorm:"type:text;not null" json:"description"`

	// Contact
	LeaveApprover    string     `gorm:"size:255" json:"leave_approver"`
	PostingDate      time.Time  `gorm:"type:date" json:"posting_date"`

	// Status & Workflow
	Status           string     `gorm:"size:50;not null;index;default:'Open'" json:"status"` // Open, Approved, Rejected, Cancelled
	WorkflowState    string     `gorm:"size:50" json:"workflow_state"`

	// Follow Up
	FollowViaEmail   bool       `gorm:"default:false" json:"follow_via_email"`

	// Color for UI
	Color            string     `gorm:"size:20" json:"color"`
}

// TableName specifies the table name
func (LeaveApplication) TableName() string {
	return "leave_applications"
}

// BeforeCreate hook
func (l *LeaveApplication) BeforeCreate(tx *gorm.DB) error {
	if err := l.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Calculate total leave days if not set
	if l.TotalLeaveDays == 0 {
		days := l.ToDate.Sub(l.FromDate).Hours() / 24
		l.TotalLeaveDays = days + 1 // Include both start and end dates

		if l.HalfDay {
			l.TotalLeaveDays -= 0.5
		}
	}

	return nil
}

// LeaveAllocation represents leave balance allocation for an employee
type LeaveAllocation struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Leave Type
	LeaveType        string     `gorm:"size:255;not null;index" json:"leave_type"`

	// Period
	FromDate         time.Time  `gorm:"type:date;not null" json:"from_date"`
	ToDate           time.Time  `gorm:"type:date;not null" json:"to_date"`

	// Allocation
	NewLeavesAllocated float64  `gorm:"type:decimal(10,2);not null" json:"new_leaves_allocated"`
	TotalLeavesAllocated float64 `gorm:"type:decimal(10,2)" json:"total_leaves_allocated"`

	// Carry Forward
	UnusedLeaves     float64    `gorm:"type:decimal(10,2);default:0" json:"unused_leaves"`
	CarryForwarded   bool       `gorm:"default:false" json:"carry_forwarded_leaves"`
	CarryForwardedLeaves float64 `gorm:"type:decimal(10,2);default:0" json:"carry_forwarded_leave_count"`
	ExpiredLeaves    float64    `gorm:"type:decimal(10,2);default:0" json:"expired_leaves"`

	// Description
	Description      string     `gorm:"type:text" json:"description"`
}

// TableName specifies the table name
func (LeaveAllocation) TableName() string {
	return "leave_allocations"
}

// LeaveEncashment represents leave encashment request
type LeaveEncashment struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Leave Details
	LeaveType        string     `gorm:"size:255;not null" json:"leave_type"`
	EncashmentDate   time.Time  `gorm:"type:date;not null" json:"encashment_date"`
	EncashableDays   float64    `gorm:"type:decimal(10,2);not null" json:"encashable_days"`
	EncashedAmount   float64    `gorm:"type:decimal(18,2)" json:"encashment_amount"`

	// Payroll
	SalarySlip       string     `gorm:"size:255" json:"salary_slip"`
	AdditionalSalary string     `gorm:"size:255" json:"additional_salary"`

	// Leave Allocation Reference
	LeaveAllocation  string     `gorm:"size:255" json:"leave_allocation"`

	// Leave Balance
	LeaveBalance     float64    `gorm:"type:decimal(10,2)" json:"leave_balance"`
}

// TableName specifies the table name
func (LeaveEncashment) TableName() string {
	return "leave_encashments"
}

// CompensatoryLeaveRequest represents request for compensatory leave
type CompensatoryLeaveRequest struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Work Details
	WorkFromDate     time.Time  `gorm:"type:date;not null" json:"work_from_date"`
	WorkEndDate      time.Time  `gorm:"type:date;not null" json:"work_end_date"`
	Reason           string     `gorm:"type:text;not null" json:"reason"`

	// Leave Allocation
	LeaveType        string     `gorm:"size:255" json:"leave_type"`
	HalfDayDate      *time.Time `gorm:"type:date" json:"half_day_date"`
	HalfDay          bool       `gorm:"default:false" json:"half_day"`

	// Workflow
	Status           string     `gorm:"size:50;default:'Open'" json:"status"` // Open, Approved, Rejected
	Approver         string     `gorm:"size:255" json:"approver"`
}

// TableName specifies the table name
func (CompensatoryLeaveRequest) TableName() string {
	return "compensatory_leave_requests"
}

// LeaveBlockList represents blocked leave periods
type LeaveBlockList struct {
	base.BaseModel
	LeaveBlockListName string `gorm:"size:255;not null;uniqueIndex" json:"leave_block_list_name"`
	Company            string `gorm:"size:255" json:"company"`
	AppliesToCompany   string `gorm:"size:255" json:"applies_to_company"`

	// Block List Days
	BlockDays          []LeaveBlockListDate `gorm:"foreignKey:ParentID" json:"leave_block_list_dates"`

	// Allowed Employees/Departments
	AllowedDepartments []LeaveBlockListAllowed `gorm:"foreignKey:ParentID" json:"leave_block_list_allowed"`
}

// TableName specifies the table name
func (LeaveBlockList) TableName() string {
	return "leave_block_lists"
}

// LeaveBlockListDate represents blocked dates
type LeaveBlockListDate struct {
	gorm.Model
	ParentID         uint      `gorm:"index" json:"parent_id"`
	BlockDate        time.Time `gorm:"type:date;not null" json:"block_date"`
	Reason           string    `gorm:"size:255" json:"reason"`
}

// TableName specifies the table name
func (LeaveBlockListDate) TableName() string {
	return "leave_block_list_dates"
}

// LeaveBlockListAllowed represents allowed employees/departments during block period
type LeaveBlockListAllowed struct {
	gorm.Model
	ParentID         uint   `gorm:"index" json:"parent_id"`
	AllowUser        string `gorm:"size:255" json:"allow_user"`
}

// TableName specifies the table name
func (LeaveBlockListAllowed) TableName() string {
	return "leave_block_list_allowed"
}

// LeavePeriod represents a leave allocation period
type LeavePeriod struct {
	base.BaseModel
	LeavePeriodName  string    `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Company          string    `gorm:"size:255;not null" json:"company"`
	FromDate         time.Time `gorm:"type:date;not null" json:"from_date"`
	ToDate           time.Time `gorm:"type:date;not null" json:"to_date"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
}

// TableName specifies the table name
func (LeavePeriod) TableName() string {
	return "leave_periods"
}
