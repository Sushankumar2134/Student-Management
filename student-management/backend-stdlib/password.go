package main

import "golang.org/x/crypto/bcrypt"

// HashPassword converts a plain-text password into a bcrypt hash.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// CheckPassword verifies a plain-text password against its bcrypt hash.
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)

	return err == nil
}
