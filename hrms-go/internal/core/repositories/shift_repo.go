package repositories

import (
	"context"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"gorm.io/gorm"
)

// ShiftRepository handles shift data operations
type ShiftRepository struct {
	db *gorm.DB
}

// NewShiftRepository creates a new shift repository
func NewShiftRepository(db *gorm.DB) *ShiftRepository {
	return &ShiftRepository{db: db}
}

// GetShiftType retrieves a shift type by name
func (r *ShiftRepository) GetShiftType(ctx context.Context, shiftName string) (*hr.ShiftType, error) {
	var shift hr.ShiftType
	err := r.db.WithContext(ctx).Where("shift_name = ?", shiftName).First(&shift).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

// ListShiftTypes retrieves all shift types
func (r *ShiftRepository) ListShiftTypes(ctx context.Context) ([]hr.ShiftType, error) {
	var shifts []hr.ShiftType
	err := r.db.WithContext(ctx).Find(&shifts).Error
	return shifts, err
}

// CreateShiftType creates a new shift type
func (r *ShiftRepository) CreateShiftType(ctx context.Context, shift *hr.ShiftType) error {
	return r.db.WithContext(ctx).Create(shift).Error
}

// UpdateShiftType updates an existing shift type
func (r *ShiftRepository) UpdateShiftType(ctx context.Context, shift *hr.ShiftType) error {
	return r.db.WithContext(ctx).Save(shift).Error
}

// DeleteShiftType deletes a shift type
func (r *ShiftRepository) DeleteShiftType(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.ShiftType{}, id).Error
}

// GetShiftAssignment retrieves a shift assignment by ID
func (r *ShiftRepository) GetShiftAssignment(ctx context.Context, id uint) (*hr.ShiftAssignment, error) {
	var assignment hr.ShiftAssignment
	err := r.db.WithContext(ctx).First(&assignment, id).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// GetCurrentShiftAssignment retrieves the current active shift for an employee
func (r *ShiftRepository) GetCurrentShiftAssignment(ctx context.Context, employee string, date time.Time) (*hr.ShiftAssignment, error) {
	var assignment hr.ShiftAssignment
	err := r.db.WithContext(ctx).
		Where("employee = ? AND status = ? AND from_date <= ? AND (to_date IS NULL OR to_date >= ?)",
			employee, "Active", date, date).
		First(&assignment).Error

	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// ListShiftAssignments retrieves shift assignments with filters
func (r *ShiftRepository) ListShiftAssignments(ctx context.Context, employee, shiftType, status string, fromDate, toDate *time.Time) ([]hr.ShiftAssignment, error) {
	var assignments []hr.ShiftAssignment
	query := r.db.WithContext(ctx)

	if employee != "" {
		query = query.Where("employee = ?", employee)
	}
	if shiftType != "" {
		query = query.Where("shift_type = ?", shiftType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if fromDate != nil {
		query = query.Where("from_date >= ?", fromDate)
	}
	if toDate != nil {
		query = query.Where("to_date <= ? OR to_date IS NULL", toDate)
	}

	err := query.Find(&assignments).Error
	return assignments, err
}

// CreateShiftAssignment creates a new shift assignment
func (r *ShiftRepository) CreateShiftAssignment(ctx context.Context, assignment *hr.ShiftAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

// UpdateShiftAssignment updates an existing shift assignment
func (r *ShiftRepository) UpdateShiftAssignment(ctx context.Context, assignment *hr.ShiftAssignment) error {
	return r.db.WithContext(ctx).Save(assignment).Error
}

// DeleteShiftAssignment deletes a shift assignment
func (r *ShiftRepository) DeleteShiftAssignment(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.ShiftAssignment{}, id).Error
}

// GetShiftRequest retrieves a shift request by ID
func (r *ShiftRepository) GetShiftRequest(ctx context.Context, id uint) (*hr.ShiftRequest, error) {
	var request hr.ShiftRequest
	err := r.db.WithContext(ctx).First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// ListShiftRequests retrieves shift requests with filters
func (r *ShiftRepository) ListShiftRequests(ctx context.Context, employee, status string) ([]hr.ShiftRequest, error) {
	var requests []hr.ShiftRequest
	query := r.db.WithContext(ctx)

	if employee != "" {
		query = query.Where("employee = ?", employee)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&requests).Error
	return requests, err
}

// CreateShiftRequest creates a new shift request
func (r *ShiftRepository) CreateShiftRequest(ctx context.Context, request *hr.ShiftRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

// UpdateShiftRequest updates an existing shift request
func (r *ShiftRepository) UpdateShiftRequest(ctx context.Context, request *hr.ShiftRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

// DeleteShiftRequest deletes a shift request
func (r *ShiftRepository) DeleteShiftRequest(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.ShiftRequest{}, id).Error
}
