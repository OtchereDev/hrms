package frappe

import (
	"reflect"
	"strings"
)

// DocTypeMeta represents metadata about a DocType (similar to Frappe's get_meta)
type DocTypeMeta struct {
	Name        string      `json:"name"`
	Module      string      `json:"module"`
	IsTable     bool        `json:"istable"`
	Editable    bool        `json:"editable_grid"`
	TrackChanges bool       `json:"track_changes"`
	Fields      []FieldMeta `json:"fields"`
}

// FieldMeta represents metadata about a field
type FieldMeta struct {
	FieldName    string `json:"fieldname"`
	FieldType    string `json:"fieldtype"`
	Label        string `json:"label"`
	Options      string `json:"options,omitempty"`
	Reqd         int    `json:"reqd"`
	ReadOnly     int    `json:"read_only"`
	Hidden       int    `json:"hidden"`
	InListView   int    `json:"in_list_view"`
	InStandardFilter int `json:"in_standard_filter"`
}

// GetMeta returns metadata for a DocType
func (s *DocTypeService) GetMeta(doctype string) (*DocTypeMeta, error) {
	model, err := s.GetModelForDocType(doctype)
	if err != nil {
		return nil, err
	}

	meta := &DocTypeMeta{
		Name:         doctype,
		Module:       getModuleFromDocType(doctype),
		IsTable:      false,
		Editable:     true,
		TrackChanges: true,
		Fields:       []FieldMeta{},
	}

	// Use reflection to extract field information
	modelType := reflect.TypeOf(model).Elem()

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldMeta := extractFieldMeta(field)
		if fieldMeta != nil {
			meta.Fields = append(meta.Fields, *fieldMeta)
		}
	}

	return meta, nil
}

// extractFieldMeta extracts metadata from a struct field
func extractFieldMeta(field reflect.StructField) *FieldMeta {
	// Get JSON tag for field name
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return nil
	}

	// Remove options from json tag (e.g., "omitempty")
	fieldName := strings.Split(jsonTag, ",")[0]
	if fieldName == "" {
		fieldName = strings.ToLower(field.Name)
	}

	// Map Go types to Frappe field types
	fieldType := mapGoTypeToFrappeType(field.Type)

	// Generate label from field name
	label := generateLabel(field.Name)

	// Check if field is required (from gorm tag)
	gormTag := field.Tag.Get("gorm")
	reqd := 0
	if strings.Contains(gormTag, "not null") {
		reqd = 1
	}

	// Determine if field should be in list view
	inListView := 0
	inStandardFilter := 0

	// Common fields that should appear in list view
	listViewFields := map[string]bool{
		"name": true, "title": true, "status": true,
		"employee": true, "employee_name": true,
		"from_date": true, "to_date": true,
		"posting_date": true, "company": true,
	}

	if listViewFields[fieldName] {
		inListView = 1
	}

	// Common filter fields
	filterFields := map[string]bool{
		"status": true, "docstatus": true,
		"employee": true, "company": true,
		"department": true, "from_date": true, "to_date": true,
	}

	if filterFields[fieldName] {
		inStandardFilter = 1
	}

	return &FieldMeta{
		FieldName:          fieldName,
		FieldType:          fieldType,
		Label:              label,
		Reqd:               reqd,
		ReadOnly:           0,
		Hidden:             0,
		InListView:         inListView,
		InStandardFilter:   inStandardFilter,
	}
}

// mapGoTypeToFrappeType maps Go types to Frappe field types
func mapGoTypeToFrappeType(t reflect.Type) string {
	// Handle pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "Data"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "Int"
	case reflect.Float32, reflect.Float64:
		return "Float"
	case reflect.Bool:
		return "Check"
	case reflect.Struct:
		// Check if it's a time.Time
		if t.Name() == "Time" {
			return "Datetime"
		}
		return "Data"
	default:
		return "Data"
	}
}

// generateLabel generates a human-readable label from a field name
func generateLabel(fieldName string) string {
	// Convert camelCase or PascalCase to Title Case with spaces
	var result []rune
	for i, r := range fieldName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, ' ')
		}
		result = append(result, r)
	}
	return strings.Title(string(result))
}

// getModuleFromDocType determines the module based on DocType
func getModuleFromDocType(doctype string) string {
	// Payroll module DocTypes
	payrollDoctypes := map[string]bool{
		"Salary Slip": true,
		"Salary Structure": true,
		"Salary Structure Assignment": true,
		"Salary Component": true,
		"Additional Salary": true,
		"Payroll Entry": true,
		"Loan": true,
		"Employee Advance": true,
		"Expense Claim": true,
		"Gratuity": true,
		"Income Tax Slab": true,
		"Employee Tax Exemption Declaration": true,
		"Employee Benefit Application": true,
	}

	if payrollDoctypes[doctype] {
		return "Payroll"
	}

	return "HR"
}
