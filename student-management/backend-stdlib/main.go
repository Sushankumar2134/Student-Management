package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	connectDatabase()
	defer closeDatabase()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /students", handleGetStudents)
	mux.HandleFunc("GET /students/{id}", handleGetStudentByID)

	mux.HandleFunc("POST /students", chain(handleCreateStudent, jwtMiddleware, adminOnly))
	mux.HandleFunc("PUT /students/{id}", chain(handleUpdateStudent, jwtMiddleware, adminOnly))
	mux.HandleFunc("DELETE /students/{id}", chain(handleDeleteStudent, jwtMiddleware, adminOnly))

	mux.HandleFunc("POST /register", handleRegister)
	mux.HandleFunc("POST /login", handleLogin)

	mux.HandleFunc("GET /profile", chain(handleGetProfile, jwtMiddleware, studentOnly))
	mux.HandleFunc("PUT /profile", chain(handleUpdateProfile, jwtMiddleware, studentOnly))

	mux.HandleFunc("PUT /change-password", chain(handleChangePassword, jwtMiddleware))
	mux.HandleFunc("POST /forgot-password", handleForgotPassword)

	mux.HandleFunc("GET /statistics", chain(handleGetStatistics, jwtMiddleware, adminOnly))

	handler := corsMiddleware(mux)

	server := &http.Server{
		Addr:         ":8081",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server running at http://localhost:8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited properly")
}

func chain(h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler = h
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		handler.ServeHTTP(w, r)
	}
}