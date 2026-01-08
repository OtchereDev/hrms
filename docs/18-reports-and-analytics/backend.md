# Reports and Analytics Module - Backend Documentation

## Overview

Comprehensive reporting and analytics for HR metrics, compliance, and business intelligence.

## Standard HR Reports

### Employee Reports
- Employee Directory
- Employee Birthday Report
- New Hire Report
- Employee Exit Report
- Employee Leave Balance
- Employees working on a holiday

### Attendance Reports
- Monthly Attendance Sheet
- Attendance Summary
- Late Entry Report
- Early Exit Report
- Absenteeism Report

### Leave Reports
- Leave Balance Report
- Leave Ledger Entry
- Leave Type wise Balance
- Employee Leave Application

### Recruitment Reports
- Job Applicant Report
- Job Opening Summary
- Interview Assessment Report
- Recruitment Analytics

### Performance Reports
- Appraisal Summary
- Goal Completion Report
- Performance Trend Analysis

## Payroll Reports

### Salary Reports
- Salary Register
- Bank Remittance Report
- Salary Slip Summary
- Payroll Component-wise summary

### Statutory Reports
- Provident Fund Report
- Professional Tax Report
- Form 16 (Tax)
- Gratuity Report

### Cost Reports
- Department-wise Salary Cost
- Grade-wise Salary Distribution
- Salary vs Budget Variance

## Analytics Dashboards

### HR Dashboard
- Headcount trends
- Attrition rate
- New hires vs exits
- Average tenure
- Department distribution
- Grade distribution

### Attendance Dashboard
- Daily attendance overview
- Shift-wise attendance
- Overtime trends
- Leave patterns

### Recruitment Dashboard
- Open positions
- Time to hire
- Source effectiveness
- Offer acceptance rate

### Payroll Dashboard
- Total payroll cost
- Average salary
- Cost per department
- Variance analysis

## Report Generation Logic

### Filters
- Date range
- Department
- Designation
- Employee grade
- Company
- Branch
- Custom fields

### Export Formats
- PDF
- Excel
- CSV
- Print

## Key Features
- Scheduled reports (daily/weekly/monthly)
- Email distribution lists
- Custom report builder
- Drill-down capability
- Comparison periods (YoY, MoM)
- Filters and grouping
- Chart visualizations

## Custom Report Builder

### Fields
- `report_name`: Report identifier
- `ref_doctype`: Source doctype
- `filters`: JSON filters
- `columns`: Selected fields
- `sort_by`: Ordering
- `group_by`: Grouping field

## Scheduled Reports

### Logic
```python
def generate_scheduled_report(report_config):
    # Apply filters
    data = frappe.get_all(
        report_config.ref_doctype,
        filters=report_config.filters,
        fields=report_config.columns
    )
    
    # Format and generate
    file = generate_excel(data)
    
    # Email to recipients
    send_email(
        recipients=report_config.recipients,
        subject=f"Scheduled Report: {report_config.report_name}",
        attachments=[file]
    )
```

## Compliance Reports
- Form 16 (Annual tax)
- PF/ESI monthly returns
- Bonus Act compliance
- Contract labor report
- Wage register
