package service

import (
	"context"

	"github.com/krlspj/go-jwt-auth/internal/user/domain"
	"github.com/krlspj/go-jwt-auth/internal/user/platform/storage/mongodb"
)

//go:generate mockery --case=snake --outpkg=mocks --output=mocks --name=UserService
type UserService interface {
	GetByID(ctx context.Context, id string) (domain.User, error)
}

type defUserService struct {
	userRepo *mongodb.UserRepository
}

func NewUserService(repo *mongodb.UserRepository) *defUserService {
	return &defUserService{
		userRepo: repo,
	}
}

func (s *defUserService) GetByID(ctx context.Context, id string) (domain.User, error) {
	return domain.User{ID: "12345", Name: "John Doe"}, nil
}
