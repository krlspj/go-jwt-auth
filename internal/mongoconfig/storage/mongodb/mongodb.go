package mongodb

import (
	"context"
	"fmt"

	"github.com/krlspj/go-jwt-auth/internal/config"
	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoConfig struct {
	app *config.AppConfig
	cli *mongo.Client
}

func NewMongoConfigRepo(a *config.AppConfig, client *mongo.Client) *MongoConfig {
	return &MongoConfig{
		app: a,
		cli: client,
	}
}

//collNames, err := cli.Database("jwt_test").ListCollectionNames(context.Background(), bson.M{})

func (m *MongoConfig) GetDatabases(ctx context.Context) ([]string, error) {
	dbList, err := m.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return dbList, nil
}

func (m *MongoConfig) GetCollections(ctx context.Context, dbName string) ([]string, error) {
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

func (m *MongoConfig) FindConfig(ctx context.Context) (domain.Config, error) {
	configColl := m.cli.Database(m.app.Database).Collection(m.app.ConfigCollection)
	filter := bson.M{}
	var mgConfig mongoConfig
	err := configColl.FindOne(ctx, filter).Decode(&mgConfig)
	if err != nil {
		return domain.Config{}, err
	}
	return mgConfig.toDomain(), nil
}

func (m *MongoConfig) CreateConfig(ctx context.Context, config domain.Config) (string, error) {
	coll := m.cli.Database(m.app.Database).Collection(m.app.ConfigCollection)
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
