package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain-text password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword compares a hashed password with a plain-text password
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
