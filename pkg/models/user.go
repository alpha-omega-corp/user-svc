package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id       int64  `json:"id" bun:"id,pk,autoincrement"`
	Email    string `json:"email" bun:"email,unique"`
	Password string `json:"-" bun:"encrypted_password"`
}

func (u *User) Verify(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)) == nil
}

func (u *User) CreateToken() (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"email":     u.Email,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (u *User) IsClaimValid(c jwt.MapClaims) bool {
	if u.Email != c["email"].(string) {
		return false
	}

	return true
}
