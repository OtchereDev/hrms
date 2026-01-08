# Automation and Notifications Module - Backend Documentation

## Overview

Automated workflows, scheduled jobs, email notifications, and business rule automation.

## Email Notifications

### Leave Notifications
- Leave application submitted → Manager
- Leave approved/rejected → Employee
- Leave balance low warning → Employee
- Leave allocation created → Employee

### Attendance Notifications
- Missing attendance → Employee & Manager
- Late arrival alert → Manager
- Weekly attendance summary → Manager
- Attendance regularization request → Manager

### Payroll Notifications
- Salary slip generated → Employee
- Payment processed → Employee
- Tax documents available → Employee

### Recruitment Notifications
- Job application received → HR
- Interview scheduled → Candidate & Interviewer
- Offer letter → Candidate
- Joining reminder → New hire

### Approval Reminders
- Pending approvals → Approver (daily digest)
- Escalation on delayed approval
- Bulk approval summary

## Scheduled Jobs

### Daily Jobs
```python
# Run at 9:00 AM daily
def send_attendance_reminders():
    employees = get_employees_without_attendance(yesterday)
    for emp in employees:
        send_email(
            recipients=[emp.email],
            subject="Missing Attendance",
            template="attendance_reminder"
        )

# Run at 11:59 PM daily
def mark_absent_for_no_checkin():
    process_auto_attendance_for_all_shifts()
```

### Weekly Jobs
```python
# Run Monday 9:00 AM
def send_weekly_reports():
    managers = get_all_managers()
    for manager in managers:
        report_data = get_team_summary(manager, last_week)
        send_email(
            recipients=[manager.email],
            subject="Weekly Team Report",
            template="weekly_summary",
            context=report_data
        )
```

### Monthly Jobs
```python
# Run 1st of month
def create_monthly_leave_allocations():
    leave_policies = get_active_leave_policies()
    for policy in leave_policies:
        create_leave_allocation(policy)

# Run 25th of month
def process_monthly_payroll():
    create_salary_slips_for_all()
    send_notification_to_hr()
```

### Yearly Jobs
```python
# Run January 1st
def annual_leave_carryforward():
    carry_forward_unused_leaves()
    reset_annual_quotas()
```

## Workflow Automation

### Auto Assignment Rules
```python
def auto_assign_leave_approver(doc):
    # Assign to reporting manager
    if doc.employee.reports_to:
        add_assignee(doc, doc.employee.reports_to)
    else:
        # Fallback to HR manager
        add_assignee(doc, get_hr_manager())
```

### Auto Status Updates
- Mark employee inactive after separation
- Update probation status on completion
- Change contract status on expiry
- Archive old documents

### Escalation Rules
```python
def check_approval_escalation():
    # Escalate if pending > 3 days
    pending_docs = get_pending_approvals(days=3)
    for doc in pending_docs:
        escalate_to_next_level(doc)
        send_escalation_notification(doc)
```

## Business Rule Engine

### Custom Validation Rules
- Salary component rules
- Leave policy rules
- Attendance rules
- Expense approval rules

### Auto-Calculation Rules
- Overtime calculation
- Leave balance updates
- Payroll component calculations
- Performance score aggregation

## Notification Channels

### Email
- Transactional emails
- Digest/summary emails
- Alert emails
- Newsletter

### In-App Notifications
- Real-time notifications
- Notification center
- Read/unread status
- Action buttons

### SMS (Optional)
- OTP for authentication
- Critical alerts
- Attendance reminders

### Push Notifications (Mobile)
- Approval requests
- Status updates
- Reminders

## Email Templates

### Template Structure
```html
<!-- Subject: Leave Application Approved -->
<div>
    <h2>Leave Approved</h2>
    <p>Dear {{ employee_name }},</p>
    <p>Your leave application has been approved.</p>
    
    <table>
        <tr><td>Leave Type:</td><td>{{ leave_type }}</td></tr>
        <tr><td>From Date:</td><td>{{ from_date }}</td></tr>
        <tr><td>To Date:</td><td>{{ to_date }}</td></tr>
        <tr><td>Total Days:</td><td>{{ total_days }}</td></tr>
    </table>
    
    <p>Approved by: {{ approved_by }}</p>
</div>
```

## Webhook Integration

### Outgoing Webhooks
```python
def send_webhook_on_employee_created(doc):
    webhook_url = frappe.db.get_single_value(
        "HR Settings", 
        "employee_webhook_url"
    )
    
    payload = {
        "event": "employee.created",
        "data": doc.as_dict()
    }
    
    requests.post(webhook_url, json=payload)
```

### Supported Events
- Employee created/updated
- Leave approved/rejected
- Attendance marked
- Salary slip created
- Expense claim submitted

## Report Scheduling

### Configuration
- Report name
- Frequency (daily/weekly/monthly)
- Recipients
- Filters
- Format (PDF/Excel)

### Execution
```python
def execute_scheduled_report(config):
    report_data = frappe.get_report(
        config.report_name,
        filters=config.filters
    )
    
    file = generate_file(report_data, config.format)
    
    send_email(
        recipients=config.recipients,
        subject=f"Scheduled Report: {config.report_name}",
        attachments=[file]
    )
```
