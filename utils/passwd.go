package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func CheckPassword(password, hash string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(hash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func ValidatePassword(password, hash string) bool {
	err := CheckPassword(password, hash)
	if err != nil {
		return false
	}
	return true
}

// RandomHex will create random n bytes and return formatted as a hex string.
// The string length will be double the byte length.
func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
