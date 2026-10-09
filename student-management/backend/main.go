package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/jackc/pgx/v5"
	
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

)
func login(c *gin.Context) {

    // Structure for login request
    var loginData struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    // Read JSON sent by frontend
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Invalid request",
        })
        return
    }

    // Variables to store user information from database
    var (
        userID       int
        username     string
        passwordHash string
        role         string
        studentID    *int
    )
// Connect to PostgreSQL
conn, err := connectDatabase()

if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "message": "Database connection failed",
    })
    return
}
// Debug: Check which database the backend uses
var currentDB string

err = conn.QueryRow(
    context.Background(),
    "SELECT current_database()",
).Scan(&currentDB)

fmt.Println("Backend database:", currentDB)
fmt.Printf("Login username received: %q\n", loginData.Username)

defer conn.Close(context.Background())

// Find user by username
err = conn.QueryRow(
    context.Background(),
    `SELECT id, username, password_hash, role, student_id
     FROM users
     WHERE username = $1`,
    loginData.Username,
).Scan(
    &userID,
    &username,
    &passwordHash,
    &role,
    &studentID,
)
    // Username doesn't exist
    // Username doesn't exist
if err != nil {
    fmt.Println("Login username lookup failed:", err)

    c.JSON(http.StatusUnauthorized, gin.H{
        "message": "Invalid username or password",
    })
    return
}

    // Compare entered password with stored bcrypt password
    err = bcrypt.CompareHashAndPassword(
        []byte(passwordHash),
        []byte(loginData.Password),
    )

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "Invalid username or password",
        })
        return
    }

    // Create JWT token
    claims := jwt.MapClaims{
        "user_id":   userID,
        "username":  username,
        "role":      role,
        "student_id": studentID,
        "exp":       time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(
        jwt.SigningMethodHS256,
        claims,
    )

    // Secret key used to sign JWT
    tokenString, err := token.SignedString(
        []byte("student-management-secret"),
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "Could not create token",
        })
        return
    }

    // Send login response
    c.JSON(http.StatusOK, gin.H{
        "message":    "Login successful",
        "token":      tokenString,
        "role":       role,
        "student_id": studentID,
    })
}

// ============================================================
// DATABASE CONNECTION
// ============================================================

func connectDatabase() (*pgx.Conn, error) {

	// Connect to PostgreSQL
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://postgres@localhost:5432/student_db",
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// ============================================================
// GET ALL STUDENTS API
// ============================================================

func getStudents(c *gin.Context) {

	// Connect to database
	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	// Close database connection when function finishes
	defer conn.Close(context.Background())

	// SQL query
	search := c.Query("search")

rows, err := conn.Query(
    context.Background(),
    `SELECT id, name, email, age, gender, city
     FROM students
     WHERE name ILIKE $1
        OR email ILIKE $1
        OR city ILIKE $1`,
    "%"+search+"%",
)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get students",
		})
		return
	}

	defer rows.Close()

	// Create a slice to store students
	var students []gin.H

	// Read every student record
	for rows.Next() {

		var id int
		var name string
		var email string
		var age int
		var gender string
		var city string

		// Copy database values into variables
		err := rows.Scan(
			&id,
			&name,
			&email,
			&age,
			&gender,
			&city,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read student",
			})
			return
		}

		// Add student to our slice
		students = append(students, gin.H{
			"id":     id,
			"name":   name,
			"email":  email,
			"age":    age,
			"gender": gender,
			"city":   city,
		})
	}

	// Send students as JSON
	c.JSON(http.StatusOK, students)
}
// ============================================================
// GET ONE STUDENT BY ID
// ============================================================

func getStudentByID(c *gin.Context) {

	// Get the ID from the URL
	// Example: /students/101
	// id will be "101"
	id := c.Param("id")

	// Connect to database
	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	defer conn.Close(context.Background())

	// Variables to store student data
	var studentID int
	var name string
	var email string
	var age int
	var gender string
	var city string

	// Find student with the given ID
	err = conn.QueryRow(
		context.Background(),
		"SELECT id, name, email, age, gender, city FROM students WHERE id = $1",
		id,
	).Scan(
		&studentID,
		&name,
		&email,
		&age,
		&gender,
		&city,
	)

	// Student not found
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}

	// Send student data as JSON
	c.JSON(http.StatusOK, gin.H{
		"id":     studentID,
		"name":   name,
		"email":  email,
		"age":    age,
		"gender": gender,
		"city":   city,
	})
}

// ============================================================
// CREATE A NEW STUDENT
// ============================================================

func createStudent(c *gin.Context) {

	// Create a structure to receive data from the request
	var student struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
		City   string `json:"city"`
	}

	// Read JSON data sent by the user
	err := c.ShouldBindJSON(&student)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid student data",
		})
		return
	}

	// Connect to PostgreSQL
	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	defer conn.Close(context.Background())

	// Insert student into database
	var id int

	err = conn.QueryRow(
		context.Background(),
		`INSERT INTO students
		(name, email, age, gender, city)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		student.Name,
		student.Email,
		student.Age,
		student.Gender,
		student.City,
	).Scan(&id)

	// Check if insertion failed
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create student",
		})
		return
	}

	// Send successful response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Student created successfully",
		"id":      id,
	})
}

// ============================================================
// UPDATE STUDENT
// ============================================================

func updateStudent(c *gin.Context) {

	// Get student ID from URL
	// Example: /students/101
	id := c.Param("id")

	// Structure to receive updated data
	var student struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
		City   string `json:"city"`
	}

	// Read JSON from request
	err := c.ShouldBindJSON(&student)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid student data",
		})
		return
	}

	// Connect to PostgreSQL
	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	defer conn.Close(context.Background())

	// Update student in database
	result, err := conn.Exec(
		context.Background(),
		`UPDATE students
		 SET name = $1,
		     email = $2,
		     age = $3,
		     gender = $4,
		     city = $5
		 WHERE id = $6`,
		student.Name,
		student.Email,
		student.Age,
		student.Gender,
		student.City,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update student",
		})
		return
	}

	// Check whether the student actually existed
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}

	// Successful response
	c.JSON(http.StatusOK, gin.H{
		"message": "Student updated successfully",
	})
}
// ============================================================
// DELETE STUDENT
// ============================================================

func deleteStudent(c *gin.Context) {

	// Get student ID from URL
	// Example: /students/115
	id := c.Param("id")

	// Connect to PostgreSQL
	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection failed",
		})
		return
	}

	defer conn.Close(context.Background())

	// Delete student from database
	result, err := conn.Exec(
		context.Background(),
		"DELETE FROM students WHERE id = $1",
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete student",
		})
		return
	}

	// Check if student existed
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}

	// Successful response
	c.JSON(http.StatusOK, gin.H{
		"message": "Student deleted successfully",
	})
}

func getMyProfile(c *gin.Context) {

	studentIDValue, exists := c.Get("student_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Student ID not found in token",
		})
		return
	}

	studentID := int(studentIDValue.(float64))

	conn, err := connectDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database connection failed",
		})
		return
	}
	defer conn.Close(context.Background())

	var student struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
		City   string `json:"city"`
	}

	err = conn.QueryRow(
		context.Background(),
		`SELECT id, name, email, age, gender, city
		 FROM students
		 WHERE id = $1`,
		studentID,
	).Scan(
		&student.ID,
		&student.Name,
		&student.Email,
		&student.Age,
		&student.Gender,
		&student.City,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Student profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, student)
}

func updateMyProfile(c *gin.Context) {

	// Get student ID from JWT
	studentIDValue, exists := c.Get("student_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Student ID not found in token",
		})
		return
	}

	studentID := int(studentIDValue.(float64))

	// Data student is allowed to update
	var profile struct {
		Email string `json:"email"`
		City  string `json:"city"`
	}

	// Read JSON request
	err := c.ShouldBindJSON(&profile)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid profile data",
		})
		return
	}

	// Connect to database
	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database connection failed",
		})
		return
	}

	defer conn.Close(context.Background())

	// Update only email and city
	result, err := conn.Exec(
		context.Background(),
		`UPDATE students
		 SET email = $1,
		     city = $2
		 WHERE id = $3`,
		profile.Email,
		profile.City,
		studentID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update profile",
		})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Student profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}
func changePassword(c *gin.Context) {

	// Get user ID from JWT
	userIDValue, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User ID not found in token",
		})
		return
	}

	userID := int(userIDValue.(float64))

	// Request data
	var passwordData struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	// Read JSON
	err := c.ShouldBindJSON(&passwordData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// Check new password and confirm password
	if passwordData.NewPassword != passwordData.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "New password and confirm password do not match",
		})
		return
	}

	// Connect database
	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database connection failed",
		})
		return
	}

	defer conn.Close(context.Background())

	// Get current password hash
	var passwordHash string

	err = conn.QueryRow(
		context.Background(),
		"SELECT password_hash FROM users WHERE id = $1",
		userID,
	).Scan(&passwordHash)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	// Check current password
	err = bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(passwordData.CurrentPassword),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Current password is incorrect",
		})
		return
	}

	// Hash new password
	newPasswordHash, err := bcrypt.GenerateFromPassword(
		[]byte(passwordData.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create password hash",
		})
		return
	}

	// Update password
	_, err = conn.Exec(
		context.Background(),
		"UPDATE users SET password_hash = $1 WHERE id = $2",
		string(newPasswordHash),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}
func forgotPassword(c *gin.Context) {

	var passwordData struct {
		Username        string `json:"username"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	// Read JSON
	err := c.ShouldBindJSON(&passwordData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// Check passwords
	if passwordData.NewPassword != passwordData.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "New password and confirm password do not match",
		})
		return
	}

	// Connect database
	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database connection failed",
		})
		return
	}

	defer conn.Close(context.Background())

	// Check whether username exists
	var userID int

	err = conn.QueryRow(
		context.Background(),
		"SELECT id FROM users WHERE username = $1",
		passwordData.Username,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Username not found",
		})
		return
	}

	// Hash new password
	newPasswordHash, err := bcrypt.GenerateFromPassword(
		[]byte(passwordData.NewPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create password hash",
		})
		return
	}

	// Update password
	_, err = conn.Exec(
		context.Background(),
		"UPDATE users SET password_hash = $1 WHERE id = $2",
		string(newPasswordHash),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to reset password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}
func getStatistics(c *gin.Context) {

	conn, err := connectDatabase()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database connection failed",
		})
		return
	}

	defer conn.Close(context.Background())

	var totalStudents int
	var maleStudents int
	var femaleStudents int

	// Total students
	err = conn.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM students",
	).Scan(&totalStudents)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get total students",
		})
		return
	}

	// Male students
	err = conn.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM students WHERE LOWER(gender) = 'male'",
	).Scan(&maleStudents)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get male students",
		})
		return
	}

	// Female students
	err = conn.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM students WHERE LOWER(gender) = 'female'",
	).Scan(&femaleStudents)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get female students",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_students":  totalStudents,
		"male_students":   maleStudents,
		"female_students": femaleStudents,
	})
}
func getProfile(c *gin.Context) {

	// Get student_id from JWT
	studentIDValue, exists := c.Get("student_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Student ID not found in token",
		})
		return
	}

	// JWT numbers are normally decoded as float64
	studentIDFloat, ok := studentIDValue.(float64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid student ID",
		})
		return
	}

	studentID := int(studentIDFloat)

	var student struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
		City   string `json:"city"`
	}

conn, err := connectDatabase()
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "message": "Database connection failed",
    })
    return
}
defer conn.Close(context.Background())

err = conn.QueryRow(
    context.Background(),
    `SELECT id, name, email, age, gender, city
     FROM students
     WHERE id = $1`,
    studentID,
).Scan(
    &student.ID,
    &student.Name,
    &student.Email,
    &student.Age,
    &student.Gender,
    &student.City,
)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Student profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, student)
}
func updateProfile(c *gin.Context) {

	// Get student_id from JWT
	studentIDValue, exists := c.Get("student_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Student ID not found in token",
		})
		return
	}

	studentIDFloat, ok := studentIDValue.(float64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid student ID",
		})
		return
	}

	studentID := int(studentIDFloat)

	var input struct {
		Email string `json:"email"`
		City  string `json:"city"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
		})
		return
	}

	conn, err := connectDatabase()
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "message": "Database connection failed",
    })
    return
}
defer conn.Close(context.Background())

_, err = conn.Exec(
    context.Background(),
    `UPDATE students
     SET email = $1, city = $2
     WHERE id = $3`,
    input.Email,
    input.City,
    studentID,
)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}
// ============================================================
// REGISTER ENDPOINT
// ============================================================

func register(c *gin.Context) {

	var registration struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=8"`
		Age      int    `json:"age" binding:"required,min=16,max=100"`
		Gender   string `json:"gender" binding:"required,oneof=Male Female"`
		City     string `json:"city" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid registration data: " + err.Error(),
		})
		return
	}

	conn, err := connectDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database connection failed",
		})
		return
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to start transaction",
		})
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	var studentID int
	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO students (name, email, age, gender, city)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		registration.Name,
		registration.Email,
		registration.Age,
		registration.Gender,
		registration.City,
	).Scan(&studentID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "email") {
			c.JSON(http.StatusConflict, gin.H{
				"message": "Email already registered",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create student record",
		})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(registration.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to hash password",
		})
		return
	}

	_, err = tx.Exec(
		context.Background(),
		`INSERT INTO users (username, password_hash, role, student_id, email)
		 VALUES ($1, $2, $3, $4, $5)`,
		registration.Username,
		string(passwordHash),
		"student",
		studentID,
		registration.Email,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "username") {
				c.JSON(http.StatusConflict, gin.H{
					"message": "Username already taken",
				})
				return
			}
			if strings.Contains(err.Error(), "email") {
				c.JSON(http.StatusConflict, gin.H{
					"message": "Email already registered",
				})
				return
			}
			if strings.Contains(err.Error(), "student_id") {
				c.JSON(http.StatusConflict, gin.H{
					"message": "Student account already exists",
				})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user account",
		})
		return
	}

	if err = tx.Commit(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to complete registration",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Registration successful",
		"student_id": studentID,
	})
}

// ============================================================
// MAIN FUNCTION
// ============================================================
func main() {

	// Test database connection
	conn, err := connectDatabase()

	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}

	conn.Close(context.Background())

	fmt.Println("Database connected successfully!")

	// Create Gin router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))

	// ========================================================
	// CREATE GET APIs
	// ========================================================

	// Get all students
	// GET http://localhost:8080/students
	//protected route, requires JWT token
	router.GET("/students", JWTMiddleware(), AdminOnly(), getStudents)

	// Get one student by ID
	// GET http://localhost:8080/students/101
	//protected route, requires JWT token
	router.GET("/students/:id", JWTMiddleware(),AdminOnly(), getStudentByID)

	// Student Profile APIs
	router.GET("/profile", JWTMiddleware(), StudentOnly(), getProfile)
	router.PUT("/profile", JWTMiddleware(), StudentOnly(), updateProfile)
	// ========================================================
	// CREATE POST API
	// ========================================================
	// Create a new student
	//// Admin-only APIs
	router.POST("/students", JWTMiddleware(), AdminOnly(), createStudent)

	// ========================================================
	// CREATE PUT API
	// ========================================================
	// Update a student by ID
	router.PUT("/students/:id", JWTMiddleware(), AdminOnly(), updateStudent)
 	// ========================================================
	// CREATE DELETE API
	// ========================================================
	// Delete a student by ID
	router.DELETE("/students/:id", JWTMiddleware(), AdminOnly(), deleteStudent)

	//router.GET("/profile", JWTMiddleware(), StudentOnly(), getMyProfile)
	//router.PUT("/profile", JWTMiddleware(), StudentOnly(), updateMyProfile)

	router.PUT("/change-password", JWTMiddleware(), changePassword)
	router.POST("/forgot-password", forgotPassword)
	router.POST("/register", register)
	router.GET("/statistics", JWTMiddleware(), AdminOnly(), getStatistics)
	// ========================================================
	// CREATE LOGIN API
	// ========================================================
	// Login endpoint
	router.POST("/login", login)
	// Start server
	fmt.Println("Server running on http://localhost:8080")

	router.Run(":8080")
}