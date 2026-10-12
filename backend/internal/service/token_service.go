package service

import (
	"fmt"
	"time"

	"live-coding-starter/backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	Generate(userID uint) (string, error)
	Validate(tokenValue string) (uint, error)
}

type tokenService struct {
	secret     []byte
	expiration time.Duration
}

type authClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewTokenService(cfg config.JWTConfig) TokenService {
	return &tokenService{
		secret:     []byte(cfg.Secret),
		expiration: time.Duration(cfg.ExpirationHours) * time.Hour,
	}
}

func (s *tokenService) Generate(userID uint) (string, error) {
	now := time.Now()
	claims := authClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *tokenService) Validate(tokenValue string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &authClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return s.secret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*authClaims)
	if !ok || !token.Valid || claims.UserID == 0 {
		return 0, ErrUnauthorized
	}

	return claims.UserID, nil
}
