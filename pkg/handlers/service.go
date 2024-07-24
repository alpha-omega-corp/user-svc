package handlers

import (
	"context"
	"fmt"
	"github.com/alpha-omega-corp/user-svc/pkg/models"
	"github.com/alpha-omega-corp/user-svc/proto"
	"github.com/uptrace/bun"
	"net/http"
	"strings"
)

type PermService interface {
	GetServices(ctx context.Context) (*proto.GetServicesResponse, error)
	CreateServicePermissions(ctx context.Context, req *proto.CreateServicePermissionsRequest) (*proto.CreateServicePermissionsResponse, error)
	GetServicePermissions(ctx context.Context, req *proto.GetServicePermissionsRequest) (*proto.GetServicePermissionsResponse, error)
	GetUserPermissions(ctx context.Context, req *proto.GetUserPermissionsRequest) (*proto.GetUserPermissionsResponse, error)
}

type permService struct {
	db *bun.DB
}

func NewPermService(db *bun.DB) PermService {
	return &permService{
		db: db,
	}
}

func (s *permService) GetServices(ctx context.Context) (*proto.GetServicesResponse, error) {
	var services []models.Service
	if err := s.db.NewSelect().Model(&services).Scan(ctx); err != nil {
		return nil, err
	}

	var resSlice []*proto.Service
	for _, service := range services {
		resSlice = append(resSlice, &proto.Service{
			Id:   service.Id,
			Name: service.Name,
		})
	}

	return &proto.GetServicesResponse{
		Services: resSlice,
	}, nil
}

func (s *permService) CreateServicePermissions(ctx context.Context, req *proto.CreateServicePermissionsRequest) (*proto.CreateServicePermissionsResponse, error) {
	permissions := &models.Permission{
		Read:      req.CanRead,
		Write:     req.CanWrite,
		Manage:    req.CanManage,
		ServiceID: req.ServiceId,
		RoleId:    req.RoleId,
	}

	_, err := s.db.NewInsert().Model(permissions).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.CreateServicePermissionsResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *permService) GetServicePermissions(ctx context.Context, req *proto.GetServicePermissionsRequest) (*proto.GetServicePermissionsResponse, error) {
	var service models.Service
	if err := s.db.NewSelect().
		Model(&service).
		Relation("Permissions").
		Where("id = ?", req.ServiceId).
		Scan(ctx); err != nil {
		return nil, err
	}

	resSlice := make([]*proto.Permission, len(service.Permissions))
	for index, permission := range service.Permissions {
		role := new(models.Role)
		if err := s.db.NewSelect().
			Model(role).
			Where("id  = ?", permission.RoleId).
			Scan(ctx); err != nil {
			return nil, err
		}

		resSlice[index] = &proto.Permission{
			Id: permission.Id,
			Service: &proto.Service{
				Id:   service.Id,
				Name: service.Name,
			},
			Role: &proto.Role{
				Id:   role.Id,
				Name: role.Name,
			},
			CanRead:   permission.Read,
			CanWrite:  permission.Write,
			CanManage: permission.Manage,
		}
	}

	return &proto.GetServicePermissionsResponse{
		Permissions: resSlice,
	}, nil
}

func (s *permService) GetUserPermissions(ctx context.Context, req *proto.GetUserPermissionsRequest) (*proto.GetUserPermissionsResponse, error) {
	user := new(models.User)
	if err := s.db.NewSelect().
		Model(user).
		Relation("Roles").
		Where("id = ?", req.UserId).
		Scan(ctx); err != nil {
		return nil, err
	}

	var permSlice []models.Permission
	for _, role := range user.Roles {
		if err := s.db.NewSelect().
			Model(&role).
			Relation("Permissions").
			Where("id = ?", role.Id).
			Scan(ctx); err != nil {
			return nil, err
		}

		permSlice = append(permSlice, role.Permissions...)
	}

	permMap := make(map[string]bool)
	for index, perm := range permSlice {
		service := new(models.Service)
		if err := s.db.NewSelect().
			Model(service).
			Where("id = ?", perm.ServiceID).
			Scan(ctx); err != nil {
			return nil, err
		}

		svc := strings.ToLower(service.Name)
		idxRead := fmt.Sprintf("%s.read", svc)
		idxWrite := fmt.Sprintf("%s.write", svc)
		idxManage := fmt.Sprintf("%s.manage", svc)

		if index > 0 {
			if permMap[idxRead] != true {
				permMap[idxRead] = perm.Read
			}
			if permMap[idxWrite] != true {
				permMap[idxWrite] = perm.Write
			}
			if permMap[idxManage] != true {
				permMap[idxManage] = perm.Manage
			}
		} else {
			permMap[idxRead] = perm.Read
			permMap[idxWrite] = perm.Write
			permMap[idxManage] = perm.Manage
		}
	}

	return &proto.GetUserPermissionsResponse{
		Matrix: permMap,
	}, nil
}
