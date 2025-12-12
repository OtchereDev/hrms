# Frappe HRMS API Compatibility Mapping

This document maps the original Frappe HRMS API endpoints to our Go implementation.

## Status Legend
- ✅ Implemented
- 🔄 Partially Implemented (different format)
- ❌ Not Implemented
- 📝 Planned

---

## Attendance Module

### Original Frappe Endpoints

| Method | Frappe Path | Status | Go Equivalent | Notes |
|--------|-------------|--------|---------------|-------|
| `get_events` | `hrms.hr.doctype.attendance.attendance.get_events` | ❌ | N/A | Calendar events for attendance |
| `mark_bulk_attendance` | `hrms.hr.doctype.attendance.attendance.mark_bulk_attendance` | ❌ | N/A | Bulk mark multiple dates |
| `get_unmarked_days` | `hrms.hr.doctype.attendance.attendance.get_unmarked_days` | ❌ | N/A | Get days without attendance |
| **DocType CRUD** | `frappe.client.get` | 🔄 | `GET /get_attendance` | Get single record |
| **DocType CRUD** | `frappe.client.get_list` | 🔄 | `GET /list_attendance` | List records |
| **DocType CRUD** | `frappe.client.insert` | 🔄 | `POST /mark_attendance` | Create record |
| **DocType CRUD** | `frappe.client.save` | ❌ | N/A | Update record |
| **DocType CRUD** | `frappe.client.submit` | ❌ | N/A | Submit (workflow) |
| **DocType CRUD** | `frappe.client.cancel` | ❌ | N/A | Cancel (workflow) |

### Employee Checkin

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `add_log_based_on_employee_field` | `hrms.hr.doctype.employee_checkin.employee_checkin.add_log_based_on_employee_field` | ❌ | N/A |
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /checkin` |
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /get_today_checkins` |

### Attendance Request

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /create_attendance_request` |
| **Approval** | Via Workflow | 🔄 | `POST /approve_attendance_request` |
| **Rejection** | Via Workflow | 🔄 | `POST /reject_attendance_request` |

---

## Leave Management Module

### Leave Application

| Method | Frappe Path | Status | Go Equivalent | Notes |
|--------|-------------|--------|---------------|-------|
| `get_leave_details` | `hrms.hr.doctype.leave_application.leave_application.get_leave_details` | ❌ | N/A | Get complete leave details |
| `get_number_of_leave_days` | `hrms.hr.doctype.leave_application.leave_application.get_number_of_leave_days` | ❌ | N/A | Calculate leave days |
| `get_leave_balance_on` | `hrms.hr.doctype.leave_application.leave_application.get_leave_balance_on` | 🔄 | `GET /get_leave_balance` | Different format |
| `get_leaves_for_period` | `hrms.hr.doctype.leave_application.leave_application.get_leaves_for_period` | ❌ | N/A | Get leaves in date range |
| `is_lwp` | `hrms.hr.doctype.leave_application.leave_application.is_lwp` | ❌ | N/A | Check if Leave Without Pay |
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /apply_leave` | Create leave application |
| **DocType CRUD** | `frappe.client.get` | ✅ | `GET /get` | Get leave application |
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /list` | List applications |
| **Approval** | Via Workflow/Button | ✅ | `POST /approve` | Approve application |
| **Rejection** | Via Workflow/Button | ✅ | `POST /reject` | Reject application |
| **Cancel** | Via Workflow/Button | ✅ | `POST /cancel` | Cancel application |

### Leave Allocation

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `get_leave_allocation_records` | `hrms.hr.doctype.leave_allocation.leave_allocation.get_leave_allocation_records` | ❌ | N/A |
| `get_carry_forwarded_leaves` | `hrms.hr.doctype.leave_allocation.leave_allocation.get_carry_forwarded_leaves` | ❌ | N/A |
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /allocate` |
| **DocType CRUD** | `frappe.client.get_list` | 🔄 | Included in `/get_leave_balance` |

### Leave Type

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /get_active` |

### Leave Encashment

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /create_encashment` |

---

## Payroll Module

### Salary Slip

| Method | Frappe Path | Status | Go Equivalent | Notes |
|--------|-------------|--------|---------------|-------|
| `get_salary_slip_details` | `hrms.payroll.doctype.salary_slip.salary_slip.get_salary_slip_details` | ❌ | N/A | Complex calculation |
| `calculate_net_pay` | `hrms.payroll.doctype.salary_slip.salary_slip.calculate_net_pay` | 🔄 | Built into `/generate` | Auto-calculated |
| `make_salary_slip` | `hrms.payroll.doctype.salary_slip.salary_slip.make_salary_slip` | ✅ | `POST /generate` | Generate slip |
| `email_salary_slip` | `hrms.payroll.doctype.salary_slip.salary_slip.email_salary_slip` | ❌ | N/A | Email functionality |
| **DocType CRUD** | `frappe.client.get` | ✅ | `GET /get` | Get salary slip |
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /list` | List slips |
| **Submit** | `frappe.client.submit` | ✅ | `POST /submit` | Submit slip |

### Salary Structure

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `make_salary_structure` | `hrms.payroll.doctype.salary_structure.salary_structure.make_salary_structure` | ❌ | N/A |
| **DocType CRUD** | `frappe.client.get_list` | ❌ | N/A |

### Salary Structure Assignment

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `create_salary_structure_assignment` | Various | ✅ | `POST /assign` | Assign structure |
| **DocType CRUD** | `frappe.client.get` | ✅ | `GET /get_active_assignment` |

### Loan

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `calculate_amounts` | `hrms.payroll.doctype.loan.loan.calculate_amounts` | ❌ | N/A | Calculate EMI |
| `request_loan` | `hrms.payroll.doctype.loan.loan.request_loan` | ✅ | `POST /create` |
| **Approval** | Via Workflow | ✅ | `POST /approve` |
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /list` |

### Employee Advance

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /create` |
| **Approval** | Via Workflow | ✅ | `POST /approve` |

### Expense Claim

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `get_expense_approver` | `hrms.payroll.doctype.expense_claim.expense_claim.get_expense_approver` | ❌ | N/A |
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /create` |
| **Approval** | Via Workflow | ✅ | `POST /approve` |
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /list` |

### Payroll Entry

| Method | Frappe Path | Status | Go Equivalent | Notes |
|--------|-------------|--------|---------------|-------|
| `get_start_end_dates` | `hrms.payroll.doctype.payroll_entry.payroll_entry.get_start_end_dates` | ❌ | N/A | Get payroll period |
| `get_employees` | `hrms.payroll.doctype.payroll_entry.payroll_entry.get_employees` | ❌ | N/A | Get employees for payroll |
| `create_salary_slips` | `hrms.payroll.doctype.payroll_entry.payroll_entry.create_salary_slips` | ❌ | N/A | Bulk create slips |
| `submit_salary_slips` | `hrms.payroll.doctype.payroll_entry.payroll_entry.submit_salary_slips` | ❌ | N/A | Bulk submit |
| **DocType CRUD** | `frappe.client.insert` | ❌ | N/A | Create payroll run |

---

## Performance Management Module

### Appraisal

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `calculate_total` | `hrms.hr.doctype.appraisal.appraisal.calculate_total` | 🔄 | Built into handlers | Auto-calculated |
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /create` |
| **DocType CRUD** | `frappe.client.get` | ✅ | `GET /get` |
| **DocType CRUD** | `frappe.client.save` | ✅ | `PUT /update` |
| **Submit** | `frappe.client.submit` | ✅ | `POST /submit` |
| **Complete** | Custom button | ✅ | `POST /complete` |
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /list` |

### Goal

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `get_goals_for_employee` | `hrms.hr.doctype.goal.goal.get_goals_for_employee` | ✅ | `GET /get_active_for_employee` |
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /create` |
| **DocType CRUD** | `frappe.client.get` | ✅ | `GET /get` |
| **DocType CRUD** | `frappe.client.save` | ✅ | `PUT /update` |
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /list` |

### Employee Performance Feedback

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| `get_feedback_history` | `hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.get_feedback_history` | 🔄 | `GET /get_for_employee` |
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /create` |

### Employee Skill Map

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| **DocType CRUD** | `frappe.client.insert` | ✅ | `POST /create_or_update` |
| **DocType CRUD** | `frappe.client.get` | ✅ | `GET /get_by_employee` |

### Appraisal Cycle

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /get_active` |

### Appraisal Template

| Method | Frappe Path | Status | Go Equivalent |
|--------|-------------|--------|---------------|
| **DocType CRUD** | `frappe.client.get_list` | ✅ | `GET /get_active` |

---

## Missing Modules (Not Implemented)

### Recruitment
- ❌ Job Opening
- ❌ Job Applicant
- ❌ Job Offer
- ❌ Interview
- ❌ Interview Round
- ❌ Interview Feedback

### Employee Lifecycle
- ❌ Employee Onboarding
- ❌ Employee Separation
- ❌ Employee Transfer
- ❌ Employee Promotion
- ❌ Exit Interview
- ❌ Full and Final Statement

### Shift Management
- ❌ Shift Type
- ❌ Shift Assignment
- ❌ Shift Request
- ❌ Shift Schedule

### Training
- ❌ Training Program
- ❌ Training Event
- ❌ Training Result
- ❌ Training Feedback

### Expense Management (Extended)
- ❌ Travel Request
- ❌ Employee Tax Exemption Declaration
- ❌ Employee Tax Exemption Proof Submission

### Reports
- ❌ Monthly Attendance Sheet
- ❌ Employee Leave Balance
- ❌ Salary Register
- ❌ Bank Remittance Report
- ❌ All standard Frappe reports

---

## Implementation Priority

### Phase 1: Critical (Core CRUD + Compatibility Layer)
1. **Frappe API Compatibility Layer** ⭐ **HIGHEST PRIORITY**
   - Create `/api/resource` endpoint (Frappe's standard)
   - Create `/api/method` endpoint (for whitelisted functions)
   - Implement Frappe response format
   - Session-based auth compatibility

2. **DocType Generic CRUD**
   - `frappe.client.get`
   - `frappe.client.get_list`
   - `frappe.client.insert`
   - `frappe.client.save`
   - `frappe.client.delete`
   - `frappe.client.submit`
   - `frappe.client.cancel`

3. **Missing Utility Functions**
   - All `get_*` helper methods
   - Calculation functions
   - Validation functions

### Phase 2: Workflows
1. Implement Frappe workflow engine equivalent
2. Approval/Rejection flows
3. State transitions

### Phase 3: Advanced Features
1. Real-time updates (Socket.IO)
2. Email notifications
3. PDF generation
4. Report generation
5. Calendar integration

### Phase 4: Missing Modules
1. Recruitment
2. Training
3. Shift Management
4. Employee Lifecycle

---

## Compatibility Layer Architecture

### Proposed Structure

```
hrms-go/
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── frappe_compat/     # NEW: Frappe compatibility layer
│   │   │   │   ├── resource.go    # Generic DocType CRUD
│   │   │   │   ├── method.go      # Whitelisted method calls
│   │   │   │   └── auth.go        # Session compatibility
│   │   │   ├── attendance.go      # Existing handlers
│   │   │   └── ...
│   │   └── routes/
│   │       └── frappe_routes.go   # NEW: Frappe-style routes
│   └── core/
│       └── frappe/                 # NEW: Frappe compatibility layer
│           ├── response.go         # Frappe response format
│           ├── session.go          # Session management
│           └── doctype.go          # DocType abstraction
```

### Example Implementation

**Frappe Response Format**:
```go
// internal/core/frappe/response.go
package frappe

type FrappeResponse struct {
    Message interface{} `json:"message"`
    Exc     interface{} `json:"exc"`
    ExcType string      `json:"exc_type,omitempty"`
}

func Success(data interface{}) *FrappeResponse {
    return &FrappeResponse{
        Message: data,
        Exc:     nil,
    }
}

func Error(err error, excType string) *FrappeResponse {
    return &FrappeResponse{
        Message: nil,
        Exc:     err.Error(),
        ExcType: excType,
    }
}
```

**Generic DocType Handler**:
```go
// internal/api/handlers/frappe_compat/resource.go
package frappe_compat

func (h *ResourceHandler) Get(c *fiber.Ctx) error {
    doctype := c.Query("doctype")
    name := c.Query("name")

    // Generic get based on doctype
    data, err := h.doctypeService.Get(c.Context(), doctype, name)
    if err != nil {
        return c.JSON(frappe.Error(err, "NotFoundError"))
    }

    return c.JSON(frappe.Success(data))
}

func (h *ResourceHandler) GetList(c *fiber.Ctx) error {
    doctype := c.Query("doctype")
    filters := c.Query("filters")
    fields := c.Query("fields")

    // Generic list based on doctype
    data, err := h.doctypeService.GetList(c.Context(), doctype, filters, fields)
    if err != nil {
        return c.JSON(frappe.Error(err, "QueryError"))
    }

    return c.JSON(frappe.Success(data))
}
```

---

## Estimated Effort

| Phase | Endpoints | Estimated Time | Priority |
|-------|-----------|----------------|----------|
| Compatibility Layer | Core framework | 2-3 weeks | **Critical** |
| Missing Utility Functions | ~50 functions | 3-4 weeks | High |
| Workflow Engine | Framework | 2 weeks | High |
| Missing CRUD Operations | ~30 operations | 2 weeks | Medium |
| Missing Modules | ~15 modules | 6-8 weeks | Low |
| Real-time & Advanced | Various | 3-4 weeks | Medium |
| **TOTAL** | **~100+ endpoints** | **18-24 weeks** | - |

---

## Next Steps

1. **Immediate**: Build Frappe compatibility layer
2. **Add missing utility functions** (get_leave_details, calculate functions)
3. **Implement generic DocType CRUD** (works for all doctypes)
4. **Build workflow engine**
5. **Add missing modules** as needed

This provides a roadmap for full Frappe HRMS API compatibility.
