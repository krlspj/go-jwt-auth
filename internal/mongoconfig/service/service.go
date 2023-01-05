package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/krlspj/go-jwt-auth/internal/config"
	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/domain"
	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/platform/storage"
)

type DBConfigService interface {
	GetDatabases()
	CreateConfig()
}

type MongoConfigService struct {
	app  *config.AppConfig
	mgdb storage.StorageConfig
}

func NewMongoConfigService(a *config.AppConfig, sc storage.StorageConfig) *MongoConfigService {
	return &MongoConfigService{
		app:  a,
		mgdb: sc,
	}
}

func (s *MongoConfigService) GetDatabases() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.mgdb.GetDatabases(ctx)
	s.mgdb.GetCollections(ctx, "jwt_test")

}

func (s *MongoConfigService) CreateConfig() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	config := new(domain.Config)
	config.SetCreatedAt(time.Now().UnixMilli())
	config.SetRefresh("false")
	config.SetRefreshB(false)
	config.SetID("63b6aa7ef7d6abd165511767")

	fmt.Println("created at time:", config.CreatedAt())

	//id, err := s.mgdb.CreateDBConfig(ctx, *config)
	err := s.mgdb.UpdateDBConfig(ctx, *config)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Config ID", id)
}

func (s *MongoConfigService) DropDatabase() error {

	return nil
}
