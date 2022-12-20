package handler

import "github.com/krlspj/go-jwt-auth/internal/auth/domain"

type userHttp struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Token    string `json:"token,omitempty"`
}

func toUserHttp(user domain.User) userHttp {
	return userHttp{
		ID:       user.ID(),
		Name:     user.Name(),
		LastName: user.Lastname(),
		Token:    user.Token(),
	}
}

func toUserRestList(users []domain.User) ([]userHttp, error) {
	uRest := make([]userHttp, 0)
	for _, u := range users {
		uRest = append(uRest, toUserHttp(u))
	}
	return uRest, nil
}
