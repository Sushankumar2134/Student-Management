package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid registration data", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Username == "" || req.Password == "" || req.Gender == "" || req.City == "" {
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

	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		writeJSONError(w, "Failed to start transaction", http.StatusInternalServerError)
		log.Println("Transaction error:", err)
		return
	}
	defer tx.Rollback()

	var studentID int
	err = tx.QueryRowContext(r.Context(),
		`INSERT INTO students (name, email, age, gender, city) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		req.Name, req.Email, req.Age, req.Gender, req.City).
		Scan(&studentID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "email") {
			writeJSONError(w, "Email already registered", http.StatusConflict)
			return
		}
		writeJSONError(w, "Failed to create student record", http.StatusInternalServerError)
		log.Println("Insert student error:", err)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, "Failed to hash password", http.StatusInternalServerError)
		log.Println("Hash error:", err)
		return
	}

	_, err = tx.ExecContext(r.Context(),
		`INSERT INTO users (username, password_hash, role, student_id, email) VALUES ($1, $2, $3, $4, $5)`,
		req.Username, string(passwordHash), "student", studentID, req.Email)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "username") {
				writeJSONError(w, "Username already taken", http.StatusConflict)
				return
			}
			if strings.Contains(err.Error(), "email") {
				writeJSONError(w, "Email already registered", http.StatusConflict)
				return
			}
			if strings.Contains(err.Error(), "student_id") {
				writeJSONError(w, "Student account already exists", http.StatusConflict)
				return
			}
		}
		writeJSONError(w, "Failed to create user account", http.StatusInternalServerError)
		log.Println("Insert user error:", err)
		return
	}

	if err = tx.Commit(); err != nil {
		writeJSONError(w, "Failed to complete registration", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	writeJSON(w, http.StatusCreated, RegisterResponse{
		Message:   "Registration successful",
		StudentID: studentID,
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRowContext(r.Context(),
		`SELECT id, username, password_hash, role, student_id FROM users WHERE username = $1`,
		req.Username).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.StudentID)

	if err != nil {
		writeJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeJSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"role":      user.Role,
		"student_id": user.StudentID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJSONError(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		Message:   "Login successful",
		Token:     tokenString,
		Role:      user.Role,
		StudentID: user.StudentID,
	})
}

func handleChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserID(r)
	if !ok {
		writeJSONError(w, "User ID not found in token", http.StatusUnauthorized)
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

	var passwordHash string
	err := db.QueryRowContext(r.Context(),
		"SELECT password_hash FROM users WHERE id = $1", userID).
		Scan(&passwordHash)

	if err != nil {
		writeJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.CurrentPassword)); err != nil {
		writeJSONError(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, "Failed to create password hash", http.StatusInternalServerError)
		return
	}

	_, err = db.ExecContext(r.Context(),
		"UPDATE users SET password_hash = $1 WHERE id = $2", string(newPasswordHash), userID)

	if err != nil {
		writeJSONError(w, "Failed to update password", http.StatusInternalServerError)
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

	var userID int
	err := db.QueryRowContext(r.Context(),
		"SELECT id FROM users WHERE username = $1", req.Username).
		Scan(&userID)

	if err != nil {
		writeJSONError(w, "Username not found", http.StatusNotFound)
		return
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, "Failed to create password hash", http.StatusInternalServerError)
		return
	}

	_, err = db.ExecContext(r.Context(),
		"UPDATE users SET password_hash = $1 WHERE id = $2", string(newPasswordHash), userID)

	if err != nil {
		writeJSONError(w, "Failed to reset password", http.StatusInternalServerError)
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