package models

import "github.com/uptrace/bun"

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r"`

	Id   int64  `json:"id" bun:"id,pk,autoincrement"`
	Name string `json:"name" bun:"name,unique"`
}
