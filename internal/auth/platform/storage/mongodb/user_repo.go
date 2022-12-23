package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrBadMethod error = errors.New("Bad method used, use FindOne instead")
)

const (
	dbname   = "jwt_test"
	userColl = "users"
)

type UserRepository struct {
	cli        *mongo.Client
	database   mongo.Database
	collection mongo.Collection
}

func NewUserRepositoryMongo(client *mongo.Client, dbName, collName string) *UserRepository {
	db := client.Database(dbName)
	return &UserRepository{
		cli:        client,
		database:   *db,
		collection: *db.Collection(collName),
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

	var allusers []userMongo
	err = cur.All(ctx, &allusers)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return nil, err
	}

	return toDomainUsers(allusers)
}

func (m *UserRepository) FindOne(ctx context.Context, id string) (domain.User, error) {
	var (
		user userMongo
		err  error
	)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	err = m.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}
	return user.toDomainUser(), nil
}

func (m *UserRepository) FindOneByField(ctx context.Context, fieldName, fieldValue string) (domain.User, error) {
	var user userMongo
	if fieldName == "_id" {
		return domain.User{}, ErrBadMethod
	}

	err := m.collection.FindOne(ctx, bson.M{fieldName: fieldValue}).Decode(&user)
	if err != nil {
		log.Println("error on fetching user:", err)
		return domain.User{}, err

	}
	return user.toDomainUser(), nil
}

func (s *UserRepository) InsertUser(ctx context.Context, user domain.User) error {
	userMg, err := toMongoUser(user)
	if err != nil {
		return err
	}
	result, err := s.collection.InsertOne(ctx, userMg)
	if err != nil {
		return err
	}
	//uid := result.InsertedID.(string)
	uid := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Println("Inserted id:", uid)
	return nil
	//return nil
}
