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
	"go.mongodb.org/mongo-driver/mongo/options"
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
}

func (s *UserRepository) ReplaceUser(ctx context.Context, user domain.User) error {
	userMg, err := toMongoUser(user)
	if err != nil {
		return err
	}
	if userMg.ID == (primitive.ObjectID{}) {
		fmt.Println("Set new ID")
		userMg.ID = primitive.NewObjectID()
	}
	filter := bson.D{{Key: "_id", Value: userMg.ID}}
	upsert := true
	options := options.ReplaceOptions{
		Upsert: &upsert,
	}
	_ = options
	result, err := s.collection.ReplaceOne(ctx, filter, userMg, &options)
	if err != nil {
		return err
	}
	fmt.Println("result", result)
	return nil
}

func (s *UserRepository) UpdateUser(ctx context.Context, user domain.User) error {
	userMg, err := toMongoUser(user)
	if err != nil {
		return err
	}
	fmt.Println("--------- userMg", userMg)
	filter := bson.D{{Key: "_id", Value: userMg.ID}}
	doc := bson.D{{Key: "$set", Value: userMg}}
	result, err := s.collection.UpdateOne(ctx, filter, doc)
	if err != nil {
		return err
	}
	fmt.Println("result", result)
	return nil
}
func (s *UserRepository) CountRecords(ctx context.Context, fieldName, fieldValue string) (int64, error) {
	filter := bson.D{{Key: fieldName, Value: fieldValue}}
	return s.collection.CountDocuments(ctx, filter)
}

func (s *UserRepository) CreatIndex() error {
	fmt.Println("============== Create Index")
	indexName, err := s.collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	fmt.Println("------------ index:", indexName)
	if err != nil {
		return err
	}
	return nil
}
