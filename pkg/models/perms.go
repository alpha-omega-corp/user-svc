package models

import "github.com/uptrace/bun"

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r"`

	Id          int64        `json:"id" bun:"id,pk,autoincrement"`
	Name        string       `json:"name" bun:"name,unique"`
	Permissions []Permission `bun:"m2m:role_to_permissions,join:Role=Permission"`
}

type Permission struct {
	Id        int64 `json:"id" bun:"id,pk,autoincrement"`
	Read      bool  `json:"read" bun:"read"`
	Write     bool  `json:"write" bun:"write"`
	Manage    bool  `json:"manage" bun:"manage"`
	ServiceID int64
}

type Service struct {
	Id          int64         `json:"id" bun:"id,pk,autoincrement"`
	Name        string        `json:"name" bun:"name,unique"`
	Permissions []*Permission `bun:"rel:has-many,join:id=service_id"`
}

type RoleToPermission struct {
	RoleID       int64       `bun:",pk"`
	Role         *Role       `bun:"rel:belongs-to,join:role_id=id"`
	PermissionID int64       `bun:",pk"`
	Permission   *Permission `bun:"rel:belongs-to,join:permission_id=id"`
}
