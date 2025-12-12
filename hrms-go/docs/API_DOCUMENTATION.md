# HRMS API Documentation

## Table of Contents
- [Overview](#overview)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [Modules](#modules)
  - [Attendance](#attendance-module)
  - [Leave Management](#leave-management-module)
  - [Payroll](#payroll-module)
  - [Performance Management](#performance-management-module)

---

## Overview

The HRMS API is a RESTful API built with Go Fiber framework that provides comprehensive Human Resource Management capabilities. The API follows Frappe-compatible URL patterns for easy migration from Frappe HRMS.

**Base URL**: `http://localhost:8080`

**API Version**: v1

**Content-Type**: `application/json`

---

## Authentication

### Login

Authenticate and receive a JWT token for subsequent requests.

**Endpoint**: `POST /api/method/hrms.api.login`

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "your_password"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "HR Manager",
      "employee_id": "EMP-001"
    }
  }
}
```

### Using the Token

Include the JWT token in the `Authorization` header for all authenticated requests:

```
Authorization: Bearer <your_jwt_token>
```

### Register

Create a new user account.

**Endpoint**: `POST /api/method/hrms.api.register`

**Request Body**:
```json
{
  "email": "newuser@example.com",
  "password": "secure_password",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "Employee"
}
```

---

## Error Handling

### Standard Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message"
  }
}
```

### Common Error Codes

| Code | Status | Description |
|------|--------|-------------|
| `UNAUTHORIZED` | 401 | Invalid or missing authentication token |
| `FORBIDDEN` | 403 | Insufficient permissions for this operation |
| `NOT_FOUND` | 404 | Requested resource not found |
| `VALIDATION_ERROR` | 400 | Invalid request data |
| `CONFLICT` | 409 | Resource conflict (e.g., overlapping leave) |
| `INTERNAL_ERROR` | 500 | Server error |

---

## Rate Limiting

- **Rate Limit**: 100 requests per minute per IP
- **Headers**: Response includes rate limit headers:
  - `X-RateLimit-Limit`: Maximum requests allowed
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Time when limit resets (Unix timestamp)

---

## Modules

## Attendance Module

Manage employee attendance, check-ins/check-outs, and attendance requests.

### Mark Attendance

Mark attendance for an employee (Manual marking).

**Endpoint**: `POST /api/method/hrms.hr.doctype.attendance.attendance.mark_attendance`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "employee_id": "uuid",
  "attendance_date": "2024-01-15",
  "status": "Present",
  "shift": "Day Shift",
  "working_hours": 8.0,
  "late_entry": false,
  "early_exit": false
}
```

**Status Options**: `Present`, `Absent`, `Half Day`, `Work From Home`, `On Leave`

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "attendance_date": "2024-01-15",
    "status": "Present",
    "working_hours": 8.0,
    "created_at": "2024-01-15T09:00:00Z"
  }
}
```

### Get Attendance

Retrieve a specific attendance record.

**Endpoint**: `GET /api/method/hrms.hr.doctype.attendance.attendance.get_attendance`

**Required Role**: Authenticated users

**Query Parameters**:
- `id` (required): Attendance ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "attendance_date": "2024-01-15",
    "status": "Present",
    "shift": "Day Shift",
    "working_hours": 8.0,
    "check_in_time": "09:00:00",
    "check_out_time": "18:00:00",
    "late_entry": false,
    "early_exit": false
  }
}
```

### List Attendance

Get paginated list of attendance records with filters.

**Endpoint**: `GET /api/method/hrms.hr.doctype.attendance.attendance.list_attendance`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (optional): Filter by employee
- `from_date` (optional): Start date (YYYY-MM-DD)
- `to_date` (optional): End date (YYYY-MM-DD)
- `status` (optional): Filter by status
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Records per page

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "records": [
      {
        "id": "uuid",
        "employee_id": "uuid",
        "attendance_date": "2024-01-15",
        "status": "Present",
        "working_hours": 8.0
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

### Get Monthly Attendance Summary

Get attendance summary for an employee for a specific month.

**Endpoint**: `GET /api/method/hrms.hr.doctype.attendance.attendance.get_monthly_attendance`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (required): Employee ID
- `year` (required): Year (e.g., 2024)
- `month` (required): Month (1-12)

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "employee_id": "uuid",
    "year": 2024,
    "month": 1,
    "summary": {
      "total_days": 31,
      "working_days": 22,
      "present": 20,
      "absent": 1,
      "half_day": 1,
      "work_from_home": 2,
      "on_leave": 5,
      "late_entries": 2,
      "early_exits": 1,
      "total_working_hours": 160.0
    },
    "daily_records": [
      {
        "date": "2024-01-15",
        "status": "Present",
        "working_hours": 8.0
      }
    ]
  }
}
```

### Check In

Record employee check-in (clock in).

**Endpoint**: `POST /api/method/hrms.hr.doctype.employee_checkin.employee_checkin.checkin`

**Required Role**: Authenticated users

**Request Body**:
```json
{
  "employee_id": "uuid",
  "log_type": "IN",
  "time": "2024-01-15T09:00:00Z",
  "device_id": "DEVICE-001"
}
```

**Log Type Options**: `IN`, `OUT`

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "log_type": "IN",
    "time": "2024-01-15T09:00:00Z",
    "device_id": "DEVICE-001",
    "created_at": "2024-01-15T09:00:00Z"
  },
  "message": "Successfully checked in. Auto-marked attendance as Present."
}
```

**Note**: The system automatically marks attendance when check-in/check-out pairs are complete.

### Get Today's Check-ins

Retrieve all check-in/check-out logs for an employee today.

**Endpoint**: `GET /api/method/hrms.hr.doctype.employee_checkin.employee_checkin.get_today_checkins`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (required): Employee ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "employee_id": "uuid",
      "log_type": "IN",
      "time": "2024-01-15T09:00:00Z",
      "device_id": "DEVICE-001"
    },
    {
      "id": "uuid",
      "employee_id": "uuid",
      "log_type": "OUT",
      "time": "2024-01-15T18:00:00Z",
      "device_id": "DEVICE-001"
    }
  ]
}
```

### Create Attendance Request

Request attendance correction or approval.

**Endpoint**: `POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.create`

**Required Role**: Authenticated users

**Request Body**:
```json
{
  "employee_id": "uuid",
  "from_date": "2024-01-15",
  "to_date": "2024-01-15",
  "reason": "Forgot to check in, was present and working",
  "half_day": false
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "from_date": "2024-01-15",
    "to_date": "2024-01-15",
    "reason": "Forgot to check in, was present and working",
    "workflow_state": "Pending",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

**Workflow States**: `Pending`, `Approved`, `Rejected`

### Approve Attendance Request

Approve a pending attendance request.

**Endpoint**: `POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.approve`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "id": "uuid"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "workflow_state": "Approved",
    "approved_by": "uuid",
    "approved_at": "2024-01-15T14:00:00Z"
  },
  "message": "Attendance request approved successfully"
}
```

### Reject Attendance Request

Reject a pending attendance request.

**Endpoint**: `POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.reject`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "id": "uuid",
  "rejection_reason": "Insufficient documentation provided"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "workflow_state": "Rejected",
    "rejected_by": "uuid",
    "rejected_at": "2024-01-15T14:00:00Z",
    "rejection_reason": "Insufficient documentation provided"
  },
  "message": "Attendance request rejected"
}
```

---

## Leave Management Module

Manage leave applications, allocations, balances, and encashments.

### Apply for Leave

Submit a new leave application.

**Endpoint**: `POST /api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave`

**Required Role**: Authenticated users

**Request Body**:
```json
{
  "employee_id": "uuid",
  "leave_type_id": "uuid",
  "from_date": "2024-02-10",
  "to_date": "2024-02-12",
  "total_leave_days": 3,
  "description": "Family vacation",
  "half_day": false
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "leave_type_id": "uuid",
    "from_date": "2024-02-10",
    "to_date": "2024-02-12",
    "total_leave_days": 3,
    "description": "Family vacation",
    "status": "Open",
    "created_at": "2024-01-15T10:00:00Z"
  },
  "message": "Leave application submitted successfully"
}
```

**Leave Status Options**: `Open`, `Approved`, `Rejected`, `Cancelled`

**Validations**:
- Checks sufficient leave balance
- Detects overlapping leave applications
- Validates date ranges

### Get Leave Application

Retrieve a specific leave application.

**Endpoint**: `GET /api/method/hrms.hr.doctype.leave_application.leave_application.get`

**Required Role**: Authenticated users

**Query Parameters**:
- `id` (required): Leave Application ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "leave_type_id": "uuid",
    "leave_type_name": "Annual Leave",
    "from_date": "2024-02-10",
    "to_date": "2024-02-12",
    "total_leave_days": 3,
    "description": "Family vacation",
    "status": "Approved",
    "approved_by": "uuid",
    "approved_at": "2024-01-16T09:00:00Z"
  }
}
```

### List Leave Applications

Get paginated list of leave applications with filters.

**Endpoint**: `GET /api/method/hrms.hr.doctype.leave_application.leave_application.list`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (optional): Filter by employee
- `leave_type_id` (optional): Filter by leave type
- `status` (optional): Filter by status
- `from_date` (optional): Start date filter
- `to_date` (optional): End date filter
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Records per page

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "records": [
      {
        "id": "uuid",
        "employee_id": "uuid",
        "leave_type_name": "Annual Leave",
        "from_date": "2024-02-10",
        "to_date": "2024-02-12",
        "total_leave_days": 3,
        "status": "Approved"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

### Approve Leave Application

Approve a leave application.

**Endpoint**: `POST /api/method/hrms.hr.doctype.leave_application.leave_application.approve`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "id": "uuid"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "Approved",
    "approved_by": "uuid",
    "approved_at": "2024-01-16T09:00:00Z"
  },
  "message": "Leave application approved successfully"
}
```

### Reject Leave Application

Reject a leave application.

**Endpoint**: `POST /api/method/hrms.hr.doctype.leave_application.leave_application.reject`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "id": "uuid",
  "rejection_reason": "Insufficient leave balance"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "Rejected",
    "rejected_by": "uuid",
    "rejected_at": "2024-01-16T09:00:00Z"
  },
  "message": "Leave application rejected"
}
```

### Cancel Leave Application

Cancel an approved leave application.

**Endpoint**: `POST /api/method/hrms.hr.doctype.leave_application.leave_application.cancel`

**Required Role**: Authenticated users (own leaves), HR Manager (any leaves)

**Request Body**:
```json
{
  "id": "uuid",
  "cancellation_reason": "Plans changed"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "Cancelled",
    "cancelled_at": "2024-01-17T10:00:00Z"
  },
  "message": "Leave application cancelled successfully"
}
```

### Allocate Leave

Allocate leave days to an employee.

**Endpoint**: `POST /api/method/hrms.hr.doctype.leave_allocation.leave_allocation.allocate`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "employee_id": "uuid",
  "leave_type_id": "uuid",
  "from_date": "2024-01-01",
  "to_date": "2024-12-31",
  "new_leaves_allocated": 20,
  "description": "Annual leave allocation for 2024"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "leave_type_id": "uuid",
    "from_date": "2024-01-01",
    "to_date": "2024-12-31",
    "new_leaves_allocated": 20,
    "total_leaves_allocated": 20,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "Leave allocated successfully"
}
```

### Get Leave Balance

Get leave balance for an employee.

**Endpoint**: `GET /api/method/hrms.hr.doctype.leave_allocation.leave_allocation.get_leave_balance`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (required): Employee ID
- `leave_type_id` (required): Leave Type ID
- `date` (optional): Balance as of date (default: today)

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "employee_id": "uuid",
    "leave_type_id": "uuid",
    "leave_type_name": "Annual Leave",
    "total_allocated": 20,
    "used": 5,
    "pending": 2,
    "available": 13,
    "expired": 0
  }
}
```

**Balance Calculation**:
- `total_allocated`: Sum of all allocations
- `used`: Approved leave applications
- `pending`: Submitted but not yet approved
- `available`: total_allocated - used - pending
- `expired`: Leaves that expired (if leave type has expiry)

### Get Active Leave Types

Get all active leave types.

**Endpoint**: `GET /api/method/hrms.hr.doctype.leave_type.leave_type.get_active`

**Required Role**: Authenticated users

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "leave_type_name": "Annual Leave",
      "max_leaves_allowed": 20,
      "max_continuous_days_allowed": 10,
      "is_carry_forward": true,
      "is_optional_leave": false,
      "allow_negative": false,
      "is_paid_leave": true
    },
    {
      "id": "uuid",
      "leave_type_name": "Sick Leave",
      "max_leaves_allowed": 12,
      "max_continuous_days_allowed": 7,
      "is_carry_forward": false,
      "is_paid_leave": true
    }
  ]
}
```

### Create Leave Encashment

Request encashment of unused leave days.

**Endpoint**: `POST /api/method/hrms.hr.doctype.leave_encashment.leave_encashment.create`

**Required Role**: Authenticated users

**Request Body**:
```json
{
  "employee_id": "uuid",
  "leave_type_id": "uuid",
  "encashment_date": "2024-12-31",
  "leave_balance": 10,
  "encashable_days": 5
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "leave_type_id": "uuid",
    "encashment_date": "2024-12-31",
    "leave_balance": 10,
    "encashable_days": 5,
    "encashment_amount": 2500.00,
    "status": "Pending",
    "created_at": "2024-12-20T10:00:00Z"
  },
  "message": "Leave encashment request created successfully"
}
```

**Note**: Encashment amount is calculated based on the employee's per-day salary.

---

## Payroll Module

Manage salary structures, salary slips, loans, advances, and expense claims.

### Generate Salary Slip

Generate a salary slip for an employee.

**Endpoint**: `POST /api/method/hrms.payroll.doctype.salary_slip.salary_slip.generate`

**Required Role**: `HR Manager`, `System Manager`, `Payroll Manager`

**Request Body**:
```json
{
  "employee_id": "uuid",
  "start_date": "2024-01-01",
  "end_date": "2024-01-31",
  "posting_date": "2024-01-31"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "employee_name": "John Doe",
    "start_date": "2024-01-01",
    "end_date": "2024-01-31",
    "posting_date": "2024-01-31",
    "salary_structure_id": "uuid",
    "earnings": [
      {
        "salary_component": "Basic Salary",
        "amount": 50000.00
      },
      {
        "salary_component": "HRA",
        "amount": 20000.00
      }
    ],
    "deductions": [
      {
        "salary_component": "Tax",
        "amount": 5000.00
      },
      {
        "salary_component": "PF",
        "amount": 1800.00
      }
    ],
    "gross_pay": 70000.00,
    "total_deduction": 6800.00,
    "net_pay": 63200.00,
    "status": "Draft",
    "created_at": "2024-01-31T10:00:00Z"
  },
  "message": "Salary slip generated successfully"
}
```

**Salary Slip Status**: `Draft`, `Submitted`, `Paid`, `Cancelled`

**Calculation Logic**:
- Gets active salary assignment for employee
- Copies earnings and deductions from salary structure
- Applies additional salaries (bonuses, incentives)
- Deducts loan EMIs and employee advances
- Calculates taxes based on tax slabs
- `gross_pay` = Sum of all earnings
- `total_deduction` = Sum of all deductions
- `net_pay` = gross_pay - total_deduction

### Get Salary Slip

Retrieve a specific salary slip.

**Endpoint**: `GET /api/method/hrms.payroll.doctype.salary_slip.salary_slip.get`

**Required Role**: Authenticated users

**Query Parameters**:
- `id` (required): Salary Slip ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "start_date": "2024-01-01",
    "end_date": "2024-01-31",
    "gross_pay": 70000.00,
    "total_deduction": 6800.00,
    "net_pay": 63200.00,
    "status": "Submitted",
    "earnings": [...],
    "deductions": [...]
  }
}
```

### Submit Salary Slip

Submit a salary slip (marks it as ready for payment).

**Endpoint**: `POST /api/method/hrms.payroll.doctype.salary_slip.salary_slip.submit`

**Required Role**: `HR Manager`, `System Manager`, `Payroll Manager`

**Request Body**:
```json
{
  "id": "uuid"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "Submitted",
    "submitted_at": "2024-02-01T10:00:00Z"
  },
  "message": "Salary slip submitted successfully"
}
```

### List Salary Slips

Get paginated list of salary slips with filters.

**Endpoint**: `GET /api/method/hrms.payroll.doctype.salary_slip.salary_slip.list`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (optional): Filter by employee
- `start_date` (optional): From date
- `end_date` (optional): To date
- `status` (optional): Filter by status
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Records per page

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "records": [
      {
        "id": "uuid",
        "employee_id": "uuid",
        "employee_name": "John Doe",
        "start_date": "2024-01-01",
        "end_date": "2024-01-31",
        "gross_pay": 70000.00,
        "net_pay": 63200.00,
        "status": "Submitted"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

### Assign Salary Structure

Assign a salary structure to an employee.

**Endpoint**: `POST /api/method/hrms.payroll.doctype.salary_structure.salary_structure.assign`

**Required Role**: `HR Manager`, `System Manager`, `Payroll Manager`

**Request Body**:
```json
{
  "employee_id": "uuid",
  "salary_structure_id": "uuid",
  "from_date": "2024-01-01",
  "base_salary": 50000.00,
  "variable_salary": 20000.00
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "salary_structure_id": "uuid",
    "from_date": "2024-01-01",
    "base_salary": 50000.00,
    "variable_salary": 20000.00,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "Salary structure assigned successfully"
}
```

### Get Active Salary Assignment

Get the active salary assignment for an employee.

**Endpoint**: `GET /api/method/hrms.payroll.doctype.salary_structure.salary_structure.get_active_assignment`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (required): Employee ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "salary_structure_id": "uuid",
    "salary_structure_name": "Standard Salary Structure",
    "from_date": "2024-01-01",
    "base_salary": 50000.00,
    "variable_salary": 20000.00,
    "is_active": true
  }
}
```

### Create Loan

Create a loan application for an employee.

**Endpoint**: `POST /api/method/hrms.payroll.doctype.loan.loan.create`

**Required Role**: Authenticated users

**Request Body**:
```json
{
  "employee_id": "uuid",
  "loan_type": "Personal Loan",
  "loan_amount": 100000.00,
  "rate_of_interest": 10.0,
  "repayment_start_date": "2024-02-01",
  "repayment_periods": 12,
  "repayment_method": "Equal Principal",
  "monthly_repayment_amount": 8333.33
}
```

**Loan Types**: `Personal Loan`, `Housing Loan`, `Education Loan`, `Vehicle Loan`

**Repayment Methods**: `Equal Principal`, `Equal Principal + Interest`

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "loan_type": "Personal Loan",
    "loan_amount": 100000.00,
    "rate_of_interest": 10.0,
    "repayment_periods": 12,
    "monthly_repayment_amount": 8333.33,
    "total_payment": 105000.00,
    "status": "Pending",
    "created_at": "2024-01-15T10:00:00Z"
  },
  "message": "Loan application created successfully"
}
```

**Loan Status**: `Pending`, `Approved`, `Rejected`, `Disbursed`, `Repaid`, `Closed`

### Approve Loan

Approve a loan application.

**Endpoint**: `POST /api/method/hrms.payroll.doctype.loan.loan.approve`

**Required Role**: `HR Manager`, `System Manager`, `Payroll Manager`

**Request Body**:
```json
{
  "id": "uuid",
  "disbursement_date": "2024-01-20"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "Approved",
    "approved_by": "uuid",
    "approved_at": "2024-01-16T10:00:00Z",
    "disbursement_date": "2024-01-20"
  },
  "message": "Loan approved successfully"
}
```

### List Loans

Get paginated list of loans with filters.

**Endpoint**: `GET /api/method/hrms.payroll.doctype.loan.loan.list`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (optional): Filter by employee
- `loan_type` (optional): Filter by loan type
- `status` (optional): Filter by status
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Records per page

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "records": [
      {
        "id": "uuid",
        "employee_id": "uuid",
        "loan_type": "Personal Loan",
        "loan_amount": 100000.00,
        "outstanding_amount": 75000.00,
        "status": "Disbursed"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 25,
      "total_pages": 2
    }
  }
}
```

### Create Employee Advance

Request an advance payment.

**Endpoint**: `POST /api/method/hrms.payroll.doctype.employee_advance.employee_advance.create`

**Required Role**: Authenticated users

**Request Body**:
```json
{
  "employee_id": "uuid",
  "purpose": "Medical Emergency",
  "advance_amount": 10000.00,
  "advance_account": "Employee Advances - Main"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "purpose": "Medical Emergency",
    "advance_amount": 10000.00,
    "paid_amount": 0.00,
    "status": "Pending",
    "created_at": "2024-01-15T10:00:00Z"
  },
  "message": "Employee advance request created"
}
```

**Advance Status**: `Pending`, `Approved`, `Rejected`, `Paid`, `Returned`

### Approve Employee Advance

Approve an advance payment request.

**Endpoint**: `POST /api/method/hrms.payroll.doctype.employee_advance.employee_advance.approve`

**Required Role**: `HR Manager`, `System Manager`, `Payroll Manager`

**Request Body**:
```json
{
  "id": "uuid"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "Approved",
    "approved_by": "uuid",
    "approved_at": "2024-01-16T10:00:00Z"
  },
  "message": "Employee advance approved"
}
```

### Create Expense Claim

Submit an expense claim for reimbursement.

**Endpoint**: `POST /api/method/hrms.payroll.doctype.expense_claim.expense_claim.create`

**Required Role**: Authenticated users

**Request Body**:
```json
{
  "employee_id": "uuid",
  "expense_date": "2024-01-10",
  "total_claimed_amount": 5000.00,
  "expenses": [
    {
      "expense_type": "Travel",
      "description": "Client meeting travel",
      "amount": 3000.00
    },
    {
      "expense_type": "Food",
      "description": "Client lunch",
      "amount": 2000.00
    }
  ]
}
```

**Common Expense Types**: `Travel`, `Food`, `Accommodation`, `Office Supplies`, `Communication`, `Entertainment`

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "expense_date": "2024-01-10",
    "total_claimed_amount": 5000.00,
    "total_sanctioned_amount": 0.00,
    "total_amount_reimbursed": 0.00,
    "status": "Draft",
    "expenses": [
      {
        "expense_type": "Travel",
        "amount": 3000.00
      },
      {
        "expense_type": "Food",
        "amount": 2000.00
      }
    ],
    "created_at": "2024-01-15T10:00:00Z"
  },
  "message": "Expense claim created successfully"
}
```

**Expense Claim Status**: `Draft`, `Submitted`, `Approved`, `Rejected`, `Paid`

### Approve Expense Claim

Approve an expense claim.

**Endpoint**: `POST /api/method/hrms.payroll.doctype.expense_claim.expense_claim.approve`

**Required Role**: `HR Manager`, `System Manager`, `Expense Approver`

**Request Body**:
```json
{
  "id": "uuid",
  "sanctioned_amount": 4500.00
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "Approved",
    "total_sanctioned_amount": 4500.00,
    "approved_by": "uuid",
    "approved_at": "2024-01-16T10:00:00Z"
  },
  "message": "Expense claim approved"
}
```

### List Expense Claims

Get paginated list of expense claims with filters.

**Endpoint**: `GET /api/method/hrms.payroll.doctype.expense_claim.expense_claim.list`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (optional): Filter by employee
- `status` (optional): Filter by status
- `from_date` (optional): From date
- `to_date` (optional): To date
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Records per page

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "records": [
      {
        "id": "uuid",
        "employee_id": "uuid",
        "expense_date": "2024-01-10",
        "total_claimed_amount": 5000.00,
        "total_sanctioned_amount": 4500.00,
        "status": "Approved"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 30,
      "total_pages": 2
    }
  }
}
```

---

## Performance Management Module

Manage employee appraisals, goals, 360-degree feedback, and skill mapping.

### Create Appraisal

Create a new performance appraisal.

**Endpoint**: `POST /api/method/hrms.hr.doctype.appraisal.appraisal.create`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "employee_id": "uuid",
  "appraisal_template_id": "uuid",
  "start_date": "2024-01-01",
  "end_date": "2024-12-31",
  "kra_based_appraisal": true
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "appraisal_template_id": "uuid",
    "start_date": "2024-01-01",
    "end_date": "2024-12-31",
    "status": "Draft",
    "total_score": 0,
    "goals": [
      {
        "kra": "Sales Target",
        "per_weightage": 30,
        "score": 0,
        "score_earned": 0
      }
    ],
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "Appraisal created successfully"
}
```

**Appraisal Status**: `Draft`, `In Progress`, `Completed`, `Cancelled`

### Get Appraisal

Retrieve a specific appraisal.

**Endpoint**: `GET /api/method/hrms.hr.doctype.appraisal.appraisal.get`

**Required Role**: Authenticated users

**Query Parameters**:
- `id` (required): Appraisal ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "employee_name": "John Doe",
    "start_date": "2024-01-01",
    "end_date": "2024-12-31",
    "status": "Completed",
    "total_score": 85,
    "goals": [
      {
        "kra": "Sales Target",
        "per_weightage": 30,
        "score": 90,
        "score_earned": 27
      },
      {
        "kra": "Customer Satisfaction",
        "per_weightage": 25,
        "score": 85,
        "score_earned": 21.25
      }
    ],
    "remarks": "Excellent performance throughout the year"
  }
}
```

### Update Appraisal

Update appraisal scores and remarks.

**Endpoint**: `PUT /api/method/hrms.hr.doctype.appraisal.appraisal.update`

**Required Role**: `HR Manager`, `System Manager`, Reporting Manager

**Request Body**:
```json
{
  "id": "uuid",
  "goals": [
    {
      "id": "goal-uuid",
      "score": 90
    }
  ],
  "remarks": "Strong performance in Q4"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "total_score": 85,
    "updated_at": "2024-12-30T10:00:00Z"
  },
  "message": "Appraisal updated successfully"
}
```

### Submit Appraisal

Submit appraisal (changes status from Draft to In Progress).

**Endpoint**: `POST /api/method/hrms.hr.doctype.appraisal.appraisal.submit`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "id": "uuid"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "In Progress",
    "submitted_at": "2024-01-05T10:00:00Z"
  },
  "message": "Appraisal submitted successfully"
}
```

### Complete Appraisal

Mark appraisal as completed.

**Endpoint**: `POST /api/method/hrms.hr.doctype.appraisal.appraisal.complete`

**Required Role**: `HR Manager`, `System Manager`

**Request Body**:
```json
{
  "id": "uuid",
  "final_score": 85,
  "final_remarks": "Excellent performance"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "Completed",
    "final_score": 85,
    "completed_at": "2024-12-31T10:00:00Z"
  },
  "message": "Appraisal completed successfully"
}
```

### List Appraisals

Get paginated list of appraisals with filters.

**Endpoint**: `GET /api/method/hrms.hr.doctype.appraisal.appraisal.list`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (optional): Filter by employee
- `status` (optional): Filter by status
- `start_date` (optional): From date
- `end_date` (optional): To date
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Records per page

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "records": [
      {
        "id": "uuid",
        "employee_id": "uuid",
        "employee_name": "John Doe",
        "start_date": "2024-01-01",
        "end_date": "2024-12-31",
        "status": "Completed",
        "total_score": 85
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

### Create Goal

Create a performance goal for an employee.

**Endpoint**: `POST /api/method/hrms.hr.doctype.goal.goal.create`

**Required Role**: Authenticated users (own goals), HR Manager (any employee)

**Request Body**:
```json
{
  "employee_id": "uuid",
  "goal_name": "Increase sales by 20%",
  "description": "Achieve 20% growth in Q1 sales compared to Q4",
  "start_date": "2024-01-01",
  "end_date": "2024-03-31",
  "progress": 0,
  "status": "Pending",
  "priority": "High"
}
```

**Goal Status**: `Pending`, `In Progress`, `Achieved`, `Partially Achieved`, `Not Achieved`

**Priority Levels**: `Low`, `Medium`, `High`, `Critical`

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "goal_name": "Increase sales by 20%",
    "description": "Achieve 20% growth in Q1 sales compared to Q4",
    "start_date": "2024-01-01",
    "end_date": "2024-03-31",
    "progress": 0,
    "status": "Pending",
    "priority": "High",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "message": "Goal created successfully"
}
```

**Auto-Status Logic**:
- Progress 0% → Status: Pending
- Progress 1-99% → Status: In Progress
- Progress 100% → Status: Achieved

### Get Goal

Retrieve a specific goal.

**Endpoint**: `GET /api/method/hrms.hr.doctype.goal.goal.get`

**Required Role**: Authenticated users

**Query Parameters**:
- `id` (required): Goal ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "goal_name": "Increase sales by 20%",
    "description": "Achieve 20% growth in Q1 sales compared to Q4",
    "start_date": "2024-01-01",
    "end_date": "2024-03-31",
    "progress": 75,
    "status": "In Progress",
    "priority": "High"
  }
}
```

### Update Goal

Update goal progress and details.

**Endpoint**: `PUT /api/method/hrms.hr.doctype.goal.goal.update`

**Required Role**: Authenticated users (own goals), HR Manager (any employee)

**Request Body**:
```json
{
  "id": "uuid",
  "progress": 75,
  "status": "In Progress",
  "description": "Updated description with Q1 mid-quarter results"
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "progress": 75,
    "status": "In Progress",
    "updated_at": "2024-02-15T10:00:00Z"
  },
  "message": "Goal updated successfully"
}
```

### List Goals

Get paginated list of goals with filters.

**Endpoint**: `GET /api/method/hrms.hr.doctype.goal.goal.list`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (optional): Filter by employee
- `status` (optional): Filter by status
- `priority` (optional): Filter by priority
- `start_date` (optional): From date
- `end_date` (optional): To date
- `page` (optional, default: 1): Page number
- `limit` (optional, default: 20): Records per page

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "records": [
      {
        "id": "uuid",
        "employee_id": "uuid",
        "goal_name": "Increase sales by 20%",
        "start_date": "2024-01-01",
        "end_date": "2024-03-31",
        "progress": 75,
        "status": "In Progress",
        "priority": "High"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 25,
      "total_pages": 2
    }
  }
}
```

### Get Active Goals for Employee

Get all active (non-completed) goals for an employee.

**Endpoint**: `GET /api/method/hrms.hr.doctype.goal.goal.get_active_for_employee`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (required): Employee ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "goal_name": "Increase sales by 20%",
      "start_date": "2024-01-01",
      "end_date": "2024-03-31",
      "progress": 75,
      "status": "In Progress",
      "priority": "High"
    },
    {
      "id": "uuid",
      "goal_name": "Complete leadership training",
      "progress": 50,
      "status": "In Progress",
      "priority": "Medium"
    }
  ]
}
```

### Create 360-Degree Feedback

Submit performance feedback for an employee.

**Endpoint**: `POST /api/method/hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.create`

**Required Role**: Authenticated users

**Request Body**:
```json
{
  "employee_id": "uuid",
  "feedback_by": "uuid",
  "feedback_date": "2024-06-30",
  "overall_rating": 4.5,
  "criteria": [
    {
      "criteria_name": "Communication",
      "rating": 4.5,
      "feedback": "Excellent communicator, clear and concise"
    },
    {
      "criteria_name": "Teamwork",
      "rating": 5.0,
      "feedback": "Outstanding team player, always helpful"
    },
    {
      "criteria_name": "Technical Skills",
      "rating": 4.0,
      "feedback": "Strong technical abilities with room for growth"
    }
  ]
}
```

**Rating Scale**: 1.0 to 5.0 (1 = Poor, 5 = Excellent)

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "feedback_by": "uuid",
    "feedback_date": "2024-06-30",
    "overall_rating": 4.5,
    "criteria": [
      {
        "criteria_name": "Communication",
        "rating": 4.5,
        "feedback": "Excellent communicator"
      },
      {
        "criteria_name": "Teamwork",
        "rating": 5.0,
        "feedback": "Outstanding team player"
      }
    ],
    "created_at": "2024-06-30T10:00:00Z"
  },
  "message": "Feedback submitted successfully"
}
```

### Get Feedback for Employee

Get all feedback received by an employee.

**Endpoint**: `GET /api/method/hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.get_for_employee`

**Required Role**: Authenticated users (own feedback), HR Manager (any employee)

**Query Parameters**:
- `employee_id` (required): Employee ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "employee_id": "uuid",
    "average_rating": 4.3,
    "total_feedback_count": 5,
    "feedback": [
      {
        "id": "uuid",
        "feedback_by": "uuid",
        "feedback_by_name": "Jane Smith",
        "feedback_date": "2024-06-30",
        "overall_rating": 4.5,
        "criteria": [...]
      }
    ]
  }
}
```

### Create or Update Skill Map

Create or update an employee's skill map.

**Endpoint**: `POST /api/method/hrms.hr.doctype.employee_skill_map.employee_skill_map.create_or_update`

**Required Role**: Authenticated users (own skill map), HR Manager (any employee)

**Request Body**:
```json
{
  "employee_id": "uuid",
  "skills": [
    {
      "skill_name": "Python",
      "proficiency": "Expert",
      "years_of_experience": 5
    },
    {
      "skill_name": "Leadership",
      "proficiency": "Intermediate",
      "years_of_experience": 2
    }
  ],
  "trainings": [
    {
      "training_name": "Advanced Python Programming",
      "completion_date": "2023-06-15"
    }
  ]
}
```

**Proficiency Levels**: `Beginner`, `Intermediate`, `Advanced`, `Expert`

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "skills": [
      {
        "skill_name": "Python",
        "proficiency": "Expert",
        "years_of_experience": 5
      },
      {
        "skill_name": "Leadership",
        "proficiency": "Intermediate",
        "years_of_experience": 2
      }
    ],
    "trainings": [
      {
        "training_name": "Advanced Python Programming",
        "completion_date": "2023-06-15"
      }
    ],
    "updated_at": "2024-01-15T10:00:00Z"
  },
  "message": "Skill map updated successfully"
}
```

### Get Skill Map by Employee

Retrieve an employee's skill map.

**Endpoint**: `GET /api/method/hrms.hr.doctype.employee_skill_map.employee_skill_map.get_by_employee`

**Required Role**: Authenticated users

**Query Parameters**:
- `employee_id` (required): Employee ID

**Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employee_id": "uuid",
    "employee_name": "John Doe",
    "skills": [
      {
        "skill_name": "Python",
        "proficiency": "Expert",
        "years_of_experience": 5
      },
      {
        "skill_name": "Go",
        "proficiency": "Advanced",
        "years_of_experience": 3
      }
    ],
    "trainings": [
      {
        "training_name": "Advanced Python Programming",
        "completion_date": "2023-06-15"
      }
    ],
    "total_skills": 2,
    "total_trainings": 1
  }
}
```

### Get Active Appraisal Cycles

Get all active appraisal cycles.

**Endpoint**: `GET /api/method/hrms.hr.doctype.appraisal_cycle.appraisal_cycle.get_active`

**Required Role**: Authenticated users

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "cycle_name": "Annual Appraisal 2024",
      "start_date": "2024-01-01",
      "end_date": "2024-12-31",
      "description": "Annual performance review cycle",
      "is_active": true
    }
  ]
}
```

### Get Active Appraisal Templates

Get all active appraisal templates.

**Endpoint**: `GET /api/method/hrms.hr.doctype.appraisal_template.appraisal_template.get_active`

**Required Role**: Authenticated users

**Response** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "template_name": "Sales Team Appraisal",
      "description": "Performance appraisal template for sales team",
      "kra_based_appraisal": true,
      "goals": [
        {
          "kra": "Sales Target",
          "per_weightage": 30
        },
        {
          "kra": "Customer Satisfaction",
          "per_weightage": 25
        }
      ],
      "rating_scale": [
        {
          "rating": "Outstanding",
          "score_from": 90,
          "score_to": 100
        },
        {
          "rating": "Exceeds Expectations",
          "score_from": 75,
          "score_to": 89
        }
      ]
    }
  ]
}
```

---

## RBAC (Role-Based Access Control)

### Role Hierarchy

1. **System Manager**: Full access to all modules and operations
2. **HR Manager**: Full access to HR, Attendance, Leave, Performance modules
3. **Payroll Manager**: Full access to Payroll module
4. **Expense Approver**: Can approve expense claims
5. **Employee**: Limited access to own records and self-service operations

### Permission Matrix

| Operation | System Manager | HR Manager | Payroll Manager | Employee |
|-----------|----------------|------------|-----------------|----------|
| Mark Attendance (any) | ✓ | ✓ | ✗ | ✗ |
| View Own Attendance | ✓ | ✓ | ✓ | ✓ |
| Apply Leave | ✓ | ✓ | ✓ | ✓ |
| Approve Leave | ✓ | ✓ | ✗ | ✗ |
| Generate Salary Slip | ✓ | ✓ | ✓ | ✗ |
| View Own Salary Slip | ✓ | ✓ | ✓ | ✓ |
| Create Appraisal | ✓ | ✓ | ✗ | ✗ |
| View Own Appraisal | ✓ | ✓ | ✓ | ✓ |
| Approve Loan | ✓ | ✓ | ✓ | ✗ |
| Apply for Loan | ✓ | ✓ | ✓ | ✓ |
| Approve Expense Claim | ✓ | ✓ | ✗ | With "Expense Approver" role |

---

## Testing with cURL

### Example: Login and Get Token

```bash
# Login
curl -X POST http://localhost:8080/api/method/hrms.api.login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'

# Save the token from response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Example: Mark Attendance

```bash
curl -X POST http://localhost:8080/api/method/hrms.hr.doctype.attendance.attendance.mark_attendance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "employee_id": "employee-uuid",
    "attendance_date": "2024-01-15",
    "status": "Present",
    "working_hours": 8.0
  }'
```

### Example: Apply for Leave

```bash
curl -X POST http://localhost:8080/api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "employee_id": "employee-uuid",
    "leave_type_id": "leave-type-uuid",
    "from_date": "2024-02-10",
    "to_date": "2024-02-12",
    "total_leave_days": 3,
    "description": "Family vacation"
  }'
```

### Example: Generate Salary Slip

```bash
curl -X POST http://localhost:8080/api/method/hrms.payroll.doctype.salary_slip.salary_slip.generate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "employee_id": "employee-uuid",
    "start_date": "2024-01-01",
    "end_date": "2024-01-31",
    "posting_date": "2024-01-31"
  }'
```

### Example: Get Leave Balance

```bash
curl -X GET "http://localhost:8080/api/method/hrms.hr.doctype.leave_allocation.leave_allocation.get_leave_balance?employee_id=employee-uuid&leave_type_id=leave-type-uuid" \
  -H "Authorization: Bearer $TOKEN"
```

---

## API Best Practices

### 1. Always Include Authorization Header
```
Authorization: Bearer <your_jwt_token>
```

### 2. Handle Rate Limiting
Check response headers and implement exponential backoff if rate limited.

### 3. Validate Request Data
Always validate required fields before making requests to avoid 400 errors.

### 4. Use Pagination
For list endpoints, use `page` and `limit` parameters to manage large datasets efficiently.

### 5. Error Handling
Always check the `success` field in responses and handle errors appropriately:

```javascript
if (!response.success) {
  console.error(response.error.message);
  // Handle error
}
```

### 6. Date Formats
- Use ISO 8601 format for dates: `YYYY-MM-DD`
- Use RFC 3339 format for timestamps: `2024-01-15T10:00:00Z`

### 7. Filtering
Use query parameters for filtering list endpoints:
- Single value: `status=Approved`
- Date ranges: `from_date=2024-01-01&to_date=2024-12-31`
- Pagination: `page=1&limit=20`

---

## Support and Contact

For API support, please contact:
- **Email**: support@hrms.example.com
- **Documentation**: https://docs.hrms.example.com
- **Status Page**: https://status.hrms.example.com

---

**Last Updated**: 2024-01-15
**API Version**: v1.0.0
