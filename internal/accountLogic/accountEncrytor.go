package accountLogic

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// GenerateSalt creates a random 16-byte hex string
func GenerateSalt() (string, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashWithSalt hashes (salt + password) using bcrypt
func HashWithSalt(password, salt string) (string, error) {
	saltedPassword := salt + password
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	return string(hash), err
}

// VerifyWithSalt compares stored hash with (salt + inputPassword)
func VerifyWithSalt(storedHash, salt, inputPassword string) bool {
	saltedInput := salt + inputPassword
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(saltedInput))
	return err == nil
}
