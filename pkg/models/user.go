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

type UserToRole struct {
	UserID int64 `bun:",pk"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`
	RoleID int64 `bun:",pk"`
	Role   *Role `bun:"rel:belongs-to,join:role_id=id"`
}
