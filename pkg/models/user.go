package models

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id       int64  `json:"id" bun:"id,pk,autoincrement"`
	Name     string `json:"name" bun:"name"`
	Email    string `json:"email" bun:"email,unique"`
	Password string `json:"-" bun:"encrypted_password"`
	Roles    []Role `bun:"m2m:user_to_roles,join:User=Role"`
}

type UserToRole struct {
	UserID int64 `bun:",pk"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`
	RoleID int64 `bun:",pk"`
	Role   *Role `bun:"rel:belongs-to,join:role_id=id"`
}
