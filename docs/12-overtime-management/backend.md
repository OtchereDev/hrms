# Overtime Management Module - Backend Documentation

## Overview

Track and compensate overtime work beyond regular working hours.

## Core Doctype: Overtime Request

### Fields
- `employee`: Employee reference
- `from_date`, `to_date`: Overtime period
- `total_hours`: Calculated overtime hours
- `overtime_type`: Compensatory/Paid
- `reason`: Justification text
- `status`: Draft/Pending/Approved/Rejected

### Business Logic
- Calculate hours beyond shift timings
- Apply overtime multipliers (1.5x, 2x)
- Create compensatory off if applicable
- Integration with payroll for paid overtime

## Integration Points
- **Attendance**: Source for actual hours worked
- **Shift Type**: Regular hours baseline
- **Salary Slip**: Additional earning component
- **Leave Ledger**: Compensatory off credits

## Key Rules
- Requires manager approval
- Cannot overlap with leave periods
- Must have corresponding attendance records
- Overtime rates configurable per company
- Maximum overtime hours per week/month limits

## Validation Logic
```python
def validate_overtime_hours(employee, from_date, to_date):
    # Get attendance records
    # Calculate actual vs scheduled hours
    # Validate overtime doesn't exceed limits
    # Check for weekend/holiday multipliers
```
