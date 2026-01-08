# HRMS Application Documentation

## Overview

This documentation provides comprehensive details about the HRMS (Human Resource Management System) application functionality. The documentation is organized by functional modules and broken down into three perspectives:

- **Designer**: UI/UX elements, layouts, visual components, and user flows
- **Frontend**: User interactions, forms, data display, and client-side behaviors
- **Backend**: Business logic, data processing, validations, and workflows

## Purpose

This documentation is designed to enable exact replication of the HRMS application. It focuses on **what the system does** rather than the technical implementation details, making it useful for rebuilding the application in any technology stack.

## Documentation Structure

Each module contains three files:
- `designer.md` - For UI/UX designers and visual design requirements
- `frontend.md` - For frontend developers implementing user interactions
- `backend.md` - For backend developers implementing business logic

## Modules

### Core HR Modules

1. **[Employee Lifecycle](./01-employee-lifecycle/)** - Employee onboarding, transfers, promotions, separations, and exit management
2. **[Leave Management](./02-leave-management/)** - Leave types, policies, allocations, applications, and encashment
3. **[Attendance & Shifts](./03-attendance-and-shifts/)** - Attendance tracking, shift management, and employee check-ins
4. **[Recruitment](./04-recruitment/)** - Job openings, applicant tracking, interviews, and job offers

### Performance & Development

5. **[Performance Management](./05-performance-management/)** - Appraisals, KRAs, feedback, and performance reviews
6. **[Skills & Competency](./16-skills-and-competency/)** - Skills mapping, assessment, and competency tracking

### Financial Modules

7. **[Expense Management](./06-expense-management/)** - Expense claims and employee advances
8. **[Travel Management](./07-travel-management/)** - Travel requests, itineraries, and vehicle logs
9. **[Payroll](./08-payroll/)** - Salary structures, salary slips, and payroll processing
10. **[Taxation & Benefits](./09-taxation-and-benefits/)** - Income tax, exemptions, and employee benefits
11. **[Overtime Management](./12-overtime-management/)** - Overtime tracking and calculations
12. **[Gratuity & Bonus](./13-gratuity-and-bonus/)** - Gratuity rules and retention bonuses
13. **[Salary Withholding](./14-salary-withholding/)** - Salary withholding cycles and management
14. **[Full & Final Settlement](./17-full-and-final-settlement/)** - Employee exit settlements

### Communication & Documentation

15. **[Daily Work Summary](./10-daily-work-summary/)** - Work summaries and team updates
16. **[Appointment Letters](./11-appointment-letters/)** - Letter templates and generation
17. **[Employee Grievances](./15-employee-grievances/)** - Grievance tracking and resolution

### Analytics & Organization

18. **[Reports & Analytics](./18-reports-and-analytics/)** - Dashboards, reports, charts, and KPIs
19. **[Organizational Management](./19-organizational-management/)** - Organizational hierarchy and structure

### Platform & Automation

20. **[Mobile App](./20-mobile-app/)** - Mobile PWA interface and functionality
21. **[Automation & Notifications](./21-automation-and-notifications/)** - Scheduled jobs, reminders, and notifications
22. **[Regional Customizations](./22-regional-customizations/)** - Country-specific features and regulations

## How to Use This Documentation

### For Project Planning
- Start with the README in each module to understand the scope
- Review all three perspectives (designer, frontend, backend) for complete understanding

### For Designers
- Focus on `designer.md` files for UI/UX specifications
- Reference `frontend.md` for interactive behavior requirements

### For Frontend Developers
- Start with `frontend.md` for user interaction requirements
- Reference `designer.md` for visual specifications
- Reference `backend.md` for understanding data requirements and business rules

### For Backend Developers
- Focus on `backend.md` files for business logic and data models
- Reference `frontend.md` to understand user-facing requirements

### For Full-Stack Developers
- Read all three files for complete context
- Start with backend, then frontend, then designer for implementation order

## Application Statistics

- **Total Modules**: 22 functional areas
- **Core Doctypes**: 145+ (97 HR + 48 Payroll)
- **Reports**: 27 (18 HR + 9 Payroll)
- **Workspaces**: 9 (7 HR + 2 Payroll)
- **Mobile Views**: 15+ screens
- **Dashboard Charts**: 33+
- **KPI Cards**: 25+

## Key Workflows

### Employee Journey
```
Job Opening → Job Applicant → Interview → Job Offer →
Employee Onboarding → Employee Record →
(Transfers/Promotions/Training) → Appraisals →
Employee Separation → Exit Interview → Full & Final Settlement
```

### Leave Workflow
```
Leave Type Definition → Leave Policy → Policy Assignment →
Leave Allocation → Leave Application → Approval → Leave Ledger
```

### Payroll Workflow
```
Salary Structure → Structure Assignment →
Payroll Entry → Salary Slip Generation →
Payment Processing
```

### Attendance Workflow
```
Shift Type → Shift Assignment → Employee Check-in →
Attendance Record → Attendance Request (if needed)
```

## Document Conventions

### Field Types Referenced
- **Link**: Reference to another document
- **Select**: Dropdown with predefined options
- **Data**: Single-line text input
- **Text**: Multi-line text input
- **Date/Datetime**: Date or date-time picker
- **Currency**: Monetary value
- **Int/Float**: Numeric values
- **Check**: Boolean checkbox
- **Table**: Child table with multiple rows
- **Attach**: File attachment
- **Read Only**: Calculated or display-only field

### Status Conventions
Most documents follow these common statuses:
- **Draft**: Editable, not submitted
- **Submitted**: Locked, officially recorded
- **Cancelled**: Reversed/voided

### Permission Levels
- **Create**: Can create new documents
- **Read**: Can view documents
- **Write**: Can edit documents
- **Submit**: Can submit documents
- **Cancel**: Can cancel documents
- **Amend**: Can amend cancelled documents

## Contributing to Documentation

This documentation should be updated whenever:
- New features are added
- Existing features are modified
- Business rules change
- UI/UX is updated

Maintain the three-perspective structure (designer/frontend/backend) for all updates.

---

**Version**: 1.0
**Last Updated**: 2025-12-18
**Application**: Frappe HRMS
