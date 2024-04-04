package repository

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/sixojke/internal/config"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
}

type Deps struct {
	Redis  *redis.Client
	Config *config.Config
}

type Repository struct {
	Cache Cache
}

func NewRepository(deps *Deps) *Repository {
	return &Repository{
		Cache: NewCacheRedis(deps.Redis),
	}
}
