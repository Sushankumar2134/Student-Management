# API Documentation

Base URL: `http://localhost:8081`

All responses are JSON with `Content-Type: application/json`.

## Error Response Format

```json
{
  "message": "Error description"
}
```

## Common HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 500 | Internal Server Error |

---

## Authentication

### POST /register

Register a new student account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "username": "johndoe",
  "password": "SecurePass123",
  "age": 20,
  "gender": "Male",
  "city": "New York"
}
```

**Validation Rules:**
- All fields required
- Username: 3-50 characters
- Password: minimum 8 characters
- Age: 16-100
- Gender: "Male" or "Female"

**Responses:**

| Status | Body |
|--------|------|
| 201 | `{"message": "Registration successful", "student_id": 123}` |
| 400 | `{"message": "Invalid registration data"}` |
| 409 | `{"message": "Username already taken"}` or `{"message": "Email already registered"}` |
| 500 | `{"message": "Failed to create student record"}` |

---

### POST /login

Authenticate and receive a JWT token.

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "SecurePass123"
}
```

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"message": "Login successful", "token": "eyJ...", "role": "student", "student_id": 123}` |
| 400 | `{"message": "Invalid request"}` |
| 401 | `{"message": "Invalid username or password"}` |

**Token Usage:**
Include in subsequent requests:
```
Authorization: Bearer eyJ...
```

---

### POST /forgot-password

Reset password using username (no email verification in current implementation).

**Request Body:**
```json
{
  "username": "johndoe",
  "new_password": "NewSecurePass123",
  "confirm_password": "NewSecurePass123"
}
```

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"message": "Password reset successfully"}` |
| 400 | `{"message": "Invalid request"}` or `{"message": "New password and confirm password do not match"}` |
| 404 | `{"message": "Username not found"}` |
| 500 | `{"message": "Failed to reset password"}` |

---

## Student Management (Admin Only)

All endpoints require `Authorization: Bearer <token>` with `role: "admin"`.

### GET /students

List all students with optional search.

**Query Parameters:**
- `search` (optional): Search by name, email, or city (case-insensitive)

**Responses:**

| Status | Body |
|--------|------|
| 200 | `[{"id": 1, "name": "John", "email": "john@example.com", "age": 20, "gender": "Male", "city": "New York"}, ...]` |
| 401 | `{"message": "Authorization token required"}` |
| 403 | `{"message": "Admin access required"}` |

**Example:**
```
GET /students?search=John
```

---

### GET /students/{id}

Get a single student by ID.

**Path Parameters:**
- `id`: Student ID

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"id": 1, "name": "John", "email": "john@example.com", "age": 20, "gender": "Male", "city": "New York"}` |
| 400 | `{"message": "Invalid student ID"}` |
| 401 | `{"message": "Authorization token required"}` |
| 403 | `{"message": "Admin access required"}` |
| 404 | `{"message": "Student not found"}` |

---

### POST /students

Create a new student (admin only).

**Request Body:**
```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "age": 22,
  "gender": "Female",
  "city": "Boston"
}
```

**Responses:**

| Status | Body |
|--------|------|
| 201 | `{"message": "Student created successfully", "id": 124}` |
| 400 | `{"message": "Invalid student data"}` |
| 401 | `{"message": "Authorization token required"}` |
| 403 | `{"message": "Admin access required"}` |
| 500 | `{"message": "Failed to create student"}` |

---

### PUT /students/{id}

Update a student (admin only).

**Path Parameters:**
- `id`: Student ID

**Request Body:**
```json
{
  "name": "Jane Smith",
  "email": "jane.smith@example.com",
  "age": 23,
  "gender": "Female",
  "city": "Chicago"
}
```

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"message": "Student updated successfully"}` |
| 400 | `{"message": "Invalid student ID"}` or `{"message": "Invalid student data"}` |
| 401 | `{"message": "Authorization token required"}` |
| 403 | `{"message": "Admin access required"}` |
| 404 | `{"message": "Student not found"}` |
| 500 | `{"message": "Failed to update student"}` |

---

### DELETE /students/{id}

Delete a student and associated user account (admin only).

**Path Parameters:**
- `id`: Student ID

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"message": "Student deleted successfully"}` |
| 400 | `{"message": "Invalid student ID"}` |
| 401 | `{"message": "Authorization token required"}` |
| 403 | `{"message": "Admin access required"}` |
| 404 | `{"message": "Student not found"}` |
| 500 | `{"message": "Failed to delete student"}` |

---

## Profile Management (Student Only)

All endpoints require `Authorization: Bearer <token>` with `role: "student"`.

### GET /profile

Get the authenticated student's profile.

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"id": 123, "name": "John", "email": "john@example.com", "age": 20, "gender": "Male", "city": "New York"}` |
| 401 | `{"message": "Authorization token required"}` or `{"message": "Student ID not found in token"}` |
| 403 | `{"message": "Student access required"}` |
| 404 | `{"message": "Student profile not found"}` |

---

### PUT /profile

Update the authenticated student's profile (email and city only).

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "city": "San Francisco"
}
```

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"message": "Profile updated successfully"}` |
| 400 | `{"message": "Invalid request data"}` |
| 401 | `{"message": "Authorization token required"}` or `{"message": "Student ID not found in token"}` |
| 403 | `{"message": "Student access required"}` |
| 404 | `{"message": "Student profile not found"}` |
| 500 | `{"message": "Failed to update profile"}` |

---

## Password Management (Authenticated)

### PUT /change-password

Change password (requires current password).

**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "current_password": "OldPass123",
  "new_password": "NewPass123",
  "confirm_password": "NewPass123"
}
```

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"message": "Password changed successfully"}` |
| 400 | `{"message": "Invalid request"}` or `{"message": "New password and confirm password do not match"}` |
| 401 | `{"message": "Authorization token required"}` or `{"message": "User ID not found in token"}` or `{"message": "Current password is incorrect"}` |
| 404 | `{"message": "User not found"}` |
| 500 | `{"message": "Failed to update password"}` |

---

## Statistics (Admin Only)

### GET /statistics

Get student statistics.

**Headers:** `Authorization: Bearer <token>` with `role: "admin"`

**Responses:**

| Status | Body |
|--------|------|
| 200 | `{"total_students": 150, "male_students": 85, "female_students": 65}` |
| 401 | `{"message": "Authorization token required"}` |
| 403 | `{"message": "Admin access required"}` |
| 500 | `{"message": "Failed to get statistics"}` |

---

## CORS

The server allows requests from `http://localhost:5173` with:
- Methods: GET, POST, PUT, DELETE, OPTIONS
- Headers: Origin, Content-Type, Authorization
- Credentials: true

Preflight OPTIONS requests return 204 No Content.

---

## Testing with cURL

### Register a new student
```bash
curl -X POST http://localhost:8081/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","username":"testuser","password":"TestPass123","age":20,"gender":"Male","city":"Test City"}'
```

### Login
```bash
curl -X POST http://localhost:8081/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"TestPass123"}'
```

### Get students (admin)
```bash
curl -X GET http://localhost:8081/students \
  -H "Authorization: Bearer <admin-token>"
```

### Get own profile (student)
```bash
curl -X GET http://localhost:8081/profile \
  -H "Authorization: Bearer <student-token>"
```

### Change password
```bash
curl -X PUT http://localhost:8081/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"current_password":"TestPass123","new_password":"NewPass123","confirm_password":"NewPass123"}'
```

---

## Postman Collection

Import the following as a Postman collection for easier testing:

1. Create a new collection
2. Set base URL variable: `{{baseUrl}} = http://localhost:8081`
3. Add requests for each endpoint above
4. For authenticated requests, add Header: `Authorization: Bearer {{token}}`
5. After login, save the token: `pm.environment.set("token", pm.response.json().token)`