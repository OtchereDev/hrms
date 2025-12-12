package hr

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// Attendance represents employee attendance record
type Attendance struct {
	base.BaseModel

	// Employee Information
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	AttendanceDate   time.Time  `gorm:"type:date;not null;index" json:"attendance_date"`

	// Status
	Status           string     `gorm:"size:50;not null;index" json:"status"` // Present, Absent, On Leave, Half Day, Work From Home

	// Timing
	InTime           *time.Time `json:"in_time"`
	OutTime          *time.Time `json:"out_time"`
	WorkingHours     float64    `gorm:"type:decimal(5,2);default:0" json:"working_hours"`

	// Shift Information
	Shift            string     `gorm:"size:255" json:"shift"`

	// Late Entry / Early Exit
	LateEntry        bool       `gorm:"default:false" json:"late_entry"`
	EarlyExit        bool       `gorm:"default:false" json:"early_exit"`

	// Leave Reference
	LeaveType        string     `gorm:"size:255" json:"leave_type"`
	LeaveApplication string     `gorm:"size:255" json:"leave_application"`

	// Organization
	Company          string     `gorm:"size:255;not null;index" json:"company"`
	Department       string     `gorm:"size:255" json:"department"`

	// Remarks
	Remarks          string     `gorm:"type:text" json:"remarks"`
}

// TableName specifies the table name
func (Attendance) TableName() string {
	return "attendances"
}

// BeforeCreate hook
func (a *Attendance) BeforeCreate(tx *gorm.DB) error {
	if err := a.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Calculate working hours if both in_time and out_time are set
	if a.InTime != nil && a.OutTime != nil {
		duration := a.OutTime.Sub(*a.InTime)
		a.WorkingHours = duration.Hours()
	}

	return nil
}

// ShiftType represents a work shift template
type ShiftType struct {
	base.BaseModel
	ShiftName        string     `gorm:"size:255;not null;uniqueIndex" json:"name"`

	// Timing
	StartTime        string     `gorm:"size:10;not null" json:"start_time"` // HH:MM format
	EndTime          string     `gorm:"size:10;not null" json:"end_time"`   // HH:MM format

	// Grace Period & Breaks
	BeginCheckInBeforeShiftStartTime int `gorm:"default:60" json:"begin_check_in_before_shift_start_time"` // minutes
	AllowCheckOutAfterShiftEndTime   int `gorm:"default:60" json:"allow_check_out_after_shift_end_time"`   // minutes

	// Working Hours
	WorkingHoursThresholdForHalfDay  float64 `gorm:"type:decimal(5,2)" json:"working_hours_threshold_for_half_day"`
	WorkingHoursThresholdForAbsent   float64 `gorm:"type:decimal(5,2)" json:"working_hours_threshold_for_absent"`

	// Auto Attendance
	EnableAutoAttendance             bool    `gorm:"default:false" json:"enable_auto_attendance"`
	DetermineCheckInAndCheckOut      string  `gorm:"size:50;default:'Alternating entries as IN and OUT during the same shift'" json:"determine_check_in_and_check_out"`
	ProcessAttendanceAfter           string  `gorm:"size:10" json:"process_attendance_after"` // HH:MM
	LastSyncOfCheckin                *time.Time `json:"last_sync_of_checkin"`

	// Holiday List
	HolidayList                      string  `gorm:"size:255" json:"holiday_list"`

	// Color Coding (for UI)
	Color                            string  `gorm:"size:20" json:"color"`
}

// TableName specifies the table name
func (ShiftType) TableName() string {
	return "shift_types"
}

// ShiftAssignment assigns a shift to an employee
type ShiftAssignment struct {
	base.BaseModel

	// Employee & Shift
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	ShiftType        string     `gorm:"size:255;not null;index" json:"shift_type"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Date Range
	FromDate         time.Time  `gorm:"type:date;not null" json:"from_date"`
	ToDate           *time.Time `gorm:"type:date" json:"to_date"`

	// Status
	Status           string     `gorm:"size:50;default:'Active'" json:"status"` // Active, Inactive

	// Approval
	ApprovedBy       string     `gorm:"size:255" json:"approved_by"`
	ApprovalDate     *time.Time `json:"approval_date"`
}

// TableName specifies the table name
func (ShiftAssignment) TableName() string {
	return "shift_assignments"
}

// EmployeeCheckin represents employee check-in/check-out logs
type EmployeeCheckin struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`

	// Timestamp
	Time             time.Time  `gorm:"not null;index" json:"time"`
	LogType          string     `gorm:"size:10;not null" json:"log_type"` // IN, OUT

	// Device Information
	DeviceID         string     `gorm:"size:255" json:"device_id"`

	// Location (if GPS tracking enabled)
	Latitude         float64    `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude        float64    `gorm:"type:decimal(11,8)" json:"longitude"`

	// Linked Attendance
	Attendance       string     `gorm:"size:255" json:"attendance"`

	// Shift
	Shift            string     `gorm:"size:255" json:"shift"`
	ShiftActualStart *time.Time `json:"shift_actual_start"`
	ShiftActualEnd   *time.Time `json:"shift_actual_end"`

	// Skip Auto Attendance
	SkipAutoAttendance bool     `gorm:"default:false" json:"skip_auto_attendance"`
}

// TableName specifies the table name
func (EmployeeCheckin) TableName() string {
	return "employee_checkins"
}

// AttendanceRequest represents a request to mark/change attendance
type AttendanceRequest struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Date Range
	FromDate         time.Time  `gorm:"type:date;not null" json:"from_date"`
	ToDate           time.Time  `gorm:"type:date;not null" json:"to_date"`

	// Request Details
	HalfDay          bool       `gorm:"default:false" json:"half_day"`
	HalfDayDate      *time.Time `gorm:"type:date" json:"half_day_date"`
	Reason           string     `gorm:"type:text;not null" json:"reason"`
	Explanation      string     `gorm:"type:text" json:"explanation"`

	// Workflow
	WorkflowState    string     `gorm:"size:50" json:"workflow_state"` // Draft, Pending, Approved, Rejected
	Approver         string     `gorm:"size:255" json:"approver"`
	ApprovalDate     *time.Time `json:"approval_date"`
}

// TableName specifies the table name
func (AttendanceRequest) TableName() string {
	return "attendance_requests"
}

// ShiftRequest represents a request for shift change
type ShiftRequest struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Shift Details
	FromDate         time.Time  `gorm:"type:date;not null" json:"from_date"`
	ToDate           time.Time  `gorm:"type:date;not null" json:"to_date"`
	ShiftType        string     `gorm:"size:255;not null" json:"shift_type"`

	// Request Details
	Reason           string     `gorm:"type:text" json:"reason"`

	// Workflow
	WorkflowState    string     `gorm:"size:50" json:"workflow_state"` // Draft, Pending, Approved, Rejected
	Approver         string     `gorm:"size:255" json:"approver"`
	ApprovalDate     *time.Time `json:"approval_date"`
}

// TableName specifies the table name
func (ShiftRequest) TableName() string {
	return "shift_requests"
}

// AttendanceDeviceSettings represents biometric device settings
type AttendanceDeviceSettings struct {
	gorm.Model
	DeviceID         string     `gorm:"uniqueIndex;size:255;not null" json:"device_id"`
	DeviceName       string     `gorm:"size:255" json:"device_name"`
	IPAddress        string     `gorm:"size:50" json:"ip_address"`
	Port             int        `gorm:"default:4370" json:"port"`
	Enabled          bool       `gorm:"default:true" json:"enabled"`
	LastSync         *time.Time `json:"last_sync"`
}

// TableName specifies the table name
func (AttendanceDeviceSettings) TableName() string {
	return "attendance_device_settings"
}
