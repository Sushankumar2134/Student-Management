package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type CreateTeacherAccountRequest struct {
	TeacherID int    `json:"teacher_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Teacher struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Subject  string         `json:"subject"`
	Phone    string         `json:"phone"`
	Address  sql.NullString `json:"address"`
	HireDate time.Time      `json:"hire_date"`
}

type TeacherProfileResponse struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Subject  string    `json:"subject"`
	Phone    string    `json:"phone"`
	Address  string    `json:"address"`
	HireDate time.Time `json:"hire_date"`
	Username string    `json:"username"`
}

type UpdateTeacherProfileRequest struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Subject string `json:"subject"`
}

type Class struct {
	ID           int    `json:"id"`
	ClassName    string `json:"class_name"`
	Section      string `json:"section"`
	AcademicYear string `json:"academic_year"`
	StudentCount int    `json:"student_count"`
}

type Subject struct {
	ID          int    `json:"id"`
	SubjectName string `json:"subject_name"`
	SubjectCode string `json:"subject_code"`
}

type TeacherAssignment struct {
	Class   Class   `json:"class"`
	Subject Subject `json:"subject"`
}

type StudentInClass struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	City   string `json:"city"`
}

type AttendanceRecord struct {
	ID             int    `json:"id"`
	StudentID      int    `json:"student_id"`
	StudentName    string `json:"student_name"`
	ClassID        int    `json:"class_id"`
	ClassName      string `json:"class_name"`
	SubjectID      int    `json:"subject_id"`
	SubjectName    string `json:"subject_name"`
	AttendanceDate string `json:"attendance_date"`
	Status         string `json:"status"`
}

type MarkAttendanceRequest struct {
	ClassID   int              `json:"class_id"`
	SubjectID int              `json:"subject_id"`
	Date      string           `json:"date"`
	Records   []AttendanceMark `json:"records"`
}

type AttendanceMark struct {
	StudentID int    `json:"student_id"`
	Status    string `json:"status"`
}

type DashboardStats struct {
	AssignedClasses  int `json:"assigned_classes"`
	AssignedSubjects int `json:"assigned_subjects"`
	AssignedStudents int `json:"assigned_students"`
	RecentAttendance int `json:"recent_attendance"`
}

func CreateTeacherAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateTeacherAccountRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON body"}`,
			http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)

	if req.TeacherID <= 0 ||
		req.Username == "" ||
		len(req.Password) < 8 {
		http.Error(w,
			`{"error":"Valid teacher_id, username and password of at least 8 characters are required"}`,
			http.StatusBadRequest)
		return
	}

	// Check whether the teacher profile exists.
	var exists bool

	err := db.QueryRowContext(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM teachers WHERE id = $1)`,
		req.TeacherID,
	).Scan(&exists)

	if err != nil {
		http.Error(w, `{"error":"Database error"}`,
			http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, `{"error":"Teacher profile not found"}`,
			http.StatusNotFound)
		return
	}

	// Hash the teacher's password using bcrypt.
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		log.Println("Teacher password hashing error:", err)
		http.Error(
			w,
			`{"error":"Failed to secure password"}`,
			http.StatusInternalServerError,
		)
		return
	}

	// Create the teacher login account with the hashed password.
	var accountID int

	err = db.QueryRowContext(r.Context(), `
        INSERT INTO teacher_accounts
            (teacher_id, username, password)
        VALUES ($1, $2, $3)
        RETURNING id
    `,
		req.TeacherID,
		req.Username,
		hashedPassword,
	).Scan(&accountID)

	if err != nil {
		log.Println("Teacher account creation error:", err)
		http.Error(
			w,
			`{"error":"Account could not be created; check for duplicate username or teacher account"}`,
			http.StatusConflict,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Teacher account created successfully",
		"account_id": accountID,
		"teacher_id": req.TeacherID,
		"username":   req.Username,
	})

}

func handleGetTeacherProfile(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	var t Teacher
	var username string
	err := db.QueryRowContext(r.Context(),
		`SELECT t.id, t.name, t.email, t.subject, t.phone, t.address, t.hire_date, ta.username
		 FROM teachers t
		 JOIN teacher_accounts ta ON ta.teacher_id = t.id
		 WHERE t.id = $1`, teacherID).
		Scan(&t.ID, &t.Name, &t.Email, &t.Subject, &t.Phone, &t.Address, &t.HireDate, &username)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "Teacher profile not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Failed to get teacher profile", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}

	response := TeacherProfileResponse{
		ID:       t.ID,
		Name:     t.Name,
		Email:    t.Email,
		Subject:  t.Subject,
		Phone:    t.Phone,
		Address:  t.Address.String,
		HireDate: t.HireDate,
		Username: username,
	}

	writeJSON(w, http.StatusOK, response)
}

func handleUpdateTeacherProfile(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	var req UpdateTeacherProfileRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	result, err := db.ExecContext(r.Context(),
		`UPDATE teachers SET email = $1, phone = $2, address = $3, subject = $4 WHERE id = $5`,
		req.Email, req.Phone, req.Address, req.Subject, teacherID)

	if err != nil {
		writeJSONError(w, "Failed to update teacher profile", http.StatusInternalServerError)
		log.Println("Update error:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		writeJSONError(w, "Teacher profile not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Teacher profile updated successfully"})
}

func handleGetTeacherDashboard(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	var t Teacher
	err := db.QueryRowContext(r.Context(),
		`SELECT id, name, email, subject, phone, address, hire_date FROM teachers WHERE id = $1`, teacherID).
		Scan(&t.ID, &t.Name, &t.Email, &t.Subject, &t.Phone, &t.Address, &t.HireDate)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "Teacher not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Failed to get teacher data", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}

	// Get some dashboard statistics
	var totalStudents int
	err = db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM students").Scan(&totalStudents)
	if err != nil {
		log.Println("Failed to get total students:", err)
		totalStudents = 0
	}

	dashboard := map[string]interface{}{
		"teacher": map[string]interface{}{
			"id":       t.ID,
			"name":     t.Name,
			"email":    t.Email,
			"subject":  t.Subject,
			"phone":    t.Phone,
			"address":  t.Address.String,
			"hireDate": t.HireDate,
		},
		"statistics": map[string]interface{}{
			"totalStudents": totalStudents,
		},
	}

	writeJSON(w, http.StatusOK, dashboard)
}

func handleGetTeacherClasses(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT c.id, c.class_name, c.section, c.academic_year,
		       COUNT(s.id) as student_count
		FROM classes c
		JOIN teacher_assignments ta ON ta.class_id = c.id
		LEFT JOIN students s ON s.class_id = c.id
		WHERE ta.teacher_id = $1
		GROUP BY c.id, c.class_name, c.section, c.academic_year
		ORDER BY c.class_name, c.section
	`, teacherID)

	if err != nil {
		writeJSONError(w, "Failed to get classes", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var c Class
		if err := rows.Scan(&c.ID, &c.ClassName, &c.Section, &c.AcademicYear, &c.StudentCount); err != nil {
			writeJSONError(w, "Failed to read class", http.StatusInternalServerError)
			log.Println("Scan error:", err)
			return
		}
		classes = append(classes, c)
	}

	writeJSON(w, http.StatusOK, classes)
}

func handleGetTeacherSubjects(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	classID := r.URL.Query().Get("class_id")
	var rows *sql.Rows
	var err error

	if classID != "" {
		rows, err = db.QueryContext(r.Context(), `
			SELECT s.id, s.subject_name, s.subject_code
			FROM subjects s
			JOIN teacher_assignments ta ON ta.subject_id = s.id
			WHERE ta.teacher_id = $1 AND ta.class_id = $2
			ORDER BY s.subject_name
		`, teacherID, classID)
	} else {
		rows, err = db.QueryContext(r.Context(), `
			SELECT DISTINCT s.id, s.subject_name, s.subject_code
			FROM subjects s
			JOIN teacher_assignments ta ON ta.subject_id = s.id
			WHERE ta.teacher_id = $1
			ORDER BY s.subject_name
		`, teacherID)
	}

	if err != nil {
		writeJSONError(w, "Failed to get subjects", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		if err := rows.Scan(&s.ID, &s.SubjectName, &s.SubjectCode); err != nil {
			writeJSONError(w, "Failed to read subject", http.StatusInternalServerError)
			log.Println("Scan error:", err)
			return
		}
		subjects = append(subjects, s)
	}

	writeJSON(w, http.StatusOK, subjects)
}

func handleGetTeacherStudents(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	classID := r.URL.Query().Get("class_id")
	if classID == "" {
		writeJSONError(w, "class_id is required", http.StatusBadRequest)
		return
	}

	// Verify teacher is assigned to this class
	var exists bool
	err := db.QueryRowContext(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM teacher_assignments WHERE teacher_id = $1 AND class_id = $2)`,
		teacherID, classID).Scan(&exists)

	if err != nil || !exists {
		writeJSONError(w, "Unauthorized to access this class", http.StatusForbidden)
		return
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, name, email, age, gender, city
		FROM students
		WHERE class_id = $1
		ORDER BY name
	`, classID)

	if err != nil {
		writeJSONError(w, "Failed to get students", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var students []StudentInClass
	for rows.Next() {
		var s StudentInClass
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Age, &s.Gender, &s.City); err != nil {
			writeJSONError(w, "Failed to read student", http.StatusInternalServerError)
			log.Println("Scan error:", err)
			return
		}
		students = append(students, s)
	}

	writeJSON(w, http.StatusOK, students)
}

func handleGetAttendance(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	classID := r.URL.Query().Get("class_id")
	subjectID := r.URL.Query().Get("subject_id")
	date := r.URL.Query().Get("date")

	if classID == "" || subjectID == "" || date == "" {
		writeJSONError(w, "class_id, subject_id, and date are required", http.StatusBadRequest)
		return
	}

	// Verify teacher is assigned to this class and subject
	var exists bool
	err := db.QueryRowContext(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM teacher_assignments WHERE teacher_id = $1 AND class_id = $2 AND subject_id = $3)`,
		teacherID, classID, subjectID).Scan(&exists)

	if err != nil || !exists {
		writeJSONError(w, "Unauthorized to access this class/subject", http.StatusForbidden)
		return
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT a.id, a.student_id, s.name as student_name, a.class_id, c.class_name,
		       a.subject_id, sub.subject_name, a.attendance_date, a.status
		FROM attendance a
		JOIN students s ON s.id = a.student_id
		JOIN classes c ON c.id = a.class_id
		JOIN subjects sub ON sub.id = a.subject_id
		WHERE a.class_id = $1 AND a.subject_id = $2 AND a.attendance_date = $3
		ORDER BY s.name
	`, classID, subjectID, date)

	if err != nil {
		writeJSONError(w, "Failed to get attendance", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var records []AttendanceRecord
	for rows.Next() {
		var rec AttendanceRecord
		if err := rows.Scan(&rec.ID, &rec.StudentID, &rec.StudentName, &rec.ClassID, &rec.ClassName,
			&rec.SubjectID, &rec.SubjectName, &rec.AttendanceDate, &rec.Status); err != nil {
			writeJSONError(w, "Failed to read attendance", http.StatusInternalServerError)
			log.Println("Scan error:", err)
			return
		}
		records = append(records, rec)
	}

	writeJSON(w, http.StatusOK, records)
}

func handleMarkAttendance(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	var req MarkAttendanceRequest
	if err := readJSON(r, &req); err != nil {
		writeJSONError(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	if req.ClassID == 0 || req.SubjectID == 0 || req.Date == "" || len(req.Records) == 0 {
		writeJSONError(w, "class_id, subject_id, date, and records are required", http.StatusBadRequest)
		return
	}

	// Verify teacher is assigned to this class and subject
	var exists bool
	err := db.QueryRowContext(r.Context(),
		`SELECT EXISTS(SELECT 1 FROM teacher_assignments WHERE teacher_id = $1 AND class_id = $2 AND subject_id = $3)`,
		teacherID, req.ClassID, req.SubjectID).Scan(&exists)

	if err != nil || !exists {
		writeJSONError(w, "Unauthorized to mark attendance for this class/subject", http.StatusForbidden)
		return
	}

	// Validate each status
	for _, rec := range req.Records {
		if rec.Status != "Present" && rec.Status != "Absent" {
			writeJSONError(w, "Status must be 'Present' or 'Absent'", http.StatusBadRequest)
			return
		}
	}

	// Use transaction for atomicity
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		writeJSONError(w, "Failed to start transaction", http.StatusInternalServerError)
		log.Println("Transaction error:", err)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(r.Context(), `
		INSERT INTO attendance (student_id, class_id, subject_id, teacher_id, attendance_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (student_id, class_id, subject_id, attendance_date)
		DO UPDATE SET status = EXCLUDED.status, teacher_id = EXCLUDED.teacher_id
	`)
	if err != nil {
		writeJSONError(w, "Failed to prepare statement", http.StatusInternalServerError)
		log.Println("Prepare error:", err)
		return
	}
	defer stmt.Close()

	for _, rec := range req.Records {
		_, err = stmt.ExecContext(r.Context(), rec.StudentID, req.ClassID, req.SubjectID, teacherID, req.Date, rec.Status)
		if err != nil {
			writeJSONError(w, "Failed to save attendance", http.StatusInternalServerError)
			log.Println("Exec error:", err)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		writeJSONError(w, "Failed to commit attendance", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Attendance saved successfully"})
}

func handleGetTeacherDashboardEnhanced(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := getTeacherID(r)
	if !ok {
		writeJSONError(w, "Teacher ID not found in token", http.StatusUnauthorized)
		return
	}

	var t Teacher
	err := db.QueryRowContext(r.Context(),
		`SELECT id, name, email, subject, phone, address, hire_date FROM teachers WHERE id = $1`, teacherID).
		Scan(&t.ID, &t.Name, &t.Email, &t.Subject, &t.Phone, &t.Address, &t.HireDate)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "Teacher not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Failed to get teacher data", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}

	// Get dashboard statistics
	var stats DashboardStats

	// Assigned classes count
	err = db.QueryRowContext(r.Context(),
		`SELECT COUNT(DISTINCT class_id) FROM teacher_assignments WHERE teacher_id = $1`,
		teacherID).Scan(&stats.AssignedClasses)
	if err != nil {
		stats.AssignedClasses = 0
	}

	// Assigned subjects count
	err = db.QueryRowContext(r.Context(),
		`SELECT COUNT(DISTINCT subject_id) FROM teacher_assignments WHERE teacher_id = $1`,
		teacherID).Scan(&stats.AssignedSubjects)
	if err != nil {
		stats.AssignedSubjects = 0
	}

	// Assigned students count (unique students in assigned classes)
	err = db.QueryRowContext(r.Context(), `
		SELECT COUNT(DISTINCT s.id)
		FROM students s
		JOIN teacher_assignments ta ON ta.class_id = s.class_id
		WHERE ta.teacher_id = $1
	`, teacherID).Scan(&stats.AssignedStudents)
	if err != nil {
		stats.AssignedStudents = 0
	}

	// Recent attendance count (last 7 days)
	err = db.QueryRowContext(r.Context(), `
		SELECT COUNT(*)
		FROM attendance a
		JOIN teacher_assignments ta ON ta.class_id = a.class_id AND ta.subject_id = a.subject_id
		WHERE ta.teacher_id = $1 AND a.attendance_date >= CURRENT_DATE - INTERVAL '7 days'
	`, teacherID).Scan(&stats.RecentAttendance)
	if err != nil {
		stats.RecentAttendance = 0
	}

	dashboard := map[string]interface{}{
		"teacher": map[string]interface{}{
			"id":       t.ID,
			"name":     t.Name,
			"email":    t.Email,
			"subject":  t.Subject,
			"phone":    t.Phone,
			"address":  t.Address.String,
			"hireDate": t.HireDate,
		},
		"statistics": stats,
	}

	writeJSON(w, http.StatusOK, dashboard)
}
