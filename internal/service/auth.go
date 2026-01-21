package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type authService struct {
	jwtSecret      string
	jwtTTL         time.Duration
	buyerPassword  string
	sellerPassword string
}

func NewAuthService(secret string, ttl time.Duration, buyerPassword string, sellerPassword string) AuthService {
	return &authService{
		jwtSecret:      secret,
		jwtTTL:         ttl,
		buyerPassword:  buyerPassword,
		sellerPassword: sellerPassword,
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	if s.jwtSecret == "" {
		return "", errors.New("JWT secret is not set (set JWT_SECRET or config jwt.secret)")
	}

	var userID int
	var role string

	switch username {
	case "buyer@test.com":
		userID = 1
		role = "buyer"
		if s.buyerPassword == "" {
			return "", errors.New("buyer password is not set (set AUTH_BUYER_PASSWORD or config auth.buyer_password)")
		}
		if password != s.buyerPassword {
			return "", errors.New("invalid credentials")
		}
	case "seller@test.com":
		userID = 2
		role = "seller"
		if s.sellerPassword == "" {
			return "", errors.New("seller password is not set (set AUTH_SELLER_PASSWORD or config auth.seller_password)")
		}
		if password != s.sellerPassword {
			return "", errors.New("invalid credentials")
		}
	default:
		return "", errors.New("user not found")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(s.jwtTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return signedToken, nil
}

