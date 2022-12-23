package bootstrap

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/server/handler"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/storage/inmemory"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/storage/mongodb"
	"github.com/krlspj/go-jwt-auth/internal/auth/service"
	"github.com/krlspj/go-jwt-auth/internal/config"
	"github.com/krlspj/go-jwt-auth/internal/dbdriver"
	jwt_uc "github.com/krlspj/go-jwt-auth/internal/jwt/usecase"
	"github.com/krlspj/go-jwt-auth/internal/server"
)

const (
	mongoDatabase       = "jwt_test"
	mongoUserCollection = "users"
	dbtype              = "mongo" // "mongo" | "inmemory"
)

func Run() error {

	// Configuration
	app := config.NewAppConfig()
	app.TokenLifetime = 1

	// check connection with database, if error -> use inmemory database
	// uncoment lines to connect to mongodb if available
	var authUserRepo domain.UserRepo
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
		authUserRepo = mongodb.NewUserRepositoryMongo(dbClient.MONGO, mongoDatabase, mongoUserCollection)

	} else {
		log.Println("\033[00;32m[NOTICE] Using inmemory database\033[0m")
		authUserRepo = inmemory.NewUserRepositoryStub()
	}

	//log.Println("\033[00;32m[NOTICE] Using inmemory database\033[0m")

	//authUserRepo = inmemory.NewUserRepositoryStub()
	hmacSampleSecret := "this is my secret"
	juc := jwt_uc.NewJwtUsecase(hmacSampleSecret)
	as := service.NewAuthService(authUserRepo)

	ah := handler.NewAuthHanlderRepo(app, as, juc, validator.New())

	ctx := context.TODO()
	s := server.NewServer(ctx, "localhost", 60002, ah)
	serverType := "native" // "gin" | "native"

	return s.Run(ctx, serverType)
}
