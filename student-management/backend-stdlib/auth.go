package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid registration data", http.StatusBadRequest)
		return
	}

	// Validate required fields.
	if req.Name == "" || req.Email == "" ||
		req.Username == "" || req.Password == "" ||
		req.Gender == "" || req.City == "" {
		writeJSONError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		writeJSONError(w, "Username must be 3-50 characters", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		writeJSONError(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if req.Age < 16 || req.Age > 100 {
		writeJSONError(w, "Age must be between 16 and 100", http.StatusBadRequest)
		return
	}

	if req.Gender != "Male" && req.Gender != "Female" {
		writeJSONError(w, "Gender must be Male or Female", http.StatusBadRequest)
		return
	}

	// Start a transaction.
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		writeJSONError(w, "Failed to start transaction", http.StatusInternalServerError)
		log.Println("Transaction error:", err)
		return
	}
	defer tx.Rollback()

	// Insert student profile first to get student_id.
	var studentID int
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO students (name, email, age, gender, city)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Name, req.Email, req.Age, req.Gender, req.City,
	).Scan(&studentID)

	if err != nil {
		log.Println("Insert student error:", err)
		if strings.Contains(err.Error(), "duplicate key") {
			writeJSONError(w, "Email already exists", http.StatusConflict)
		} else {
			writeJSONError(w, "Failed to create student", http.StatusInternalServerError)
		}
		return
	}

	// Hash the student's password.
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Println("Password hashing error:", err)
		writeJSONError(w, "Failed to secure password", http.StatusInternalServerError)
		return
	}

	// Insert hashed credentials into student_accounts.
	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO student_accounts (username, email, password, student_id)
		 VALUES ($1, $2, $3, $4)`,
		req.Username, req.Email, hashedPassword, studentID,
	)

	if err != nil {
		log.Println("Insert student account error:", err)
		if strings.Contains(err.Error(), "duplicate key") {
			writeJSONError(w, "Username or email already exists", http.StatusConflict)
		} else {
			writeJSONError(w, "Failed to create student account", http.StatusInternalServerError)
		}
		return
	}

	// Commit both inserts together.
	if err = tx.Commit(); err != nil {
		log.Println("Registration commit error:", err)
		writeJSONError(w, "Failed to complete registration", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, RegisterResponse{
		Message: "Registration successful",
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest

	// Read JSON request.
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSONError(w, "Username/email and password are required",
			http.StatusBadRequest)
		return
	}

	var (
		accountID    int
		username     string
		password     string
		role         string
		studentID    *int
		teacherID    *int
		teacherIDVal int
		isActive     bool
		email        string
	)

	// Try to find account by username or email across all three tables
	found := false

	// STEP 1: Check admin account (by username or email)
	err := db.QueryRowContext(r.Context(),
		`SELECT id, username, password, email
         FROM admins
         WHERE username = $1 OR email = $1`,
		req.Username,
	).Scan(&accountID, &username, &password, &email)

	if err == nil {
		role = "admin"
		found = true
	} else if err != sql.ErrNoRows {
		log.Println("Admin login query error:", err)
		writeJSONError(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// STEP 2: Check student account (by username or email)
	if !found {
		err = db.QueryRowContext(r.Context(),
			`SELECT id, username, password, student_id, email
             FROM student_accounts
             WHERE username = $1 OR email = $1`,
			req.Username,
		).Scan(&accountID, &username, &password, &studentID, &email)

		if err == nil {
			role = "student"
			found = true
		} else if err != sql.ErrNoRows {
			log.Println("Student login query error:", err)
			writeJSONError(w, "Login failed", http.StatusInternalServerError)
			return
		}
	}

	// STEP 3: Check teacher account (by username or email)
	if !found {
		err = db.QueryRowContext(r.Context(),
			`SELECT id, username, password, teacher_id, is_active, email
             FROM teacher_accounts
             WHERE username = $1 OR email = $1`,
			req.Username,
		).Scan(
			&accountID,
			&username,
			&password,
			&teacherIDVal,
			&isActive,
			&email,
		)
		teacherID = &teacherIDVal

		if err == sql.ErrNoRows {
			writeJSONError(w,
				"Invalid username/email or password",
				http.StatusUnauthorized)
			return
		}

		if err != nil {
			log.Println("Teacher login query error:", err)
			writeJSONError(w,
				"Login failed",
				http.StatusInternalServerError)
			return
		}

		if !isActive {
			writeJSONError(w,
				"Teacher account is inactive",
				http.StatusForbidden)
			return
		}

		role = "teacher"
		found = true
	}

	if !found {
		writeJSONError(w,
			"Invalid username/email or password",
			http.StatusUnauthorized)
		return
	}

	// STEP 4: Verify password using bcrypt for all roles.
	if !CheckPassword(req.Password, password) {
		writeJSONError(
			w,
			"Invalid username/email or password",
			http.StatusUnauthorized,
		)
		return
	}

	// STEP 5: Generate JWT with role information.
	claims := jwt.MapClaims{
		"user_id":    accountID,
		"username":   username,
		"email":      email,
		"role":       role,
		"student_id": studentID,
		"teacher_id": teacherID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}
	//contains information about the logged-in user, such as user_id, username, role, and teacher_id.
	//jwt.SigningMethodHS256 tells the JWT library to use the HS256 signing algorithm
	//jwt.NewWithClaims(...) creates a JWT using those claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("JWT signing error:", err)
		writeJSONError(w,
			"Could not create token",
			http.StatusInternalServerError)
		return
	}

	// STEP 6: Return login response.
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":    "Login successful",
		"token":      tokenString,
		"role":       role,
		"user_id":    accountID,
		"username":   username,
		"email":      email,
		"student_id": studentID,
		"teacher_id": teacherID,
	})
}

func handleCreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		writeJSONError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		writeJSONError(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Println("Password hashing error:", err)
		writeJSONError(
			w,
			"Failed to secure password",
			http.StatusInternalServerError,
		)
		return
	}

	_, err = db.ExecContext(
		r.Context(),
		`INSERT INTO admins (username, email, password)
         VALUES ($1, $2, $3)`,
		req.Username,
		req.Email,
		hashedPassword,
	)

	if err != nil {
		log.Println("Create admin error:", err)
		writeJSONError(
			w,
			"Failed to create admin. Username or email may already exist.",
			http.StatusConflict,
		)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"message": "Admin created successfully",
	})
}

func handleChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}

	role, ok := getRole(r)
	if !ok {
		writeJSONError(w, "Role not found in token", http.StatusUnauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		writeJSONError(w, "New password and confirm password do not match", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 8 {
		writeJSONError(w, "New password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	var password string
	var tableName string

	switch role {
	case "admin":
		tableName = "admins"
	case "student":
		tableName = "student_accounts"
	case "teacher":
		tableName = "teacher_accounts"
	default:
		writeJSONError(w, "Invalid role", http.StatusBadRequest)
		return
	}

	err := db.QueryRowContext(r.Context(),
		"SELECT password FROM "+tableName+" WHERE id = $1", userID).
		Scan(&password)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "User not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Failed to get user", http.StatusInternalServerError)
			log.Println("Query error:", err)
		}
		return
	}

	if !CheckPassword(req.CurrentPassword, password) {
		writeJSONError(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		log.Println("Password hashing error:", err)
		writeJSONError(w, "Failed to secure password", http.StatusInternalServerError)
		return
	}

	_, err = db.ExecContext(
		r.Context(),
		"UPDATE "+tableName+" SET password = $1 WHERE id = $2",
		hashedPassword,
		userID,
	)
	if err != nil {
		writeJSONError(w, "Failed to update password", http.StatusInternalServerError)
		log.Println("Update error:", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Password changed successfully"})
}

func handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		writeJSONError(w, "New password and confirm password do not match", http.StatusBadRequest)
		return
	}

	if len(req.NewPassword) < 8 {
		writeJSONError(w, "New password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Hash the new password before updating any account.
	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		log.Println("Password hashing error:", err)
		writeJSONError(
			w,
			"Failed to secure password",
			http.StatusInternalServerError,
		)
		return
	}
	// Try to update in each table
	tables := []string{"admins", "student_accounts", "teacher_accounts"}
	updated := false

	for _, table := range tables {
		result, err := db.ExecContext(
			r.Context(),
			"UPDATE "+table+" SET password= $1 WHERE username = $2",
			hashedPassword,
			req.Username,
		)

		if err != nil {
			log.Println("Forgot password update error for "+table+":", err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()

		if rowsAffected > 0 {
			updated = true
			break
		}
	}
	if !updated {
		writeJSONError(w, "Username not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

func handleGetProfile(w http.ResponseWriter, r *http.Request) {
	studentID, ok := getStudentID(r)
	if !ok {
		writeJSONError(w, "Student ID not found in token", http.StatusUnauthorized)
		return
	}

	var s Student
	err := db.QueryRowContext(r.Context(),
		`SELECT id, name, email, age, gender, city FROM students WHERE id = $1`, studentID).
		Scan(&s.ID, &s.Name, &s.Email, &s.Age, &s.Gender, &s.City)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "Student profile not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Failed to get profile", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}

	writeJSON(w, http.StatusOK, s)
}

func handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	studentID, ok := getStudentID(r)
	if !ok {
		writeJSONError(w, "Student ID not found in token", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	result, err := db.ExecContext(r.Context(),
		`UPDATE students SET email = $1, city = $2 WHERE id = $3`,
		req.Email, req.City, studentID)

	if err != nil {
		writeJSONError(w, "Failed to update profile", http.StatusInternalServerError)
		log.Println("Update error:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		writeJSONError(w, "Student profile not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

func handleGetStatistics(w http.ResponseWriter, r *http.Request) {
	var totalStudents, maleStudents, femaleStudents int

	err := db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM students").Scan(&totalStudents)
	if err != nil {
		writeJSONError(w, "Failed to get total students", http.StatusInternalServerError)
		return
	}

	err = db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM students WHERE LOWER(gender) = 'male'").Scan(&maleStudents)
	if err != nil {
		writeJSONError(w, "Failed to get male students", http.StatusInternalServerError)
		return
	}

	err = db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM students WHERE LOWER(gender) = 'female'").Scan(&femaleStudents)
	if err != nil {
		writeJSONError(w, "Failed to get female students", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, StatisticsResponse{
		TotalStudents:  totalStudents,
		MaleStudents:   maleStudents,
		FemaleStudents: femaleStudents,
	})
}
