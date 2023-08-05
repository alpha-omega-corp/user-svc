package utils

import (
	"errors"
	"github.com/alpha-omega-corp/auth-svc/pkg/models"
	"github.com/golang-jwt/jwt"
	"time"
)

type AuthWrapper struct {
	secretKey string
	expiresAt int64

	provider string
}

func NewAuthWrapper(key string) *AuthWrapper {
	return &AuthWrapper{
		secretKey: key,
		expiresAt: 24,

		provider: "auth-svc",
	}
}

type AuthClaims struct {
	jwt.StandardClaims
	Id    int64
	Email string
}

func (w *AuthWrapper) GenerateToken(user models.User) (signedToken string, err error) {
	claims := &AuthClaims{
		Id:    user.Id,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(w.expiresAt)).Unix(),
			Issuer:    w.provider,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(w.secretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (w *AuthWrapper) ValidateToken(signedToken string) (claims *AuthClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(w.secretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*AuthClaims)

	if !ok {
		return nil, errors.New("unable to parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}

	return claims, nil

}
