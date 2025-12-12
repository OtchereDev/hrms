# HRMS API Testing Guide

## Table of Contents
- [Overview](#overview)
- [Setup](#setup)
- [Testing with Postman](#testing-with-postman)
- [Testing with cURL](#testing-with-curl)
- [Testing Workflows](#testing-workflows)
- [Automated Testing](#automated-testing)
- [Troubleshooting](#troubleshooting)

---

## Overview

This guide provides comprehensive instructions for testing the HRMS API endpoints. We'll cover:
- Setting up the testing environment
- Testing with Postman collection
- Manual testing with cURL
- Common workflows and scenarios
- Automated testing approaches

---

## Setup

### Prerequisites

1. **HRMS API Server Running**
   ```bash
   cd hrms-go
   go run cmd/api/main.go
   ```

   Expected output:
   ```
   Server starting on localhost:8080
   ```

2. **PostgreSQL Database Running**
   - Ensure PostgreSQL is running
   - Database credentials configured in `.env`

3. **Environment Variables**
   Create `.env` file in `hrms-go/` directory:
   ```env
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=hrms_db
   DB_SSLMODE=disable

   # JWT
   JWT_SECRET=your-super-secret-jwt-key-min-32-characters
   JWT_EXPIRATION_HOURS=24

   # Server
   SERVER_HOST=localhost
   SERVER_PORT=8080
   SERVER_ENV=development
   ```

---

## Testing with Postman

### Step 1: Import Collection

1. Open Postman
2. Click **Import** button
3. Select `docs/HRMS_API.postman_collection.json`
4. Collection will be imported with all endpoints organized by module

### Step 2: Configure Variables

The collection uses variables for easy testing:

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `base_url` | `http://localhost:8080` | API base URL |
| `access_token` | (empty) | JWT token (auto-filled after login) |
| `employee_id` | (empty) | Employee ID (auto-filled after register/login) |

**To modify variables:**
1. Click on the collection name
2. Go to **Variables** tab
3. Update the values as needed

### Step 3: Test Authentication Flow

#### 1. Register a Test User

**Endpoint**: `Authentication → Register`

**Pre-filled Request**:
```json
{
  "email": "test.user@example.com",
  "password": "SecureP@ssw0rd123",
  "first_name": "Test",
  "last_name": "User",
  "role": "Employee"
}
```

**Expected Response** (201 Created):
```json
{
  "success": true,
  "data": {
    "id": "550e8400-...",
    "email": "test.user@example.com",
    "first_name": "Test",
    "last_name": "User",
    "role": "Employee"
  }
}
```

**Auto-actions**:
- `employee_id` variable is automatically set from response

#### 2. Login

**Endpoint**: `Authentication → Login`

**Pre-filled Request**:
```json
{
  "email": "test.user@example.com",
  "password": "SecureP@ssw0rd123"
}
```

**Expected Response** (200 OK):
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "550e8400-...",
      "email": "test.user@example.com",
      "role": "Employee"
    }
  }
}
```

**Auto-actions**:
- `access_token` variable is automatically set
- All subsequent requests will use this token automatically

### Step 4: Test Module Endpoints

After successful login, test endpoints in this order:

#### **Attendance Module**

1. **Check In** (Morning)
   - Endpoint: `Attendance → Check In`
   - Change `log_type` to `"IN"`
   - Update `time` to current timestamp

2. **Get Today's Check-ins**
   - Endpoint: `Attendance → Get Today's Check-ins`
   - Should show the check-in record

3. **Check Out** (Evening)
   - Endpoint: `Attendance → Check In`
   - Change `log_type` to `"OUT"`
   - Update `time` to current timestamp

4. **List Attendance**
   - Endpoint: `Attendance → List Attendance`
   - Should show auto-marked attendance from check-in/out

#### **Leave Management Module**

1. **Get Active Leave Types**
   - Endpoint: `Leave Management → Get Active Leave Types`
   - Note a `leave_type_id` for next steps

2. **Apply for Leave**
   - Endpoint: `Leave Management → Apply for Leave`
   - Update `leave_type_id` with actual ID
   - Set future dates

3. **Get Leave Balance**
   - Endpoint: `Leave Management → Get Leave Balance`
   - Update `leave_type_id`

4. **List Leave Applications**
   - Endpoint: `Leave Management → List Leave Applications`
   - Should show your application

#### **Payroll Module** (Requires Payroll Manager Role)

1. **Get Active Salary Assignment**
   - Endpoint: `Payroll → Get Active Salary Assignment`

2. **Create Expense Claim**
   - Endpoint: `Payroll → Create Expense Claim`
   - Update dates and amounts

3. **List Expense Claims**
   - Endpoint: `Payroll → List Expense Claims`

#### **Performance Module**

1. **Create Goal**
   - Endpoint: `Performance Management → Create Goal`
   - Set realistic dates and targets

2. **Update Goal**
   - Endpoint: `Performance Management → Update Goal`
   - Update progress to 50%

3. **Get Active Goals**
   - Endpoint: `Performance Management → Get Active Goals for Employee`
   - Should show your goals

### Step 5: Test RBAC (Role-Based Access Control)

To test RBAC, you need users with different roles:

1. **Create HR Manager User**
   ```json
   {
     "email": "hr.manager@example.com",
     "password": "SecureP@ssw0rd123",
     "first_name": "HR",
     "last_name": "Manager",
     "role": "HR Manager"
   }
   ```

2. **Login as HR Manager**
   - Use HR Manager credentials
   - Test approval endpoints (Leave, Attendance Requests)

3. **Try Restricted Endpoints as Employee**
   - Login as regular employee
   - Try to approve leave → Should get 403 Forbidden

---

## Testing with cURL

### Authentication Flow

#### 1. Register
```bash
curl -X POST http://localhost:8080/api/method/hrms.api.register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecureP@ssw0rd123",
    "first_name": "Test",
    "last_name": "User",
    "role": "Employee"
  }'
```

#### 2. Login and Save Token
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/method/hrms.api.login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecureP@ssw0rd123"
  }' | jq -r '.data.token')

echo "Token: $TOKEN"
```

#### 3. Save Employee ID
```bash
EMPLOYEE_ID=$(curl -s -X POST http://localhost:8080/api/method/hrms.api.login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecureP@ssw0rd123"
  }' | jq -r '.data.user.employee_id')

echo "Employee ID: $EMPLOYEE_ID"
```

### Attendance Operations

#### Check In
```bash
curl -X POST http://localhost:8080/api/method/hrms.hr.doctype.employee_checkin.employee_checkin.checkin \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"employee_id\": \"$EMPLOYEE_ID\",
    \"log_type\": \"IN\",
    \"time\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",
    \"device_id\": \"CLI-001\"
  }"
```

#### List Attendance
```bash
curl -X GET "http://localhost:8080/api/method/hrms.hr.doctype.attendance.attendance.list_attendance?employee_id=$EMPLOYEE_ID&page=1&limit=20" \
  -H "Authorization: Bearer $TOKEN"
```

### Leave Operations

#### Get Leave Types
```bash
curl -X GET http://localhost:8080/api/method/hrms.hr.doctype.leave_type.leave_type.get_active \
  -H "Authorization: Bearer $TOKEN"
```

#### Apply for Leave
```bash
# Save a leave type ID first
LEAVE_TYPE_ID="your-leave-type-uuid"

curl -X POST http://localhost:8080/api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"employee_id\": \"$EMPLOYEE_ID\",
    \"leave_type_id\": \"$LEAVE_TYPE_ID\",
    \"from_date\": \"2024-02-10\",
    \"to_date\": \"2024-02-12\",
    \"total_leave_days\": 3,
    \"description\": \"Family vacation\"
  }"
```

---

## Testing Workflows

### Workflow 1: Complete Attendance Day

**Scenario**: Employee completes a full working day with check-in and check-out.

```bash
#!/bin/bash

# Setup
TOKEN="your-token"
EMPLOYEE_ID="your-employee-id"

# Morning check-in (9 AM)
echo "1. Checking in..."
curl -s -X POST http://localhost:8080/api/method/hrms.hr.doctype.employee_checkin.employee_checkin.checkin \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"employee_id\": \"$EMPLOYEE_ID\",
    \"log_type\": \"IN\",
    \"time\": \"2024-01-15T09:00:00Z\",
    \"device_id\": \"TEST-001\"
  }" | jq

# Get today's check-ins
echo -e "\n2. Getting today's check-ins..."
curl -s -X GET "http://localhost:8080/api/method/hrms.hr.doctype.employee_checkin.employee_checkin.get_today_checkins?employee_id=$EMPLOYEE_ID" \
  -H "Authorization: Bearer $TOKEN" | jq

# Evening check-out (6 PM)
echo -e "\n3. Checking out..."
curl -s -X POST http://localhost:8080/api/method/hrms.hr.doctype.employee_checkin.employee_checkin.checkin \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"employee_id\": \"$EMPLOYEE_ID\",
    \"log_type\": \"OUT\",
    \"time\": \"2024-01-15T18:00:00Z\",
    \"device_id\": \"TEST-001\"
  }" | jq

# Verify attendance was auto-marked
echo -e "\n4. Checking attendance..."
curl -s -X GET "http://localhost:8080/api/method/hrms.hr.doctype.attendance.attendance.list_attendance?employee_id=$EMPLOYEE_ID&from_date=2024-01-15&to_date=2024-01-15" \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Results**:
1. Check-in successful
2. Today's check-ins shows IN log
3. Check-out successful
4. Attendance auto-marked as "Present" with 9 working hours

### Workflow 2: Leave Application & Approval

**Scenario**: Employee applies for leave, HR Manager approves it.

```bash
#!/bin/bash

# Employee applies for leave
echo "1. Employee applying for leave..."
LEAVE_APP=$(curl -s -X POST http://localhost:8080/api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"employee_id\": \"$EMPLOYEE_ID\",
    \"leave_type_id\": \"$LEAVE_TYPE_ID\",
    \"from_date\": \"2024-02-10\",
    \"to_date\": \"2024-02-12\",
    \"total_leave_days\": 3,
    \"description\": \"Family vacation\"
  }")

LEAVE_ID=$(echo $LEAVE_APP | jq -r '.data.id')
echo "Leave application created: $LEAVE_ID"

# HR Manager approves leave
echo -e "\n2. HR Manager approving leave..."
curl -s -X POST http://localhost:8080/api/method/hrms.hr.doctype.leave_application.leave_application.approve \
  -H "Authorization: Bearer $HR_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"id\": \"$LEAVE_ID\"}" | jq

# Check updated leave balance
echo -e "\n3. Checking leave balance..."
curl -s -X GET "http://localhost:8080/api/method/hrms.hr.doctype.leave_allocation.leave_allocation.get_leave_balance?employee_id=$EMPLOYEE_ID&leave_type_id=$LEAVE_TYPE_ID" \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq
```

**Expected Results**:
1. Leave application created with status "Open"
2. HR Manager approves → status changes to "Approved"
3. Leave balance decreases by 3 days

### Workflow 3: Salary Slip Generation

**Scenario**: Payroll Manager generates salary slip for an employee.

```bash
#!/bin/bash

# Generate salary slip
echo "1. Generating salary slip..."
SALARY_SLIP=$(curl -s -X POST http://localhost:8080/api/method/hrms.payroll.doctype.salary_slip.salary_slip.generate \
  -H "Authorization: Bearer $PAYROLL_MANAGER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"employee_id\": \"$EMPLOYEE_ID\",
    \"start_date\": \"2024-01-01\",
    \"end_date\": \"2024-01-31\",
    \"posting_date\": \"2024-01-31\"
  }")

SLIP_ID=$(echo $SALARY_SLIP | jq -r '.data.id')
echo "Salary slip generated: $SLIP_ID"
echo $SALARY_SLIP | jq '.data | {gross_pay, total_deduction, net_pay}'

# Submit salary slip
echo -e "\n2. Submitting salary slip..."
curl -s -X POST http://localhost:8080/api/method/hrms.payroll.doctype.salary_slip.salary_slip.submit \
  -H "Authorization: Bearer $PAYROLL_MANAGER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"id\": \"$SLIP_ID\"}" | jq

# Employee views salary slip
echo -e "\n3. Employee viewing salary slip..."
curl -s -X GET "http://localhost:8080/api/method/hrms.payroll.doctype.salary_slip.salary_slip.get?id=$SLIP_ID" \
  -H "Authorization: Bearer $EMPLOYEE_TOKEN" | jq
```

**Expected Results**:
1. Salary slip generated with earnings, deductions, and net pay calculated
2. Status changed from "Draft" to "Submitted"
3. Employee can view their salary slip

### Workflow 4: Performance Goal Tracking

**Scenario**: Employee creates goal, updates progress, marks as achieved.

```bash
#!/bin/bash

# Create goal
echo "1. Creating performance goal..."
GOAL=$(curl -s -X POST http://localhost:8080/api/method/hrms.hr.doctype.goal.goal.create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"employee_id\": \"$EMPLOYEE_ID\",
    \"goal_name\": \"Increase sales by 20%\",
    \"description\": \"Q1 sales target\",
    \"start_date\": \"2024-01-01\",
    \"end_date\": \"2024-03-31\",
    \"progress\": 0,
    \"status\": \"Pending\",
    \"priority\": \"High\"
  }")

GOAL_ID=$(echo $GOAL | jq -r '.data.id')
echo "Goal created: $GOAL_ID"

# Update progress to 50%
echo -e "\n2. Updating goal progress to 50%..."
curl -s -X PUT http://localhost:8080/api/method/hrms.hr.doctype.goal.goal.update \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": \"$GOAL_ID\",
    \"progress\": 50,
    \"status\": \"In Progress\"
  }" | jq

# Update progress to 100%
echo -e "\n3. Marking goal as achieved (100%)..."
curl -s -X PUT http://localhost:8080/api/method/hrms.hr.doctype.goal.goal.update \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": \"$GOAL_ID\",
    \"progress\": 100,
    \"status\": \"Achieved\"
  }" | jq

# View updated goal
echo -e "\n4. Viewing final goal status..."
curl -s -X GET "http://localhost:8080/api/method/hrms.hr.doctype.goal.goal.get?id=$GOAL_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Results**:
1. Goal created with status "Pending"
2. Progress updated to 50% → status "In Progress"
3. Progress updated to 100% → status "Achieved"

---

## Automated Testing

### Unit Tests (Go)

Create test files for services:

**Example**: `hrms-go/internal/core/services/attendance/attendance_service_test.go`

```go
package attendance_test

import (
    "context"
    "testing"
    "time"

    "github.com/OtchereDev/hrms-go/internal/core/services/attendance"
    "github.com/stretchr/testify/assert"
)

func TestCheckin(t *testing.T) {
    // Setup
    service := setupTestService()
    ctx := context.Background()

    // Test check-in
    req := &attendance.CheckinRequest{
        EmployeeID: "test-employee-id",
        LogType:    "IN",
        Time:       time.Now(),
        DeviceID:   "TEST-001",
    }

    checkin, err := service.Checkin(ctx, req)

    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, checkin)
    assert.Equal(t, "IN", checkin.LogType)
}
```

Run tests:
```bash
cd hrms-go
go test ./internal/core/services/... -v
```

### Integration Tests

Create integration test suite:

**Example**: `hrms-go/tests/integration/attendance_test.go`

```go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAttendanceFlow(t *testing.T) {
    // Setup test server
    app := setupTestApp()

    // 1. Login
    loginBody := map[string]string{
        "email":    "test@example.com",
        "password": "password",
    }
    token := login(t, app, loginBody)

    // 2. Check in
    checkinBody := map[string]interface{}{
        "employee_id": "test-id",
        "log_type":    "IN",
        "time":        time.Now(),
    }
    resp := makeAuthRequest(t, app, "POST", "/api/method/hrms.hr.doctype.employee_checkin.employee_checkin.checkin", token, checkinBody)
    assert.Equal(t, http.StatusOK, resp.Code)

    // 3. Verify attendance created
    resp = makeAuthRequest(t, app, "GET", "/api/method/hrms.hr.doctype.attendance.attendance.list_attendance?employee_id=test-id", token, nil)
    assert.Equal(t, http.StatusOK, resp.Code)

    var result map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &result)
    assert.True(t, result["success"].(bool))
}
```

Run integration tests:
```bash
go test ./tests/integration/... -v
```

### Load Testing (with Apache Bench)

Test API performance:

```bash
# Test login endpoint
ab -n 1000 -c 10 -p login.json -T application/json http://localhost:8080/api/method/hrms.api.login

# Test authenticated endpoint
ab -n 1000 -c 10 -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/method/hrms.hr.doctype.attendance.attendance.list_attendance
```

### API Contract Testing (with Postman/Newman)

Run Postman collection from CLI:

```bash
# Install Newman
npm install -g newman

# Run collection
newman run docs/HRMS_API.postman_collection.json \
  --environment production.postman_environment.json \
  --reporters cli,json \
  --reporter-json-export results.json
```

---

## Troubleshooting

### Common Issues

#### 1. 401 Unauthorized

**Symptom**: All authenticated requests return 401

**Possible Causes**:
- Token not included in request
- Token expired (> 24 hours old)
- Invalid token format

**Solutions**:
```bash
# Check if token is set
echo $TOKEN

# Re-login to get new token
TOKEN=$(curl -s -X POST http://localhost:8080/api/method/hrms.api.login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password"}' \
  | jq -r '.data.token')

# Verify token format (should start with "eyJ")
echo $TOKEN | cut -c1-3
```

#### 2. 403 Forbidden

**Symptom**: Request returns 403

**Cause**: User role doesn't have permission for endpoint

**Solution**: Check required roles in API documentation and use appropriate user

#### 3. 400 Bad Request

**Symptom**: Request validation fails

**Causes**:
- Missing required fields
- Invalid data types
- Invalid date formats

**Solutions**:
```bash
# Check request body format
echo '{"employee_id": "uuid", ...}' | jq

# Verify date format (YYYY-MM-DD)
date +%Y-%m-%d
```

#### 4. 500 Internal Server Error

**Symptom**: Server returns 500

**Causes**:
- Database connection issues
- Missing environment variables
- Application bugs

**Solutions**:
1. Check server logs
2. Verify database connection
3. Check `.env` file

#### 5. Connection Refused

**Symptom**: Cannot connect to API

**Solutions**:
```bash
# Check if server is running
curl http://localhost:8080/health

# Check server logs
tail -f logs/app.log

# Restart server
cd hrms-go && go run cmd/api/main.go
```

### Debugging Tips

1. **Enable Detailed Logging**
   ```bash
   export LOG_LEVEL=debug
   go run cmd/api/main.go
   ```

2. **Use jq for Pretty JSON**
   ```bash
   curl ... | jq '.'
   ```

3. **Save Responses to File**
   ```bash
   curl ... > response.json
   cat response.json | jq '.'
   ```

4. **Check HTTP Headers**
   ```bash
   curl -v ... 2>&1 | grep "^<"
   ```

5. **Test Database Directly**
   ```bash
   psql -h localhost -U postgres -d hrms_db -c "SELECT * FROM users LIMIT 5;"
   ```

---

## Best Practices

1. **Use Test Data**
   - Create separate test accounts
   - Use test email domains (e.g., `test@example.com`)
   - Clean up test data after testing

2. **Test Edge Cases**
   - Empty fields
   - Invalid UUIDs
   - Past dates
   - Very large numbers

3. **Test Error Scenarios**
   - Duplicate data
   - Conflicting dates
   - Insufficient permissions

4. **Performance Testing**
   - Test with realistic data volumes
   - Monitor response times
   - Check memory usage

5. **Security Testing**
   - Test without authentication
   - Test with expired tokens
   - Test SQL injection attempts
   - Test XSS attempts

---

**Last Updated**: 2024-01-15
**Version**: 1.0.0
