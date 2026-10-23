
// package main

// import (
// 	"fmt"

// 	"golang.org/x/crypto/bcrypt"
// )

// func main() {
// 	// Replace this with the exact password you are testing.
// 	password := "Sneha@123"

// 	// Paste the stored hash from PostgreSQL here.
// 	hash := "PASTE_TEACHER_PASSWORD_HASH_HERE"

// 	err := bcrypt.CompareHashAndPassword(
// 		[]byte(hash),
// 		[]byte(password),
// 	)

// 	if err == nil {
// 		fmt.Println("PASSWORD MATCHES")
// 	} else {
// 		fmt.Println("PASSWORD DOES NOT MATCH")
// 	}
// }
