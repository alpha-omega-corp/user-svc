package services

import (
	"context"
	"github.com/alpha-omega-corp/user-svc/pkg/models"
	"github.com/alpha-omega-corp/user-svc/proto"
	"github.com/uptrace/bun"
	"net/http"
)

type RoleService interface {
	GetAll(ctx context.Context) (*proto.GetRolesResponse, error)
	Create(ctx context.Context, req *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error)
}

type roleService struct {
	RoleService
	db *bun.DB
}

func NewRoleService(db *bun.DB) RoleService {
	return &roleService{
		db: db,
	}
}

func (s *roleService) GetAll(ctx context.Context) (*proto.GetRolesResponse, error) {
	var roles []*models.Role

	err := s.db.NewSelect().Model(&roles).Scan(ctx)
	if err != nil {
		return nil, err
	}

	var resSlice []*proto.Role
	for _, role := range roles {
		resSlice = append(resSlice, &proto.Role{
			Id:   role.Id,
			Name: role.Name,
		})
	}

	return &proto.GetRolesResponse{
		Roles: resSlice,
	}, nil
}

func (s *roleService) Create(ctx context.Context, req *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error) {
	role := new(models.Role)
	role.Name = req.Name

	_, err := s.db.NewInsert().Model(role).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &proto.CreateRoleResponse{
		Status: http.StatusCreated,
	}, nil
}
