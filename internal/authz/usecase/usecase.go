package usecase

import (
	"context"
	"errors"

	"github.com/krlspj/go-jwt-auth/internal/authz/domain"
	"github.com/krlspj/go-jwt-auth/internal/jwt/service"
)

var (
	ErrNoAdmin error = errors.New("User is not Admin")
)

type AuthzUsecase interface {
	CheckValidToken(token string) error
	IsAdmin(ctx context.Context, token string) error
}

type defAuthzUsecase struct {
	jwtService service.JwtService
	userRepo   domain.AuthzUserRepo
}

func NewAuthzUsecase(jwts service.JwtService, aur domain.AuthzUserRepo) *defAuthzUsecase {
	return &defAuthzUsecase{
		jwtService: jwts,
		userRepo:   aur,
	}
}

func (uc *defAuthzUsecase) CheckValidToken(token string) error {
	//claims, err := service.ValidateToken(token)
	_, err := uc.jwtService.ValidateToken(token)

	//fmt.Println("claims", claims)
	if err != nil {
		return err
	}
	return nil
}

func (uc *defAuthzUsecase) IsAdmin(ctx context.Context, token string) error {
	claims, err := uc.jwtService.ValidateToken(token)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindOne(ctx, claims.UserId)
	if user.IsAdmin() == "true" {
		return nil
	} else {
		return ErrNoAdmin
	}
}
