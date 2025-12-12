# Authentication & Security Guide

## Table of Contents
- [Overview](#overview)
- [Authentication Flow](#authentication-flow)
- [JWT Token Structure](#jwt-token-structure)
- [Token Management](#token-management)
- [Role-Based Access Control (RBAC)](#role-based-access-control-rbac)
- [Security Best Practices](#security-best-practices)
- [API Endpoints](#api-endpoints)
- [Code Examples](#code-examples)

---

## Overview

The HRMS API uses **JWT (JSON Web Tokens)** for authentication. All authenticated endpoints require a valid JWT token in the Authorization header.

**Security Features**:
- Password hashing with bcrypt
- JWT-based stateless authentication
- Role-based access control (RBAC)
- Token expiration (24 hours default)
- Rate limiting (100 requests/minute)
- CORS protection
- Request ID tracking

---

## Authentication Flow

### 1. User Registration

```
┌─────────┐                    ┌─────────┐                    ┌──────────┐
│ Client  │                    │   API   │                    │ Database │
└────┬────┘                    └────┬────┘                    └────┬─────┘
     │                              │                              │
     │  POST /api/method/          │                              │
     │  hrms.api.register          │                              │
     ├────────────────────────────>│                              │
     │  {email, password, ...}     │                              │
     │                              │                              │
     │                              │  Hash password with bcrypt   │
     │                              ├──────────┐                   │
     │                              │          │                   │
     │                              │<─────────┘                   │
     │                              │                              │
     │                              │  INSERT user record          │
     │                              ├─────────────────────────────>│
     │                              │                              │
     │                              │  User created                │
     │                              │<─────────────────────────────┤
     │                              │                              │
     │  201 Created                 │                              │
     │  {success: true, data: {...}}│                              │
     │<─────────────────────────────┤                              │
     │                              │                              │
```

**Request**:
```http
POST /api/method/hrms.api.register
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "SecureP@ssw0rd",
  "first_name": "John",
  "last_name": "Doe",
  "role": "Employee"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "Employee",
    "created_at": "2024-01-15T10:00:00Z"
  },
  "message": "User registered successfully"
}
```

**Notes**:
- Password must be at least 8 characters
- Password is hashed using bcrypt (cost factor: 10)
- Email must be unique
- Default role is "Employee" if not specified

### 2. User Login

```
┌─────────┐                    ┌─────────┐                    ┌──────────┐
│ Client  │                    │   API   │                    │ Database │
└────┬────┘                    └────┬────┘                    └────┬─────┘
     │                              │                              │
     │  POST /api/method/          │                              │
     │  hrms.api.login             │                              │
     ├────────────────────────────>│                              │
     │  {email, password}           │                              │
     │                              │                              │
     │                              │  SELECT user by email        │
     │                              ├─────────────────────────────>│
     │                              │                              │
     │                              │  User record                 │
     │                              │<─────────────────────────────┤
     │                              │                              │
     │                              │  Compare password hash       │
     │                              ├──────────┐                   │
     │                              │          │                   │
     │                              │<─────────┘                   │
     │                              │                              │
     │                              │  Generate JWT token          │
     │                              ├──────────┐                   │
     │                              │          │                   │
     │                              │<─────────┘                   │
     │                              │                              │
     │  200 OK                      │                              │
     │  {token, user}               │                              │
     │<─────────────────────────────┤                              │
     │                              │                              │
     │  Store token in              │                              │
     │  localStorage/cookie         │                              │
     ├──────────┐                   │                              │
     │          │                   │                              │
     │<─────────┘                   │                              │
```

**Request**:
```http
POST /api/method/hrms.api.login
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "password": "SecureP@ssw0rd"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwiZW1haWwiOiJqb2huLmRvZUBleGFtcGxlLmNvbSIsInJvbGUiOiJFbXBsb3llZSIsImV4cCI6MTcwNTQwNjQwMH0.signature",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john.doe@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "Employee",
      "employee_id": "EMP-001"
    }
  }
}
```

**Error Response (Invalid Credentials)**:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid email or password"
  }
}
```

### 3. Authenticated Requests

```
┌─────────┐                    ┌─────────┐                    ┌──────────┐
│ Client  │                    │   API   │                    │ Database │
└────┬────┘                    └────┬────┘                    └────┬─────┘
     │                              │                              │
     │  GET /api/method/...         │                              │
     │  Authorization: Bearer <JWT> │                              │
     ├────────────────────────────>│                              │
     │                              │                              │
     │                              │  Extract & Verify JWT        │
     │                              ├──────────┐                   │
     │                              │          │                   │
     │                              │<─────────┘                   │
     │                              │                              │
     │                              │  Check token expiration      │
     │                              ├──────────┐                   │
     │                              │          │                   │
     │                              │<─────────┘                   │
     │                              │                              │
     │                              │  Verify signature            │
     │                              ├──────────┐                   │
     │                              │          │                   │
     │                              │<─────────┘                   │
     │                              │                              │
     │                              │  Check RBAC permissions      │
     │                              ├──────────┐                   │
     │                              │          │                   │
     │                              │<─────────┘                   │
     │                              │                              │
     │                              │  Execute business logic      │
     │                              ├─────────────────────────────>│
     │                              │                              │
     │                              │  Query result                │
     │                              │<─────────────────────────────┤
     │                              │                              │
     │  200 OK                      │                              │
     │  {success: true, data: {...}}│                              │
     │<─────────────────────────────┤                              │
     │                              │                              │
```

**Request**:
```http
GET /api/method/hrms.hr.doctype.attendance.attendance.list_attendance?employee_id=uuid
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response**:
```json
{
  "success": true,
  "data": {
    "records": [...],
    "pagination": {...}
  }
}
```

**Error Response (Missing Token)**:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing or invalid authorization token"
  }
}
```

**Error Response (Expired Token)**:
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Token has expired"
  }
}
```

**Error Response (Insufficient Permissions)**:
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Insufficient permissions to perform this action"
  }
}
```

---

## JWT Token Structure

### Token Components

A JWT consists of three parts separated by dots (`.`):

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwiZW1haWwiOiJqb2huLmRvZUBleGFtcGxlLmNvbSIsInJvbGUiOiJFbXBsb3llZSIsImV4cCI6MTcwNTQwNjQwMH0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

│                                │                                                                                                                                       │                                      │
│          HEADER (Base64)        │                                             PAYLOAD (Base64)                                                                          │      SIGNATURE (HMAC-SHA256)         │
```

### 1. Header

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

- `alg`: Algorithm used for signing (HMAC-SHA256)
- `typ`: Token type (JWT)

### 2. Payload (Claims)

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john.doe@example.com",
  "role": "Employee",
  "employee_id": "EMP-001",
  "iat": 1705320000,
  "exp": 1705406400
}
```

**Standard Claims**:
- `iat` (Issued At): Timestamp when token was created
- `exp` (Expiration): Timestamp when token expires

**Custom Claims**:
- `user_id`: User's UUID
- `email`: User's email address
- `role`: User's role (for RBAC)
- `employee_id`: Associated employee ID (if applicable)

### 3. Signature

The signature is created by:

```
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret_key
)
```

This ensures:
- Token integrity (cannot be tampered with)
- Token authenticity (issued by our server)

---

## Token Management

### Token Storage (Client-Side)

**Recommended Approaches**:

#### 1. Local Storage (Web Apps)
```javascript
// Store token
localStorage.setItem('hrms_token', token);

// Retrieve token
const token = localStorage.getItem('hrms_token');

// Remove token (logout)
localStorage.removeItem('hrms_token');
```

**Pros**: Simple, persists across sessions
**Cons**: Vulnerable to XSS attacks

#### 2. HTTP-Only Cookies (More Secure)
```javascript
// Server sets cookie (in response headers)
Set-Cookie: hrms_token=<jwt>; HttpOnly; Secure; SameSite=Strict; Max-Age=86400

// Browser automatically sends cookie with requests
// No JavaScript access (XSS protection)
```

**Pros**: Protected from XSS, automatic inclusion
**Cons**: Requires CORS configuration, CSRF protection needed

#### 3. Memory Storage (Most Secure, SPA)
```javascript
// Store in JavaScript variable/Redux store
let authToken = token;

// Clear on page refresh (requires re-authentication)
```

**Pros**: Most secure, immune to XSS/CSRF
**Cons**: Lost on page refresh

### Token Expiration

**Default Expiration**: 24 hours

**Handling Expiration**:

```javascript
// Check if token is expired before making requests
function isTokenExpired(token) {
  const payload = JSON.parse(atob(token.split('.')[1]));
  const exp = payload.exp * 1000; // Convert to milliseconds
  return Date.now() >= exp;
}

// Automatically redirect to login on 401
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('hrms_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### Token Refresh Strategy

**Option 1: Re-login** (Current Implementation)
- Token expires after 24 hours
- User must log in again
- Simple, secure

**Option 2: Refresh Tokens** (Future Enhancement)
- Issue short-lived access tokens (15 minutes)
- Issue long-lived refresh tokens (30 days)
- Use refresh token to get new access token
- More complex but better UX

---

## Role-Based Access Control (RBAC)

### Role Hierarchy

```
System Manager (Full Access)
    │
    ├── HR Manager (HR, Attendance, Leave, Performance)
    │   │
    │   └── Employee (Self-service operations)
    │
    └── Payroll Manager (Payroll operations)
        │
        └── Expense Approver (Expense approval)
```

### Permission Checking

The API uses middleware to check permissions:

```go
// Example: Only HR Managers and System Managers can approve leave
attendanceRoutes.Post(
    "/hrms.hr.doctype.leave_application.leave_application.approve",
    permMiddleware.RequireRole("HR Manager", "System Manager"),
    leaveHandler.ApproveLeaveApplication,
)
```

### How It Works

1. **JWT Token Extraction**: Middleware extracts user info from JWT
2. **Role Verification**: Checks if user's role matches required roles
3. **Access Decision**:
   - ✅ If authorized → Continue to handler
   - ❌ If not authorized → Return 403 Forbidden

### Role Definitions

| Role | Code | Description |
|------|------|-------------|
| System Manager | `System Manager` | Full system access |
| HR Manager | `HR Manager` | HR operations |
| Payroll Manager | `Payroll Manager` | Payroll operations |
| Expense Approver | `Expense Approver` | Expense approval |
| Employee | `Employee` | Self-service access |

### Special Permission Rules

**Employee Access**:
- Can view own attendance, leave, salary, appraisals
- Can apply for leave, check in/out, submit expense claims
- Cannot view or modify other employees' data

**HR Manager Access**:
- Can view and modify all employee data
- Can approve/reject leave, attendance requests
- Can create appraisals, manage performance

**System Manager Access**:
- Full access to all modules
- Can manage users and roles
- Can configure system settings

---

## Security Best Practices

### For API Developers

1. **Never Log Tokens**
   ```go
   // ❌ Bad
   log.Printf("User token: %s", token)

   // ✅ Good
   log.Printf("User authenticated: %s", userID)
   ```

2. **Use Secure Secret Keys**
   ```go
   // ❌ Bad
   jwtSecret := "secret123"

   // ✅ Good - Use environment variable with strong random key
   jwtSecret := os.Getenv("JWT_SECRET") // At least 32 characters
   ```

3. **Hash Passwords Properly**
   ```go
   // ✅ Using bcrypt with appropriate cost
   hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
   ```

4. **Validate Input**
   ```go
   // ✅ Always validate and sanitize input
   if len(email) == 0 || !isValidEmail(email) {
       return errors.New("invalid email")
   }
   ```

5. **Use HTTPS in Production**
   - Always use TLS/SSL certificates
   - Redirect HTTP to HTTPS
   - Use HSTS headers

### For API Consumers

1. **Store Tokens Securely**
   - Use HTTP-only cookies when possible
   - Never log tokens
   - Clear tokens on logout

2. **Handle Token Expiration**
   - Check token expiration before requests
   - Implement automatic logout on 401
   - Show user-friendly expiration messages

3. **Use HTTPS Only**
   ```javascript
   // ❌ Bad
   const API_URL = 'http://api.hrms.com';

   // ✅ Good
   const API_URL = 'https://api.hrms.com';
   ```

4. **Implement Request Timeout**
   ```javascript
   axios.get(url, {
     timeout: 10000 // 10 seconds
   })
   ```

5. **Validate SSL Certificates**
   - Don't disable SSL verification in production
   - Use certificate pinning if needed

### Security Headers

The API includes these security headers:

```http
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1705406400
```

---

## API Endpoints

### Public Endpoints (No Authentication Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/method/hrms.api.register` | Register new user |
| POST | `/api/method/hrms.api.login` | Login and get JWT token |

### Protected Endpoints (Authentication Required)

All other endpoints require a valid JWT token in the Authorization header.

**Format**: `Authorization: Bearer <jwt_token>`

Examples:
- `/api/method/hrms.hr.doctype.attendance.*`
- `/api/method/hrms.hr.doctype.leave_application.*`
- `/api/method/hrms.payroll.doctype.salary_slip.*`
- `/api/method/hrms.hr.doctype.appraisal.*`

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) for complete endpoint list.

---

## Code Examples

### JavaScript/Node.js

#### Register & Login

```javascript
const axios = require('axios');

const API_BASE_URL = 'https://api.hrms.com';

// Register
async function register(email, password, firstName, lastName) {
  try {
    const response = await axios.post(
      `${API_BASE_URL}/api/method/hrms.api.register`,
      {
        email,
        password,
        first_name: firstName,
        last_name: lastName,
        role: 'Employee'
      }
    );
    console.log('Registration successful:', response.data);
    return response.data;
  } catch (error) {
    console.error('Registration failed:', error.response?.data);
    throw error;
  }
}

// Login
async function login(email, password) {
  try {
    const response = await axios.post(
      `${API_BASE_URL}/api/method/hrms.api.login`,
      { email, password }
    );

    const { token, user } = response.data.data;

    // Store token
    localStorage.setItem('hrms_token', token);
    localStorage.setItem('hrms_user', JSON.stringify(user));

    console.log('Login successful');
    return { token, user };
  } catch (error) {
    console.error('Login failed:', error.response?.data);
    throw error;
  }
}
```

#### Authenticated Requests

```javascript
// Create axios instance with interceptors
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
});

// Request interceptor - Add auth token
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('hrms_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => Promise.reject(error)
);

// Response interceptor - Handle 401
api.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('hrms_token');
      localStorage.removeItem('hrms_user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Example: Get attendance
async function getAttendance(employeeId, fromDate, toDate) {
  try {
    const response = await api.get(
      '/api/method/hrms.hr.doctype.attendance.attendance.list_attendance',
      {
        params: {
          employee_id: employeeId,
          from_date: fromDate,
          to_date: toDate,
          page: 1,
          limit: 20
        }
      }
    );
    return response.data.data;
  } catch (error) {
    console.error('Failed to fetch attendance:', error);
    throw error;
  }
}

// Example: Apply for leave
async function applyLeave(leaveData) {
  try {
    const response = await api.post(
      '/api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave',
      leaveData
    );
    return response.data.data;
  } catch (error) {
    console.error('Failed to apply for leave:', error);
    throw error;
  }
}
```

### Python

```python
import requests
from datetime import datetime, timedelta

API_BASE_URL = 'https://api.hrms.com'

class HRMSClient:
    def __init__(self):
        self.base_url = API_BASE_URL
        self.token = None
        self.user = None

    def register(self, email, password, first_name, last_name):
        """Register a new user"""
        response = requests.post(
            f'{self.base_url}/api/method/hrms.api.register',
            json={
                'email': email,
                'password': password,
                'first_name': first_name,
                'last_name': last_name,
                'role': 'Employee'
            }
        )
        response.raise_for_status()
        return response.json()['data']

    def login(self, email, password):
        """Login and store token"""
        response = requests.post(
            f'{self.base_url}/api/method/hrms.api.login',
            json={'email': email, 'password': password}
        )
        response.raise_for_status()

        data = response.json()['data']
        self.token = data['token']
        self.user = data['user']

        return data

    def _headers(self):
        """Get headers with authorization"""
        if not self.token:
            raise Exception('Not authenticated. Please login first.')
        return {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }

    def get_attendance(self, employee_id, from_date, to_date, page=1, limit=20):
        """Get attendance records"""
        response = requests.get(
            f'{self.base_url}/api/method/hrms.hr.doctype.attendance.attendance.list_attendance',
            params={
                'employee_id': employee_id,
                'from_date': from_date,
                'to_date': to_date,
                'page': page,
                'limit': limit
            },
            headers=self._headers()
        )
        response.raise_for_status()
        return response.json()['data']

    def apply_leave(self, employee_id, leave_type_id, from_date, to_date,
                    total_days, description):
        """Apply for leave"""
        response = requests.post(
            f'{self.base_url}/api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave',
            json={
                'employee_id': employee_id,
                'leave_type_id': leave_type_id,
                'from_date': from_date,
                'to_date': to_date,
                'total_leave_days': total_days,
                'description': description
            },
            headers=self._headers()
        )
        response.raise_for_status()
        return response.json()['data']

# Usage
client = HRMSClient()

# Login
client.login('john.doe@example.com', 'password123')

# Get attendance for last 30 days
today = datetime.now().date()
thirty_days_ago = today - timedelta(days=30)

attendance = client.get_attendance(
    employee_id='employee-uuid',
    from_date=thirty_days_ago.isoformat(),
    to_date=today.isoformat()
)

print(f"Found {len(attendance['records'])} attendance records")
```

### cURL

```bash
#!/bin/bash

API_BASE_URL="https://api.hrms.com"

# Login and save token
TOKEN=$(curl -s -X POST "$API_BASE_URL/api/method/hrms.api.login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }' | jq -r '.data.token')

echo "Token: $TOKEN"

# Get attendance
curl -X GET "$API_BASE_URL/api/method/hrms.hr.doctype.attendance.attendance.list_attendance?employee_id=uuid&page=1&limit=20" \
  -H "Authorization: Bearer $TOKEN"

# Apply for leave
curl -X POST "$API_BASE_URL/api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave" \
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

---

## Troubleshooting

### Common Issues

#### 1. "Missing or invalid authorization token"

**Cause**: Token not included in request or wrong format

**Solution**:
```javascript
// ✅ Correct
headers: {
  'Authorization': 'Bearer ' + token
}

// ❌ Wrong
headers: {
  'Authorization': token  // Missing "Bearer "
}
```

#### 2. "Token has expired"

**Cause**: Token older than 24 hours

**Solution**: Log in again to get a new token

```javascript
if (error.response.status === 401) {
  // Re-login
  await login(email, password);
  // Retry request
}
```

#### 3. "Insufficient permissions"

**Cause**: User role doesn't have access to endpoint

**Solution**: Check required roles in API documentation

#### 4. CORS Errors

**Cause**: Frontend domain not allowed

**Solution**: Configure CORS in API server or use same domain

---

## Security Considerations

### Password Requirements

- Minimum 8 characters
- Recommended: Mix of uppercase, lowercase, numbers, symbols
- Never send passwords in URL parameters
- Always use HTTPS

### Token Security

- Never share tokens
- Never commit tokens to version control
- Rotate JWT secret keys periodically
- Use short expiration times for sensitive operations

### Rate Limiting

- 100 requests per minute per IP
- Helps prevent brute force attacks
- Returns 429 status code when exceeded
- Reset time provided in headers

---

**Last Updated**: 2024-01-15
**Version**: 1.0.0
