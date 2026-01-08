# Full and Final Settlement Module - Backend Documentation

## Overview

Process final settlement for employees leaving the organization, including dues calculation, clearance, and final payment.

## Core Doctype: Full and Final Statement

### Fields
- `employee`: Employee reference
- `resignation_letter`: Reference to resignation
- `relieving_date`: Last working day
- `gross_pay`: Final month prorated salary
- `leave_encashment`: Unused leave payment
- `gratuity_amount`: Gratuity due
- `bonus_amount`: Pro-rata bonus
- `deductions`: Recoveries/dues
- `net_payable`: Final settlement amount
- `status`: Draft/Pending Clearance/Approved/Paid

### Components Child Table
- `component_name`: Earning/Deduction
- `amount`: Component value
- `remarks`: Explanation

## Business Logic

### Earnings Calculation
- Salary for notice period worked
- Salary in lieu of notice (if applicable)
- Unused leave encashment
- Pro-rata bonus
- Gratuity (if eligible)
- Outstanding reimbursements

### Deductions Calculation
- Notice period shortfall recovery
- Loans/advances outstanding
- Company property not returned
- Damages/losses
- Tax deductions

## Core Doctype: Employee Clearance

### Fields
- `employee`: Employee reference
- `clearance_date`: Date initiated
- `departments`: Child table with:
  - `department`: IT/Admin/Finance/HR
  - `status`: Pending/Cleared/On Hold
  - `cleared_by`: Approver
  - `cleared_on`: Clearance date
  - `remarks`: Notes/pending items

### Clearance Checklist
- IT: Laptop, access cards, software licenses
- Admin: ID card, parking pass, locker keys
- Finance: Credit card, fuel card, advances
- HR: Documents, exit interview
- Project: Knowledge transfer, handover

## Workflow
1. Employee submits resignation
2. F&F statement auto-created
3. Clearance process initiated
4. Department-wise clearance
5. Final dues calculated
6. Approval by Finance/HR
7. Payment processing
8. Full and final letter issued

## Key Rules
- Cannot process until all clearances complete
- Gratuity only if eligible (5+ years)
- Leave encashment as per policy
- Notice period recovery if applicable
- Tax compliance on all payments
- Final statement immutable once paid

## Settlement Calculation
```python
def calculate_fnf_settlement(employee, relieving_date):
    earnings = {
        'salary': calculate_prorated_salary(),
        'leave_encashment': calculate_leave_balance() * per_day_salary,
        'gratuity': calculate_gratuity() if eligible else 0,
        'bonus': calculate_prorata_bonus(),
        'reimbursements': get_pending_reimbursements()
    }
    
    deductions = {
        'notice_recovery': calculate_notice_shortfall(),
        'advances': get_outstanding_advances(),
        'property': get_property_charges(),
        'tax': calculate_tds()
    }
    
    net_payable = sum(earnings.values()) - sum(deductions.values())
    return net_payable
```
