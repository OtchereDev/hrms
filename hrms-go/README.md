# HRMS Go - Human Resource Management System

A complete rewrite of the Frappe HRMS system in Go using Fiber framework, GORM, and PostgreSQL.

## Features

- ✅ **Generic File Storage** - Interface supporting both local filesystem and S3-compatible storage
- ✅ **JWT Authentication** - Secure token-based authentication with refresh tokens
- ✅ **RBAC Permission System** - Role-based access control with document-level permissions
- ✅ **RESTful API** - Frappe-compatible API endpoints
- 🚧 **HR Module** - Employee management, attendance, leave, recruitment (In Progress)
- 🚧 **Payroll Module** - Salary processing, tax calculations (In Progress)
- 🚧 **Workflow Engine** - Approval workflows (Planned)
- 🚧 **Ghana Regional Features** - Tax, SSNIT, statutory compliance (Planned)
- 🚧 **Reports** - PDF/Excel export (Planned)
- 🚧 **Real-time Updates** - WebSocket support (Planned)

## Technology Stack

### Backend
- **Go** 1.21+
- **Fiber v2** - Express-like web framework
- **GORM** - ORM for PostgreSQL
- **JWT** - Authentication
- **AWS SDK v2** - S3 storage support
- **PostgreSQL 15+** - Primary database
- **Redis** - Caching and job queue

### Frontend (Existing)
- **Vue 3** - Progressive framework
- **Ionic** - Mobile UI components
- **TypeScript** - Type safety

## Project Structure

```
hrms-go/
├── cmd/
│   ├── api/              # API server entry point
│   ├── scheduler/        # Background job scheduler
│   └── migrate/          # Database migrations
├── internal/
│   ├── config/           # Configuration management
│   ├── core/
│   │   ├── models/       # GORM models (base, hr, payroll)
│   │   ├── services/     # Business logic
│   │   └── repositories/ # Data access layer
│   ├── api/
│   │   ├── handlers/     # HTTP handlers
│   │   ├── middleware/   # Authentication, permissions
│   │   └── routes/       # Route definitions
│   ├── storage/          # File storage (S3 + Local)
│   ├── workflow/         # Workflow engine
│   └── scheduler/        # Background jobs
├── pkg/
│   ├── logger/           # Logging utilities
│   ├── errors/           # Error definitions
│   └── response/         # Standard API responses
└── docs/                 # Documentation
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Redis (for caching and jobs)
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**
   ```bash
   cd /home/user/hrms/hrms-go
   ```

2. **Install dependencies**
   ```bash
   make install
   ```

3. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start dependencies with Docker**
   ```bash
   make docker-up
   ```

   This will start:
   - PostgreSQL on port 5432
   - Redis on port 6379
   - MinIO (S3) on ports 9000/9001 (optional)

5. **Run database migrations**
   ```bash
   make migrate
   ```

6. **Start the API server**
   ```bash
   make run
   ```

   The server will start on http://localhost:8080

### Default Admin Credentials

When running in development mode with seeded data:
- **Email:** admin@hrms.local
- **Password:** admin123

⚠️ **Important:** Change these credentials in production!

## API Documentation

### Authentication Endpoints

#### Login
```http
POST /api/method/login
Content-Type: application/json

{
  "usr": "admin@hrms.local",
  "pwd": "admin123"
}
```

Response:
```json
{
  "success": true,
  "message": "Logged In",
  "data": {
    "user": {
      "id": 1,
      "email": "admin@hrms.local",
      "full_name": "System Administrator",
      "roles": ["System Manager"]
    },
    "token": {
      "access_token": "eyJhbGciOiJIUzI1NiIs...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
      "token_type": "Bearer",
      "expires_in": 900
    }
  }
}
```

#### Get Current User Info
```http
GET /api/method/hrms.api.get_current_user_info
Authorization: Bearer {access_token}
```

#### Refresh Token
```http
POST /api/method/refresh_token
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Logout
```http
POST /api/method/logout
Authorization: Bearer {access_token}
```

### Health Check

```http
GET /health
```

Response:
```json
{
  "status": "ok",
  "service": "hrms-api"
}
```

## Development

### Running in Development Mode

```bash
# Start with hot reload (requires air)
make dev

# Or manually
make run
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Formatting

```bash
# Format code
make fmt

# Run linter
make lint
```

### Database Operations

```bash
# Run migrations
make migrate

# Seed test data
make seed
```

## Docker Deployment

### Build Docker Image

```bash
make docker-build
```

### Start All Services

```bash
make docker-up
```

This starts:
- PostgreSQL database
- Redis cache
- HRMS API server
- MinIO (S3-compatible storage)

### Stop All Services

```bash
make docker-down
```

## Configuration

### Environment Variables

See `.env.example` for all available configuration options.

Key settings:

#### Application
- `APP_ENV` - Environment (development/production)
- `APP_DEBUG` - Enable debug mode
- `SERVER_PORT` - API server port (default: 8080)

#### Database
- `DB_HOST` - PostgreSQL host
- `DB_PORT` - PostgreSQL port
- `DB_NAME` - Database name
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password

#### JWT
- `JWT_SECRET` - Secret key for signing tokens (**Change in production!**)
- `JWT_ACCESS_TTL` - Access token lifetime (default: 15m)
- `JWT_REFRESH_TTL` - Refresh token lifetime (default: 7d)

#### Storage
- `STORAGE_TYPE` - Storage backend (local/s3)
- `STORAGE_LOCAL_PATH` - Local storage path
- `S3_BUCKET` - S3 bucket name (if using S3)
- `S3_REGION` - AWS region
- `S3_ACCESS_KEY_ID` - AWS access key
- `S3_SECRET_ACCESS_KEY` - AWS secret key

## File Storage

The system supports two storage backends:

### Local Storage

```env
STORAGE_TYPE=local
STORAGE_LOCAL_PATH=./storage
```

Files are stored in the local filesystem.

### S3 Storage

```env
STORAGE_TYPE=s3
S3_BUCKET=my-hrms-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY_ID=your-key-id
S3_SECRET_ACCESS_KEY=your-secret-key
```

Compatible with:
- AWS S3
- MinIO
- DigitalOcean Spaces
- Any S3-compatible service

## Authentication

### JWT Tokens

The system uses JWT tokens for authentication:

1. **Access Token** - Short-lived (default: 15 minutes)
   - Used for API requests
   - Must be included in Authorization header

2. **Refresh Token** - Long-lived (default: 7 days)
   - Used to obtain new access tokens
   - Stored securely with session

### Session Management

Sessions are tracked in the database with:
- Session ID
- User information
- IP address and user agent
- Last activity timestamp
- Expiration time

### Permissions

Role-based access control (RBAC) with:
- User roles (HR Manager, HR User, Employee, etc.)
- Document-type permissions (read, write, create, delete, submit, cancel)
- Ownership checks (users can only access their own data)

## Migration from Frappe HRMS

This project maintains API compatibility with the original Frappe HRMS:

### Compatible API Endpoints

```
POST /api/method/login                              ✅
POST /api/method/logout                             ✅
GET  /api/method/hrms.api.get_current_user_info    ✅
GET  /api/method/hrms.api.get_current_employee_info 🚧
```

More endpoints will be added as modules are ported.

### Data Migration

For migrating data from Frappe HRMS:

1. Export data from Frappe
2. Transform to GORM-compatible format
3. Import using migration scripts

(Migration tools coming soon)

## Roadmap

### Phase 1: Foundation ✅ (Completed)
- [x] Project structure
- [x] Configuration system
- [x] Database setup (GORM + PostgreSQL)
- [x] JWT authentication
- [x] RBAC middleware
- [x] Generic file storage (S3 + Local)

### Phase 2: Core HR Module 🚧 (In Progress)
- [ ] Employee models and APIs
- [ ] Attendance system
- [ ] Leave management
- [ ] Recruitment module
- [ ] Performance tracking

### Phase 3: Payroll Module 🚧 (Planned)
- [ ] Salary structure
- [ ] Salary slip calculation
- [ ] Tax computation (Ghana)
- [ ] Payroll processing

### Phase 4: Advanced Features 📅 (Planned)
- [ ] Workflow engine
- [ ] Background job scheduler
- [ ] WebSocket for real-time updates
- [ ] Report generation (PDF/Excel)
- [ ] Email notifications

### Phase 5: Regional Features 📅 (Planned)
- [ ] Ghana tax and SSNIT
- [ ] Nigeria regional features
- [ ] Multi-currency support

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

[Add your license here]

## Support

For issues and questions:
- GitHub Issues: [Create an issue](../../issues)
- Documentation: [See docs/](docs/)

## Credits

Based on the original [Frappe HRMS](https://github.com/frappe/hrms) by Frappe Technologies.

Rewritten in Go by OtchereDev with focus on performance, scalability, and maintainability.
