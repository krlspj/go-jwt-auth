package service

import (
	"context"
	"log"

	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	FindUsers(ctx context.Context) ([]domain.User, error)
	FindUser(ctx context.Context, id string) (domain.User, error)
	FindUserByField(ctx context.Context, field, value string) (domain.User, error)
	LoginUser(ctx context.Context, username, password string) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User) error
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

func (s *authService) FindUserByField(ctx context.Context, field, value string) (domain.User, error) {
	return s.userRepo.FindOneByField(ctx, field, value)
}

func (s *authService) FindUser(ctx context.Context, id string) (domain.User, error) {
	return s.userRepo.FindOne(ctx, id)
}

func (s *authService) LoginUser(ctx context.Context, username, password string) (domain.User, error) {
	// find user in db
	user, err := s.userRepo.FindOneByField(ctx, "username", username)
	if err != nil {
		return domain.User{}, err
	}

	// verify passwords
	if err := user.VerifyPassword(password); err != nil {
		return domain.User{}, err
	}

	// (TODO) create a refresh token and return it

	return user, nil

}

func (s *authService) CreateUser(ctx context.Context, user domain.User) error {
	// user domain validation (disabled) need to create public struct fields to user validate on them
	//if err := user.ValidateUser(); err != nil {
	//	return err
	//}

	// Verify user

	// config new user
	if err := user.HashPassword(); err != nil {
		return err
	}

	// Insert user in database
	if err := s.userRepo.InsertUser(ctx, user); err != nil {
		return err
	}
	return nil
}

/*_____________ PRIVATE FUNCTIONS ____________*/

// HashPassword returns an encrypted passwor
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // this hash generation already uses a salt in it
	if err != nil {
		//log.Panic(err)
		return "", err
	}

	return string(bytes), nil
}

// VerifyPassword returns true if passwords maches. When not, false and a message
func verifyPassword(userPassword string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		log.Println("VerifyPassword:", err.Error())
		return err
	}

	return nil
}
