package repositories

import (
	"context"
	"fmt"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"gorm.io/gorm"
)

// EmployeeRepository handles employee data access
type EmployeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// Create creates a new employee
func (r *EmployeeRepository) Create(ctx context.Context, employee *hr.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

// GetByID retrieves an employee by ID
func (r *EmployeeRepository) GetByID(ctx context.Context, id uint) (*hr.Employee, error) {
	var employee hr.Employee
	err := r.db.WithContext(ctx).First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByEmployeeNumber retrieves an employee by employee number
func (r *EmployeeRepository) GetByEmployeeNumber(ctx context.Context, employeeNumber string) (*hr.Employee, error) {
	var employee hr.Employee
	err := r.db.WithContext(ctx).Where("employee_number = ?", employeeNumber).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByUserID retrieves an employee by user ID
func (r *EmployeeRepository) GetByUserID(ctx context.Context, userID string) (*hr.Employee, error) {
	var employee hr.Employee
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByEmail retrieves an employee by company email
func (r *EmployeeRepository) GetByEmail(ctx context.Context, email string) (*hr.Employee, error) {
	var employee hr.Employee
	err := r.db.WithContext(ctx).Where("company_email = ?", email).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Update updates an employee
func (r *EmployeeRepository) Update(ctx context.Context, employee *hr.Employee) error {
	return r.db.WithContext(ctx).Save(employee).Error
}

// Delete soft deletes an employee
func (r *EmployeeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.Employee{}, id).Error
}

// List retrieves employees with pagination and filters
func (r *EmployeeRepository) List(ctx context.Context, filters EmployeeFilters, page, pageSize int) ([]hr.Employee, int64, error) {
	var employees []hr.Employee
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.Employee{})

	// Apply filters
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Department != "" {
		query = query.Where("department = ?", filters.Department)
	}
	if filters.Branch != "" {
		query = query.Where("branch = ?", filters.Branch)
	}
	if filters.Designation != "" {
		query = query.Where("designation = ?", filters.Designation)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.EmploymentType != "" {
		query = query.Where("employment_type = ?", filters.EmploymentType)
	}
	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		query = query.Where(
			"employee_name LIKE ? OR employee_number LIKE ? OR company_email LIKE ?",
			search, search, search,
		)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute query
	if err := query.Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

// GetActiveEmployees retrieves all active employees
func (r *EmployeeRepository) GetActiveEmployees(ctx context.Context) ([]hr.Employee, error) {
	var employees []hr.Employee
	err := r.db.WithContext(ctx).Where("status = ?", "Active").Find(&employees).Error
	return employees, err
}

// GetEmployeesByDepartment retrieves employees in a department
func (r *EmployeeRepository) GetEmployeesByDepartment(ctx context.Context, department string) ([]hr.Employee, error) {
	var employees []hr.Employee
	err := r.db.WithContext(ctx).Where("department = ? AND status = ?", department, "Active").Find(&employees).Error
	return employees, err
}

// GetEmployeesByManager retrieves employees reporting to a manager
func (r *EmployeeRepository) GetEmployeesByManager(ctx context.Context, managerEmployeeNumber string) ([]hr.Employee, error) {
	var employees []hr.Employee
	err := r.db.WithContext(ctx).Where("reports_to = ? AND status = ?", managerEmployeeNumber, "Active").Find(&employees).Error
	return employees, err
}

// ExistsByEmployeeNumber checks if an employee exists by employee number
func (r *EmployeeRepository) ExistsByEmployeeNumber(ctx context.Context, employeeNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&hr.Employee{}).Where("employee_number = ?", employeeNumber).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail checks if an employee exists by email
func (r *EmployeeRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&hr.Employee{}).Where("company_email = ?", email).Count(&count).Error
	return count > 0, err
}

// UpdateStatus updates employee status
func (r *EmployeeRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&hr.Employee{}).Where("id = ?", id).Update("status", status).Error
}

// EmployeeFilters represents filters for employee listing
type EmployeeFilters struct {
	Company        string
	Department     string
	Branch         string
	Designation    string
	Status         string
	EmploymentType string
	Search         string // Search in name, employee number, email
}

// AddEducation adds educational qualification for an employee
func (r *EmployeeRepository) AddEducation(ctx context.Context, education *hr.EmployeeEducation) error {
	return r.db.WithContext(ctx).Create(education).Error
}

// GetEducation retrieves educational qualifications for an employee
func (r *EmployeeRepository) GetEducation(ctx context.Context, employeeID uint) ([]hr.EmployeeEducation, error) {
	var education []hr.EmployeeEducation
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Find(&education).Error
	return education, err
}

// AddExternalWorkHistory adds external work history for an employee
func (r *EmployeeRepository) AddExternalWorkHistory(ctx context.Context, history *hr.EmployeeExternalWorkHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetExternalWorkHistory retrieves external work history for an employee
func (r *EmployeeRepository) GetExternalWorkHistory(ctx context.Context, employeeID uint) ([]hr.EmployeeExternalWorkHistory, error) {
	var history []hr.EmployeeExternalWorkHistory
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Order("from_date DESC").Find(&history).Error
	return history, err
}

// AddInternalWorkHistory adds internal work history for an employee
func (r *EmployeeRepository) AddInternalWorkHistory(ctx context.Context, history *hr.EmployeeInternalWorkHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetInternalWorkHistory retrieves internal work history for an employee
func (r *EmployeeRepository) GetInternalWorkHistory(ctx context.Context, employeeID uint) ([]hr.EmployeeInternalWorkHistory, error) {
	var history []hr.EmployeeInternalWorkHistory
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Order("from_date DESC").Find(&history).Error
	return history, err
}

// AddSkill adds a skill for an employee
func (r *EmployeeRepository) AddSkill(ctx context.Context, skill *hr.EmployeeSkill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

// GetSkills retrieves skills for an employee
func (r *EmployeeRepository) GetSkills(ctx context.Context, employeeID uint) ([]hr.EmployeeSkill, error) {
	var skills []hr.EmployeeSkill
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Find(&skills).Error
	return skills, err
}

// GetEmployeeWithDetails retrieves employee with all related data
func (r *EmployeeRepository) GetEmployeeWithDetails(ctx context.Context, id uint) (*EmployeeDetails, error) {
	employee, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	education, err := r.GetEducation(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get education: %w", err)
	}

	externalHistory, err := r.GetExternalWorkHistory(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get external work history: %w", err)
	}

	internalHistory, err := r.GetInternalWorkHistory(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get internal work history: %w", err)
	}

	skills, err := r.GetSkills(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	return &EmployeeDetails{
		Employee:            *employee,
		Education:           education,
		ExternalWorkHistory: externalHistory,
		InternalWorkHistory: internalHistory,
		Skills:              skills,
	}, nil
}

// EmployeeDetails contains employee with all related information
type EmployeeDetails struct {
	Employee            hr.Employee
	Education           []hr.EmployeeEducation
	ExternalWorkHistory []hr.EmployeeExternalWorkHistory
	InternalWorkHistory []hr.EmployeeInternalWorkHistory
	Skills              []hr.EmployeeSkill
}
