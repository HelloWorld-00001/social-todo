package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func HashPasswordWithSalt(password, salt string) (string, error) {
	combined := password + salt
	// Pre-hash with SHA256 to stay within bcrypt's 72-byte limit
	sha := sha256.Sum256([]byte(combined))

	// Hash with bcrypt
	hashed, err := bcrypt.GenerateFromPassword(sha[:], bcrypt.DefaultCost)
	return string(hashed), err
}

func GenerateSalt(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// ComparePasswordWithSalt checks if password + salt matches stored hash
func ComparePasswordWithSalt(password, salt, hashedPassword string) bool {
	// Recreate the same SHA-256 hash as during password creation
	combined := password + salt
	sha := sha256.Sum256([]byte(combined))

	// Compare with bcrypt-hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), sha[:])
	return err == nil
}
