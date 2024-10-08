package service

import (
	"cloud/internal/root/service/user_service"
	"go.uber.org/fx"
)

func NewService() fx.Option {
	return fx.Module(
		"service",
		fx.Options(
			user_service.NewUser(),
		),
	)
}
