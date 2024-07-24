package handlers

import (
	"context"
	"github.com/alpha-omega-corp/user-svc/pkg/models"
	"github.com/alpha-omega-corp/user-svc/proto"
	"github.com/uptrace/bun"
	"net/http"
)

type UserService interface {
	GetAll(ctx context.Context) (*proto.GetUsersResponse, error)
	GetOne(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error)
	Create(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error)
	Update(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error)
	Delete(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error)
	Assign(ctx context.Context, req *proto.AssignUserRequest) (*proto.AssignUserResponse, error)
}

type userService struct {
	UserService

	db *bun.DB
}

func NewUserService(db *bun.DB) UserService {
	return &userService{
		db: db,
	}
}

func (s *userService) GetOne(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user := new(models.User)

	err := s.db.NewSelect().Model(&user).Relation("Roles").Where("id = ?", req.Id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.GetUserResponse{
		User: &proto.User{},
	}, nil
}

func (s *userService) GetAll(ctx context.Context) (*proto.GetUsersResponse, error) {
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

func (s *userService) Create(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	_, err := s.db.NewInsert().Model(&models.User{
		Name:  req.Name,
		Email: req.Email,
	}).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.CreateUserResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *userService) Update(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {

	if err := s.db.NewSelect().Model(&models.User{
		Name: req.Name,
	}).
		Where("id = ?", req.Id).Scan(ctx); err != nil {
		return nil, err
	}

	return &proto.UpdateUserResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *userService) Delete(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	_, err := s.db.NewDelete().Model(&models.User{}).Where("id = ?", req.Id).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteUserResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *userService) Assign(ctx context.Context, req *proto.AssignUserRequest) (*proto.AssignUserResponse, error) {
	userRoles := new([]models.UserToRole)

	if err := s.db.NewSelect().Model(userRoles).Where("user_id = ?", req.UserId).Scan(ctx); err != nil {
		return nil, err
	}

	requestRoles := make(map[int64]int64, len(req.Roles))
	currentRoles := make(map[int64]int64, len(*userRoles))

	for idx, userRole := range *userRoles {
		currentRoles[userRole.RoleID] = int64(idx)
	}

	for idx, reqRole := range req.Roles {
		requestRoles[reqRole] = int64(idx)
	}

	// Add roles that are in the request
	for _, roleId := range req.Roles {
		if _, ok := currentRoles[roleId]; !ok {
			_, err := s.db.NewInsert().Model(&models.UserToRole{
				UserID: req.UserId,
				RoleID: roleId,
			}).Exec(ctx)

			if err != nil {
				return nil, err
			}
		}
	}

	// Delete user's roles that are not in the request
	for roleId := range currentRoles {
		if _, ok := requestRoles[roleId]; !ok {
			_, err := s.db.NewDelete().Model(&models.UserToRole{}).
				Where("user_id = ?", req.UserId).
				Where("role_id = ?", roleId).
				Exec(ctx)

			if err != nil {
				return nil, err
			}
		}
	}

	return &proto.AssignUserResponse{
		Status: http.StatusCreated,
	}, nil
}
