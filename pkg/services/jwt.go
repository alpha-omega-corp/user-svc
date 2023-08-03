package services

import (
	"context"
	"fmt"
	"github.com/alpha-omega-corp/authentication-svc/pkg/models"
	"github.com/alpha-omega-corp/authentication-svc/pkg/services/pb"
	"github.com/golang-jwt/jwt/v5"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	db *bun.DB
}

func NewServer(db *bun.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	encPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, errDb := s.db.NewInsert().Model(&models.User{
		Email:    req.Email,
		Password: string(encPw),
	}).Exec(ctx)

	if errDb != nil {
		return nil, errDb
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user := new(models.User)
	err := s.db.NewSelect().Model(user).Where("email = ?", req.Email).Scan(ctx, user)
	if err != nil {
		return nil, err
	}

	if !user.Verify(req.Password) {
		return nil, fmt.Errorf("authentication failed for user %s", req.Email)
	}

	jwtToken, err := user.CreateToken()
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  jwtToken,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	token, err := parseToken(req.Token)
	if err != nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	user := new(models.User)
	errDb := s.db.NewSelect().Model(user).Where("email = ?", claims["email"].(string)).Scan(ctx, user)
	if errDb != nil {
		return nil, errDb
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id,
	}, nil
}

func parseToken(ts string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		signingMethod, isValid := token.Method.(*jwt.SigningMethodHMAC)

		if isValid {
			return nil, fmt.Errorf("unexpected signing method: %v", signingMethod)
		}

		return []byte(secret), nil
	})
}
