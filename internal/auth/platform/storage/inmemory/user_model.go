package inmemory

import (
	"github.com/google/uuid"
	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
)

// userDB is the user type in the database
type userDB struct {
	id          string
	username    string
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
	user.SetName(u.username)
	user.SetLastname(u.lastname)
	user.SetPassword(u.password)
	user.SetToken(u.token)

	return *user
}

func (u *userDB) GenerateUUID() {
	u.id = uuid.New().String()
}

// toDomainUsers convert a slice of usersDB to domain user slice
func toDomainUsers(usersDB []userDB) ([]domain.User, error) {
	users := make([]domain.User, 0)
	for _, u := range usersDB {
		users = append(users, u.toDomainUser())
	}
	return users, nil
}

func toUserDB(user domain.User) userDB {
	/*
		udb := new(userDB)
		if user.ID() != "" {
			udb.id = user.ID()
		}
		if user.Name() != "" {
			udb.name = user.Name()
		}
		if user.Lastname() != "" {
			udb.lastname = user.Lastname()
		}
		if user.Token() != "" {
			udb.token = user.Token()
		}
		if user.Password() != "" {
			udb.password = user.Password()
		}
	*/
	return userDB{
		id:       user.ID(),
		username: user.Name(),
		lastname: user.Lastname(),
		password: user.Password(),
		token:    user.Token(),
	}

}
