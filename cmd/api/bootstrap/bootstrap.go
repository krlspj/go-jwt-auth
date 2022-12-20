package bootstrap

import (
	"context"
	"log"

	"github.com/krlspj/go-jwt-auth/internal/auth/domain"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/server/handler"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/storage/inmemory"
	"github.com/krlspj/go-jwt-auth/internal/auth/platform/storage/mongodb"
	"github.com/krlspj/go-jwt-auth/internal/auth/service"
	"github.com/krlspj/go-jwt-auth/internal/dbdriver"
	"github.com/krlspj/go-jwt-auth/internal/server"
)

func Run() error {
	// check connection with database, if error -> use inmemory database
	// uncoment lines to connect to mongodb if available
	var authUserRepo domain.UserRepo

	dbClient, err := dbdriver.ConnectMongo()
	if err != nil {
		// use inmemory database
		log.Printf("\033[00;33m[WARNING] error connecting to external database: %s\033[0m\n", err.Error())
		log.Println("\033[00;32m[NOTICE] Using inmemory database\033[0m")
		authUserRepo = inmemory.NewUserRepositoryStub()
	} else {
		log.Println("\033[00;32m[NOTICE] Using external database\033[0m")
		authUserRepo = mongodb.NewUserRepositoryMongo(dbClient.MONGO)
	}

	//log.Println("\033[00;32m[NOTICE] Using inmemory database\033[0m")

	//authUserRepo = inmemory.NewUserRepositoryStub()

	as := service.NewAuthService(authUserRepo)

	ah := handler.NewAuthHanlderRepo(as)

	ctx := context.TODO()
	s := server.NewServer(ctx, "localhost", 60002, ah)

	return s.Run(ctx)
}
