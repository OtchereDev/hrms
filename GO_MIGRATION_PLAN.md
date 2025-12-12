# HRMS Go Migration Plan

## Project: Frappe HRMS → Go Fiber + GORM + PostgreSQL

**Timeline:** 16 weeks
**Start Date:** 2025-12-12
**Scope:** Complete port of 165 models, 210 API endpoints, all workflows & reports

---

## Technology Stack

### Backend
- **Framework:** Go Fiber v2 (Express-like for Go)
- **ORM:** GORM (PostgreSQL)
- **Authentication:** JWT (go-jwt)
- **Background Jobs:** Asynq (Redis-backed)
- **WebSocket:** Fiber WebSocket
- **File Storage:** Generic interface (S3 + Local)
- **Reports:** go-pdf, excelize
- **Email:** gomail
- **Cache:** Redis
- **Validation:** go-playground/validator

### Frontend (Existing)
- Vue 3 + Ionic (PWA)
- Vue 3 + TypeScript (Roster App)
- **No changes required** - API compatible

### Database
- **Production:** PostgreSQL 15+
- **Development:** PostgreSQL (Docker)

---

## Project Structure

```
hrms-go/
├── cmd/
│   ├── api/                    # Main API server
│   │   └── main.go
│   ├── scheduler/              # Background job worker
│   │   └── main.go
│   └── migrate/                # Database migration tool
│       └── main.go
│
├── internal/
│   ├── config/                 # Configuration management
│   │   ├── config.go
│   │   └── env.go
│   │
│   ├── core/                   # Core domain models & business logic
│   │   ├── models/             # GORM models organized by module
│   │   │   ├── base/           # Base models (User, Role, Session)
│   │   │   ├── hr/             # HR module models
│   │   │   │   ├── employee.go
│   │   │   │   ├── attendance.go
│   │   │   │   ├── leave.go
│   │   │   │   ├── recruitment.go
│   │   │   │   ├── performance.go
│   │   │   │   └── training.go
│   │   │   ├── payroll/        # Payroll module models
│   │   │   │   ├── salary.go
│   │   │   │   ├── payroll_entry.go
│   │   │   │   ├── tax.go
│   │   │   │   └── benefits.go
│   │   │   └── regional/       # Regional customizations
│   │   │       ├── ghana/
│   │   │       └── nigeria/
│   │   │
│   │   ├── services/           # Business logic services
│   │   │   ├── auth/
│   │   │   ├── employee/
│   │   │   ├── attendance/
│   │   │   ├── leave/
│   │   │   ├── payroll/
│   │   │   └── workflow/
│   │   │
│   │   └── repositories/       # Data access layer
│   │       ├── employee_repo.go
│   │       ├── attendance_repo.go
│   │       └── leave_repo.go
│   │
│   ├── api/                    # HTTP handlers & routes
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── employee.go
│   │   │   ├── attendance.go
│   │   │   ├── leave.go
│   │   │   └── payroll.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── permission.go
│   │   │   ├── logger.go
│   │   │   └── cors.go
│   │   ├── routes/
│   │   │   └── routes.go
│   │   └── validators/
│   │       └── validators.go
│   │
│   ├── workflow/               # Workflow engine
│   │   ├── engine.go
│   │   ├── state_machine.go
│   │   └── transitions.go
│   │
│   ├── storage/                # Generic file storage
│   │   ├── storage.go          # Interface
│   │   ├── local.go            # Local filesystem
│   │   └── s3.go               # S3/MinIO
│   │
│   ├── scheduler/              # Background jobs
│   │   ├── jobs/
│   │   │   ├── attendance.go   # Auto-attendance processing
│   │   │   ├── reminders.go    # Birthday, leave reminders
│   │   │   └── payroll.go      # Payroll calculations
│   │   └── scheduler.go
│   │
│   ├── reports/                # Report generation
│   │   ├── generator.go
│   │   ├── pdf.go
│   │   └── excel.go
│   │
│   ├── websocket/              # Real-time updates
│   │   ├── hub.go
│   │   └── client.go
│   │
│   ├── email/                  # Email service
│   │   ├── mailer.go
│   │   └── templates/
│   │
│   └── utils/                  # Utility functions
│       ├── date.go
│       ├── formula.go          # Formula evaluation engine
│       └── helpers.go
│
├── pkg/                        # Public reusable packages
│   ├── logger/
│   ├── errors/
│   └── response/
│
├── migrations/                 # SQL migrations
│   └── postgres/
│
├── scripts/                    # Utility scripts
│   ├── generate_models.go      # Generate GORM from DocType JSON
│   └── seed_data.go
│
├── docs/                       # API documentation
│   └── api.md
│
├── .env.example
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

---

## Implementation Phases

### **Phase 1: Foundation (Weeks 1-2)**

#### Week 1: Project Setup
- [x] Initialize Go module
- [x] Set up project structure
- [x] Create base configuration system
- [x] Set up PostgreSQL with GORM
- [x] Create base models (User, Role, Permission, Session)
- [x] Implement JWT authentication
- [x] Build RBAC middleware

#### Week 2: Core Infrastructure
- [x] Generic file storage interface (S3 + Local)
- [x] Workflow engine framework
- [x] Email service setup
- [x] Logger & error handling
- [x] WebSocket hub setup
- [x] Background job scheduler (Asynq)

### **Phase 2: Core HR Module (Weeks 3-6)**

#### Week 3: Employee Management
- [ ] Employee model & lifecycle
- [ ] Department, Designation, Branch models
- [ ] Employee Onboarding workflow
- [ ] Employee Transfer
- [ ] Employee Separation
- [ ] Employee Promotion
- [ ] APIs for employee CRUD

#### Week 4: Attendance System
- [ ] Attendance, Shift Type, Shift Assignment models
- [ ] Employee Checkin model
- [ ] Auto-attendance processing
- [ ] Late entry, early exit, overtime calculation
- [ ] Attendance Request workflow
- [ ] APIs for attendance tracking
- [ ] Scheduled job: Process auto-attendance

#### Week 5: Leave Management
- [ ] Leave Type, Leave Policy models
- [ ] Leave Application, Leave Allocation models
- [ ] Leave Encashment, Compensatory Leave
- [ ] Leave balance calculation
- [ ] Leave approval workflow
- [ ] APIs for leave management
- [ ] Scheduled job: Allocate earned leaves

#### Week 6: Recruitment & Performance
- [ ] Job Opening, Job Applicant, Job Offer models
- [ ] Interview, Interview Round models
- [ ] Appraisal, Appraisal Cycle models
- [ ] Goal, Employee Performance Feedback
- [ ] Training Program, Training Event models
- [ ] APIs for recruitment & performance
- [ ] Scheduled job: Close expired job openings

### **Phase 3: Payroll Module (Weeks 7-10)**

#### Week 7: Salary Structure
- [ ] Salary Component model
- [ ] Salary Structure, Salary Structure Assignment
- [ ] Additional Salary model
- [ ] Formula evaluation engine
- [ ] APIs for salary structure management

#### Week 8: Salary Slip Processing
- [ ] Salary Slip model
- [ ] Earnings & deductions calculation
- [ ] Pro-rated salary calculation
- [ ] Loan repayment integration
- [ ] Employee Benefits calculation
- [ ] Salary slip generation API

#### Week 9: Tax & Statutory
- [ ] Income Tax Slab model (Ghana)
- [ ] Employee Tax Exemption Declaration
- [ ] SSNIT (Social Security) calculation
- [ ] Tier 2 & Tier 3 pension calculations
- [ ] Tax calculation engine
- [ ] APIs for tax management

#### Week 10: Payroll Entry & Benefits
- [ ] Payroll Entry model
- [ ] Payroll Period model
- [ ] Bulk salary slip generation
- [ ] Gratuity calculation
- [ ] Retention Bonus, Employee Incentive
- [ ] Expense Claim, Employee Advance
- [ ] APIs for payroll processing

### **Phase 4: Additional Modules (Weeks 11-12)**

#### Week 11: Document Management & Workflows
- [ ] Document lifecycle hooks
- [ ] Workflow state machine
- [ ] Approval routing
- [ ] Notification system
- [ ] File attachment handling
- [ ] Document submission & cancellation

#### Week 12: Reports & Analytics
- [ ] Report generator framework
- [ ] PDF export (go-pdf)
- [ ] Excel export (excelize)
- [ ] Port 50+ custom reports:
  - Monthly Attendance Sheet
  - Salary Register
  - Leave Balance Report
  - Employee Birthday Report
  - Payroll Summary
  - etc.
- [ ] Dashboard charts data APIs

### **Phase 5: Regional & Integration (Weeks 13-14)**

#### Week 13: Ghana Regional Features
- [ ] Ghana tax slabs & rates
- [ ] SSNIT contribution rules
- [ ] GRA (Ghana Revenue Authority) compliance
- [ ] Bank payment formats (Ghana banks)
- [ ] Statutory reports for Ghana
- [ ] Holiday calendar for Ghana
- [ ] Ghana-specific validations

#### Week 14: West Africa Expansion Framework
- [ ] Nigeria tax & pension models
- [ ] Ivory Coast statutory models
- [ ] Senegal regional features
- [ ] Regional configuration system
- [ ] Multi-currency support
- [ ] Regional report templates

### **Phase 6: Frontend Integration (Week 15)**

#### Week 15: Vue App Integration
- [ ] Update API base URLs
- [ ] JWT token storage & refresh
- [ ] Update authentication flow
- [ ] Test all PWA features
- [ ] Test roster management app
- [ ] Fix any API compatibility issues
- [ ] Update file upload/download
- [ ] Test WebSocket notifications

### **Phase 7: Testing & Deployment (Week 16)**

#### Week 16: Final Testing
- [ ] Integration testing
- [ ] API endpoint testing (all 210 endpoints)
- [ ] Workflow testing
- [ ] Performance testing (load testing)
- [ ] Security testing (auth, permissions)
- [ ] Report generation testing
- [ ] Background job testing
- [ ] Docker deployment setup
- [ ] CI/CD pipeline
- [ ] Documentation

---

## Model Porting Strategy

### Total Models: 165

#### Base Models (10)
- User, Role, Permission, Session
- Company, Department, Branch, Designation
- Holiday List, Holiday

#### HR Module (90+)
**Employee Lifecycle:**
- Employee, Employee Onboarding, Employee Separation
- Employee Transfer, Employee Promotion, Employee Grade
- Exit Interview, Retention Bonus

**Attendance (12):**
- Attendance, Attendance Request
- Employee Checkin, Shift Type, Shift Assignment, Shift Request
- Operational Hours, Attendance Device Settings

**Leave (15):**
- Leave Type, Leave Policy, Leave Policy Assignment
- Leave Application, Leave Allocation, Leave Block List
- Leave Encashment, Compensatory Leave Request
- Leave Period, Leave Control Panel

**Recruitment (10):**
- Job Opening, Job Applicant, Job Offer
- Staffing Plan, Job Requisition
- Interview, Interview Round, Interview Feedback

**Performance (8):**
- Appraisal, Appraisal Template, Appraisal Cycle
- Goal, Employee Performance Feedback
- Employee Skill Map, Skill Assessment

**Training (6):**
- Training Program, Training Event
- Training Result, Training Feedback

**Other HR (30+):**
- Employee Education, Employee External Work History
- Employee Health Insurance, Employee Group Insurance
- Employee Grievance, Employee Referral
- Daily Work Summary, Employee Engagement Survey
- etc.

#### Payroll Module (45+)
**Salary (15):**
- Salary Component, Salary Structure
- Salary Structure Assignment, Salary Slip
- Additional Salary, Salary Detail

**Payroll (8):**
- Payroll Entry, Payroll Period
- Payroll Employee Detail, Payroll Settings

**Tax (12):**
- Income Tax Slab, Employee Tax Exemption Declaration
- Employee Tax Exemption Proof Submission
- Employee Other Income, Tax Withholding Category
- Employee Benefit Application, Employee Benefit Claim

**Benefits & Deductions (10):**
- Gratuity, Gratuity Slab
- Retention Bonus, Employee Incentive
- Loan, Loan Type, Loan Application
- Loan Repayment, Loan Security

**Expenses (5):**
- Expense Claim, Expense Claim Type
- Employee Advance, Employee Advance Payment

---

## API Endpoint Mapping

### Frappe → Fiber Route Pattern

**Frappe:**
```
POST /api/method/hrms.api.get_leave_applications
```

**Fiber (maintain compatibility):**
```
POST /api/method/hrms.api.get_leave_applications
```

### All Endpoints by Module

#### Authentication (5)
- `/api/method/login`
- `/api/method/logout`
- `/api/method/hrms.api.get_current_user_info`
- `/api/method/hrms.api.get_current_employee_info`
- `/api/method/refresh_token`

#### Employee (15)
- `/api/method/hrms.api.get_employees`
- `/api/method/hrms.api.create_employee`
- `/api/method/hrms.api.update_employee`
- etc.

#### Attendance (20)
- `/api/method/hrms.api.get_attendance_calendar_events`
- `/api/method/hrms.api.mark_attendance`
- `/api/method/hrms.hr.doctype.attendance.attendance.mark_bulk_attendance`
- etc.

#### Leave (25)
- `/api/method/hrms.api.get_leave_applications`
- `/api/method/hrms.api.apply_leave`
- `/api/method/hrms.api.get_leave_balance`
- etc.

#### Payroll (30)
- `/api/method/hrms.api.get_salary_slips`
- `/api/method/hrms.payroll.doctype.salary_slip.salary_slip.make_salary_slip`
- etc.

#### Reports (40)
- `/api/method/hrms.api.get_monthly_attendance`
- `/api/method/hrms.api.get_salary_register`
- etc.

#### File Upload (5)
- `/api/method/upload_file`
- `/api/method/download_file`
- etc.

#### WebSocket (5)
- `/ws/realtime`
- `/ws/notifications`

---

## Database Design

### Base Tables
```sql
-- Users & Auth
users
roles
permissions
role_permissions
user_roles
sessions

-- Organization
companies
departments
branches
designations
holiday_lists
holidays

-- Files
files
file_attachments
```

### HR Tables
```sql
-- Employee
employees
employee_onboardings
employee_separations
employee_transfers
employee_promotions

-- Attendance
attendances
attendance_requests
employee_checkins
shift_types
shift_assignments
shift_requests

-- Leave
leave_types
leave_policies
leave_policy_assignments
leave_applications
leave_allocations
leave_encashments

-- Recruitment
job_openings
job_applicants
job_offers
interviews
interview_rounds

-- Performance
appraisals
appraisal_cycles
goals
employee_performance_feedbacks

-- Training
training_programs
training_events
training_results
```

### Payroll Tables
```sql
-- Salary
salary_components
salary_structures
salary_structure_assignments
salary_slips
salary_details
additional_salaries

-- Payroll
payroll_entries
payroll_periods
payroll_employee_details

-- Tax
income_tax_slabs
employee_tax_exemption_declarations
employee_tax_exemption_proofs

-- Benefits
gratuities
retention_bonuses
employee_incentives

-- Expenses
expense_claims
employee_advances

-- Loans
loans
loan_types
loan_applications
loan_repayments
```

---

## Ghana Regional Customizations

### Tax System
- **Personal Income Tax:** Progressive rates (0% - 30%)
- **SSNIT (Social Security):** 5.5% employee + 13% employer
- **Tier 2 Pension:** 5% employee contribution
- **Tier 3 Pension:** Voluntary

### Tax Slabs (2024)
| Annual Income (GHS) | Rate |
|---------------------|------|
| 0 - 4,380           | 0%   |
| 4,381 - 120         | 5%   |
| 121 - 480           | 10%  |
| 481 - 3,840         | 17.5%|
| 3,841 - 240,000     | 25%  |
| Above 240,000       | 30%  |

### Statutory Reports
- Monthly SSNIT contribution report
- GRA tax remittance report
- Annual tax returns (P9 form)
- Employer annual return

---

## Development Workflow

### Local Development
```bash
# Start PostgreSQL & Redis
docker-compose up -d

# Run migrations
make migrate

# Seed data
make seed

# Start API server
make run-api

# Start scheduler
make run-scheduler

# Run tests
make test
```

### Code Generation
```bash
# Generate GORM models from DocType JSON
go run scripts/generate_models.go

# Generate API routes
go run scripts/generate_routes.go
```

### Testing
```bash
# Unit tests
make test-unit

# Integration tests
make test-integration

# E2E tests
make test-e2e

# Coverage
make coverage
```

---

## Performance Targets

- **API Response Time:** < 100ms (p95)
- **Report Generation:** < 5s for 1000 records
- **Salary Slip Calculation:** < 2s per employee
- **Concurrent Users:** 500+
- **Database Queries:** < 50ms (p95)

---

## Security Considerations

1. **JWT Security:**
   - Short-lived access tokens (15 min)
   - Long-lived refresh tokens (7 days)
   - Token rotation
   - Secure storage

2. **Permission System:**
   - Role-based access control
   - Document-level permissions
   - Field-level permissions (sensitive data)
   - Audit logs

3. **Data Protection:**
   - Encryption at rest (PostgreSQL)
   - Encryption in transit (TLS)
   - PII data masking in logs
   - GDPR compliance

4. **Rate Limiting:**
   - API rate limits (100 req/min per user)
   - Login attempt limits
   - File upload limits

---

## Deployment Architecture

```
┌─────────────────┐
│   Load Balancer │
│     (Nginx)     │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼──┐  ┌──▼───┐
│ API  │  │ API  │  (Horizontal scaling)
│Server│  │Server│
└───┬──┘  └──┬───┘
    │         │
    └────┬────┘
         │
    ┌────▼─────┐
    │PostgreSQL│
    │(Primary) │
    └────┬─────┘
         │
    ┌────▼─────┐
    │PostgreSQL│
    │ (Replica)│
    └──────────┘

┌─────────────┐
│   Redis     │
│  (Cache &   │
│   Queue)    │
└──────┬──────┘
       │
┌──────▼──────┐
│  Asynq      │
│  Worker     │
└─────────────┘

┌─────────────┐
│     S3      │
│  (Files)    │
└─────────────┘
```

---

## Next Steps

1. ✅ Review and approve this plan
2. Set up Go project structure
3. Initialize Git branch for Go migration
4. Start Phase 1: Foundation
5. Weekly progress reviews

---

**Estimated Completion:** 2025-04-11 (16 weeks from start)
**Developer:** Claude AI + Human Review
**Status:** Planning Phase
