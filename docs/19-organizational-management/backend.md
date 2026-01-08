# Organizational Management Module - Backend Documentation

## Overview

Manage organizational structure, departments, designations, branches, and hierarchies.

## Core Doctype: Company

### Fields
- `company_name`: Legal entity name
- `abbr`: Abbreviation for naming
- `country`, `default_currency`
- `parent_company`: For group structure
- `date_of_establishment`
- `default_holiday_list`
- `default_letter_head`

### Business Logic
- Multi-company support
- Inter-company transactions
- Company-specific settings
- Consolidated reporting

## Core Doctype: Department

### Fields
- `department_name`: Department identifier
- `parent_department`: For nested departments
- `company`: Company reference
- `is_group`: Group vs leaf node
- `leave_approver`, `expense_approver`
- `payroll_cost_center`

### Structure
- Uses nested set model for hierarchy
- Tree structure with parent-child relationships
- Support for multi-level departments

## Core Doctype: Designation

### Fields
- `designation_name`: Job title
- `description`: Role description
- `required_skills`: Child table of skills
- `required_experience_years`

### Usage
- Job role definition
- Career ladder
- Salary grade mapping
- Skill requirements

## Core Doctype: Branch

### Fields
- `branch`: Branch name
- `company`: Parent company
- `address`, `contact_details`
- `branch_head`: Employee reference

### Business Logic
- Geographic locations
- Branch-wise reporting
- Local compliance
- Branch transfers

## Core Doctype: Employment Type

### Fields
- `employment_type_name`: Full-time/Part-time/Contract/Intern
- `description`: Type details

### Usage
- Contract terms
- Benefits eligibility
- Leave policies
- Payroll rules

## Core Doctype: Employee Grade

### Fields
- `grade_name`: L1/L2/Manager/Director
- `default_leave_policy`
- `default_salary_structure`

### Business Logic
- Compensation bands
- Benefits tier
- Leave entitlements
- Reporting structure

## Organizational Chart

### Hierarchy Logic
```python
def get_org_chart(company):
    # Get all employees with reporting relationships
    employees = frappe.get_all(
        "Employee",
        filters={"company": company, "status": "Active"},
        fields=["name", "employee_name", "reports_to", "designation"]
    )
    
    # Build tree structure
    org_tree = build_tree(employees, root="CEO")
    return org_tree

def build_tree(employees, root):
    tree = {
        "name": root,
        "children": []
    }
    
    for emp in employees:
        if emp.reports_to == root:
            tree["children"].append(build_tree(employees, emp.name))
    
    return tree
```

## Department Hierarchy

Uses nested set for efficient tree operations:
- `lft`, `rgt` fields for tree positioning
- Fast subtree queries
- Efficient ancestor/descendant lookups

## Key Rules
- Department hierarchy enforced
- Reporting relationships validated
- Cross-company transfers require approval
- Branch-specific policies
- Grade-based access control
