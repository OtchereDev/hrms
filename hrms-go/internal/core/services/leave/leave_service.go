package leave

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
	ErrLeaveApplicationNotFound = errors.New("leave application not found")
	ErrLeaveAllocationNotFound  = errors.New("leave allocation not found")
	ErrLeaveTypeNotFound        = errors.New("leave type not found")
	ErrInsufficientLeaveBalance = errors.New("insufficient leave balance")
	ErrOverlappingLeaves        = errors.New("overlapping leave applications exist")
	ErrInvalidDates             = errors.New("invalid date range")
	ErrEmployeeRequired         = errors.New("employee is required")
	ErrLeaveTypeRequired        = errors.New("leave type is required")
	ErrInvalidStatus            = errors.New("invalid leave status")
)

// LeaveService handles leave management business logic
type LeaveService struct {
	leaveAppRepo    *repositories.LeaveApplicationRepository
	leaveAllocRepo  *repositories.LeaveAllocationRepository
	leaveTypeRepo   *repositories.LeaveTypeRepository
	leaveBalanceRepo *repositories.LeaveBalanceRepository
	encashmentRepo  *repositories.LeaveEncashmentRepository
	db              *gorm.DB
}

// NewLeaveService creates a new leave service
func NewLeaveService(db *gorm.DB) *LeaveService {
	return &LeaveService{
		leaveAppRepo:     repositories.NewLeaveApplicationRepository(db),
		leaveAllocRepo:   repositories.NewLeaveAllocationRepository(db),
		leaveTypeRepo:    repositories.NewLeaveTypeRepository(db),
		leaveBalanceRepo: repositories.NewLeaveBalanceRepository(db),
		encashmentRepo:   repositories.NewLeaveEncashmentRepository(db),
		db:               db,
	}
}

// ApplyLeaveRequest represents a leave application request
type ApplyLeaveRequest struct {
	Employee         string     `json:"employee" validate:"required"`
	LeaveType        string     `json:"leave_type" validate:"required"`
	FromDate         time.Time  `json:"from_date" validate:"required"`
	ToDate           time.Time  `json:"to_date" validate:"required"`
	HalfDay          bool       `json:"half_day"`
	HalfDayDate      *time.Time `json:"half_day_date"`
	TotalLeaveDays   float64    `json:"total_leave_days"`
	Company          string     `json:"company" validate:"required"`
	Description      string     `json:"description"`
	LeaveApprover    string     `json:"leave_approver"`
	FollowViaEmail   bool       `json:"follow_via_email"`
}

// AllocateLeaveRequest represents a leave allocation request
type AllocateLeaveRequest struct {
	Employee              string    `json:"employee" validate:"required"`
	LeaveType             string    `json:"leave_type" validate:"required"`
	FromDate              time.Time `json:"from_date" validate:"required"`
	ToDate                time.Time `json:"to_date" validate:"required"`
	NewLeavesAllocated    float64   `json:"new_leaves_allocated" validate:"required"`
	Company               string    `json:"company" validate:"required"`
	Description           string    `json:"description"`
	AddUnusedLeaves       bool      `json:"add_unused_leaves"`
	UnusedLeaves          float64   `json:"unused_leaves"`
}

// UpdateLeaveApplicationRequest represents an update request
type UpdateLeaveApplicationRequest struct {
	LeaveApprover  *string    `json:"leave_approver"`
	Description    *string    `json:"description"`
	TotalLeaveDays *float64   `json:"total_leave_days"`
}

// ApplyLeave creates a new leave application
func (s *LeaveService) ApplyLeave(ctx context.Context, req *ApplyLeaveRequest) (*hr.LeaveApplication, error) {
	// Validate request
	if err := s.validateApplyLeaveRequest(req); err != nil {
		return nil, err
	}

	// Check leave balance
	balance, err := s.leaveBalanceRepo.GetLeaveBalance(ctx, req.Employee, req.LeaveType)
	if err != nil {
		return nil, fmt.Errorf("failed to get leave balance: %w", err)
	}

	if balance < req.TotalLeaveDays {
		return nil, ErrInsufficientLeaveBalance
	}

	// Check for overlapping leaves
	overlapping, err := s.leaveAppRepo.GetOverlappingLeaves(ctx, req.Employee, req.FromDate, req.ToDate, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to check overlapping leaves: %w", err)
	}
	if len(overlapping) > 0 {
		return nil, ErrOverlappingLeaves
	}

	// Create leave application
	leave := &hr.LeaveApplication{
		Employee:       req.Employee,
		LeaveType:      req.LeaveType,
		FromDate:       req.FromDate,
		ToDate:         req.ToDate,
		HalfDay:        req.HalfDay,
		HalfDayDate:    req.HalfDayDate,
		TotalLeaveDays: req.TotalLeaveDays,
		Company:        req.Company,
		Description:    req.Description,
		LeaveApprover:  req.LeaveApprover,
		FollowViaEmail: req.FollowViaEmail,
		Status:         "Open",
	}

	if err := s.leaveAppRepo.Create(ctx, leave); err != nil {
		return nil, fmt.Errorf("failed to create leave application: %w", err)
	}

	return leave, nil
}

// GetLeaveApplicationByID retrieves a leave application by ID
func (s *LeaveService) GetLeaveApplicationByID(ctx context.Context, id uint) (*hr.LeaveApplication, error) {
	leave, err := s.leaveAppRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeaveApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get leave application: %w", err)
	}
	return leave, nil
}

// UpdateLeaveApplication updates a leave application
func (s *LeaveService) UpdateLeaveApplication(ctx context.Context, id uint, req *UpdateLeaveApplicationRequest) (*hr.LeaveApplication, error) {
	leave, err := s.leaveAppRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeaveApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get leave application: %w", err)
	}

	// Can only update if status is Open or Rejected
	if leave.Status != "Open" && leave.Status != "Rejected" {
		return nil, fmt.Errorf("cannot update leave application with status: %s", leave.Status)
	}

	// Update fields if provided
	if req.LeaveApprover != nil {
		leave.LeaveApprover = *req.LeaveApprover
	}
	if req.Description != nil {
		leave.Description = *req.Description
	}
	if req.TotalLeaveDays != nil {
		leave.TotalLeaveDays = *req.TotalLeaveDays
	}

	if err := s.leaveAppRepo.Update(ctx, leave); err != nil {
		return nil, fmt.Errorf("failed to update leave application: %w", err)
	}

	return leave, nil
}

// ApproveLeaveApplication approves a leave application
func (s *LeaveService) ApproveLeaveApplication(ctx context.Context, id uint, approver string) (*hr.LeaveApplication, error) {
	leave, err := s.leaveAppRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeaveApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get leave application: %w", err)
	}

	if leave.Status != "Open" {
		return nil, fmt.Errorf("cannot approve leave application with status: %s", leave.Status)
	}

	// Verify leave balance again
	balance, err := s.leaveBalanceRepo.GetLeaveBalance(ctx, leave.Employee, leave.LeaveType)
	if err != nil {
		return nil, fmt.Errorf("failed to get leave balance: %w", err)
	}

	if balance < leave.TotalLeaveDays {
		return nil, ErrInsufficientLeaveBalance
	}

	leave.Status = "Approved"

	if err := s.leaveAppRepo.Update(ctx, leave); err != nil {
		return nil, fmt.Errorf("failed to approve leave application: %w", err)
	}

	return leave, nil
}

// RejectLeaveApplication rejects a leave application
func (s *LeaveService) RejectLeaveApplication(ctx context.Context, id uint, approver string, reason string) (*hr.LeaveApplication, error) {
	leave, err := s.leaveAppRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeaveApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get leave application: %w", err)
	}

	if leave.Status != "Open" {
		return nil, fmt.Errorf("cannot reject leave application with status: %s", leave.Status)
	}

	leave.Status = "Rejected"

	if err := s.leaveAppRepo.Update(ctx, leave); err != nil {
		return nil, fmt.Errorf("failed to reject leave application: %w", err)
	}

	return leave, nil
}

// CancelLeaveApplication cancels a leave application
func (s *LeaveService) CancelLeaveApplication(ctx context.Context, id uint) (*hr.LeaveApplication, error) {
	leave, err := s.leaveAppRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeaveApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get leave application: %w", err)
	}

	if leave.Status == "Cancelled" {
		return nil, fmt.Errorf("leave application already cancelled")
	}

	leave.Status = "Cancelled"

	if err := s.leaveAppRepo.Update(ctx, leave); err != nil {
		return nil, fmt.Errorf("failed to cancel leave application: %w", err)
	}

	return leave, nil
}

// ListLeaveApplications retrieves leave applications with filters
func (s *LeaveService) ListLeaveApplications(ctx context.Context, filters repositories.LeaveApplicationFilters, page, pageSize int) ([]hr.LeaveApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.leaveAppRepo.List(ctx, filters, page, pageSize)
}

// AllocateLeave allocates leaves to an employee
func (s *LeaveService) AllocateLeave(ctx context.Context, req *AllocateLeaveRequest) (*hr.LeaveAllocation, error) {
	// Validate request
	if err := s.validateAllocateLeaveRequest(req); err != nil {
		return nil, err
	}

	// Check if leave type exists
	_, err := s.leaveTypeRepo.GetByName(ctx, req.LeaveType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLeaveTypeNotFound
		}
		return nil, fmt.Errorf("failed to get leave type: %w", err)
	}

	// Calculate total leaves
	totalLeaves := req.NewLeavesAllocated
	if req.AddUnusedLeaves {
		totalLeaves += req.UnusedLeaves
	}

	// Create leave allocation
	allocation := &hr.LeaveAllocation{
		Employee:             req.Employee,
		LeaveType:            req.LeaveType,
		FromDate:             req.FromDate,
		ToDate:               req.ToDate,
		NewLeavesAllocated:   req.NewLeavesAllocated,
		Company:              req.Company,
		Description:          req.Description,
		AddUnusedLeaves:      req.AddUnusedLeaves,
		UnusedLeaves:         req.UnusedLeaves,
		TotalLeavesAllocated: totalLeaves,
	}

	if err := s.leaveAllocRepo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("failed to create leave allocation: %w", err)
	}

	return allocation, nil
}

// GetLeaveBalance retrieves leave balance for an employee
func (s *LeaveService) GetLeaveBalance(ctx context.Context, employee, leaveType string) (float64, error) {
	if employee == "" {
		return 0, ErrEmployeeRequired
	}
	if leaveType == "" {
		return 0, ErrLeaveTypeRequired
	}

	return s.leaveBalanceRepo.GetLeaveBalance(ctx, employee, leaveType)
}

// GetAllLeaveBalances retrieves all leave balances for an employee
func (s *LeaveService) GetAllLeaveBalances(ctx context.Context, employee string) (map[string]float64, error) {
	if employee == "" {
		return nil, ErrEmployeeRequired
	}

	return s.leaveBalanceRepo.GetAllLeaveBalances(ctx, employee)
}

// ListLeaveAllocations retrieves leave allocations with filters
func (s *LeaveService) ListLeaveAllocations(ctx context.Context, filters repositories.LeaveAllocationFilters, page, pageSize int) ([]hr.LeaveAllocation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.leaveAllocRepo.List(ctx, filters, page, pageSize)
}

// GetActiveLeaveTypes retrieves all active leave types
func (s *LeaveService) GetActiveLeaveTypes(ctx context.Context) ([]hr.LeaveType, error) {
	return s.leaveTypeRepo.ListActive(ctx)
}

// CreateLeaveEncashment creates a leave encashment request
func (s *LeaveService) CreateLeaveEncashment(ctx context.Context, employee, leaveType string, encashableDays float64) (*hr.LeaveEncashment, error) {
	if employee == "" {
		return nil, ErrEmployeeRequired
	}
	if leaveType == "" {
		return nil, ErrLeaveTypeRequired
	}

	// Verify leave balance
	balance, err := s.leaveBalanceRepo.GetLeaveBalance(ctx, employee, leaveType)
	if err != nil {
		return nil, fmt.Errorf("failed to get leave balance: %w", err)
	}

	if balance < encashableDays {
		return nil, ErrInsufficientLeaveBalance
	}

	// Get leave type to check if encashment is allowed
	lt, err := s.leaveTypeRepo.GetByName(ctx, leaveType)
	if err != nil {
		return nil, fmt.Errorf("failed to get leave type: %w", err)
	}

	if !lt.AllowEncashment {
		return nil, fmt.Errorf("leave type does not allow encashment")
	}

	encashment := &hr.LeaveEncashment{
		Employee:        employee,
		LeaveType:       leaveType,
		EncashmentDate:  time.Now().UTC(),
		EncashableDays:  encashableDays,
	}

	if err := s.encashmentRepo.Create(ctx, encashment); err != nil {
		return nil, fmt.Errorf("failed to create leave encashment: %w", err)
	}

	return encashment, nil
}

// Helper methods

func (s *LeaveService) validateApplyLeaveRequest(req *ApplyLeaveRequest) error {
	if req.Employee == "" {
		return ErrEmployeeRequired
	}
	if req.LeaveType == "" {
		return ErrLeaveTypeRequired
	}
	if req.FromDate.IsZero() || req.ToDate.IsZero() {
		return ErrInvalidDates
	}
	if req.FromDate.After(req.ToDate) {
		return ErrInvalidDates
	}
	if req.TotalLeaveDays <= 0 {
		return fmt.Errorf("total leave days must be greater than zero")
	}
	if req.Company == "" {
		return fmt.Errorf("company is required")
	}
	return nil
}

func (s *LeaveService) validateAllocateLeaveRequest(req *AllocateLeaveRequest) error {
	if req.Employee == "" {
		return ErrEmployeeRequired
	}
	if req.LeaveType == "" {
		return ErrLeaveTypeRequired
	}
	if req.FromDate.IsZero() || req.ToDate.IsZero() {
		return ErrInvalidDates
	}
	if req.FromDate.After(req.ToDate) {
		return ErrInvalidDates
	}
	if req.NewLeavesAllocated <= 0 {
		return fmt.Errorf("new leaves allocated must be greater than zero")
	}
	if req.Company == "" {
		return fmt.Errorf("company is required")
	}
	return nil
}

// GetLeaveDetails returns leave details for an employee on a specific date
func (s *LeaveService) GetLeaveDetails(ctx context.Context, employee string, date time.Time, forSalarySlip bool) (*hr.LeaveApplication, error) {
	filter := repositories.LeaveApplicationFilter{
		Employee:  employee,
		StartDate: &date,
		EndDate:   &date,
		Status:    "Approved",
	}

	leaves, _, err := s.leaveRepo.ListApplications(ctx, filter, 1, 1)
	if err != nil {
		return nil, err
	}

	if len(leaves) == 0 {
		return nil, nil
	}

	return leaves[0], nil
}

// GetNumberOfLeaveDays calculates the number of leave days between dates
func (s *LeaveService) GetNumberOfLeaveDays(ctx context.Context, employee, leaveType string, fromDate, toDate time.Time, halfDay bool, halfDayDate *time.Time) (float64, error) {
	if fromDate.After(toDate) {
		return 0, fmt.Errorf("from_date cannot be after to_date")
	}

	// Calculate total days
	totalDays := toDate.Sub(fromDate).Hours()/24 + 1

	// If half day is specified
	if halfDay && halfDayDate != nil {
		// Check if half day date is within the range
		if (halfDayDate.Equal(fromDate) || halfDayDate.After(fromDate)) &&
		   (halfDayDate.Equal(toDate) || halfDayDate.Before(toDate)) {
			totalDays -= 0.5
		}
	}

	// TODO: Exclude holidays based on employee's holiday list
	// This would require a holiday repository/service

	return totalDays, nil
}

// GetLeaveBalanceOn returns the leave balance for an employee on a specific date
func (s *LeaveService) GetLeaveBalanceOn(ctx context.Context, employee, leaveType string, date time.Time) (float64, error) {
	// Get all allocations for this leave type up to the date
	filter := repositories.LeaveAllocationFilter{
		Employee:  employee,
		LeaveType: leaveType,
		EndDate:   &date,
	}

	allocations, _, err := s.leaveRepo.ListAllocations(ctx, filter, 1, 10000)
	if err != nil {
		return 0, err
	}

	// Sum up total allocated
	var totalAllocated float64
	for _, allocation := range allocations {
		totalAllocated += allocation.NewLeavesAllocated
	}

	// Get all approved leave applications up to the date
	appFilter := repositories.LeaveApplicationFilter{
		Employee:  employee,
		LeaveType: leaveType,
		Status:    "Approved",
		EndDate:   &date,
	}

	applications, _, err := s.leaveRepo.ListApplications(ctx, appFilter, 1, 10000)
	if err != nil {
		return 0, err
	}

	// Calculate total used
	var totalUsed float64
	for _, app := range applications {
		days, _ := s.GetNumberOfLeaveDays(ctx, employee, leaveType, app.FromDate, app.ToDate, app.HalfDay, app.HalfDayDate)
		totalUsed += days
	}

	return totalAllocated - totalUsed, nil
}

// GetLeavesForPeriod returns all approved leaves for an employee in a date range
func (s *LeaveService) GetLeavesForPeriod(ctx context.Context, employee string, fromDate, toDate time.Time) ([]*hr.LeaveApplication, error) {
	filter := repositories.LeaveApplicationFilter{
		Employee:  employee,
		StartDate: &fromDate,
		EndDate:   &toDate,
		Status:    "Approved",
	}

	leaves, _, err := s.leaveRepo.ListApplications(ctx, filter, 1, 10000)
	if err != nil {
		return nil, err
	}

	return leaves, nil
}
