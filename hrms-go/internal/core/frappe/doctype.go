package frappe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"github.com/OtchereDev/hrms-go/internal/core/models/payroll"
)

// DocTypeService provides generic CRUD operations for any DocType
type DocTypeService struct {
	db *gorm.DB
}

// NewDocTypeService creates a new DocType service
func NewDocTypeService(db *gorm.DB) *DocTypeService {
	return &DocTypeService{db: db}
}

// DocTypeModel maps Frappe DocType names to Go model types
var DocTypeModelMap = map[string]interface{}{
	// HR Module
	"Attendance":                      &hr.Attendance{},
	"Employee Checkin":                &hr.EmployeeCheckin{},
	"Attendance Request":              &hr.AttendanceRequest{},
	"Leave Application":               &hr.LeaveApplication{},
	"Leave Allocation":                &hr.LeaveAllocation{},
	"Leave Type":                      &hr.LeaveType{},
	"Leave Encashment":                &hr.LeaveEncashment{},
	"Appraisal":                       &hr.Appraisal{},
	"Appraisal Template":              &hr.AppraisalTemplate{},
	"Appraisal Cycle":                 &hr.AppraisalCycle{},
	"Goal":                            &hr.Goal{},
	"Employee Performance Feedback":   &hr.EmployeePerformanceFeedback{},
	"Employee Skill Map":              &hr.EmployeeSkillMap{},

	// Payroll Module
	"Salary Slip":                     &payroll.SalarySlip{},
	"Salary Structure":                &payroll.SalaryStructure{},
	"Salary Structure Assignment":     &payroll.SalaryStructureAssignment{},
	"Salary Component":                &payroll.SalaryComponent{},
	"Additional Salary":               &payroll.AdditionalSalary{},
	"Payroll Entry":                   &payroll.PayrollEntry{},
	"Loan":                            &payroll.Loan{},
	"Employee Advance":                &payroll.EmployeeAdvance{},
	"Expense Claim":                   &payroll.ExpenseClaim{},
	"Gratuity":                        &payroll.Gratuity{},
	"Income Tax Slab":                 &payroll.IncomeTaxSlab{},
	"Employee Tax Exemption Declaration": &payroll.EmployeeTaxExemptionDeclaration{},
	"Employee Benefit Application":    &payroll.EmployeeBenefitApplication{},
}

// DocTypeTableMap maps DocType names to database table names
var DocTypeTableMap = map[string]string{
	"Attendance":                      "attendances",
	"Employee Checkin":                "employee_checkins",
	"Attendance Request":              "attendance_requests",
	"Leave Application":               "leave_applications",
	"Leave Allocation":                "leave_allocations",
	"Leave Type":                      "leave_types",
	"Leave Encashment":                "leave_encashments",
	"Appraisal":                       "appraisals",
	"Appraisal Template":              "appraisal_templates",
	"Appraisal Cycle":                 "appraisal_cycles",
	"Goal":                            "goals",
	"Employee Performance Feedback":   "employee_performance_feedbacks",
	"Employee Skill Map":              "employee_skill_maps",
	"Salary Slip":                     "salary_slips",
	"Salary Structure":                "salary_structures",
	"Salary Structure Assignment":     "salary_structure_assignments",
	"Salary Component":                "salary_components",
	"Additional Salary":               "additional_salaries",
	"Payroll Entry":                   "payroll_entries",
	"Loan":                            "loans",
	"Employee Advance":                "employee_advances",
	"Expense Claim":                   "expense_claims",
	"Gratuity":                        "gratuities",
	"Income Tax Slab":                 "income_tax_slabs",
	"Employee Tax Exemption Declaration": "employee_tax_exemption_declarations",
	"Employee Benefit Application":    "employee_benefit_applications",
}

// GetModelForDocType returns a new instance of the model for a DocType
func (s *DocTypeService) GetModelForDocType(doctype string) (interface{}, error) {
	model, ok := DocTypeModelMap[doctype]
	if !ok {
		return nil, fmt.Errorf("unknown doctype: %s", doctype)
	}

	// Create a new instance of the model type
	modelType := reflect.TypeOf(model).Elem()
	return reflect.New(modelType).Interface(), nil
}

// GetTableName returns the database table name for a DocType
func (s *DocTypeService) GetTableName(doctype string) (string, error) {
	tableName, ok := DocTypeTableMap[doctype]
	if !ok {
		return "", fmt.Errorf("unknown doctype: %s", doctype)
	}
	return tableName, nil
}

// Get retrieves a single document by name (ID)
func (s *DocTypeService) Get(ctx context.Context, doctype, name string) (interface{}, error) {
	model, err := s.GetModelForDocType(doctype)
	if err != nil {
		return nil, err
	}

	result := s.db.WithContext(ctx).First(model, "id = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("document not found: %s %s", doctype, name)
		}
		return nil, result.Error
	}

	return model, nil
}

// GetList retrieves a list of documents with filters
func (s *DocTypeService) GetList(ctx context.Context, doctype string, filters map[string]interface{}, fields []string, limit, offset int) ([]interface{}, int64, error) {
	tableName, err := s.GetTableName(doctype)
	if err != nil {
		return nil, 0, err
	}

	model, err := s.GetModelForDocType(doctype)
	if err != nil {
		return nil, 0, err
	}

	// Create slice of the model type
	modelType := reflect.TypeOf(model).Elem()
	sliceType := reflect.SliceOf(reflect.PtrTo(modelType))
	results := reflect.MakeSlice(sliceType, 0, 0)
	resultsPtr := reflect.New(sliceType)
	resultsPtr.Elem().Set(results)

	// Build query
	query := s.db.WithContext(ctx).Table(tableName)

	// Apply simple filters (for backward compatibility)
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Select specific fields if provided
	if len(fields) > 0 {
		query = query.Select(fields)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Execute query
	if err := query.Find(resultsPtr.Interface()).Error; err != nil {
		return nil, 0, err
	}

	// Convert to []interface{}
	resultSlice := resultsPtr.Elem()
	resultList := make([]interface{}, resultSlice.Len())
	for i := 0; i < resultSlice.Len(); i++ {
		resultList[i] = resultSlice.Index(i).Interface()
	}

	return resultList, total, nil
}

// GetListWithFilters retrieves a list of documents with advanced Frappe-style filters
func (s *DocTypeService) GetListWithFilters(ctx context.Context, doctype string, filters []Filter, orFilters []Filter, fields []string, orderBy string, limit, offset int) ([]interface{}, int64, error) {
	tableName, err := s.GetTableName(doctype)
	if err != nil {
		return nil, 0, err
	}

	model, err := s.GetModelForDocType(doctype)
	if err != nil {
		return nil, 0, err
	}

	// Create slice of the model type
	modelType := reflect.TypeOf(model).Elem()
	sliceType := reflect.SliceOf(reflect.PtrTo(modelType))
	results := reflect.MakeSlice(sliceType, 0, 0)
	resultsPtr := reflect.New(sliceType)
	resultsPtr.Elem().Set(results)

	// Build query
	query := s.db.WithContext(ctx).Table(tableName)

	// Apply AND filters
	query = ApplyFilters(query, filters)

	// Apply OR filters
	query = ApplyOrFilters(query, orFilters)

	// Select specific fields if provided
	if len(fields) > 0 {
		query = query.Select(fields)
	}

	// Get total count before pagination
	var total int64
	countQuery := query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering
	if orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("id DESC")
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Execute query
	if err := query.Find(resultsPtr.Interface()).Error; err != nil {
		return nil, 0, err
	}

	// Convert to []interface{}
	resultSlice := resultsPtr.Elem()
	resultList := make([]interface{}, resultSlice.Len())
	for i := 0; i < resultSlice.Len(); i++ {
		resultList[i] = resultSlice.Index(i).Interface()
	}

	return resultList, total, nil
}

// Insert creates a new document
func (s *DocTypeService) Insert(ctx context.Context, doctype string, data map[string]interface{}) (interface{}, error) {
	model, err := s.GetModelForDocType(doctype)
	if err != nil {
		return nil, err
	}

	// Convert map to struct
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonData, model); err != nil {
		return nil, err
	}

	// Insert into database
	if err := s.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

// Save updates an existing document
func (s *DocTypeService) Save(ctx context.Context, doctype, name string, data map[string]interface{}) (interface{}, error) {
	model, err := s.GetModelForDocType(doctype)
	if err != nil {
		return nil, err
	}

	// First get the existing record
	result := s.db.WithContext(ctx).First(model, "id = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("document not found: %s %s", doctype, name)
		}
		return nil, result.Error
	}

	// Update with new data
	if err := s.db.WithContext(ctx).Model(model).Updates(data).Error; err != nil {
		return nil, err
	}

	// Reload the model to get updated values
	if err := s.db.WithContext(ctx).First(model, "id = ?", name).Error; err != nil {
		return nil, err
	}

	return model, nil
}

// Delete deletes a document
func (s *DocTypeService) Delete(ctx context.Context, doctype, name string) error {
	model, err := s.GetModelForDocType(doctype)
	if err != nil {
		return err
	}

	result := s.db.WithContext(ctx).Delete(model, "id = ?", name)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("document not found: %s %s", doctype, name)
	}

	return nil
}

// Submit sets a document's status to submitted (workflow action)
func (s *DocTypeService) Submit(ctx context.Context, doctype, name string) (interface{}, error) {
	// Update docstatus to 1 (submitted in Frappe)
	data := map[string]interface{}{
		"docstatus": 1,
	}
	return s.Save(ctx, doctype, name, data)
}

// Cancel sets a document's status to cancelled (workflow action)
func (s *DocTypeService) Cancel(ctx context.Context, doctype, name string) (interface{}, error) {
	// Update docstatus to 2 (cancelled in Frappe)
	data := map[string]interface{}{
		"docstatus": 2,
	}
	return s.Save(ctx, doctype, name, data)
}
