package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func createUsers() {

	// Connect to PostgreSQL
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://postgres@localhost:5432/student_db",
	)

	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}

	defer conn.Close(context.Background())

	// =========================
	// ADMIN USER
	// =========================

	adminPassword := "Admin@123"

	// Convert admin password into bcrypt hash
	adminHash, err := bcrypt.GenerateFromPassword(
		[]byte(adminPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		fmt.Println("Password hashing failed:", err)
		return
	}

	// Insert admin user
	_, err = conn.Exec(
		context.Background(),
		`INSERT INTO users
		(username, password_hash, role, student_id, email)
		VALUES ($1, $2, $3, $4, $5)`,
		"akarsh",
		string(adminHash),
		"admin",
		nil,
		"aakarsh88@gmail.com",
	)

	if err != nil {
		fmt.Println("Admin creation failed:", err)
		return
	}

	fmt.Println("Admin user created successfully!")


	// =========================
	// STUDENT USER
	// =========================

	studentPassword := "Satish@123"

	// Convert student password into bcrypt hash
	studentHash, err := bcrypt.GenerateFromPassword(
		[]byte(studentPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		fmt.Println("Password hashing failed:", err)
		return
	}

	// Insert student user
	_, err = conn.Exec(
		context.Background(),
		`INSERT INTO users
		(username, password_hash, role, student_id, email)
		VALUES ($1, $2, $3, $4, $5)`,
		"satish",
		string(studentHash),
		"student",
		101,
		"satish@gmail.com",
	)

	if err != nil {
		fmt.Println("Student creation failed:", err)
		return
	}

	fmt.Println("Student user created successfully!")
}