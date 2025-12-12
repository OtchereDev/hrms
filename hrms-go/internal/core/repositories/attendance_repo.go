package repositories

import (
	"context"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"gorm.io/gorm"
)

// AttendanceRepository handles attendance data access
type AttendanceRepository struct {
	db *gorm.DB
}

// NewAttendanceRepository creates a new attendance repository
func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

// Create creates a new attendance record
func (r *AttendanceRepository) Create(ctx context.Context, attendance *hr.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

// GetByID retrieves an attendance record by ID
func (r *AttendanceRepository) GetByID(ctx context.Context, id uint) (*hr.Attendance, error) {
	var attendance hr.Attendance
	err := r.db.WithContext(ctx).First(&attendance, id).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

// GetByEmployeeAndDate retrieves attendance for an employee on a specific date
func (r *AttendanceRepository) GetByEmployeeAndDate(ctx context.Context, employee string, date time.Time) (*hr.Attendance, error) {
	var attendance hr.Attendance
	err := r.db.WithContext(ctx).
		Where("employee = ? AND attendance_date = ?", employee, date).
		First(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

// Update updates an attendance record
func (r *AttendanceRepository) Update(ctx context.Context, attendance *hr.Attendance) error {
	return r.db.WithContext(ctx).Save(attendance).Error
}

// Delete soft deletes an attendance record
func (r *AttendanceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.Attendance{}, id).Error
}

// List retrieves attendance records with pagination and filters
func (r *AttendanceRepository) List(ctx context.Context, filters AttendanceFilters, page, pageSize int) ([]hr.Attendance, int64, error) {
	var attendances []hr.Attendance
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.Attendance{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Department != "" {
		query = query.Where("department = ?", filters.Department)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if !filters.FromDate.IsZero() {
		query = query.Where("attendance_date >= ?", filters.FromDate)
	}
	if !filters.ToDate.IsZero() {
		query = query.Where("attendance_date <= ?", filters.ToDate)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("attendance_date DESC").Find(&attendances).Error

	return attendances, total, err
}

// GetMonthlyAttendance retrieves attendance summary for an employee for a month
func (r *AttendanceRepository) GetMonthlyAttendance(ctx context.Context, employee string, year int, month int) ([]hr.Attendance, error) {
	var attendances []hr.Attendance
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	err := r.db.WithContext(ctx).
		Where("employee = ? AND attendance_date BETWEEN ? AND ?", employee, startDate, endDate).
		Order("attendance_date ASC").
		Find(&attendances).Error

	return attendances, err
}

// AttendanceFilters represents filters for attendance queries
type AttendanceFilters struct {
	Employee   string
	Company    string
	Department string
	Status     string
	FromDate   time.Time
	ToDate     time.Time
}

// EmployeeCheckinRepository handles employee check-in/check-out data access
type EmployeeCheckinRepository struct {
	db *gorm.DB
}

// NewEmployeeCheckinRepository creates a new checkin repository
func NewEmployeeCheckinRepository(db *gorm.DB) *EmployeeCheckinRepository {
	return &EmployeeCheckinRepository{db: db}
}

// Create creates a new checkin record
func (r *EmployeeCheckinRepository) Create(ctx context.Context, checkin *hr.EmployeeCheckin) error {
	return r.db.WithContext(ctx).Create(checkin).Error
}

// GetByID retrieves a checkin record by ID
func (r *EmployeeCheckinRepository) GetByID(ctx context.Context, id uint) (*hr.EmployeeCheckin, error) {
	var checkin hr.EmployeeCheckin
	err := r.db.WithContext(ctx).First(&checkin, id).Error
	if err != nil {
		return nil, err
	}
	return &checkin, nil
}

// GetTodayCheckinsForEmployee retrieves all checkins for an employee today
func (r *EmployeeCheckinRepository) GetTodayCheckinsForEmployee(ctx context.Context, employee string) ([]hr.EmployeeCheckin, error) {
	var checkins []hr.EmployeeCheckin
	today := time.Now().UTC().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Where("employee = ? AND time >= ? AND time < ?", employee, today, tomorrow).
		Order("time DESC").
		Find(&checkins).Error

	return checkins, err
}

// GetLastCheckinForEmployee retrieves the last checkin for an employee
func (r *EmployeeCheckinRepository) GetLastCheckinForEmployee(ctx context.Context, employee string) (*hr.EmployeeCheckin, error) {
	var checkin hr.EmployeeCheckin
	err := r.db.WithContext(ctx).
		Where("employee = ?", employee).
		Order("time DESC").
		First(&checkin).Error

	if err != nil {
		return nil, err
	}
	return &checkin, nil
}

// List retrieves checkin records with filters
func (r *EmployeeCheckinRepository) List(ctx context.Context, filters CheckinFilters, page, pageSize int) ([]hr.EmployeeCheckin, int64, error) {
	var checkins []hr.EmployeeCheckin
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.EmployeeCheckin{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if !filters.FromDate.IsZero() {
		query = query.Where("time >= ?", filters.FromDate)
	}
	if !filters.ToDate.IsZero() {
		query = query.Where("time <= ?", filters.ToDate)
	}
	if filters.LogType != "" {
		query = query.Where("log_type = ?", filters.LogType)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("time DESC").Find(&checkins).Error

	return checkins, total, err
}

// CheckinFilters represents filters for checkin queries
type CheckinFilters struct {
	Employee string
	FromDate time.Time
	ToDate   time.Time
	LogType  string
}

// AttendanceRequestRepository handles attendance request data access
type AttendanceRequestRepository struct {
	db *gorm.DB
}

// NewAttendanceRequestRepository creates a new attendance request repository
func NewAttendanceRequestRepository(db *gorm.DB) *AttendanceRequestRepository {
	return &AttendanceRequestRepository{db: db}
}

// Create creates a new attendance request
func (r *AttendanceRequestRepository) Create(ctx context.Context, request *hr.AttendanceRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

// GetByID retrieves an attendance request by ID
func (r *AttendanceRequestRepository) GetByID(ctx context.Context, id uint) (*hr.AttendanceRequest, error) {
	var request hr.AttendanceRequest
	err := r.db.WithContext(ctx).First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// Update updates an attendance request
func (r *AttendanceRequestRepository) Update(ctx context.Context, request *hr.AttendanceRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

// List retrieves attendance requests with filters
func (r *AttendanceRequestRepository) List(ctx context.Context, filters AttendanceRequestFilters, page, pageSize int) ([]hr.AttendanceRequest, int64, error) {
	var requests []hr.AttendanceRequest
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.AttendanceRequest{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Status != "" {
		query = query.Where("workflow_state = ?", filters.Status)
	}
	if !filters.FromDate.IsZero() {
		query = query.Where("from_date >= ?", filters.FromDate)
	}
	if !filters.ToDate.IsZero() {
		query = query.Where("to_date <= ?", filters.ToDate)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&requests).Error

	return requests, total, err
}

// AttendanceRequestFilters represents filters for attendance request queries
type AttendanceRequestFilters struct {
	Employee string
	Company  string
	Status   string
	FromDate time.Time
	ToDate   time.Time
}
