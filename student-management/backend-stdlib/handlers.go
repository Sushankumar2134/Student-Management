package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func handleGetStudents(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	query := `
		SELECT id, name, email, age, gender, city
		FROM students
	`
	args := []interface{}{}

	if search != "" {
		query += ` WHERE name ILIKE $1 OR email ILIKE $1 OR city ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY id`

	rows, err := db.QueryContext(r.Context(), query, args...)
	if err != nil {
		writeJSONError(w, "Failed to get students", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	students := []Student{}

	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Age, &s.Gender, &s.City); err != nil {
			writeJSONError(w, "Failed to read student", http.StatusInternalServerError)
			log.Println("Scan error:", err)
			return
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		writeJSONError(w, "Failed to read students", http.StatusInternalServerError)
		log.Println("Rows error:", err)
		return
	}

	writeJSON(w, http.StatusOK, students)
}

func handleGetStudentByID(w http.ResponseWriter, r *http.Request) {
	idStr := getPathParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var s Student
	err = db.QueryRowContext(r.Context(),
		`SELECT id, name, email, age, gender, city FROM students WHERE id = $1`, id).
		Scan(&s.ID, &s.Name, &s.Email, &s.Age, &s.Gender, &s.City)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "Student not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Failed to get student", http.StatusInternalServerError)
		log.Println("Query error:", err)
		return
	}

	writeJSON(w, http.StatusOK, s)
}

func handleCreateStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	if err := readJSON(r, &student); err != nil {
		writeJSONError(w, "Invalid student data", http.StatusBadRequest)
		return
	}

	if student.Name == "" {
		writeJSONError(w, "Name is required", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRowContext(r.Context(),
		`INSERT INTO students (name, email, age, gender, city) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		student.Name, student.Email, student.Age, student.Gender, student.City).
		Scan(&id)

	if err != nil {
		writeJSONError(w, "Failed to create student", http.StatusInternalServerError)
		log.Println("Insert error:", err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Student created successfully",
		"id":      id,
	})
}

func handleUpdateStudent(w http.ResponseWriter, r *http.Request) {
	idStr := getPathParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var student Student
	if err := readJSON(r, &student); err != nil {
		writeJSONError(w, "Invalid student data", http.StatusBadRequest)
		return
	}

	result, err := db.ExecContext(r.Context(),
		`UPDATE students SET name = $1, email = $2, age = $3, gender = $4, city = $5 WHERE id = $6`,
		student.Name, student.Email, student.Age, student.Gender, student.City, id)

	if err != nil {
		writeJSONError(w, "Failed to update student", http.StatusInternalServerError)
		log.Println("Update error:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		writeJSONError(w, "Student not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Student updated successfully"})
}

func handleDeleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := getPathParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		writeJSONError(w, "Failed to start transaction", http.StatusInternalServerError)
		log.Println("Transaction error:", err)
		return
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(r.Context(), "DELETE FROM users WHERE student_id = $1", id)
	if err != nil {
		writeJSONError(w, "Failed to delete user account", http.StatusInternalServerError)
		log.Println("Delete user error:", err)
		return
	}

	result, err = tx.ExecContext(r.Context(), "DELETE FROM students WHERE id = $1", id)
	if err != nil {
		writeJSONError(w, "Failed to delete student", http.StatusInternalServerError)
		log.Println("Delete student error:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		writeJSONError(w, "Student not found", http.StatusNotFound)
		return
	}

	if err = tx.Commit(); err != nil {
		writeJSONError(w, "Failed to complete deletion", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Student deleted successfully"})
}