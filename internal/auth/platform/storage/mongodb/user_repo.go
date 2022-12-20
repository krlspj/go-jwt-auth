package mongodb

import (
	"context"
	"fmt"

	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbname   = "jwt_test"
	userColl = "users"
)

type UserRepository struct {
	cli *mongo.Client
}

func NewUserRepositoryMongo(client *mongo.Client) *UserRepository {
	return &UserRepository{
		cli: client,
	}
}

func (m *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	coll := m.cli.Database(dbname).Collection(userColl)

	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Error:", err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	var allusers []userDB
	err = cur.All(ctx, &allusers)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return nil, err
	}
	fmt.Println("allusers ->", allusers)

	return toDomainUsers(allusers)
}
