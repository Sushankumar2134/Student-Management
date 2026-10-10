package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func connectDatabase() {
	dsn := getDSN()
	var err error

	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Database configuration error:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err = db.Ping(); err != nil {
		log.Fatal("PostgreSQL connection failed:", err)
	}

	log.Println("PostgreSQL connected successfully!")
}

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

func closeDatabase() {
	if db != nil {
		db.Close()
	}
}