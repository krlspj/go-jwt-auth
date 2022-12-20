package inmemory

import (
	"context"

	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
)

type UserRepositoryStub struct {
	//users []domain.User
	users []userDB
}

func NewUserRepositoryStub() *UserRepositoryStub {
	newUsers := []userDB{
		{id: "1000", name: "John", lastname: "Doe", password: "example"},
		{id: "10001", name: "Mickye", lastname: "Mouse", password: "exmaple"},
	}
	return &UserRepositoryStub{
		users: newUsers,
	}
}

// FindAll return all users stored in memory database
func (s *UserRepositoryStub) FindAll(ctx context.Context) ([]domain.User, error) {
	return toDomainUsers(s.users)
}
