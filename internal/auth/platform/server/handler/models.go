package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
)

type userReq struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"username" validate:"required,min=3,max=15"`
	LastName string `json:"lastname,omitemtpy"`
	Password string `json:"password" validate:"required,min=8,max=30"`
	Email    string `json:"email,omitempty"`
}

func (m *userReq) toDomainUser() domain.User {
	var u domain.User
	u.SetID(m.ID)
	u.SetName(m.Name)
	u.SetLastname(m.LastName)
	u.SetPassword(m.Password)

	return u
}
func (m *userReq) ValidateUser() error {
	validate := validator.New()
	if err := validate.Struct(m); err != nil {
		return err
	}
	return nil
}

type userResp struct {
	ID       string `json:"id"`
	Name     string `json:"username"`
	LastName string `json:"lastname"`
	//Password string `json:"password"`
	//Token    string `json:"token,omitempty"`
}

func toUserResp(user domain.User) userResp {
	return userResp{
		ID:       user.ID(),
		Name:     user.Name(),
		LastName: user.Lastname(),
		//Token:    user.Token(),
		//Password: user.Password(),
	}
}

func toUserRestList(users []domain.User) ([]userResp, error) {
	uRest := make([]userResp, 0)
	for _, u := range users {
		uRest = append(uRest, toUserResp(u))
	}
	return uRest, nil
}
