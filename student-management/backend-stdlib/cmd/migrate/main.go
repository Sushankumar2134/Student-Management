package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

// This is a ONE-TIME migration script to convert existing plain-text passwords
// to bcrypt hashes in all three account tables.
// Run this ONCE after deploying the bcrypt changes.
// BACKUP YOUR DATABASE BEFORE RUNNING THIS SCRIPT.
//
// Usage:
//   go run ./cmd/migrate/main.go
//
// The script:
// 1. Connects to the database using the same DSN as the main app
// 2. For each table (admins, student_accounts, teacher_accounts):
//    - Selects all rows with password column
//    - Checks if password is already a valid bcrypt hash (starts with $2a$, $2b$, or $2y$)
//    - If not a bcrypt hash, hashes it and updates the row
// 3. Logs progress and summary

func getDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "student_db")

	if password != "" {
		return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname
	}
	return "postgres://" + user + "@" + host + ":" + port + "/" + dbname
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func isBcryptHash(password string) bool {
	// Bcrypt hashes start with $2a$, $2b$, or $2y$ followed by cost parameter
	return len(password) >= 4 && password[0] == '$' && password[1] == '2' && (password[2] == 'a' || password[2] == 'b' || password[2] == 'y') && password[3] == '$'
}

func main() {
	dsn := getDSN()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Connected to database. Starting password migration...")

	tables := []string{"admins", "student_accounts", "teacher_accounts"}
	totalUpdated := 0
	totalSkipped := 0
	totalErrors := 0

	for _, table := range tables {
		log.Printf("Processing table: %s", table)

		rows, err := db.Query(`SELECT id, password FROM ` + table)
		if err != nil {
			log.Printf("ERROR querying %s: %v", table, err)
			totalErrors++
			continue
		}

		type account struct {
			id       int
			password string
		}
		var accounts []account

		for rows.Next() {
			var a account
			if err := rows.Scan(&a.id, &a.password); err != nil {
				log.Printf("ERROR scanning row in %s: %v", table, err)
				totalErrors++
				continue
			}
			accounts = append(accounts, a)
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			log.Printf("ERROR iterating rows in %s: %v", table, err)
			totalErrors++
			continue
		}

		for _, a := range accounts {
			if isBcryptHash(a.password) {
				log.Printf("  ID %d: Already bcrypt hash, skipping", a.id)
				totalSkipped++
				continue
			}

			// Hash the plain-text password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("  ID %d: ERROR hashing password: %v", a.id, err)
				totalErrors++
				continue
			}

			// Update the row
			_, err = db.Exec(`UPDATE `+table+` SET password = $1 WHERE id = $2`, string(hashedPassword), a.id)
			if err != nil {
				log.Printf("  ID %d: ERROR updating password: %v", a.id, err)
				totalErrors++
				continue
			}

			log.Printf("  ID %d: Password hashed and updated", a.id)
			totalUpdated++
		}

		log.Printf("Table %s done. Processed: %d", table, len(accounts))
	}

	log.Println("=== Migration Summary ===")
	log.Printf("Total passwords updated: %d", totalUpdated)
	log.Printf("Total passwords skipped (already hashed): %d", totalSkipped)
	log.Printf("Total errors: %d", totalErrors)

	if totalErrors > 0 {
		log.Println("WARNING: Some errors occurred. Please review the logs above.")
	} else {
		log.Println("Migration completed successfully!")
	}
}
