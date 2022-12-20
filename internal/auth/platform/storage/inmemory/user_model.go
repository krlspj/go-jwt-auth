package inmemory

import "github.com/krlspj/go-jwt-auth/internal/auth/domain"

// userDB is the user type in the database
type userDB struct {
	id          string
	name        string
	lastname    string
	password    string
	token       string
	createdTime string
	updatedTime string
}

// toDomainUser converts the userDB type to domain.User type
func (u *userDB) toDomainUser() domain.User {
	user := new(domain.User)

	user.SetID(u.id)
	user.SetName(u.name)
	user.SetLastname(u.lastname)
	user.SetPassword(u.password)
	user.SetToken(u.token)

	return *user
}

// toDomainUsers convert a slice of usersDB to domain user slice
func toDomainUsers(usersDB []userDB) ([]domain.User, error) {
	users := make([]domain.User, 0)
	for _, u := range usersDB {
		users = append(users, u.toDomainUser())
	}
	return users, nil
}
