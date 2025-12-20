# Salary Withholding Module - Backend Documentation

## Overview

Manage deductions and withholdings from employee salaries for loans, advances, and other obligations.

## Core Doctype: Retention Bonus

### Fields
- `employee`: Employee reference
- `bonus_amount`: Total bonus amount
- `bonus_payment_date`: Original payment date
- `retention_period`: Months to retain
- `status`: Active/Paid/Forfeited

### Business Logic
- Hold bonus payment subject to retention period
- Release proportionally or lump sum after period
- Forfeit if employee leaves early
- Tax treatment on actual payment

## Core Doctype: Employee Advance

### Fields
- `employee`: Employee reference
- `advance_amount`: Amount advanced
- `repayment_method`: Monthly installments
- `monthly_repayment_amount`: Installment amount
- `status`: Draft/Paid/Returned
- `repayment_schedule`: Child table with dates/amounts

### Business Logic
- Deduct from salary slip automatically
- Track outstanding balance
- Handle early repayment
- Interest calculation if applicable

## Salary Slip Integration

### Deduction Components
- Loan repayment
- Advance recovery
- Damage/loss recovery
- Legal obligations (garnishment)
- Voluntary deductions (charity, savings)

## Key Rules
- Cannot exceed maximum deduction percentage (typically 50%)
- Court-ordered deductions take priority
- Minimum take-home salary must be maintained
- Withholding continues until fully recovered
- Employee notification required

## Recovery Logic
```python
def calculate_deductions(salary_slip):
    max_deductible = salary_slip.gross_pay * 0.5
    priority_deductions = get_priority_deductions()
    remaining = max_deductible
    for deduction in sorted_by_priority:
        if remaining > 0:
            deduct_amount = min(deduction.amount, remaining)
            remaining -= deduct_amount
```
