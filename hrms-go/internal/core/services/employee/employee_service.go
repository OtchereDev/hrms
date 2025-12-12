package employee

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
	ErrEmployeeNotFound       = errors.New("employee not found")
	ErrEmployeeExists         = errors.New("employee already exists")
	ErrInvalidEmployeeData    = errors.New("invalid employee data")
	ErrEmployeeNumberRequired = errors.New("employee number is required")
	ErrEmployeeNameRequired   = errors.New("employee name is required")
	ErrCompanyRequired        = errors.New("company is required")
	ErrDateOfJoiningRequired  = errors.New("date of joining is required")
)

// EmployeeService handles employee business logic
type EmployeeService struct {
	repo *repositories.EmployeeRepository
	db   *gorm.DB
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(db *gorm.DB) *EmployeeService {
	return &EmployeeService{
		repo: repositories.NewEmployeeRepository(db),
		db:   db,
	}
}

// CreateEmployeeRequest represents a request to create an employee
type CreateEmployeeRequest struct {
	// Basic Information
	EmployeeNumber   string     `json:"employee_number" validate:"required"`
	FirstName        string     `json:"first_name" validate:"required"`
	MiddleName       string     `json:"middle_name"`
	LastName         string     `json:"last_name"`
	Gender           string     `json:"gender"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	DateOfJoining    *time.Time `json:"date_of_joining" validate:"required"`

	// Contact Information
	CellNumber       string     `json:"cell_number"`
	PersonalEmail    string     `json:"personal_email"`
	CompanyEmail     string     `json:"company_email" validate:"required,email"`

	// Organization
	Company          string     `json:"company" validate:"required"`
	Department       string     `json:"department"`
	Designation      string     `json:"designation"`
	Branch           string     `json:"branch"`
	Grade            string     `json:"grade"`
	EmploymentType   string     `json:"employment_type"`
	ReportsTo        string     `json:"reports_to"`

	// Additional fields
	UserID           string     `json:"user_id"`
	DefaultShift     string     `json:"default_shift"`
	HolidayList      string     `json:"holiday_list"`
}

// UpdateEmployeeRequest represents a request to update an employee
type UpdateEmployeeRequest struct {
	FirstName        *string    `json:"first_name"`
	MiddleName       *string    `json:"middle_name"`
	LastName         *string    `json:"last_name"`
	Gender           *string    `json:"gender"`
	DateOfBirth      *time.Time `json:"date_of_birth"`

	CellNumber       *string    `json:"cell_number"`
	PersonalEmail    *string    `json:"personal_email"`
	CompanyEmail     *string    `json:"company_email"`

	Department       *string    `json:"department"`
	Designation      *string    `json:"designation"`
	Branch           *string    `json:"branch"`
	Grade            *string    `json:"grade"`
	EmploymentType   *string    `json:"employment_type"`
	ReportsTo        *string    `json:"reports_to"`

	DefaultShift     *string    `json:"default_shift"`
	HolidayList      *string    `json:"holiday_list"`
	Status           *string    `json:"status"`
}

// CreateEmployee creates a new employee
func (s *EmployeeService) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*hr.Employee, error) {
	// Validate required fields
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check if employee number already exists
	exists, err := s.repo.ExistsByEmployeeNumber(ctx, req.EmployeeNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check employee existence: %w", err)
	}
	if exists {
		return nil, ErrEmployeeExists
	}

	// Check if email already exists
	if req.CompanyEmail != "" {
		exists, err := s.repo.ExistsByEmail(ctx, req.CompanyEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, errors.New("email already in use")
		}
	}

	// Create employee
	employee := &hr.Employee{
		EmployeeNumber:   req.EmployeeNumber,
		FirstName:        req.FirstName,
		MiddleName:       req.MiddleName,
		LastName:         req.LastName,
		Gender:           req.Gender,
		DateOfBirth:      req.DateOfBirth,
		DateOfJoining:    req.DateOfJoining,
		CellNumber:       req.CellNumber,
		PersonalEmail:    req.PersonalEmail,
		CompanyEmail:     req.CompanyEmail,
		Company:          req.Company,
		Department:       req.Department,
		Designation:      req.Designation,
		Branch:           req.Branch,
		Grade:            req.Grade,
		EmploymentType:   req.EmploymentType,
		ReportsTo:        req.ReportsTo,
		UserID:           req.UserID,
		DefaultShift:     req.DefaultShift,
		HolidayList:      req.HolidayList,
		Status:           "Active",
	}

	// Calculate retirement date if date of birth is provided
	if req.DateOfBirth != nil {
		retirementDate := s.calculateRetirementDate(*req.DateOfBirth)
		employee.DateOfRetirement = &retirementDate
	}

	if err := s.repo.Create(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	return employee, nil
}

// GetEmployee retrieves an employee by ID
func (s *EmployeeService) GetEmployee(ctx context.Context, id uint) (*hr.Employee, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}
	return employee, nil
}

// GetEmployeeByNumber retrieves an employee by employee number
func (s *EmployeeService) GetEmployeeByNumber(ctx context.Context, employeeNumber string) (*hr.Employee, error) {
	employee, err := s.repo.GetByEmployeeNumber(ctx, employeeNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}
	return employee, nil
}

// GetEmployeeByUserID retrieves an employee by user ID
func (s *EmployeeService) GetEmployeeByUserID(ctx context.Context, userID string) (*hr.Employee, error) {
	employee, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}
	return employee, nil
}

// UpdateEmployee updates an employee
func (s *EmployeeService) UpdateEmployee(ctx context.Context, id uint, req UpdateEmployeeRequest) (*hr.Employee, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		employee.FirstName = *req.FirstName
	}
	if req.MiddleName != nil {
		employee.MiddleName = *req.MiddleName
	}
	if req.LastName != nil {
		employee.LastName = *req.LastName
	}
	if req.Gender != nil {
		employee.Gender = *req.Gender
	}
	if req.DateOfBirth != nil {
		employee.DateOfBirth = req.DateOfBirth
		// Recalculate retirement date
		retirementDate := s.calculateRetirementDate(*req.DateOfBirth)
		employee.DateOfRetirement = &retirementDate
	}
	if req.CellNumber != nil {
		employee.CellNumber = *req.CellNumber
	}
	if req.PersonalEmail != nil {
		employee.PersonalEmail = *req.PersonalEmail
	}
	if req.CompanyEmail != nil {
		// Check if email is already in use by another employee
		if *req.CompanyEmail != employee.CompanyEmail {
			exists, err := s.repo.ExistsByEmail(ctx, *req.CompanyEmail)
			if err != nil {
				return nil, fmt.Errorf("failed to check email existence: %w", err)
			}
			if exists {
				return nil, errors.New("email already in use")
			}
			employee.CompanyEmail = *req.CompanyEmail
		}
	}
	if req.Department != nil {
		employee.Department = *req.Department
	}
	if req.Designation != nil {
		employee.Designation = *req.Designation
	}
	if req.Branch != nil {
		employee.Branch = *req.Branch
	}
	if req.Grade != nil {
		employee.Grade = *req.Grade
	}
	if req.EmploymentType != nil {
		employee.EmploymentType = *req.EmploymentType
	}
	if req.ReportsTo != nil {
		employee.ReportsTo = *req.ReportsTo
	}
	if req.DefaultShift != nil {
		employee.DefaultShift = *req.DefaultShift
	}
	if req.HolidayList != nil {
		employee.HolidayList = *req.HolidayList
	}
	if req.Status != nil {
		employee.Status = *req.Status
	}

	if err := s.repo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	return employee, nil
}

// DeleteEmployee soft deletes an employee
func (s *EmployeeService) DeleteEmployee(ctx context.Context, id uint) error {
	// Check if employee exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEmployeeNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

// ListEmployees retrieves employees with pagination and filters
func (s *EmployeeService) ListEmployees(ctx context.Context, filters repositories.EmployeeFilters, page, pageSize int) ([]hr.Employee, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.List(ctx, filters, page, pageSize)
}

// GetEmployeeWithDetails retrieves employee with all related data
func (s *EmployeeService) GetEmployeeWithDetails(ctx context.Context, id uint) (*repositories.EmployeeDetails, error) {
	details, err := s.repo.GetEmployeeWithDetails(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, err
	}
	return details, nil
}

// UpdateEmployeeStatus updates employee status (Active, Inactive, Left, etc.)
func (s *EmployeeService) UpdateEmployeeStatus(ctx context.Context, id uint, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"Active":    true,
		"Inactive":  true,
		"Suspended": true,
		"Left":      true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// validateCreateRequest validates the create employee request
func (s *EmployeeService) validateCreateRequest(req CreateEmployeeRequest) error {
	if req.EmployeeNumber == "" {
		return ErrEmployeeNumberRequired
	}
	if req.FirstName == "" {
		return ErrEmployeeNameRequired
	}
	if req.Company == "" {
		return ErrCompanyRequired
	}
	if req.DateOfJoining == nil {
		return ErrDateOfJoiningRequired
	}
	if req.CompanyEmail == "" {
		return errors.New("company email is required")
	}

	return nil
}

// calculateRetirementDate calculates retirement date based on date of birth
// Default retirement age is 60 years (can be made configurable)
func (s *EmployeeService) calculateRetirementDate(dateOfBirth time.Time) time.Time {
	retirementAge := 60 // TODO: Make this configurable from HR Settings
	return dateOfBirth.AddDate(retirementAge, 0, 0)
}

// GetActiveEmployees retrieves all active employees
func (s *EmployeeService) GetActiveEmployees(ctx context.Context) ([]hr.Employee, error) {
	return s.repo.GetActiveEmployees(ctx)
}

// GetEmployeesByDepartment retrieves employees in a department
func (s *EmployeeService) GetEmployeesByDepartment(ctx context.Context, department string) ([]hr.Employee, error) {
	return s.repo.GetEmployeesByDepartment(ctx, department)
}

// GetEmployeesByManager retrieves employees reporting to a manager
func (s *EmployeeService) GetEmployeesByManager(ctx context.Context, managerEmployeeNumber string) ([]hr.Employee, error) {
	return s.repo.GetEmployeesByManager(ctx, managerEmployeeNumber)
}

// GetEmployeeDetails retrieves detailed employee information including related data
func (s *EmployeeService) GetEmployeeDetails(ctx context.Context, employeeNumber string) (map[string]interface{}, error) {
	employee, err := s.repo.GetByEmployeeNumber(ctx, employeeNumber)
	if err != nil {
		return nil, err
	}

	// Build detailed response with all employee information
	details := map[string]interface{}{
		"employee_number":    employee.EmployeeNumber,
		"first_name":         employee.FirstName,
		"middle_name":        employee.MiddleName,
		"last_name":          employee.LastName,
		"full_name":          employee.FirstName + " " + employee.LastName,
		"gender":             employee.Gender,
		"date_of_birth":      employee.DateOfBirth,
		"date_of_joining":    employee.DateOfJoining,
		"company":            employee.Company,
		"department":         employee.Department,
		"designation":        employee.Designation,
		"branch":             employee.Branch,
		"employment_type":    employee.EmploymentType,
		"status":             employee.Status,
		"reports_to":         employee.ReportsTo,
		"grade":              employee.Grade,
		"personal_email":     employee.PersonalEmail,
		"company_email":      employee.CompanyEmail,
		"cell_number":        employee.CellNumber,
		"emergency_contact":  employee.EmergencyContactName,
		"emergency_phone":    employee.EmergencyContactPhone,
		"current_address":    employee.CurrentAddress,
		"permanent_address":  employee.PermanentAddress,
		"user_id":            employee.UserID,
	}

	return details, nil
}

// SearchEmployees searches for employees based on various criteria
func (s *EmployeeService) SearchEmployees(ctx context.Context, query string, department, designation, status string, limit int) ([]hr.Employee, error) {
	if limit <= 0 {
		limit = 20
	}

	// Use repository's list method with filters
	var employees []hr.Employee
	db := s.repo.GetDB().WithContext(ctx)

	// Apply search query
	if query != "" {
		db = db.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(employee_number) LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	// Apply filters
	if department != "" {
		db = db.Where("department = ?", department)
	}
	if designation != "" {
		db = db.Where("designation = ?", designation)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	} else {
		db = db.Where("status = ?", "Active") // Default to active employees
	}

	err := db.Limit(limit).Find(&employees).Error
	return employees, err
}

// GetReportingStructure retrieves the complete reporting hierarchy for an employee
func (s *EmployeeService) GetReportingStructure(ctx context.Context, employeeNumber string) (map[string]interface{}, error) {
	employee, err := s.repo.GetByEmployeeNumber(ctx, employeeNumber)
	if err != nil {
		return nil, err
	}

	// Get reporting chain (upward)
	reportingChain := []map[string]interface{}{}
	currentEmployee := employee
	maxDepth := 10 // Prevent infinite loops
	depth := 0

	for currentEmployee.ReportsTo != "" && depth < maxDepth {
		manager, err := s.repo.GetByEmployeeNumber(ctx, currentEmployee.ReportsTo)
		if err != nil {
			break
		}

		reportingChain = append(reportingChain, map[string]interface{}{
			"employee_number": manager.EmployeeNumber,
			"name":            manager.FirstName + " " + manager.LastName,
			"designation":     manager.Designation,
			"level":           depth + 1,
		})

		currentEmployee = manager
		depth++
	}

	// Get direct reports (downward)
	directReports, _ := s.repo.GetEmployeesByManager(ctx, employeeNumber)
	reportsList := []map[string]interface{}{}
	for _, rep := range directReports {
		reportsList = append(reportsList, map[string]interface{}{
			"employee_number": rep.EmployeeNumber,
			"name":            rep.FirstName + " " + rep.LastName,
			"designation":     rep.Designation,
		})
	}

	structure := map[string]interface{}{
		"employee": map[string]interface{}{
			"employee_number": employee.EmployeeNumber,
			"name":            employee.FirstName + " " + employee.LastName,
			"designation":     employee.Designation,
		},
		"reporting_to":   reportingChain,
		"direct_reports": reportsList,
	}

	return structure, nil
}

// GetEmployeeFieldValue retrieves a specific field value for an employee
func (s *EmployeeService) GetEmployeeFieldValue(ctx context.Context, employeeNumber, fieldName string) (interface{}, error) {
	employee, err := s.repo.GetByEmployeeNumber(ctx, employeeNumber)
	if err != nil {
		return nil, err
	}

	// Map field names to employee struct fields
	fieldMap := map[string]interface{}{
		"employee_number":    employee.EmployeeNumber,
		"first_name":         employee.FirstName,
		"middle_name":        employee.MiddleName,
		"last_name":          employee.LastName,
		"gender":             employee.Gender,
		"date_of_birth":      employee.DateOfBirth,
		"date_of_joining":    employee.DateOfJoining,
		"company":            employee.Company,
		"department":         employee.Department,
		"designation":        employee.Designation,
		"branch":             employee.Branch,
		"employment_type":    employee.EmploymentType,
		"status":             employee.Status,
		"reports_to":         employee.ReportsTo,
		"grade":              employee.Grade,
		"personal_email":     employee.PersonalEmail,
		"company_email":      employee.CompanyEmail,
		"cell_number":        employee.CellNumber,
		"user_id":            employee.UserID,
	}

	value, exists := fieldMap[fieldName]
	if !exists {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	return value, nil
}
