# Regional Customizations Module - Backend Documentation

## Overview

Country-specific compliance, statutory requirements, and localization features for different regions.

## India Customizations

### Provident Fund (PF)
- Employee PF contribution: 12% of basic
- Employer PF contribution: 12% of basic  
- PF ceiling: Rs. 15,000 per month
- ECR (Electronic Challan Return) generation
- UAN (Universal Account Number) tracking

### Employee State Insurance (ESI)
- Employee contribution: 0.75% of gross
- Employer contribution: 3.25% of gross
- ESI ceiling: Rs. 21,000 per month
- ESI registration and returns

### Professional Tax
- State-specific PT slabs
- Monthly deduction
- PT registration per state

### Gratuity
- Formula: (15 × Last salary × Years) / 26
- Minimum service: 5 years
- Tax exemption limits

### Income Tax
- Tax slabs (Old vs New regime)
- Section 80C, 80D exemptions
- HRA exemption calculation
- TDS deduction
- Form 16 generation

### Leave Encashment
- Tax treatment rules
- Maximum encashment limits
- Formula as per service years

## United States Customizations

### Federal Tax (W-2)
- Federal income tax withholding
- Social Security (6.2%)
- Medicare (1.45%)
- Additional Medicare (0.9% if > $200k)

### State Tax
- State-specific tax rates
- State unemployment insurance
- Local taxes

### 401(k) Retirement
- Employee contribution (pre-tax)
- Employer matching
- Contribution limits

### Benefits
- Health insurance
- Dental/Vision
- Life insurance
- Disability insurance

### Paid Time Off (PTO)
- Combined leave pool
- Accrual rules
- Carryover policies

## United Kingdom Customizations

### PAYE (Pay As You Earn)
- Income tax deduction
- National Insurance (NI)
  - Employee NI: 12%
  - Employer NI: 13.8%
- Tax codes

### Pension Schemes
- Auto-enrollment
- Minimum contributions
  - Employee: 5%
  - Employer: 3%

### Statutory Payments
- Statutory Sick Pay (SSP)
- Statutory Maternity Pay (SMP)
- Statutory Paternity Pay (SPP)

### P45/P60 Forms
- P45: Leaving employment
- P60: End of tax year

## Middle East (UAE/Saudi) Customizations

### End of Service Benefits
- Gratuity calculation based on service years
- Different formulas for unlimited vs limited contracts

### Work Permit & Visa
- Visa expiry tracking
- Labor card management
- Medical insurance mandatory

### Leave Policies
- Annual leave: 30 days after 1 year
- Ticket allowance
- Ramadan working hours

## Singapore Customizations

### Central Provident Fund (CPF)
- Employee contribution: up to 20%
- Employer contribution: up to 17%
- Age-based contribution rates
- CPF ceilings (Ordinary Wage + Additional Wage)

### Leave Entitlements
- Annual leave: 7-14 days based on service
- Sick leave: 14 days outpatient, 60 days hospitalization
- Childcare leave
- Maternity/Paternity leave

### IR8A Form
- Annual tax form
- Auto-population from salary data

## Common Regional Features

### Multi-Currency Support
- Base currency per company
- Exchange rate handling
- Salary payment in local currency

### Statutory Reports
- Country-specific compliance reports
- Monthly/quarterly returns
- Annual filings

### Holiday Calendars
- Country/region-specific holidays
- Religious holidays
- State/province holidays

### Labor Law Compliance
- Working hours limits
- Overtime regulations
- Notice period requirements
- Termination rules

## Configuration

### Regional Settings
```python
{
    "country": "India",
    "currency": "INR",
    "tax_regime": "New",
    "pf_applicable": True,
    "esi_applicable": True,
    "pt_states": ["Maharashtra", "Karnataka"],
    "gratuity_applicable": True
}
```

### Salary Component Mapping
```python
# Map components to statutory requirements
STATUTORY_COMPONENTS = {
    "India": {
        "PF": ["Basic", "DA"],  # PF calculated on these
        "ESI": ["Gross"],
        "PT": ["Gross"]
    },
    "US": {
        "Social Security": ["Gross"],
        "Medicare": ["Gross"]
    }
}
```

## Localization

### Date Formats
- US: MM/DD/YYYY
- UK/India: DD/MM/YYYY
- ISO: YYYY-MM-DD

### Number Formats
- Indian: 1,00,000 (lakhs)
- Western: 100,000
- Decimal separator: . or ,

### Language Support
- English (default)
- Hindi, Spanish, Arabic, Chinese
- RTL support for Arabic

## Compliance Tracking

### Audit Trail
- All statutory calculations logged
- Report generation history
- Form submissions tracked
- Changes to tax settings recorded

### Deadlines
- Monthly return deadlines
- Quarterly filings
- Annual reports
- Automated reminders
