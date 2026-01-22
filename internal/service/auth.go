package service

import (
	"context"
	"errors"
	"time"

	"paraklitshop/internal/model"
	"paraklitshop/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type authService struct {
	userRepository repository.UserRepository
	jwtSecret      string
	jwtTTL         time.Duration
}

func NewAuthService(userRepo repository.UserRepository, secret string, ttl time.Duration) AuthService {
	return &authService{
		userRepository: userRepo,
		jwtSecret:      secret,
		jwtTTL:         ttl,
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	if s.userRepository == nil {
		return "", errors.New("user repository is not initialized")
	}

	// Пытаемся найти пользователя по email (username используется как email)
	user, err := s.userRepository.GetByEmail(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(s.jwtTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return signedToken, nil
}

func (s *authService) Register(ctx context.Context, email, password, role string) (int, error) {
	if role != "buyer" && role != "seller" {
		return 0, errors.New("invalid role")
	}
	existing, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return s.userRepository.Create(ctx, &model.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	})
}
