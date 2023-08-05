package models

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id       int64  `json:"id" bun:"id,pk,autoincrement"`
	Email    string `json:"email" bun:"email,unique"`
	Password string `json:"-" bun:"encrypted_password"`
}
