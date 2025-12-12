package attendance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"gorm.io/gorm"
)

var (
	ErrAttendanceNotFound       = errors.New("attendance record not found")
	ErrAttendanceExists         = errors.New("attendance already exists for this date")
	ErrInvalidAttendanceData    = errors.New("invalid attendance data")
	ErrEmployeeRequired         = errors.New("employee is required")
	ErrAttendanceDateRequired   = errors.New("attendance date is required")
	ErrInvalidStatus            = errors.New("invalid attendance status")
	ErrCheckinNotFound          = errors.New("checkin record not found")
	ErrAlreadyCheckedIn         = errors.New("employee already checked in")
	ErrNotCheckedIn             = errors.New("employee not checked in")
	ErrAttendanceRequestNotFound = errors.New("attendance request not found")
	ErrInvalidWorkflowState     = errors.New("invalid workflow state")
)

// AttendanceService handles attendance business logic
type AttendanceService struct {
	attendanceRepo *repositories.AttendanceRepository
	checkinRepo    *repositories.EmployeeCheckinRepository
	requestRepo    *repositories.AttendanceRequestRepository
	db             *gorm.DB
}

// NewAttendanceService creates a new attendance service
func NewAttendanceService(db *gorm.DB) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: repositories.NewAttendanceRepository(db),
		checkinRepo:    repositories.NewEmployeeCheckinRepository(db),
		requestRepo:    repositories.NewAttendanceRequestRepository(db),
		db:             db,
	}
}

// MarkAttendanceRequest represents a request to mark attendance
type MarkAttendanceRequest struct {
	Employee        string     `json:"employee" validate:"required"`
	AttendanceDate  time.Time  `json:"attendance_date" validate:"required"`
	Status          string     `json:"status" validate:"required"` // Present, Absent, Half Day, Work From Home, On Leave
	Company         string     `json:"company" validate:"required"`
	Department      string     `json:"department"`
	Shift           string     `json:"shift"`
	InTime          *time.Time `json:"in_time"`
	OutTime         *time.Time `json:"out_time"`
	LateEntry       bool       `json:"late_entry"`
	EarlyExit       bool       `json:"early_exit"`
	WorkingHours    float64    `json:"working_hours"`
	OvertimeHours   float64    `json:"overtime_hours"`
	LeaveApplication string    `json:"leave_application"`
}

// UpdateAttendanceRequest represents a request to update attendance
type UpdateAttendanceRequest struct {
	Status        *string    `json:"status"`
	InTime        *time.Time `json:"in_time"`
	OutTime       *time.Time `json:"out_time"`
	LateEntry     *bool      `json:"late_entry"`
	EarlyExit     *bool      `json:"early_exit"`
	WorkingHours  *float64   `json:"working_hours"`
	OvertimeHours *float64   `json:"overtime_hours"`
}

// CheckinRequest represents a check-in/check-out request
type CheckinRequest struct {
	Employee   string  `json:"employee" validate:"required"`
	LogType    string  `json:"log_type" validate:"required"` // IN, OUT
	Time       time.Time `json:"time"`
	DeviceID   string  `json:"device_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

// AttendanceRequestRequest represents a request to create attendance request
type AttendanceRequestRequest struct {
	Employee    string    `json:"employee" validate:"required"`
	Company     string    `json:"company" validate:"required"`
	FromDate    time.Time `json:"from_date" validate:"required"`
	ToDate      time.Time `json:"to_date" validate:"required"`
	HalfDay     bool      `json:"half_day"`
	HalfDayDate *time.Time `json:"half_day_date"`
	Reason      string    `json:"reason" validate:"required"`
	Explanation string    `json:"explanation"`
}

// MarkAttendance marks attendance for an employee
func (s *AttendanceService) MarkAttendance(ctx context.Context, req *MarkAttendanceRequest) (*hr.Attendance, error) {
	// Validate request
	if err := s.validateMarkAttendanceRequest(req); err != nil {
		return nil, err
	}

	// Check if attendance already exists for this date
	existing, err := s.attendanceRepo.GetByEmployeeAndDate(ctx, req.Employee, req.AttendanceDate)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing attendance: %w", err)
	}
	if existing != nil {
		return nil, ErrAttendanceExists
	}

	// Create attendance record
	attendance := &hr.Attendance{
		Employee:         req.Employee,
		AttendanceDate:   req.AttendanceDate,
		Status:           req.Status,
		Company:          req.Company,
		Department:       req.Department,
		Shift:            req.Shift,
		InTime:           req.InTime,
		OutTime:          req.OutTime,
		LateEntry:        req.LateEntry,
		EarlyExit:        req.EarlyExit,
		WorkingHours:     req.WorkingHours,
		OvertimeHours:    req.OvertimeHours,
		LeaveApplication: req.LeaveApplication,
	}

	if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
		return nil, fmt.Errorf("failed to create attendance: %w", err)
	}

	return attendance, nil
}

// GetAttendanceByID retrieves an attendance record by ID
func (s *AttendanceService) GetAttendanceByID(ctx context.Context, id uint) (*hr.Attendance, error) {
	attendance, err := s.attendanceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}
	return attendance, nil
}

// UpdateAttendance updates an attendance record
func (s *AttendanceService) UpdateAttendance(ctx context.Context, id uint, req *UpdateAttendanceRequest) (*hr.Attendance, error) {
	attendance, err := s.attendanceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceNotFound
		}
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}

	// Update fields if provided
	if req.Status != nil {
		if err := s.validateStatus(*req.Status); err != nil {
			return nil, err
		}
		attendance.Status = *req.Status
	}
	if req.InTime != nil {
		attendance.InTime = req.InTime
	}
	if req.OutTime != nil {
		attendance.OutTime = req.OutTime
	}
	if req.LateEntry != nil {
		attendance.LateEntry = *req.LateEntry
	}
	if req.EarlyExit != nil {
		attendance.EarlyExit = *req.EarlyExit
	}
	if req.WorkingHours != nil {
		attendance.WorkingHours = *req.WorkingHours
	}
	if req.OvertimeHours != nil {
		attendance.OvertimeHours = *req.OvertimeHours
	}

	if err := s.attendanceRepo.Update(ctx, attendance); err != nil {
		return nil, fmt.Errorf("failed to update attendance: %w", err)
	}

	return attendance, nil
}

// ListAttendance retrieves attendance records with filters
func (s *AttendanceService) ListAttendance(ctx context.Context, filters repositories.AttendanceFilters, page, pageSize int) ([]hr.Attendance, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.attendanceRepo.List(ctx, filters, page, pageSize)
}

// GetMonthlyAttendance retrieves monthly attendance summary for an employee
func (s *AttendanceService) GetMonthlyAttendance(ctx context.Context, employee string, year, month int) ([]hr.Attendance, error) {
	if employee == "" {
		return nil, ErrEmployeeRequired
	}
	if year < 2000 || year > 2100 {
		return nil, fmt.Errorf("invalid year: %d", year)
	}
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month: %d", month)
	}

	return s.attendanceRepo.GetMonthlyAttendance(ctx, employee, year, month)
}

// Checkin performs employee check-in or check-out
func (s *AttendanceService) Checkin(ctx context.Context, req *CheckinRequest) (*hr.EmployeeCheckin, error) {
	// Validate request
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}
	if req.LogType != "IN" && req.LogType != "OUT" {
		return nil, fmt.Errorf("invalid log type: must be IN or OUT")
	}

	// Set time if not provided
	if req.Time.IsZero() {
		req.Time = time.Now().UTC()
	}

	// If checking out, verify employee is checked in
	if req.LogType == "OUT" {
		lastCheckin, err := s.checkinRepo.GetLastCheckinForEmployee(ctx, req.Employee)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrNotCheckedIn
			}
			return nil, fmt.Errorf("failed to get last checkin: %w", err)
		}
		if lastCheckin.LogType == "OUT" {
			return nil, ErrNotCheckedIn
		}
	}

	// Create checkin record
	checkin := &hr.EmployeeCheckin{
		Employee:       req.Employee,
		Time:           req.Time,
		LogType:        req.LogType,
		DeviceID:       req.DeviceID,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
	}

	if err := s.checkinRepo.Create(ctx, checkin); err != nil {
		return nil, fmt.Errorf("failed to create checkin: %w", err)
	}

	// Auto-mark attendance based on checkins
	if req.LogType == "OUT" {
		if err := s.autoMarkAttendanceFromCheckins(ctx, req.Employee, req.Time); err != nil {
			// Log error but don't fail the checkin
			fmt.Printf("Warning: failed to auto-mark attendance: %v\n", err)
		}
	}

	return checkin, nil
}

// GetTodayCheckinsForEmployee retrieves today's checkins for an employee
func (s *AttendanceService) GetTodayCheckinsForEmployee(ctx context.Context, employee string) ([]hr.EmployeeCheckin, error) {
	if employee == "" {
		return nil, ErrEmployeeRequired
	}

	return s.checkinRepo.GetTodayCheckinsForEmployee(ctx, employee)
}

// CreateAttendanceRequest creates a new attendance request
func (s *AttendanceService) CreateAttendanceRequest(ctx context.Context, req *AttendanceRequestRequest) (*hr.AttendanceRequest, error) {
	// Validate request
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}
	if req.FromDate.IsZero() || req.ToDate.IsZero() {
		return nil, fmt.Errorf("from_date and to_date are required")
	}
	if req.FromDate.After(req.ToDate) {
		return nil, fmt.Errorf("from_date cannot be after to_date")
	}

	// Create attendance request
	attendanceReq := &hr.AttendanceRequest{
		Employee:    req.Employee,
		Company:     req.Company,
		FromDate:    req.FromDate,
		ToDate:      req.ToDate,
		HalfDay:     req.HalfDay,
		HalfDayDate: req.HalfDayDate,
		Reason:      req.Reason,
		Explanation: req.Explanation,
		WorkflowState: "Draft",
	}

	if err := s.requestRepo.Create(ctx, attendanceReq); err != nil {
		return nil, fmt.Errorf("failed to create attendance request: %w", err)
	}

	return attendanceReq, nil
}

// ApproveAttendanceRequest approves an attendance request
func (s *AttendanceService) ApproveAttendanceRequest(ctx context.Context, id uint, approver string) (*hr.AttendanceRequest, error) {
	req, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceRequestNotFound
		}
		return nil, fmt.Errorf("failed to get attendance request: %w", err)
	}

	if req.WorkflowState != "Pending" && req.WorkflowState != "Draft" {
		return nil, ErrInvalidWorkflowState
	}

	req.WorkflowState = "Approved"

	if err := s.requestRepo.Update(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to approve attendance request: %w", err)
	}

	return req, nil
}

// RejectAttendanceRequest rejects an attendance request
func (s *AttendanceService) RejectAttendanceRequest(ctx context.Context, id uint, approver string) (*hr.AttendanceRequest, error) {
	req, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceRequestNotFound
		}
		return nil, fmt.Errorf("failed to get attendance request: %w", err)
	}

	if req.WorkflowState != "Pending" && req.WorkflowState != "Draft" {
		return nil, ErrInvalidWorkflowState
	}

	req.WorkflowState = "Rejected"

	if err := s.requestRepo.Update(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to reject attendance request: %w", err)
	}

	return req, nil
}

// Helper methods

func (s *AttendanceService) validateMarkAttendanceRequest(req *MarkAttendanceRequest) error {
	if req.Employee == "" {
		return ErrEmployeeRequired
	}
	if req.AttendanceDate.IsZero() {
		return ErrAttendanceDateRequired
	}
	if req.Company == "" {
		return fmt.Errorf("company is required")
	}
	return s.validateStatus(req.Status)
}

func (s *AttendanceService) validateStatus(status string) error {
	validStatuses := []string{"Present", "Absent", "Half Day", "Work From Home", "On Leave"}
	for _, valid := range validStatuses {
		if status == valid {
			return nil
		}
	}
	return ErrInvalidStatus
}

// autoMarkAttendanceFromCheckins automatically marks attendance based on check-ins/outs
func (s *AttendanceService) autoMarkAttendanceFromCheckins(ctx context.Context, employee string, checkoutTime time.Time) error {
	// Get today's checkins
	checkins, err := s.checkinRepo.GetTodayCheckinsForEmployee(ctx, employee)
	if err != nil {
		return err
	}

	if len(checkins) < 2 {
		return nil // Need at least check-in and check-out
	}

	// Find first IN and last OUT
	var firstIn, lastOut *hr.EmployeeCheckin
	for i := range checkins {
		if checkins[i].LogType == "IN" && (firstIn == nil || checkins[i].Time.Before(firstIn.Time)) {
			firstIn = &checkins[i]
		}
		if checkins[i].LogType == "OUT" && (lastOut == nil || checkins[i].Time.After(lastOut.Time)) {
			lastOut = &checkins[i]
		}
	}

	if firstIn == nil || lastOut == nil {
		return nil
	}

	// Calculate working hours
	workingHours := lastOut.Time.Sub(firstIn.Time).Hours()

	// Check if attendance already exists
	attendanceDate := checkoutTime.UTC().Truncate(24 * time.Hour)
	existing, err := s.attendanceRepo.GetByEmployeeAndDate(ctx, employee, attendanceDate)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existing != nil {
		// Update existing attendance
		existing.InTime = &firstIn.Time
		existing.OutTime = &lastOut.Time
		existing.WorkingHours = workingHours
		existing.Status = "Present"
		return s.attendanceRepo.Update(ctx, existing)
	}

	// Create new attendance
	attendance := &hr.Attendance{
		Employee:       employee,
		AttendanceDate: attendanceDate,
		Status:         "Present",
		InTime:         &firstIn.Time,
		OutTime:        &lastOut.Time,
		WorkingHours:   workingHours,
	}

	return s.attendanceRepo.Create(ctx, attendance)
}
