# HRMS API Documentation

Complete documentation for the HRMS (Human Resource Management System) API built with Go and Fiber.

## 📚 Documentation Index

### **[API Documentation](./API_DOCUMENTATION.md)**
Comprehensive API reference covering all 60+ endpoints across 4 modules:
- **Attendance Management**: Check-in/out, attendance marking, requests
- **Leave Management**: Applications, allocations, balances, encashment
- **Payroll**: Salary slips, structures, loans, advances, expense claims
- **Performance Management**: Appraisals, goals, 360° feedback, skill mapping

**Includes**:
- Complete endpoint reference with request/response examples
- Query parameter documentation
- Error codes and handling
- RBAC permission matrix
- cURL examples for quick testing

---

### **[Authentication Guide](./AUTHENTICATION.md)**
Deep dive into authentication and security:
- JWT token structure and lifecycle
- Registration and login flows with diagrams
- Token management best practices
- Role-Based Access Control (RBAC) implementation
- Security best practices for developers and API consumers
- Code examples in JavaScript, Python, and cURL

**Topics Covered**:
- Password hashing with bcrypt
- Token storage strategies (localStorage, cookies, memory)
- Token expiration and refresh strategies
- RBAC role hierarchy and permissions
- Security headers and CORS
- Rate limiting (100 req/min)

---

### **[Testing Guide](./TESTING_GUIDE.md)**
Step-by-step guide for testing all API endpoints:
- Environment setup instructions
- Postman collection usage
- cURL command examples
- Complete workflow testing scenarios
- Automated testing approaches
- Troubleshooting common issues

**Testing Scenarios**:
1. Complete attendance day (check-in to check-out)
2. Leave application and approval workflow
3. Salary slip generation process
4. Performance goal tracking lifecycle

---

### **[Postman Collection](./HRMS_API.postman_collection.json)**
Ready-to-import Postman collection with:
- All 60+ endpoints organized by module
- Pre-configured authentication
- Auto-token management (saves JWT after login)
- Environment variables for easy configuration
- Example request bodies for all endpoints

**Quick Start**:
1. Import `HRMS_API.postman_collection.json` into Postman
2. Update `base_url` variable (default: `http://localhost:8080`)
3. Run `Authentication → Register` to create test user
4. Run `Authentication → Login` (token auto-saved)
5. Test any endpoint (authentication handled automatically)

---

## 🚀 Quick Start

### 1. Start the API Server

```bash
cd hrms-go

# Setup environment
cp .env.example .env
# Edit .env with your database credentials

# Install dependencies (if needed)
go mod tidy

# Run server
go run cmd/api/main.go
```

Expected output:
```
Server starting on localhost:8080
```

### 2. Test with cURL

```bash
# Register
curl -X POST http://localhost:8080/api/method/hrms.api.register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecureP@ssw0rd123",
    "first_name": "Test",
    "last_name": "User",
    "role": "Employee"
  }'

# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/method/hrms.api.login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecureP@ssw0rd123"
  }' | jq -r '.data.token')

# Make authenticated request
curl -X GET "http://localhost:8080/api/method/hrms.hr.doctype.leave_type.leave_type.get_active" \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Test with Postman

1. **Import Collection**: `File → Import → Select HRMS_API.postman_collection.json`
2. **Run**: `Authentication → Login`
3. **Test**: Any endpoint (token auto-applied)

---

## 📊 API Overview

### Modules

| Module | Endpoints | Description |
|--------|-----------|-------------|
| **Attendance** | 9 | Employee check-in/out, attendance marking, requests, monthly summaries |
| **Leave** | 10 | Leave applications, approvals, balance tracking, encashment |
| **Payroll** | 14 | Salary slips, structures, loans, advances, expense claims |
| **Performance** | 17 | Appraisals, goals, 360° feedback, skill mapping |
| **Authentication** | 2 | Registration, login (JWT tokens) |

**Total Endpoints**: 60+

### Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP + JWT
       ▼
┌─────────────┐
│   Routes    │  (Fiber router + RBAC middleware)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Handlers   │  (HTTP layer - request/response)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Services   │  (Business logic)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│Repositories │  (Data access layer)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Database   │  (PostgreSQL + GORM)
└─────────────┘
```

### Technology Stack

- **Framework**: Go Fiber v2
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Middleware**: CORS, Rate Limiting, Request ID, Compression

---

## 🔐 Authentication

All endpoints (except `/register` and `/login`) require JWT authentication.

### Get Token

```bash
# Login
curl -X POST http://localhost:8080/api/method/hrms.api.login \
  -H "Content-Type: application/json" \
  -d '{"email": "your@email.com", "password": "yourpassword"}'
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {...}
  }
}
```

### Use Token

Include in `Authorization` header:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Token Expiration**: 24 hours (configurable via `JWT_EXPIRATION_HOURS`)

---

## 🔒 Role-Based Access Control (RBAC)

### Roles

1. **System Manager**: Full access
2. **HR Manager**: HR, Attendance, Leave, Performance
3. **Payroll Manager**: Payroll operations
4. **Expense Approver**: Expense claim approvals
5. **Employee**: Self-service operations

### Permission Examples

| Operation | System Manager | HR Manager | Payroll Manager | Employee |
|-----------|----------------|------------|-----------------|----------|
| Apply Leave | ✅ | ✅ | ✅ | ✅ |
| Approve Leave | ✅ | ✅ | ❌ | ❌ |
| Generate Salary Slip | ✅ | ✅ | ✅ | ❌ |
| View Own Salary Slip | ✅ | ✅ | ✅ | ✅ |
| Create Appraisal | ✅ | ✅ | ❌ | ❌ |
| Submit Expense Claim | ✅ | ✅ | ✅ | ✅ |

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md#rbac-role-based-access-control) for complete matrix.

---

## 📖 Usage Examples

### Example 1: Mark Attendance

```bash
curl -X POST http://localhost:8080/api/method/hrms.hr.doctype.attendance.attendance.mark_attendance \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "employee-uuid",
    "attendance_date": "2024-01-15",
    "status": "Present",
    "working_hours": 8.0
  }'
```

### Example 2: Apply for Leave

```bash
curl -X POST http://localhost:8080/api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "employee-uuid",
    "leave_type_id": "leave-type-uuid",
    "from_date": "2024-02-10",
    "to_date": "2024-02-12",
    "total_leave_days": 3,
    "description": "Family vacation"
  }'
```

### Example 3: Get Leave Balance

```bash
curl -X GET "http://localhost:8080/api/method/hrms.hr.doctype.leave_allocation.leave_allocation.get_leave_balance?employee_id=employee-uuid&leave_type_id=leave-type-uuid" \
  -H "Authorization: Bearer $TOKEN"
```

### Example 4: Create Performance Goal

```bash
curl -X POST http://localhost:8080/api/method/hrms.hr.doctype.goal.goal.create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "employee-uuid",
    "goal_name": "Increase sales by 20%",
    "start_date": "2024-01-01",
    "end_date": "2024-03-31",
    "progress": 0,
    "status": "Pending",
    "priority": "High"
  }'
```

---

## 🐛 Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Missing/invalid token | Re-login to get new token |
| 403 Forbidden | Insufficient permissions | Check user role requirements |
| 400 Bad Request | Invalid request data | Validate required fields and formats |
| 500 Internal Server | Server error | Check server logs and database connection |
| Rate Limited (429) | Too many requests | Wait for rate limit reset |

See [TESTING_GUIDE.md#troubleshooting](./TESTING_GUIDE.md#troubleshooting) for detailed solutions.

---

## 📝 API Response Format

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data
  },
  "message": "Operation successful" // Optional
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message"
  }
}
```

### Pagination Response

```json
{
  "success": true,
  "data": {
    "records": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

---

## 🔄 API Conventions

### Date Formats

- **Date**: `YYYY-MM-DD` (e.g., `2024-01-15`)
- **Timestamp**: ISO 8601 / RFC 3339 (e.g., `2024-01-15T10:00:00Z`)

### Query Parameters

- **Pagination**: `page` (default: 1), `limit` (default: 20)
- **Date Range**: `from_date`, `to_date`
- **Filters**: `status`, `employee_id`, `leave_type_id`, etc.

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (auth required) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 409 | Conflict (e.g., overlapping leave) |
| 429 | Too Many Requests (rate limited) |
| 500 | Internal Server Error |

---

## 📚 Additional Resources

### Project Structure

```
hrms-go/
├── cmd/api/                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/       # HTTP handlers
│   │   ├── middleware/     # Auth, RBAC, etc.
│   │   └── routes/         # Route definitions
│   ├── core/
│   │   ├── models/         # Database models
│   │   ├── repositories/   # Data access layer
│   │   └── services/       # Business logic
│   ├── config/             # Configuration
│   └── database.go         # Database setup
├── docs/                    # API documentation (you are here!)
└── .env                     # Environment variables
```

### Environment Variables

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

## 🤝 Contributing

When updating documentation:

1. **API Changes**: Update `API_DOCUMENTATION.md` with new endpoints
2. **Authentication Changes**: Update `AUTHENTICATION.md`
3. **Testing**: Add examples to `TESTING_GUIDE.md`
4. **Postman**: Update `HRMS_API.postman_collection.json`
5. **README**: Update this file with new sections/links

---

## 📞 Support

For questions or issues:

1. Check the relevant documentation file
2. Review the [Troubleshooting](#troubleshooting) section
3. Check server logs for error details
4. Verify environment configuration

---

## 📅 Documentation Version

- **Last Updated**: 2024-01-15
- **API Version**: v1.0.0
- **Go Version**: 1.21+
- **Fiber Version**: v2.52+

---

## 📑 Quick Links

- **[Complete API Reference](./API_DOCUMENTATION.md)** - All endpoints with examples
- **[Authentication & Security](./AUTHENTICATION.md)** - JWT, RBAC, security best practices
- **[Testing Guide](./TESTING_GUIDE.md)** - Step-by-step testing instructions
- **[Postman Collection](./HRMS_API.postman_collection.json)** - Import and test immediately

---

**Happy Coding! 🚀**
