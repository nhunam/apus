package cache

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gitlab.com/apus-backend/base-service/config"
	"time"
)

// deprecated
type Cache struct {
	RedisCache *cache.Cache
}

func NewCache(c config.Config) (*Cache, error) {
	return &Cache{RedisCache: setupRedis(c)}, nil
}

func setupRedis(c config.Config) *cache.Cache {
	ttl, _ := time.ParseDuration(c.RedisCache.Ttl)
	redisCache := cache.New(&cache.Options{
		Redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // use default Addr
			Password: "",               // no password set
			DB:       0,                // use default DB
		}),
		LocalCache: cache.NewTinyLFU(c.RedisCache.Size, ttl),
	})

	return redisCache
}
