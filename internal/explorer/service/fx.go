package service

import (
	"cloud/internal/explorer/service/file_service"
	"go.uber.org/fx"
)

func NewService() fx.Option {
	return fx.Module(
		"service",
		fx.Options(
			file_service.NewFile(),
		),
	)
}
