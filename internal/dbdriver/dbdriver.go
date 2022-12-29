package dbdriver

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	MONGO *mongo.Client
	//	MDB   *mongo.Database
}

func ConnectMongo() (*DB, error) {
	log.Println("[NOTICE] [dbdriver] Connecting to MongoDB...")

	mongodbURI := "mongodb://root:example@192.168.1.35:27017"
	//mongodbURI := os.Getenv("MONGO_URI")
	//mongodbURI := "mongodb://" + os.Getenv("MONGO_IP") + ":27017"
	//mongodbURI := "mongodb://" + os.Getenv("MONGO_USER") +
	//	":" + os.Getenv("MONGO_PWD") + "@" +
	//	os.Getenv("MONGO_IP") + ":27017"

	fmt.Println("-> mongodbURI", mongodbURI)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		return nil, err
		//log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
		//log.Fatal(err)
	}

	// do a ping to ensure connection..
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
		//panic(err)
	}

	log.Println("\033[0;32m[NOTICE] [DATABASE] Successfully connected\033[0m")
	return &DB{
		MONGO: client,
	}, nil
}

func OpenCollection(client *mongo.Client, collectionName, databaseName string) *mongo.Collection {
	//var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	//databaseName := os.Getenv("DATABASE_NAME")
	return client.Database(databaseName).Collection(collectionName)
}

func CreateUserCollection(cli *mongo.Client) error {
	// err = m.DB.Database("uskytest").CreateCollection(ctx, user_col, &userCollOpts)
	//collNames, err := cli.Database("jwt_test").ListCollectionNames(context.Background(), bson.M{})
	// err := Ocli.Database("jwt_test").Drop(context.Background())
	//if err != nil {
	//	return err
	//}
	//fmt.Println("collNames", collNames)

	err := cli.Database("jwt_test").Collection("users").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}

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

	err = cli.Database("jwt_test").CreateCollection(context.Background(), "users", &userCollOpts)
	if err != nil {
		return err
	}

	// Create collection index
	//indexModel := mongo.IndexModel{
	//	Keys:    bson.D{{Key: "username", Value: 1}},
	//	Options: options.Index().SetUnique(true),
	//}
	//name, err := cli.Database("jwt_test").Collection("users").Indexes().CreateOne(context.TODO(), indexModel)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Name of Index Created: " + name)

	return nil
}
