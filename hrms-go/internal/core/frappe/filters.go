package frappe

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Filter represents a Frappe filter
// Filters can be in multiple formats:
// 1. Simple: {"status": "Active"}
// 2. Array: [["status", "=", "Active"], ["docstatus", "<", 2]]
// 3. Dict with operators: {"status": ["=", "Active"]}
type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

// ParseFilters parses Frappe-style filters
func ParseFilters(filtersJSON string) ([]Filter, error) {
	if filtersJSON == "" || filtersJSON == "{}" || filtersJSON == "[]" {
		return []Filter{}, nil
	}

	// Try to parse as array of arrays first
	var arrayFilters [][]interface{}
	if err := json.Unmarshal([]byte(filtersJSON), &arrayFilters); err == nil {
		return parseArrayFilters(arrayFilters)
	}

	// Try to parse as map
	var mapFilters map[string]interface{}
	if err := json.Unmarshal([]byte(filtersJSON), &mapFilters); err == nil {
		return parseMapFilters(mapFilters)
	}

	return nil, fmt.Errorf("invalid filter format")
}

// parseArrayFilters parses filters in array format
// [["field", "operator", "value"], ...]
func parseArrayFilters(arrayFilters [][]interface{}) ([]Filter, error) {
	filters := make([]Filter, 0, len(arrayFilters))

	for _, filterArray := range arrayFilters {
		if len(filterArray) < 3 {
			continue
		}

		field, ok := filterArray[0].(string)
		if !ok {
			continue
		}

		operator, ok := filterArray[1].(string)
		if !ok {
			operator = "="
		}

		value := filterArray[2]

		filters = append(filters, Filter{
			Field:    field,
			Operator: operator,
			Value:    value,
		})
	}

	return filters, nil
}

// parseMapFilters parses filters in map format
// {"field": "value"} or {"field": ["operator", "value"]}
func parseMapFilters(mapFilters map[string]interface{}) ([]Filter, error) {
	filters := make([]Filter, 0, len(mapFilters))

	for field, value := range mapFilters {
		// Check if value is an array [operator, value]
		if arrayValue, ok := value.([]interface{}); ok && len(arrayValue) >= 2 {
			operator, ok := arrayValue[0].(string)
			if !ok {
				operator = "="
			}

			filters = append(filters, Filter{
				Field:    field,
				Operator: operator,
				Value:    arrayValue[1],
			})
		} else {
			// Simple case: {"field": "value"}
			filters = append(filters, Filter{
				Field:    field,
				Operator: "=",
				Value:    value,
			})
		}
	}

	return filters, nil
}

// ApplyFilters applies filters to a GORM query
func ApplyFilters(query *gorm.DB, filters []Filter) *gorm.DB {
	for _, filter := range filters {
		query = applyFilter(query, filter)
	}
	return query
}

// applyFilter applies a single filter to a GORM query
func applyFilter(query *gorm.DB, filter Filter) *gorm.DB {
	field := filter.Field
	operator := strings.ToLower(filter.Operator)
	value := filter.Value

	switch operator {
	case "=", "equals":
		return query.Where(fmt.Sprintf("%s = ?", field), value)

	case "!=", "not equals", "<>":
		return query.Where(fmt.Sprintf("%s != ?", field), value)

	case ">", "greater than":
		return query.Where(fmt.Sprintf("%s > ?", field), value)

	case ">=", "greater than or equals":
		return query.Where(fmt.Sprintf("%s >= ?", field), value)

	case "<", "less than":
		return query.Where(fmt.Sprintf("%s < ?", field), value)

	case "<=", "less than or equals":
		return query.Where(fmt.Sprintf("%s <= ?", field), value)

	case "in":
		return query.Where(fmt.Sprintf("%s IN ?", field), value)

	case "not in":
		return query.Where(fmt.Sprintf("%s NOT IN ?", field), value)

	case "like":
		return query.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v%%", value))

	case "not like":
		return query.Where(fmt.Sprintf("%s NOT LIKE ?", field), fmt.Sprintf("%%%v%%", value))

	case "is":
		if value == nil || value == "null" {
			return query.Where(fmt.Sprintf("%s IS NULL", field))
		}
		return query.Where(fmt.Sprintf("%s = ?", field), value)

	case "is not":
		if value == nil || value == "null" {
			return query.Where(fmt.Sprintf("%s IS NOT NULL", field))
		}
		return query.Where(fmt.Sprintf("%s != ?", field), value)

	case "between":
		// Expect value to be an array with 2 elements
		if arr, ok := value.([]interface{}); ok && len(arr) == 2 {
			return query.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), arr[0], arr[1])
		}
		return query

	default:
		// Default to equality
		return query.Where(fmt.Sprintf("%s = ?", field), value)
	}
}

// ParseOrFilters parses Frappe-style or_filters
// Same format as regular filters but combined with OR instead of AND
func ParseOrFilters(filtersJSON string) ([]Filter, error) {
	return ParseFilters(filtersJSON)
}

// ApplyOrFilters applies OR filters to a GORM query
func ApplyOrFilters(query *gorm.DB, filters []Filter) *gorm.DB {
	if len(filters) == 0 {
		return query
	}

	return query.Where(func(q *gorm.DB) *gorm.DB {
		for i, filter := range filters {
			if i == 0 {
				q = applyFilter(q, filter)
			} else {
				q = q.Or(func(subQ *gorm.DB) *gorm.DB {
					return applyFilter(subQ, filter)
				})
			}
		}
		return q
	})
}
