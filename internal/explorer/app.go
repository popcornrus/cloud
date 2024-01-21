package explorer

import (
	"cloud/external/config"
	"cloud/external/db"
	"cloud/internal/explorer/http/handler"
	"cloud/internal/explorer/http/middleware"
	"cloud/internal/explorer/repository"
	"cloud/internal/explorer/service"
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
			NewServer,
		),
		fx.Invoke(RunServer),
	)
}

func NewCache() *cache.Cache {
	return cache.New(cache.NoExpiration, cache.NoExpiration)
}
