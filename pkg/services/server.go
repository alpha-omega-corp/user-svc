package services

import (
	"context"
	"github.com/alpha-omega-corp/auth-svc/pkg/models"
	"github.com/alpha-omega-corp/auth-svc/pkg/utils"
	"github.com/alpha-omega-corp/auth-svc/proto"
	"github.com/uptrace/bun"
	"net/http"
)

type Server struct {
	proto.UnimplementedAuthServiceServer

	db          *bun.DB
	authWrapper *utils.AuthWrapper
}

func NewServer(db *bun.DB, authWrapper *utils.AuthWrapper) *Server {
	return &Server{
		db:          db,
		authWrapper: authWrapper,
	}
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

	token, _ := s.authWrapper.GenerateToken(user)

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
	claims, err := s.authWrapper.ValidateToken(req.Token)

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
