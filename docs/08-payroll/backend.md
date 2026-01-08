# Payroll Module - Backend Documentation

## Overview

Complete payroll processing system with salary structures, salary slips, payroll entries, tax calculations, and accounting integration.

---

## Core Doctypes

### 1. Salary Structure

**Purpose**: Defines salary components and calculation formulas for employees.

#### Key Fields
- `name`: Structure name
- `company`: Company
- `payroll_frequency`: Monthly, Bi-monthly, Weekly, etc.
- `earnings`: Child table of earning components (Basic, HRA, etc.)
- `deductions`: Child table of deduction components (Tax, PF, etc.)
- `mode_of_payment`: Cash, Bank, Cheque
- `payment_account`: Bank account

#### Business Logic
- Validates total percentage allocation
- Auto-assigns to employees based on designation/grade
- Calculates amounts using formulas
- Supports condition-based components

---

### 2. Salary Slip

**Purpose**: Individual employee's monthly salary calculation.

#### Key Fields
- `employee`: Employee reference
- `salary_structure`: Applied structure
- `start_date`, `end_date`: Pay period
- `payment_days`: Working days - absences
- `earnings`: Calculated earning components
- `deductions`: Calculated deduction components
- `gross_pay`, `total_deduction`, `net_pay`: Totals
- `loan_repayment`: EMI deductions
- `income_tax`: Tax deduction

#### Business Logic
```python
def calculate():
    # Get base salary from structure
    # Calculate based on payment days
    # Add additional salaries (bonus, incentives)
    # Deduct loans, advances
    # Calculate tax
    # Create GL entries
```

---

### 3. Payroll Entry

**Purpose**: Bulk salary slip generation and processing.

#### Key Fields
- `company`, `department`, `branch`, `designation`: Filters
- `start_date`, `end_date`: Payroll period
- `employees`: Selected employees list
- `posting_date`: Accounting date
- `payment_account`: Bank account for payment

#### Business Logic
- Creates salary slips for all employees
- Submits slips in bulk
- Creates bank entries
- Generates payment vouchers

---

### 4. Salary Component

**Purpose**: Individual salary elements (Basic, HRA, Tax, etc.).

#### Key Fields
- `salary_component_name`: Component name
- `type`: Earning or Deduction
- `is_tax_applicable`: Include in taxable income
- `formula`: Calculation formula
- `amount_based_on_formula`: Auto-calculate or manual
- `accounts`: GL accounts per company

---

### 5. Additional Salary

**Purpose**: One-time additions (Bonus, Arrears, Deductions).

#### Key Fields
- `employee`: Employee reference
- `salary_component`: Component to add
- `amount`: Amount to add/deduct
- `payroll_date`: Effective date
- `overwrite_salary_structure_amount`: Override default

---

## Tax Calculations

### Income Tax Slab
- Progressive tax rates
- Annual income calculation
- Monthly TDS deduction
- Proof submission handling

---

## Key Business Rules

1. Payment days calculation includes holidays and weekends
2. Loans auto-deducted from salary
3. Tax calculated based on annual projection
4. Benefits (Flexi-benefits) allocated
5. Gratuity calculated on separation
6. Salary withholding for recoveries

---

This backend documentation covers core payroll functionality.
