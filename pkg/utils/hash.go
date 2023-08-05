package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(pw string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pw), 5)

	return string(bytes)
}

func CheckPasswordHash(pw string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))

	return err == nil
}
