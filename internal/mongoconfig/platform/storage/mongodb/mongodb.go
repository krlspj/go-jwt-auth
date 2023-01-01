package mongodb

import (
	"context"
	"fmt"

	"github.com/krlspj/go-jwt-auth/internal/config"
	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfigRepo struct {
	app *config.AppConfig
	cli *mongo.Client
}

func NewMongoConfigRepo(a *config.AppConfig, client *mongo.Client) *MongoConfigRepo {
	return &MongoConfigRepo{
		app: a,
		cli: client,
	}
}

//collNames, err := cli.Database("jwt_test").ListCollectionNames(context.Background(), bson.M{})

func (m *MongoConfigRepo) GetDatabases(ctx context.Context) ([]string, error) {
	dbList, err := m.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return dbList, nil
}

func (m *MongoConfigRepo) GetCollections(ctx context.Context, dbName string) ([]string, error) {
	//filter := bson.D{{Key: "name", Value: "users"}}
	//filter := bson.D{{Key: "$in", Value: bson.A{bson.D{{Key: "name", Value: "users"}}}}}
	collections := bson.A{}
	for _, coll := range m.app.Collections {
		collections = append(collections, coll)
	}
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$in", Value: collections}}}}
	//filter := bson.M{"name": bson.M{"$in": collections}}
	collList, err := m.cli.Database(dbName).ListCollectionNames(ctx, filter)
	if err != nil {
		return nil, err
	}
	fmt.Println("collList", collList)
	return collList, nil
}

func (m *MongoConfigRepo) FindConfig(ctx context.Context) (domain.Config, error) {
	configColl := m.cli.Database(m.app.DatabaseName).Collection(m.app.ConfigCollection)
	filter := bson.M{}
	var mgConfig mongoConfig
	err := configColl.FindOne(ctx, filter).Decode(&mgConfig)
	if err != nil {
		return domain.Config{}, err
	}

	return mgConfig.toDomain()
}

func (m *MongoConfigRepo) CreateDBConfig(ctx context.Context, config domain.Config) (string, error) {
	coll := m.cli.Database(m.app.DatabaseName).Collection(m.app.ConfigCollection)
	mgConfig, err := toMongoConfig(config)
	if err != nil {
		return "", err
	}
	result, err := coll.InsertOne(ctx, mgConfig)
	if err != nil {
		return "", err
	}
	fmt.Println("Restult:", result)
	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (m *MongoConfigRepo) UpdateDBConfig(ctx context.Context, config domain.Config) error {
	coll := m.cli.Database(m.app.DatabaseName).Collection(m.app.ConfigCollection)
	mgConfig, err := toMongoConfig(config)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: mgConfig.ID}}
	result, err := coll.UpdateOne(ctx, filter, mgConfig)
	if err != nil {
		return err
	}
	fmt.Println("Restult:", result)

	return nil
}

func (m *MongoConfigRepo) CreateUserCollection(ctx context.Context) error {
	userCollOpts := options.CreateCollectionOptions{
		Validator: bson.M{
			"$jsonSchema": bson.M{
				"bsonType":    "object",
				"description": "user definition",
				"required":    bson.A{"username", "password"},
				"properties": bson.M{
					"username": bson.M{
						"bsonType":    "string",
						"description": "Ability name",
						"maxLength":   30,
					},
					"password": bson.M{
						"bsonType":    "string",
						"description": "password is required",
					},
					"email": bson.M{
						"bsonType":    "string",
						"description": "user's email",
						"minLength":   6,
						"maxLength":   127,
					},
					"roleId": bson.M{
						"bsonType": "objectId",
					},
					"refreshPassword": bson.M{
						"enum":        bson.A{"true", "false"},
						"description": "this flag forces user to renew its password",
					},
				},
			},
		},
	}

	err := m.createCollection(ctx, m.app.UsersCollection, &userCollOpts)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoConfigRepo) createCollection(ctx context.Context, collectionName string, opts *options.CreateCollectionOptions) error {
	err := m.cli.Database(m.app.DatabaseName).CreateCollection(ctx, collectionName, opts)
	if err != nil {
		return err
	}

	return nil
}
