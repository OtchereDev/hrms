# Mobile App Module - Backend Documentation

## Overview

Mobile application functionality for employee self-service via iOS and Android apps.

## API Endpoints

### Authentication
- `POST /api/method/login`: Employee login with credentials
- `POST /api/method/logout`: Session logout
- `GET /api/method/frappe.auth.get_logged_user`: Current user info

### Employee Self-Service
- `GET /api/resource/Employee/{id}`: Employee profile
- `PUT /api/resource/Employee/{id}`: Update profile
- `GET /api/method/hrms.mobile.get_employee_dashboard`: Dashboard data

### Attendance & Check-in
- `POST /api/resource/Employee Checkin`: Create check-in/out
  - Fields: `employee`, `time`, `device_id`, `log_type` (IN/OUT)
  - Geolocation: `latitude`, `longitude`
- `GET /api/method/hrms.mobile.get_attendance_summary`: Monthly summary

### Leave Management
- `GET /api/method/hrms.mobile.get_leave_balance`: Leave balances
- `GET /api/resource/Leave Application`: Leave applications list
- `POST /api/resource/Leave Application`: Apply for leave
- `PUT /api/resource/Leave Application/{id}`: Update/cancel leave

### Expense & Travel
- `POST /api/resource/Expense Claim`: Create expense claim
  - Support for image uploads (receipts)
- `GET /api/method/hrms.mobile.get_expense_claims`: List claims
- `POST /api/resource/Travel Request`: Create travel request

### Approvals
- `GET /api/method/hrms.mobile.get_pending_approvals`: Items requiring approval
- `POST /api/method/hrms.mobile.approve_document`: Approve leave/expense
- `POST /api/method/hrms.mobile.reject_document`: Reject with comments

### Payroll
- `GET /api/method/hrms.mobile.get_salary_slips`: Salary slip history
- `GET /api/resource/Salary Slip/{id}`: Download salary slip PDF

### Team Management (for managers)
- `GET /api/method/hrms.mobile.get_team_members`: Direct reports
- `GET /api/method/hrms.mobile.get_team_attendance`: Team attendance
- `GET /api/method/hrms.mobile.get_team_leaves`: Team leave calendar

## Mobile-Specific Features

### Geolocation Validation
```python
@frappe.whitelist()
def validate_checkin_location(latitude, longitude, employee):
    employee_doc = frappe.get_doc("Employee", employee)
    office_location = get_office_location(employee_doc.branch)
    
    distance = calculate_distance(
        (latitude, longitude),
        (office_location.latitude, office_location.longitude)
    )
    
    max_distance = frappe.db.get_single_value(
        "HR Settings", 
        "max_checkin_distance_meters"
    )
    
    if distance > max_distance:
        frappe.throw("You are too far from office location")
    
    return True
```

### Push Notifications
- Leave approval/rejection
- Expense claim status
- Salary slip available
- Attendance regularization reminders
- Birthday/anniversary wishes
- Company announcements

### Offline Support
- Cache employee data
- Queue check-ins when offline
- Sync when connectivity restored

## Mobile App Settings

### Fields
- `enable_geofencing`: Require location for check-in
- `geofence_radius`: Allowed distance in meters
- `allow_photo_checkin`: Require selfie on check-in
- `biometric_authentication`: Face ID/Touch ID support
- `session_timeout`: Auto-logout after inactivity

## File Upload Handling
- Receipt images for expense claims
- Profile photo updates
- Document attachments
- Compression before upload
- Thumbnail generation

## Security
- JWT token authentication
- Role-based access control
- API rate limiting
- Device fingerprinting
- Encrypted local storage
