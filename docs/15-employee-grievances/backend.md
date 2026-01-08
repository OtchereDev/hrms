# Employee Grievances Module - Backend Documentation

## Overview

Systematic handling of employee complaints, grievances, and conflict resolution.

## Core Doctype: Employee Grievance

### Fields
- `employee`: Employee reference
- `grievance_against`: Person/Department
- `raised_on`: Submission date
- `grievance_type`: Harassment/Discrimination/Workload/Compensation/Other
- `description`: Detailed complaint
- `status`: Open/Under Investigation/Resolved/Closed
- `priority`: Low/Medium/High/Critical
- `resolution`: Final outcome text
- `resolution_date`: Date resolved

### Business Logic
- Auto-assign to HR department
- Escalation based on type and priority
- Investigation workflow
- Anonymous submission option
- Confidentiality flags

## Related Doctype: Grievance Resolution

### Fields
- `grievance`: Parent reference
- `action_taken`: Description
- `action_date`: When action performed
- `responsible_person`: Who took action
- `status_update`: Investigation findings

### Workflow
1. Employee submits grievance
2. HR acknowledges (24 hours SLA)
3. Investigation assigned
4. Evidence collection
5. Resolution proposed
6. Employee feedback
7. Final closure

## Escalation Rules
- High priority: Escalate to senior management after 3 days
- Harassment cases: Immediate escalation
- Unresolved after 30 days: Committee review
- Anonymous: Extra confidentiality protocols

## Key Rules
- All grievances must be acknowledged
- Investigation timeline tracking
- Multiple status updates allowed
- Email notifications at each stage
- Audit trail for compliance
- Whistleblower protection

## Notification Logic
```python
def send_notifications(grievance):
    # Notify HR on submission
    # Notify investigator on assignment
    # Remind if SLA breach imminent
    # Update employee on status changes
```
