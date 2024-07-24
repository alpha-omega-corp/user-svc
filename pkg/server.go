package pkg

import (
	"context"
	"github.com/alpha-omega-corp/user-svc/pkg/handlers"
	"github.com/alpha-omega-corp/user-svc/pkg/utils"
	"github.com/alpha-omega-corp/user-svc/proto"
	"github.com/uptrace/bun"
)

type Server struct {
	proto.UnimplementedUserServiceServer

	authService handlers.AuthService
	permService handlers.PermService
	roleService handlers.RoleService
	userService handlers.UserService
}

func NewServer(db *bun.DB, w *utils.AuthWrapper) *Server {
	return &Server{
		authService: handlers.NewAuthService(w, db),
		permService: handlers.NewPermService(db),
		roleService: handlers.NewRoleService(db),
		userService: handlers.NewUserService(db),
	}
}

func (s *Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	return s.userService.Create(ctx, req)
}
func (s *Server) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	return s.userService.GetAll(ctx)
}
func (s *Server) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	return s.userService.Update(ctx, req)
}
func (s *Server) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	return s.userService.Delete(ctx, req)
}
func (s *Server) AssignUser(ctx context.Context, req *proto.AssignUserRequest) (*proto.AssignUserResponse, error) {
	return s.userService.Assign(ctx, req)
}
func (s *Server) GetServices(ctx context.Context, req *proto.GetServicesRequest) (*proto.GetServicesResponse, error) {
	return s.permService.GetServices(ctx)
}
func (s *Server) CreateServicePermission(ctx context.Context, req *proto.CreateServicePermissionsRequest) (*proto.CreateServicePermissionsResponse, error) {
	return s.permService.CreateServicePermissions(ctx, req)
}
func (s *Server) GetServicePermissions(ctx context.Context, req *proto.GetServicePermissionsRequest) (*proto.GetServicePermissionsResponse, error) {
	return s.permService.GetServicePermissions(ctx, req)
}
func (s *Server) GetUserPermissions(ctx context.Context, req *proto.GetUserPermissionsRequest) (*proto.GetUserPermissionsResponse, error) {
	return s.permService.GetUserPermissions(ctx, req)
}
func (s *Server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return s.authService.Login(ctx, req)
}
func (s *Server) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return s.authService.Register(ctx, req)
}
func (s *Server) Validate(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	return s.authService.Validate(ctx, req)
}
func (s *Server) GetRoles(ctx context.Context, req *proto.GetRolesRequest) (*proto.GetRolesResponse, error) {
	return s.roleService.GetAll(ctx)
}
func (s *Server) CreateRole(ctx context.Context, req *proto.CreateRoleRequest) (*proto.CreateRoleResponse, error) {
	return s.roleService.Create(ctx, req)
}
