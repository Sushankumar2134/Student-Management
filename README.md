# Student Management System

A full-stack web application developed using **Go, Gin, PostgreSQL, and React** to manage student information and implement secure user authentication.

## 🚀 Technologies Used

- **Frontend:** React, Vite, JavaScript
- **Backend:** Go, Gin
- **Database:** PostgreSQL
- **Authentication:** JWT (JSON Web Token)
- **Password Security:** bcrypt
- **Version Control:** Git and GitHub

## ✨ Features

- User login and authentication
- Student and admin roles
- JWT-based authentication
- Role-Based Access Control (RBAC)
- Password hashing using bcrypt
- Student profile management
- REST API development
- PostgreSQL database integration
- Student self-registration *(if implemented)*

## 🔐 Authentication Flow

1. User enters their username and password.
2. The backend verifies the credentials against PostgreSQL.
3. bcrypt validates the entered password against the stored hash.
4. A JWT is generated after successful authentication.
5. The frontend uses the token for protected API requests.
6. Backend middleware validates the token.
7. Role-based authorization controls access to permitted resources.

## 👥 User Roles

### Admin
- Manage student records.
- Access admin-specific features.
- Perform authorized student management operations.

### Student
- Access the student dashboard.
- View their own profile.
- Perform permitted student-specific operations.

## 🗄️ Database Structure

The application uses PostgreSQL with the following main tables:

- `users` — Stores usernames, password hashes, roles, and linked student IDs.
- `students` — Stores student information such as name, email, age, gender, and city.

The `users.student_id` field links a student login account to the corresponding student record.

## ⚙️ Getting Started

### Prerequisites

- Go installed
- Node.js and npm installed
- PostgreSQL installed and running
- Git installed

### Backend Setup

```bash
cd backend
go mod tidy
go run .
```

Configure your PostgreSQL connection according to your local environment before running the backend.

### Frontend Setup

Navigate to your frontend directory and run:

```bash
npm install
npm run dev
```

Open the local URL displayed by Vite in your terminal.

## 📚 Learning Objectives

This project helps me understand:

- Go programming and REST API development
- Gin routing and middleware
- PostgreSQL queries and relationships
- Password hashing with bcrypt
- JWT generation and validation
- Authentication and authorization
- Role-Based Access Control
- React-to-backend API integration
- API testing and documentation

## 📖 API Documentation

Detailed API documentation, including endpoints, request bodies, response examples, authentication requirements, and HTTP status codes, will be maintained in `API_DOCUMENTATION.md`.

## 🔒 Security Notes

- Passwords should never be stored in plaintext.
- JWT secrets must not be committed to GitHub.
- Protected APIs must validate tokens on the backend.
- Admin-only operations must enforce role checks on the backend.
- Database credentials and real student information must remain private.

## 🎯 Project Purpose

This project is developed for learning and practical experience in full-stack development, backend APIs, database management, authentication, and application security.

The implementation and documentation will be updated as the project evolves.
