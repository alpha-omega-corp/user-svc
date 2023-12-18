package models

import "github.com/uptrace/bun"

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r"`

	Id          int64        `json:"id" bun:"id,pk,autoincrement"`
	Name        string       `json:"name" bun:"name,unique"`
	Permissions []Permission `bun:"rel:has-many,join:id=role_id"`
}

type Service struct {
	Id          int64        `json:"id" bun:"id,pk,autoincrement"`
	Name        string       `json:"name" bun:"name,unique"`
	Permissions []Permission `bun:"rel:has-many,join:id=service_id"`
}

type Permission struct {
	Id        int64 `json:"id" bun:"id,pk,autoincrement"`
	Read      bool  `json:"read" bun:"read"`
	Write     bool  `json:"write" bun:"write"`
	Manage    bool  `json:"manage" bun:"manage"`
	RoleId    int64
	ServiceID int64
}
