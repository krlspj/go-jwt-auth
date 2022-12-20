package service

import (
	"context"

	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
)

type AuthService interface {
	FindUsers(ctx context.Context) ([]domain.User, error)
}

type authService struct {
	userRepo domain.UserRepo
}

func NewAuthService(ur domain.UserRepo) *authService {
	return &authService{
		userRepo: ur,
	}
}

func (s *authService) FindUsers(ctx context.Context) ([]domain.User, error) {
	return s.userRepo.FindAll(ctx)
}
