package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware checks whether the user has a valid JWT token.
func JWTMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Get Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check whether Authorization header exists
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization token required",
			})
			c.Abort()
			return
		}

		// Expected format:
		// Authorization: Bearer TOKEN
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Check and decode JWT
		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				// Make sure the token uses HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}

				// Same secret used when creating JWT
				return []byte("student-management-secret"), nil
			},
		)

		// Invalid token
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Get claims from token
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Store user information in Gin context
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		c.Set("student_id", claims["student_id"])

		// Continue to the requested API
		c.Next()
	}
}
// AdminOnly allows only users with the admin role.
func AdminOnly() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Get role stored by JWTMiddleware
		role, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User role not found",
			})
			c.Abort()
			return
		}

		// Check whether user is admin
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Admin access required",
			})
			c.Abort()
			return
		}

		// User is admin → continue
		c.Next()
	}
}
// StudentOnly allows only student users.
func StudentOnly() gin.HandlerFunc {
	return func(c *gin.Context) {

		role, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User role not found",
			})
			c.Abort()
			return
		}

		if role != "student" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Student access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}