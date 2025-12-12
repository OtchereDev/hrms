package repositories

import (
	"context"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/payroll"
	"gorm.io/gorm"
)

// SalaryComponentRepository handles salary component data access
type SalaryComponentRepository struct {
	db *gorm.DB
}

// NewSalaryComponentRepository creates a new salary component repository
func NewSalaryComponentRepository(db *gorm.DB) *SalaryComponentRepository {
	return &SalaryComponentRepository{db: db}
}

// Create creates a new salary component
func (r *SalaryComponentRepository) Create(ctx context.Context, component *payroll.SalaryComponent) error {
	return r.db.WithContext(ctx).Create(component).Error
}

// GetByID retrieves a salary component by ID
func (r *SalaryComponentRepository) GetByID(ctx context.Context, id uint) (*payroll.SalaryComponent, error) {
	var component payroll.SalaryComponent
	err := r.db.WithContext(ctx).Preload("Accounts").First(&component, id).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}

// GetByName retrieves a salary component by name
func (r *SalaryComponentRepository) GetByName(ctx context.Context, name string) (*payroll.SalaryComponent, error) {
	var component payroll.SalaryComponent
	err := r.db.WithContext(ctx).
		Preload("Accounts").
		Where("salary_component_name = ?", name).
		First(&component).Error
	if err != nil {
		return nil, err
	}
	return &component, nil
}

// Update updates a salary component
func (r *SalaryComponentRepository) Update(ctx context.Context, component *payroll.SalaryComponent) error {
	return r.db.WithContext(ctx).Save(component).Error
}

// ListActive retrieves all active salary components
func (r *SalaryComponentRepository) ListActive(ctx context.Context, componentType string) ([]payroll.SalaryComponent, error) {
	var components []payroll.SalaryComponent
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	if componentType != "" {
		query = query.Where("type = ?", componentType)
	}

	err := query.Order("salary_component_name ASC").Find(&components).Error
	return components, err
}

// SalaryStructureRepository handles salary structure data access
type SalaryStructureRepository struct {
	db *gorm.DB
}

// NewSalaryStructureRepository creates a new salary structure repository
func NewSalaryStructureRepository(db *gorm.DB) *SalaryStructureRepository {
	return &SalaryStructureRepository{db: db}
}

// Create creates a new salary structure
func (r *SalaryStructureRepository) Create(ctx context.Context, structure *payroll.SalaryStructure) error {
	return r.db.WithContext(ctx).Create(structure).Error
}

// GetByID retrieves a salary structure by ID
func (r *SalaryStructureRepository) GetByID(ctx context.Context, id uint) (*payroll.SalaryStructure, error) {
	var structure payroll.SalaryStructure
	err := r.db.WithContext(ctx).
		Preload("Earnings").
		Preload("Deductions").
		Preload("OtherBenefits").
		First(&structure, id).Error
	if err != nil {
		return nil, err
	}
	return &structure, nil
}

// GetByName retrieves a salary structure by name
func (r *SalaryStructureRepository) GetByName(ctx context.Context, name string) (*payroll.SalaryStructure, error) {
	var structure payroll.SalaryStructure
	err := r.db.WithContext(ctx).
		Preload("Earnings").
		Preload("Deductions").
		Preload("OtherBenefits").
		Where("salary_structure_name = ?", name).
		First(&structure).Error
	if err != nil {
		return nil, err
	}
	return &structure, nil
}

// Update updates a salary structure
func (r *SalaryStructureRepository) Update(ctx context.Context, structure *payroll.SalaryStructure) error {
	return r.db.WithContext(ctx).Save(structure).Error
}

// ListActive retrieves all active salary structures
func (r *SalaryStructureRepository) ListActive(ctx context.Context, company string) ([]payroll.SalaryStructure, error) {
	var structures []payroll.SalaryStructure
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	if company != "" {
		query = query.Where("company = ?", company)
	}

	err := query.Order("salary_structure_name ASC").Find(&structures).Error
	return structures, err
}

// SalaryStructureAssignmentRepository handles salary structure assignments
type SalaryStructureAssignmentRepository struct {
	db *gorm.DB
}

// NewSalaryStructureAssignmentRepository creates a new assignment repository
func NewSalaryStructureAssignmentRepository(db *gorm.DB) *SalaryStructureAssignmentRepository {
	return &SalaryStructureAssignmentRepository{db: db}
}

// Create creates a new salary structure assignment
func (r *SalaryStructureAssignmentRepository) Create(ctx context.Context, assignment *payroll.SalaryStructureAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

// GetByID retrieves an assignment by ID
func (r *SalaryStructureAssignmentRepository) GetByID(ctx context.Context, id uint) (*payroll.SalaryStructureAssignment, error) {
	var assignment payroll.SalaryStructureAssignment
	err := r.db.WithContext(ctx).First(&assignment, id).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// GetActiveForEmployee retrieves the active assignment for an employee
func (r *SalaryStructureAssignmentRepository) GetActiveForEmployee(ctx context.Context, employee string) (*payroll.SalaryStructureAssignment, error) {
	var assignment payroll.SalaryStructureAssignment
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).
		Where("employee = ? AND from_date <= ?", employee, now).
		Where("to_date IS NULL OR to_date >= ?", now).
		Order("from_date DESC").
		First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// Update updates an assignment
func (r *SalaryStructureAssignmentRepository) Update(ctx context.Context, assignment *payroll.SalaryStructureAssignment) error {
	return r.db.WithContext(ctx).Save(assignment).Error
}

// List retrieves assignments with filters
func (r *SalaryStructureAssignmentRepository) List(ctx context.Context, filters SalaryAssignmentFilters, page, pageSize int) ([]payroll.SalaryStructureAssignment, int64, error) {
	var assignments []payroll.SalaryStructureAssignment
	var total int64

	query := r.db.WithContext(ctx).Model(&payroll.SalaryStructureAssignment{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.SalaryStructure != "" {
		query = query.Where("salary_structure = ?", filters.SalaryStructure)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("from_date DESC").Find(&assignments).Error

	return assignments, total, err
}

// SalaryAssignmentFilters represents filters for assignment queries
type SalaryAssignmentFilters struct {
	Employee        string
	Company         string
	SalaryStructure string
}

// SalarySlipRepository handles salary slip data access
type SalarySlipRepository struct {
	db *gorm.DB
}

// NewSalarySlipRepository creates a new salary slip repository
func NewSalarySlipRepository(db *gorm.DB) *SalarySlipRepository {
	return &SalarySlipRepository{db: db}
}

// Create creates a new salary slip
func (r *SalarySlipRepository) Create(ctx context.Context, slip *payroll.SalarySlip) error {
	return r.db.WithContext(ctx).Create(slip).Error
}

// GetByID retrieves a salary slip by ID
func (r *SalarySlipRepository) GetByID(ctx context.Context, id uint) (*payroll.SalarySlip, error) {
	var slip payroll.SalarySlip
	err := r.db.WithContext(ctx).
		Preload("Earnings").
		Preload("Deductions").
		First(&slip, id).Error
	if err != nil {
		return nil, err
	}
	return &slip, nil
}

// GetByEmployeeAndPeriod retrieves salary slip for employee in a period
func (r *SalarySlipRepository) GetByEmployeeAndPeriod(ctx context.Context, employee string, startDate, endDate time.Time) (*payroll.SalarySlip, error) {
	var slip payroll.SalarySlip
	err := r.db.WithContext(ctx).
		Preload("Earnings").
		Preload("Deductions").
		Where("employee = ? AND start_date = ? AND end_date = ?", employee, startDate, endDate).
		First(&slip).Error
	if err != nil {
		return nil, err
	}
	return &slip, nil
}

// Update updates a salary slip
func (r *SalarySlipRepository) Update(ctx context.Context, slip *payroll.SalarySlip) error {
	return r.db.WithContext(ctx).Save(slip).Error
}

// List retrieves salary slips with filters
func (r *SalarySlipRepository) List(ctx context.Context, filters SalarySlipFilters, page, pageSize int) ([]payroll.SalarySlip, int64, error) {
	var slips []payroll.SalarySlip
	var total int64

	query := r.db.WithContext(ctx).Model(&payroll.SalarySlip{})

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
	if !filters.StartDate.IsZero() {
		query = query.Where("start_date >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("end_date <= ?", filters.EndDate)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("start_date DESC").Find(&slips).Error

	return slips, total, err
}

// SalarySlipFilters represents filters for salary slip queries
type SalarySlipFilters struct {
	Employee   string
	Company    string
	Department string
	Status     string
	StartDate  time.Time
	EndDate    time.Time
}

// AdditionalSalaryRepository handles additional salary data access
type AdditionalSalaryRepository struct {
	db *gorm.DB
}

// NewAdditionalSalaryRepository creates a new additional salary repository
func NewAdditionalSalaryRepository(db *gorm.DB) *AdditionalSalaryRepository {
	return &AdditionalSalaryRepository{db: db}
}

// Create creates a new additional salary
func (r *AdditionalSalaryRepository) Create(ctx context.Context, salary *payroll.AdditionalSalary) error {
	return r.db.WithContext(ctx).Create(salary).Error
}

// GetByID retrieves an additional salary by ID
func (r *AdditionalSalaryRepository) GetByID(ctx context.Context, id uint) (*payroll.AdditionalSalary, error) {
	var salary payroll.AdditionalSalary
	err := r.db.WithContext(ctx).First(&salary, id).Error
	if err != nil {
		return nil, err
	}
	return &salary, nil
}

// GetForEmployeeAndDate retrieves additional salaries for an employee in a period
func (r *AdditionalSalaryRepository) GetForEmployeeAndDate(ctx context.Context, employee string, payrollDate time.Time) ([]payroll.AdditionalSalary, error) {
	var salaries []payroll.AdditionalSalary
	err := r.db.WithContext(ctx).
		Where("employee = ? AND payroll_date = ?", employee, payrollDate).
		Find(&salaries).Error
	return salaries, err
}

// Update updates an additional salary
func (r *AdditionalSalaryRepository) Update(ctx context.Context, salary *payroll.AdditionalSalary) error {
	return r.db.WithContext(ctx).Save(salary).Error
}

// List retrieves additional salaries with filters
func (r *AdditionalSalaryRepository) List(ctx context.Context, filters AdditionalSalaryFilters, page, pageSize int) ([]payroll.AdditionalSalary, int64, error) {
	var salaries []payroll.AdditionalSalary
	var total int64

	query := r.db.WithContext(ctx).Model(&payroll.AdditionalSalary{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("payroll_date DESC").Find(&salaries).Error

	return salaries, total, err
}

// AdditionalSalaryFilters represents filters for additional salary queries
type AdditionalSalaryFilters struct {
	Employee string
	Company  string
	Type     string
}
