package bootstrap

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	auth_domain "github.com/krlspj/go-jwt-auth/internal/auth/domain"
	auth_handler "github.com/krlspj/go-jwt-auth/internal/auth/platform/server/handler"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/storage/inmemory"
	auth_mongo "github.com/krlspj/go-jwt-auth/internal/auth/platform/storage/mongodb"
	"github.com/krlspj/go-jwt-auth/internal/auth/service"
	authz_domain "github.com/krlspj/go-jwt-auth/internal/authz/domain"
	authz_handler "github.com/krlspj/go-jwt-auth/internal/authz/platform/server/handler"
	authz_mongo "github.com/krlspj/go-jwt-auth/internal/authz/platform/storage/mongodb"
	authz_uc "github.com/krlspj/go-jwt-auth/internal/authz/usecase"
	"github.com/krlspj/go-jwt-auth/internal/config"
	"github.com/krlspj/go-jwt-auth/internal/dbdriver"
	jwt_service "github.com/krlspj/go-jwt-auth/internal/jwt/service"
	mongo_config_service "github.com/krlspj/go-jwt-auth/internal/mongoconfig/service"
	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/storage/mongodb"
	"github.com/krlspj/go-jwt-auth/internal/server"
)

const (
	// TODO -> set as enviroment variables (or flags)
	mongoDatabase         = "jwt_test"
	mongoUserCollection   = "users"
	mongoConfigCollection = "config"
	dbtype                = "mongo" // "mongo" | "inmemory"
	hmacSampleSecret      = "myTokenSecret"
)

func Run() error {

	// Configuration
	app := config.NewAppConfig()
	app.TokenLifetime = 30
	app.Database = mongoDatabase
	app.ConfigCollection = mongoConfigCollection
	app.UsersCollection = mongoUserCollection
	app.Collections = []string{app.UsersCollection, app.ConfigCollection}

	// check connection with database, if error -> use inmemory database
	// uncoment lines to connect to mongodb if available
	var (
		authUserRepo  auth_domain.UserRepo
		authzUserRepo authz_domain.AuthzUserRepo
	)
	const DB_TYPE = dbtype

	if DB_TYPE == "inmemory" {
		log.Println("\033[00;32m[NOTICE] Using inmemory database\033[0m")
		authUserRepo = inmemory.NewUserRepositoryStub()

	} else if DB_TYPE == "mongo" {
		dbClient, err := dbdriver.ConnectMongo()
		if err != nil {
			// use inmemory database
			log.Printf("\033[00;33m[WARNING] error connecting to external database: %s\033[0m\n", err.Error())
			log.Println("\033[00;32m[NOTICE] Using inmemory database\033[0m")
			authUserRepo = inmemory.NewUserRepositoryStub()
		}
		log.Println("\033[00;32m[NOTICE] Using mongo database\033[0m")
		// config Database
		configMsg, err := ConfigDatabase(app, dbClient)
		// Setup mongo collections
		//err = dbdriver.CreateUserCollection(dbClient.MONGO)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("\033[00;32m[NOTICE]", configMsg, "\033[0m")

		authUserRepo = auth_mongo.NewUserRepositoryMongo(dbClient.MONGO, mongoDatabase, mongoUserCollection)
		authzUserRepo = authz_mongo.NewAuthzUserRepositoryMongo(dbClient.MONGO, mongoDatabase, mongoUserCollection)

	} else {
		log.Println("\033[00;32m[NOTICE] Using inmemory database\033[0m")
		authUserRepo = inmemory.NewUserRepositoryStub()
	}

	// Usecases / services
	jwtService := jwt_service.NewJwtUsecase(hmacSampleSecret)
	authService := service.NewAuthService(authUserRepo)
	authzUsecase := authz_uc.NewAuthzUsecase(jwtService, authzUserRepo)

	// Handler / controller
	ah := auth_handler.NewAuthHanlderRepo(app, authService, jwtService, validator.New())
	az := authz_handler.NewAuthz(app, authzUsecase)

	ctx := context.TODO()
	s := server.NewServer(ctx, "localhost", 60002, ah, az)
	serverType := "native" // "gin" | "native"

	return s.Run(ctx, serverType)
}

func ConfigDatabase(app *config.AppConfig, client *dbdriver.DB) (string, error) {
	mongoConfigRepo := mongodb.NewMongoConfigRepo(app, client.MONGO)
	mongoSetupService := mongo_config_service.NewMongoConfigService(app, mongoConfigRepo)
	mongoSetupService.GetDatabases()
	mongoSetupService.CreateConfig()
	mongoSetupService.GetDatabases()

	return "database already configured", nil
}
