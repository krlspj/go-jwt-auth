package mongodb

import "github.com/krlspj/go-jwt-auth/internal/auth/domain"

// userDB is the user type in the database
type userDB struct {
	Id          string `bson:"_id"`
	Name        string `bson:"name"`
	Lastname    string `bson:"lastname"`
	Password    string `bson:"password"`
	Token       string `bson:"token"`
	CreatedTime string `bson:"createdTime"`
	UpdatedTime string `bson:"updatedTime"`
}

// toDomainUser converts the userDB type to domain.User type
func (u *userDB) toDomainUser() domain.User {
	user := new(domain.User)

	user.SetID(u.Id)
	user.SetName(u.Name)
	user.SetLastname(u.Lastname)
	user.SetPassword(u.Password)
	user.SetToken(u.Token)

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
