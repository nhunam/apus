package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab.com/apus-backend/base-service/config"
	"time"
)

type Redis struct {
	Rdb       *redis.Client
	CommonDur time.Duration
}

func NewRedis(c config.Config) (*Redis, error) {
	commonDur, err := time.ParseDuration(c.Redis.Ttl.Common)
	if err != nil {
		return nil, err
	}
	return &Redis{
		Rdb: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
			Password: c.Redis.Password,
			DB:       c.Redis.Database,
		}),
		CommonDur: commonDur,
	}, nil
}

func (r *Redis) Close() error {
	if err := r.Rdb.Close(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) GetStruct(ctx context.Context, key string, val interface{}) error {
	result, err := r.Rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	if len(result) > 0 {
		err = json.Unmarshal([]byte(result), val)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (r *Redis) SetStruct(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = r.Rdb.Set(ctx, key, string(jsonBytes), exp).Err()
	if err != nil {
		return err
	}
	return nil
}
