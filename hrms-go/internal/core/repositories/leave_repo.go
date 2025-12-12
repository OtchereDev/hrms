package repositories

import (
	"context"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"gorm.io/gorm"
)

// LeaveApplicationRepository handles leave application data access
type LeaveApplicationRepository struct {
	db *gorm.DB
}

// NewLeaveApplicationRepository creates a new leave application repository
func NewLeaveApplicationRepository(db *gorm.DB) *LeaveApplicationRepository {
	return &LeaveApplicationRepository{db: db}
}

// Create creates a new leave application
func (r *LeaveApplicationRepository) Create(ctx context.Context, leave *hr.LeaveApplication) error {
	return r.db.WithContext(ctx).Create(leave).Error
}

// GetByID retrieves a leave application by ID
func (r *LeaveApplicationRepository) GetByID(ctx context.Context, id uint) (*hr.LeaveApplication, error) {
	var leave hr.LeaveApplication
	err := r.db.WithContext(ctx).First(&leave, id).Error
	if err != nil {
		return nil, err
	}
	return &leave, nil
}

// Update updates a leave application
func (r *LeaveApplicationRepository) Update(ctx context.Context, leave *hr.LeaveApplication) error {
	return r.db.WithContext(ctx).Save(leave).Error
}

// Delete soft deletes a leave application
func (r *LeaveApplicationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.LeaveApplication{}, id).Error
}

// List retrieves leave applications with pagination and filters
func (r *LeaveApplicationRepository) List(ctx context.Context, filters LeaveApplicationFilters, page, pageSize int) ([]hr.LeaveApplication, int64, error) {
	var leaves []hr.LeaveApplication
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.LeaveApplication{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.LeaveType != "" {
		query = query.Where("leave_type = ?", filters.LeaveType)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.LeaveApprover != "" {
		query = query.Where("leave_approver = ?", filters.LeaveApprover)
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
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&leaves).Error

	return leaves, total, err
}

// GetOverlappingLeaves checks for overlapping approved leaves
func (r *LeaveApplicationRepository) GetOverlappingLeaves(ctx context.Context, employee string, fromDate, toDate time.Time, excludeID uint) ([]hr.LeaveApplication, error) {
	var leaves []hr.LeaveApplication

	query := r.db.WithContext(ctx).
		Where("employee = ? AND status = ?", employee, "Approved").
		Where("(from_date <= ? AND to_date >= ?) OR (from_date <= ? AND to_date >= ?) OR (from_date >= ? AND to_date <= ?)",
			toDate, toDate, fromDate, fromDate, fromDate, toDate)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Find(&leaves).Error
	return leaves, err
}

// LeaveApplicationFilters represents filters for leave application queries
type LeaveApplicationFilters struct {
	Employee      string
	Company       string
	LeaveType     string
	Status        string
	LeaveApprover string
	FromDate      time.Time
	ToDate        time.Time
}

// LeaveAllocationRepository handles leave allocation data access
type LeaveAllocationRepository struct {
	db *gorm.DB
}

// NewLeaveAllocationRepository creates a new leave allocation repository
func NewLeaveAllocationRepository(db *gorm.DB) *LeaveAllocationRepository {
	return &LeaveAllocationRepository{db: db}
}

// Create creates a new leave allocation
func (r *LeaveAllocationRepository) Create(ctx context.Context, allocation *hr.LeaveAllocation) error {
	return r.db.WithContext(ctx).Create(allocation).Error
}

// GetByID retrieves a leave allocation by ID
func (r *LeaveAllocationRepository) GetByID(ctx context.Context, id uint) (*hr.LeaveAllocation, error) {
	var allocation hr.LeaveAllocation
	err := r.db.WithContext(ctx).First(&allocation, id).Error
	if err != nil {
		return nil, err
	}
	return &allocation, nil
}

// GetByEmployeeAndLeaveType retrieves allocation for employee and leave type
func (r *LeaveAllocationRepository) GetByEmployeeAndLeaveType(ctx context.Context, employee, leaveType string) (*hr.LeaveAllocation, error) {
	var allocation hr.LeaveAllocation
	err := r.db.WithContext(ctx).
		Where("employee = ? AND leave_type = ?", employee, leaveType).
		Order("to_date DESC").
		First(&allocation).Error
	if err != nil {
		return nil, err
	}
	return &allocation, nil
}

// Update updates a leave allocation
func (r *LeaveAllocationRepository) Update(ctx context.Context, allocation *hr.LeaveAllocation) error {
	return r.db.WithContext(ctx).Save(allocation).Error
}

// List retrieves leave allocations with filters
func (r *LeaveAllocationRepository) List(ctx context.Context, filters LeaveAllocationFilters, page, pageSize int) ([]hr.LeaveAllocation, int64, error) {
	var allocations []hr.LeaveAllocation
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.LeaveAllocation{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.LeaveType != "" {
		query = query.Where("leave_type = ?", filters.LeaveType)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("to_date DESC").Find(&allocations).Error

	return allocations, total, err
}

// LeaveAllocationFilters represents filters for leave allocation queries
type LeaveAllocationFilters struct {
	Employee  string
	Company   string
	LeaveType string
}

// LeaveTypeRepository handles leave type data access
type LeaveTypeRepository struct {
	db *gorm.DB
}

// NewLeaveTypeRepository creates a new leave type repository
func NewLeaveTypeRepository(db *gorm.DB) *LeaveTypeRepository {
	return &LeaveTypeRepository{db: db}
}

// Create creates a new leave type
func (r *LeaveTypeRepository) Create(ctx context.Context, leaveType *hr.LeaveType) error {
	return r.db.WithContext(ctx).Create(leaveType).Error
}

// GetByID retrieves a leave type by ID
func (r *LeaveTypeRepository) GetByID(ctx context.Context, id uint) (*hr.LeaveType, error) {
	var leaveType hr.LeaveType
	err := r.db.WithContext(ctx).First(&leaveType, id).Error
	if err != nil {
		return nil, err
	}
	return &leaveType, nil
}

// GetByName retrieves a leave type by name
func (r *LeaveTypeRepository) GetByName(ctx context.Context, name string) (*hr.LeaveType, error) {
	var leaveType hr.LeaveType
	err := r.db.WithContext(ctx).Where("leave_type_name = ?", name).First(&leaveType).Error
	if err != nil {
		return nil, err
	}
	return &leaveType, nil
}

// Update updates a leave type
func (r *LeaveTypeRepository) Update(ctx context.Context, leaveType *hr.LeaveType) error {
	return r.db.WithContext(ctx).Save(leaveType).Error
}

// Delete soft deletes a leave type
func (r *LeaveTypeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.LeaveType{}, id).Error
}

// ListActive retrieves all active leave types
func (r *LeaveTypeRepository) ListActive(ctx context.Context) ([]hr.LeaveType, error) {
	var leaveTypes []hr.LeaveType
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("leave_type_name ASC").
		Find(&leaveTypes).Error
	return leaveTypes, err
}

// LeaveBalanceRepository handles leave balance queries
type LeaveBalanceRepository struct {
	db *gorm.DB
}

// NewLeaveBalanceRepository creates a new leave balance repository
func NewLeaveBalanceRepository(db *gorm.DB) *LeaveBalanceRepository {
	return &LeaveBalanceRepository{db: db}
}

// GetLeaveBalance calculates leave balance for an employee
func (r *LeaveBalanceRepository) GetLeaveBalance(ctx context.Context, employee, leaveType string) (float64, error) {
	// Get allocation
	var allocation hr.LeaveAllocation
	err := r.db.WithContext(ctx).
		Where("employee = ? AND leave_type = ?", employee, leaveType).
		Order("to_date DESC").
		First(&allocation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	// Get total used leaves (approved + submitted)
	var totalUsed float64
	err = r.db.WithContext(ctx).
		Model(&hr.LeaveApplication{}).
		Where("employee = ? AND leave_type = ? AND status IN ?", employee, leaveType, []string{"Approved", "Submitted"}).
		Select("COALESCE(SUM(total_leave_days), 0)").
		Scan(&totalUsed).Error
	if err != nil {
		return 0, err
	}

	balance := allocation.TotalLeavesAllocated - totalUsed
	return balance, nil
}

// GetAllLeaveBalances retrieves all leave balances for an employee
func (r *LeaveBalanceRepository) GetAllLeaveBalances(ctx context.Context, employee string) (map[string]float64, error) {
	// Get all allocations
	var allocations []hr.LeaveAllocation
	err := r.db.WithContext(ctx).
		Where("employee = ?", employee).
		Find(&allocations).Error
	if err != nil {
		return nil, err
	}

	balances := make(map[string]float64)

	for _, allocation := range allocations {
		// Get used leaves for this type
		var totalUsed float64
		err = r.db.WithContext(ctx).
			Model(&hr.LeaveApplication{}).
			Where("employee = ? AND leave_type = ? AND status IN ?",
				employee, allocation.LeaveType, []string{"Approved", "Submitted"}).
			Select("COALESCE(SUM(total_leave_days), 0)").
			Scan(&totalUsed).Error
		if err != nil {
			return nil, err
		}

		balance := allocation.TotalLeavesAllocated - totalUsed
		balances[allocation.LeaveType] = balance
	}

	return balances, nil
}

// LeaveEncashmentRepository handles leave encashment data access
type LeaveEncashmentRepository struct {
	db *gorm.DB
}

// NewLeaveEncashmentRepository creates a new leave encashment repository
func NewLeaveEncashmentRepository(db *gorm.DB) *LeaveEncashmentRepository {
	return &LeaveEncashmentRepository{db: db}
}

// Create creates a new leave encashment record
func (r *LeaveEncashmentRepository) Create(ctx context.Context, encashment *hr.LeaveEncashment) error {
	return r.db.WithContext(ctx).Create(encashment).Error
}

// GetByID retrieves a leave encashment by ID
func (r *LeaveEncashmentRepository) GetByID(ctx context.Context, id uint) (*hr.LeaveEncashment, error) {
	var encashment hr.LeaveEncashment
	err := r.db.WithContext(ctx).First(&encashment, id).Error
	if err != nil {
		return nil, err
	}
	return &encashment, nil
}

// Update updates a leave encashment
func (r *LeaveEncashmentRepository) Update(ctx context.Context, encashment *hr.LeaveEncashment) error {
	return r.db.WithContext(ctx).Save(encashment).Error
}

// List retrieves leave encashments with filters
func (r *LeaveEncashmentRepository) List(ctx context.Context, filters LeaveEncashmentFilters, page, pageSize int) ([]hr.LeaveEncashment, int64, error) {
	var encashments []hr.LeaveEncashment
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.LeaveEncashment{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.LeaveType != "" {
		query = query.Where("leave_type = ?", filters.LeaveType)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&encashments).Error

	return encashments, total, err
}

// LeaveEncashmentFilters represents filters for leave encashment queries
type LeaveEncashmentFilters struct {
	Employee  string
	Company   string
	LeaveType string
}
