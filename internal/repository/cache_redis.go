package repository

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type CacheRedis struct {
	db *redis.Client
}

func NewCacheRedis(db *redis.Client) *CacheRedis {
	return &CacheRedis{
		db: db,
	}
}

func (r *CacheRedis) Get(key string) (string, error) {
	value, err := r.db.Get(key).Result()
	if err != nil {
		if err != redis.Nil {
			return "", fmt.Errorf("cache get error: %v", err)
		} else {
			return "", nil
		}
	}

	return value, nil
}

func (r *CacheRedis) Set(key string, value string, expiration time.Duration) error {
	if err := r.db.Set(key, value, expiration).Err(); err != nil {
		return fmt.Errorf("cache set error: %v", err)
	}

	return nil
}
