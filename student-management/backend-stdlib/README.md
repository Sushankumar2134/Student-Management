# Student Management System - Standard Library Backend

A REST API for student management built using Go's standard library (`net/http`) without any external HTTP frameworks.

## Features

- Student CRUD operations (Create, Read, Update, Delete)
- User registration and authentication with JWT
- Role-based access control (Admin, Student)
- Password management (change password, forgot password)
- Profile management
- Statistics dashboard (admin only)
- PostgreSQL database with connection pooling
- CORS support for React frontend

## Tech Stack

- **Language**: Go 1.27+
- **Database**: PostgreSQL with `github.com/jackc/pgx/v5/stdlib`
- **Authentication**: JWT with `github.com/golang-jwt/jwt/v5`
- **Password Hashing**: bcrypt with `golang.org/x/crypto/bcrypt`
- **Environment**: `github.com/joho/godotenv`

## Project Structure

```
backend-stdlib/
├── main.go           # Server initialization and route registration
├── database.go       # PostgreSQL connection and configuration
├── models.go         # Data structures and request/response types
├── handlers.go       # HTTP handlers for student CRUD
├── auth.go           # Authentication, registration, password handlers
├── middleware.go     # JWT, RBAC, and CORS middleware
├── helpers.go        # JSON response helpers
├── go.mod            # Go module dependencies
├── go.sum            # Dependency checksums
├── .env.example      # Environment variable template
├── .gitignore        # Git ignore rules
├── README.md         # This file
└── API_DOCUMENTATION.md  # API endpoint documentation
```

## Getting Started

### Prerequisites

- Go 1.27 or later
- PostgreSQL 14+
- Existing `student_db` database with `students` and `users` tables

### Installation

1. Clone the repository
2. Navigate to the backend-stdlib directory:
   ```bash
   cd backend-stdlib
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Copy `.env.example` to `.env` and configure:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

5. Run the server:
   ```bash
   go run .
   ```

The server will start on `http://localhost:8081`

### Database Schema

The backend expects the following tables in PostgreSQL:

```sql
-- Students table
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100),
    age INTEGER,
    gender VARCHAR(20),
    city VARCHAR(20)
);

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    student_id INTEGER UNIQUE REFERENCES students(id),
    email VARCHAR(100) UNIQUE
);
```

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register a new student |
| POST | `/login` | Login and get JWT token |
| POST | `/forgot-password` | Reset password (username only) |

### Protected Endpoints (Require JWT)

| Method | Endpoint | Roles | Description |
|--------|----------|-------|-------------|
| GET | `/students` | Admin | List all students with optional search |
| GET | `/students/{id}` | Admin | Get student by ID |
| POST | `/students` | Admin | Create new student |
| PUT | `/students/{id}` | Admin | Update student |
| DELETE | `/students/{id}` | Admin | Delete student |
| GET | `/profile` | Student | Get own profile |
| PUT | `/profile` | Student | Update own profile (email, city) |
| PUT | `/change-password` | Student | Change password |
| GET | `/statistics` | Admin | Get student statistics |

## Authentication

Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Frontend Integration

The React frontend (in `../frontend`) uses `http://localhost:8080` by default. To use this backend:

1. Update `frontend/src/services/api.js`:
   ```javascript
   const api = axios.create({
       baseURL: "http://localhost:8081",
   });
   ```

2. Or set environment variable:
   ```bash
   VITE_API_URL=http://localhost:8081
   ```

## Development

Run with auto-reload (requires `air`):
```bash
air
```

Run tests:
```bash
go test ./...
```

## License

MIT