//+build wireinject

package main

import (
	"apuscorp.com/core/auth-service/api"
	"apuscorp.com/core/auth-service/cache"
	"apuscorp.com/core/auth-service/client"
	"apuscorp.com/core/auth-service/dao"
	"apuscorp.com/core/auth-service/database"
	"apuscorp.com/core/auth-service/i18n"
	"apuscorp.com/core/auth-service/log"
	"apuscorp.com/core/auth-service/redis"
	"apuscorp.com/core/auth-service/router"
	"github.com/google/wire"
	"gitlab.com/apus-backend/base-service/config"
	"gitlab.com/apus-backend/base-service/rabbitmq"
)

func InitApp() (App, error) {
	panic(wire.Build(config.LoadConfig, log.NewAppLog,
		database.NewDatabase, cache.NewCache, router.NewRouter, i18n.NewI18n,
		client.NewClient, redis.NewRedis, rabbitmq.NewRabbitMQ,
		wire.Struct(new(dao.UserCredentialsDao), "*"),
		wire.Struct(new(api.AuthV1Api), "*"),
		wire.Struct(new(App), "*")))
	return App{}, nil
}
