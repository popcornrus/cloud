package root

import (
	"cloud/external/config"
	"cloud/external/db"
	"cloud/grpc/users"
	"cloud/internal/root/http/handler"
	"cloud/internal/root/http/middleware"
	"cloud/internal/root/repository"
	"cloud/internal/root/service"
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Options(
			repository.NewRepository(),
			service.NewService(),
			handler.NewHandler(),
			middleware.NewMiddleware(),
			db.NewDataBase(),
		),
		fx.Provide(
			config.NewConfig,
			validator.New,
			NewCache,
			NewLogger,
			NewRouter,
			users.NewUserGRPC,
			NewServer,
		),
		fx.Invoke(users.RunUserGRPCServer),
		fx.Invoke(RunServer),
	)
}

func NewCache() *cache.Cache {
	return cache.New(cache.NoExpiration, cache.NoExpiration)
}
