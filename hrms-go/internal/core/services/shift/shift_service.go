package shift

import (
	"context"
	"fmt"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"gorm.io/gorm"
)

// ShiftService handles shift business logic
type ShiftService struct {
	shiftRepo *repositories.ShiftRepository
	db        *gorm.DB
}

// NewShiftService creates a new shift service
func NewShiftService(db *gorm.DB) *ShiftService {
	return &ShiftService{
		shiftRepo: repositories.NewShiftRepository(db),
		db:        db,
	}
}

// GetShiftType retrieves a shift type by name
func (s *ShiftService) GetShiftType(ctx context.Context, shiftName string) (*hr.ShiftType, error) {
	return s.shiftRepo.GetShiftType(ctx, shiftName)
}

// ListShiftTypes retrieves all shift types
func (s *ShiftService) ListShiftTypes(ctx context.Context) ([]hr.ShiftType, error) {
	return s.shiftRepo.ListShiftTypes(ctx)
}

// CreateShiftType creates a new shift type
func (s *ShiftService) CreateShiftType(ctx context.Context, shift *hr.ShiftType) error {
	// Validate shift times
	if shift.StartTime == "" || shift.EndTime == "" {
		return fmt.Errorf("start_time and end_time are required")
	}
	return s.shiftRepo.CreateShiftType(ctx, shift)
}

// UpdateShiftType updates an existing shift type
func (s *ShiftService) UpdateShiftType(ctx context.Context, shift *hr.ShiftType) error {
	return s.shiftRepo.UpdateShiftType(ctx, shift)
}

// DeleteShiftType deletes a shift type
func (s *ShiftService) DeleteShiftType(ctx context.Context, id uint) error {
	return s.shiftRepo.DeleteShiftType(ctx, id)
}

// AssignShift assigns a shift to an employee
func (s *ShiftService) AssignShift(ctx context.Context, employee, shiftType, company string, fromDate time.Time, toDate *time.Time, approvedBy string) (*hr.ShiftAssignment, error) {
	// Validate input
	if employee == "" {
		return nil, fmt.Errorf("employee is required")
	}
	if shiftType == "" {
		return nil, fmt.Errorf("shift_type is required")
	}
	if company == "" {
		return nil, fmt.Errorf("company is required")
	}

	// Check if shift type exists
	_, err := s.shiftRepo.GetShiftType(ctx, shiftType)
	if err != nil {
		return nil, fmt.Errorf("shift type not found")
	}

	// Create shift assignment
	assignment := &hr.ShiftAssignment{
		Employee:   employee,
		ShiftType:  shiftType,
		Company:    company,
		FromDate:   fromDate,
		ToDate:     toDate,
		Status:     "Active",
		ApprovedBy: approvedBy,
	}

	if approvedBy != "" {
		now := time.Now()
		assignment.ApprovalDate = &now
	}

	err = s.shiftRepo.CreateShiftAssignment(ctx, assignment)
	if err != nil {
		return nil, err
	}

	return assignment, nil
}

// GetCurrentShift retrieves the current active shift for an employee on a specific date
func (s *ShiftService) GetCurrentShift(ctx context.Context, employee string, date time.Time) (*hr.ShiftAssignment, *hr.ShiftType, error) {
	assignment, err := s.shiftRepo.GetCurrentShiftAssignment(ctx, employee, date)
	if err != nil {
		return nil, nil, err
	}

	shiftType, err := s.shiftRepo.GetShiftType(ctx, assignment.ShiftType)
	if err != nil {
		return assignment, nil, err
	}

	return assignment, shiftType, nil
}

// GetShiftDetails retrieves detailed information about a shift including timing and rules
func (s *ShiftService) GetShiftDetails(ctx context.Context, shiftName string) (map[string]interface{}, error) {
	shift, err := s.shiftRepo.GetShiftType(ctx, shiftName)
	if err != nil {
		return nil, err
	}

	details := map[string]interface{}{
		"shift_name":                        shift.ShiftName,
		"start_time":                        shift.StartTime,
		"end_time":                          shift.EndTime,
		"begin_check_in_before":             shift.BeginCheckInBeforeShiftStartTime,
		"allow_check_out_after":             shift.AllowCheckOutAfterShiftEndTime,
		"working_hours_threshold_half_day":  shift.WorkingHoursThresholdForHalfDay,
		"working_hours_threshold_absent":    shift.WorkingHoursThresholdForAbsent,
		"enable_auto_attendance":            shift.EnableAutoAttendance,
		"determine_check_in_and_check_out":  shift.DetermineCheckInAndCheckOut,
		"process_attendance_after":          shift.ProcessAttendanceAfter,
		"holiday_list":                      shift.HolidayList,
		"color":                             shift.Color,
	}

	return details, nil
}

// GetShiftAssignmentsForEmployee retrieves all shift assignments for an employee
func (s *ShiftService) GetShiftAssignmentsForEmployee(ctx context.Context, employee string, status string) ([]hr.ShiftAssignment, error) {
	return s.shiftRepo.ListShiftAssignments(ctx, employee, "", status, nil, nil)
}

// GetShiftAssignmentsByShiftType retrieves all shift assignments for a specific shift type
func (s *ShiftService) GetShiftAssignmentsByShiftType(ctx context.Context, shiftType string, status string) ([]hr.ShiftAssignment, error) {
	return s.shiftRepo.ListShiftAssignments(ctx, "", shiftType, status, nil, nil)
}

// UpdateShiftAssignmentStatus updates the status of a shift assignment
func (s *ShiftService) UpdateShiftAssignmentStatus(ctx context.Context, id uint, status string) error {
	assignment, err := s.shiftRepo.GetShiftAssignment(ctx, id)
	if err != nil {
		return err
	}

	assignment.Status = status
	return s.shiftRepo.UpdateShiftAssignment(ctx, assignment)
}

// CreateShiftRequest creates a new shift request
func (s *ShiftService) CreateShiftRequest(ctx context.Context, request *hr.ShiftRequest) error {
	// Validate request
	if request.Employee == "" {
		return fmt.Errorf("employee is required")
	}
	if request.ShiftType == "" {
		return fmt.Errorf("shift_type is required")
	}

	// Check if shift type exists
	_, err := s.shiftRepo.GetShiftType(ctx, request.ShiftType)
	if err != nil {
		return fmt.Errorf("shift type not found")
	}

	request.Status = "Pending"
	return s.shiftRepo.CreateShiftRequest(ctx, request)
}

// ApproveShiftRequest approves a shift request and creates a shift assignment
func (s *ShiftService) ApproveShiftRequest(ctx context.Context, id uint, approvedBy string) (*hr.ShiftRequest, error) {
	request, err := s.shiftRepo.GetShiftRequest(ctx, id)
	if err != nil {
		return nil, err
	}

	if request.Status != "Pending" {
		return nil, fmt.Errorf("cannot approve shift request with status: %s", request.Status)
	}

	// Update request status
	request.Status = "Approved"
	request.ApprovedBy = approvedBy
	now := time.Now()
	request.ApprovalDate = &now

	err = s.shiftRepo.UpdateShiftRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// Create shift assignment
	_, err = s.AssignShift(ctx, request.Employee, request.ShiftType, request.Company, request.FromDate, request.ToDate, approvedBy)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// RejectShiftRequest rejects a shift request
func (s *ShiftService) RejectShiftRequest(ctx context.Context, id uint, rejectedBy, reason string) (*hr.ShiftRequest, error) {
	request, err := s.shiftRepo.GetShiftRequest(ctx, id)
	if err != nil {
		return nil, err
	}

	if request.Status != "Pending" {
		return nil, fmt.Errorf("cannot reject shift request with status: %s", request.Status)
	}

	request.Status = "Rejected"
	request.ApprovedBy = rejectedBy
	now := time.Now()
	request.ApprovalDate = &now

	err = s.shiftRepo.UpdateShiftRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// ListShiftRequests retrieves shift requests with filters
func (s *ShiftService) ListShiftRequests(ctx context.Context, employee, status string) ([]hr.ShiftRequest, error) {
	return s.shiftRepo.ListShiftRequests(ctx, employee, status)
}
