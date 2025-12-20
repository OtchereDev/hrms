# Gratuity and Bonus Module - Backend Documentation

## Overview

Manage gratuity calculations, retention bonuses, and employee incentives.

## Core Doctype: Gratuity

### Fields
- `employee`: Employee reference
- `gratuity_rule`: Applicable rule reference
- `current_work_experience`: Years of service
- `amount`: Calculated gratuity amount
- `gratuity_status`: Draft/Submitted/Paid

### Business Logic
- Calculate based on last drawn salary
- Formula: (15 days salary × years of service) / 26
- Minimum service period: 5 years
- Pro-rata for partial years
- Tax implications

## Core Doctype: Additional Salary

### Fields
- `employee`: Employee reference
- `salary_component`: Bonus/Incentive component
- `amount`: Bonus amount
- `payroll_date`: Month to pay
- `recurring`: One-time or recurring
- `overwrite_salary_structure_amount`: Override flag

### Types
- Performance bonus
- Retention bonus
- Referral bonus
- Festival bonus
- Project completion bonus

## Business Logic
- Auto-create salary slip component
- Prorated for partial months
- Recurring vs one-time handling
- Tax deduction at source (TDS)

## Key Rules
- Gratuity only after 5 years of continuous service
- Bonus subject to company performance
- Requires HR/Finance approval
- Must integrate with payroll cycle
