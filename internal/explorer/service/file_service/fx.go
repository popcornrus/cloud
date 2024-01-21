package file_service

import (
	"go.uber.org/fx"
)

func NewFile() fx.Option {
	return fx.Module(
		"user-service",
		fx.Provide(
			fx.Annotate(
				NewFileService,
				fx.As(new(FileServiceInterface)),
			),
		),
	)
}
