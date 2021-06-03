package main

import (
	"apuscorp.com/core/auth-service/api"
	"apuscorp.com/core/auth-service/cache"
	"apuscorp.com/core/auth-service/database"
	"apuscorp.com/core/auth-service/i18n"
	appLog "apuscorp.com/core/auth-service/log"
	"apuscorp.com/core/auth-service/redis"
	"apuscorp.com/core/auth-service/router"
	"apuscorp.com/core/auth-service/util"
	"context"
	"fmt"
	"gitlab.com/apus-backend/base-service/config"
	"gitlab.com/apus-backend/base-service/rabbitmq"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type App struct {
	Config    config.Config
	AppLog    *appLog.AppLog
	Database  *database.Database
	Redis     *redis.Redis
	Cache     *cache.Cache
	Router    *router.Router
	RabbitMQ  *rabbitmq.RabbitMQ
	I18n      *i18n.I18n
	AuthV1Api *api.AuthV1Api
}

func (r App) Start() error {
	// setup routers
	r.SetupRouters()

	// run Gin engine
	util.CheckError(r.Router.Engine.Run(fmt.Sprintf(":%s", r.Config.Server.Port)))

	gracefulShutdown(&http.Server{
		Addr:    fmt.Sprintf(":%s", r.Config.Server.Port),
		Handler: r.Router.Engine,
	})

	return nil
}

func (r App) Stop() {
	if err := r.Database.Close(); err != nil {
		panic(err)
	}
	if err := r.Redis.Close(); err != nil {
		panic(err)
	}
	if err := r.RabbitMQ.Close(); err != nil {
		panic(err)
	}
}

func (r App) SetupRouters() {
	// auth group V1
	authGroupV1 := r.Router.Engine.Group("/api/v1/auth")
	{
		authGroupV1.POST("/login", r.AuthV1Api.Login)
	}
	authGroupV1.Use(r.Router.AuthMiddleware.MiddlewareFunc())
	{
		authGroupV1.POST("/refresh_token", r.Router.AuthMiddleware.RefreshHandler)
	}

	// init swagger
	r.Router.InitSwagger(r.Config)
}

// @title Auth Service API
// @version 1.0
// @description This is Auth Service API.

// @contact.name Mai Tien Hai
// @contact.email maihai86@gmail.com

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// wire
	app, err := InitApp()
	util.CheckError(err)

	err = app.Start()
	util.CheckError(err)

	defer app.Stop()

	log.Info().Msg("App started")
}

func gracefulShutdown(srv *http.Server) {
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
