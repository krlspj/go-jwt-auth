package mongodb

import (
	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// userDB is the user type in the database
type userMongo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"username,omitempty"`
	Lastname    string             `bson:"lastname,omitempty"`
	Password    string             `bson:"password,omitempty"`
	Token       string             `bson:"token,omitempty"`
	CreatedTime string             `bson:"createdTime,omitempty"`
	UpdatedTime string             `bson:"updatedTime,omitempty"`
}

// toDomainUser converts the userDB type to domain.User type
func (u *userMongo) toDomainUser() domain.User {
	user := new(domain.User)
	user.SetID(u.ID.Hex())
	user.SetName(u.Name)
	user.SetLastname(u.Lastname)
	user.SetPassword(u.Password)
	user.SetToken(u.Token)
	user.SetCreatedTime(u.CreatedTime)
	user.SetUpdatedTime(u.UpdatedTime)

	return *user
}

// toDomainUsers convert a slice of usersDB to domain user slice
func toDomainUsers(usersMg []userMongo) ([]domain.User, error) {
	users := make([]domain.User, 0)
	for _, u := range usersMg {
		users = append(users, u.toDomainUser())
	}
	return users, nil
}

func toMongoUser(user domain.User) (userMongo, error) {
	var oid primitive.ObjectID
	if user.ID() != "" {
		var err error
		oid, err = primitive.ObjectIDFromHex(user.ID())
		if err != nil {
			return userMongo{}, err
		}
	}
	return userMongo{
		ID:       oid,
		Name:     user.Name(),
		Lastname: user.Lastname(),
		Password: user.Password(),
		Token:    user.Token(),
	}, nil
}
