package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// []byte(...) converts the
//
//	string into bytes because the JWT library expects a key in this form.
//
// jwtSecret stores the key in memory for signing and verification.
var jwtSecret = []byte(getJWTSecret())

// getJWTSecret() gets the secret value.
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET not set, using default (insecure for production)")
		return "your-super-secret-jwt-key-change-in-production"
	}
	return secret
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims := jwt.MapClaims{}
		//verifies the JWT
		//jwt.ParseWithClaims(...) reads the token and checks its signature.

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			writeJSONError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(jwt.MapClaims)
		if !ok {
			writeJSONError(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			writeJSONError(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func studentOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(jwt.MapClaims)
		if !ok {
			writeJSONError(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "student" {
			writeJSONError(w, "Student access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func teacherOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(jwt.MapClaims)
		if !ok {
			writeJSONError(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "teacher" {
			writeJSONError(w, "Teacher access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getUserID(r *http.Request) (int, bool) {
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, false
	}

	return int(userIDFloat), true
}

func getStudentID(r *http.Request) (int, bool) {
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	studentIDVal, ok := claims["student_id"]
	if !ok || studentIDVal == nil {
		return 0, false
	}

	studentIDFloat, ok := studentIDVal.(float64)
	if !ok {
		return 0, false
	}

	return int(studentIDFloat), true
}

func getTeacherID(r *http.Request) (int, bool) {
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	teacherIDVal, ok := claims["teacher_id"]
	if !ok || teacherIDVal == nil {
		return 0, false
	}

	teacherIDFloat, ok := teacherIDVal.(float64)
	if !ok {
		return 0, false
	}

	return int(teacherIDFloat), true
}

func getRole(r *http.Request) (string, bool) {
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		return "", false
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", false
	}

	return role, true
}
