package storage

import (
	"context"

	"github.com/krlspj/go-jwt-auth/internal/mongoconfig/domain"
)

type StorageConfig interface {
	GetDatabases(ctx context.Context) ([]string, error)
	GetCollections(ctx context.Context, dbName string) ([]string, error)
	FindConfig(ctx context.Context) (domain.Config, error)
	CreateDBConfig(ctx context.Context, config domain.Config) (string, error)
	UpdateDBConfig(ctx context.Context, config domain.Config) error
}
