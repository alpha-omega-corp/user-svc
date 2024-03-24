package server

import (
	"context"
	"fmt"
	"github.com/alpha-omega-corp/user-svc/pkg/models"
	"github.com/alpha-omega-corp/user-svc/pkg/utils"
	"github.com/alpha-omega-corp/user-svc/proto"
	"github.com/uptrace/bun"
	"net/http"
	"strings"
)

type Server struct {
	proto.UnimplementedUserServiceServer

	auth *utils.AuthWrapper
	db   *bun.DB
}

func NewServer(db *bun.DB, w *utils.AuthWrapper) *Server {
	return &Server{
		auth: w,
		db:   db,
	}
}

func (s *Server) CreateRole(ctx context.Context, req *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error) {
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

func (s *Server) GetRoles(ctx context.Context, req *proto.GetRolesRequest) (*proto.GetRolesResponse, error) {
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

func (s *Server) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	var users []*models.User

	err := s.db.NewSelect().Model(&users).Relation("Roles").Scan(ctx)
	if err != nil {
		return nil, err
	}

	var resSlice []*proto.User
	for _, user := range users {
		rolesSlice := make([]*proto.Role, len(user.Roles))

		for index, role := range user.Roles {
			rolesSlice[index] = &proto.Role{
				Id:   role.Id,
				Name: role.Name,
			}
		}

		protoUser := &proto.User{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Roles: rolesSlice,
		}

		resSlice = append(resSlice, protoUser)
	}

	return &proto.GetUsersResponse{
		Users: resSlice,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	user := new(models.User)
	if err := s.db.NewSelect().Model(user).Where("id = ?", req.Id).Scan(ctx); err != nil {
		return nil, err
	}

	user.Name = req.Name
	_, err := s.db.NewUpdate().Model(user).Where("id = ?", req.Id).Exec(ctx)
	if err != nil {
		return nil, err
	}

	for _, roleId := range req.Roles {
		_, err := s.db.NewDelete().Model(&models.UserToRole{}).
			Where("user_id = ?", req.Id).
			Where("role_id = ?", roleId).
			Exec(ctx)

		if err != nil {
			return nil, err
		}
	}

	return &proto.UpdateUserResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) GetPermServices(ctx context.Context, req *proto.GetPermServicesRequest) (*proto.GetPermServicesResponse, error) {
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

	return &proto.GetPermServicesResponse{
		Services: resSlice,
	}, nil
}

func (s *Server) CreatePermissions(ctx context.Context, req *proto.CreatePermissionRequest) (*proto.CreatePermissionResponse, error) {
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

	return &proto.CreatePermissionResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) GetPermissions(ctx context.Context, req *proto.GetPermissionsRequest) (*proto.GetPermissionsResponse, error) {
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

	return &proto.GetPermissionsResponse{
		Permissions: resSlice,
	}, nil
}

func (s *Server) GetUserPermissions(ctx context.Context, req *proto.GetUserPermissionsRequest) (*proto.GetUserPermissionsResponse, error) {
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

func (s *Server) AssignRole(ctx context.Context, req *proto.AssignRoleRequest) (*proto.AssignRoleResponse, error) {
	_, err := s.db.NewInsert().Model(&models.UserToRole{
		UserID: req.UserId,
		RoleID: req.RoleId,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &proto.AssignRoleResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	user := new(models.User)

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)

	_, err := s.db.NewInsert().Model(user).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &proto.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	var user models.User
	err := s.db.NewSelect().Model(&user).Where("email = ?", req.Email).Scan(ctx, &user)
	if err != nil {
		return nil, err
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		return &proto.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := s.auth.GenerateToken(user)

	return &proto.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
		User: &proto.User{
			Id:    user.Id,
			Email: user.Email,
		},
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, err := s.auth.ValidateToken(req.Token)

	if err != nil {
		return &proto.ValidateResponse{
			Status: http.StatusForbidden,
			Error:  err.Error(),
		}, nil
	}

	var user models.User
	err = s.db.NewSelect().Model(&user).Where("email = ?", claims.Email).Scan(ctx, &user)
	if err != nil {
		return &proto.ValidateResponse{
			Status: http.StatusForbidden,
			Error:  "User not found",
		}, nil
	}

	return &proto.ValidateResponse{
		Status: http.StatusOK,
		User: &proto.User{
			Id:    user.Id,
			Email: user.Email,
		},
	}, nil
}
