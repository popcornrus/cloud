package service

import (
	"go.uber.org/fx"
)

func NewService() fx.Option {
	return fx.Module(
		"service",
		fx.Provide(
			fx.Annotate(
				NewFileService,
				fx.As(new(FileServiceInterface)),
			),
			fx.Annotate(
				NewFolderService,
				fx.As(new(FolderServiceInterface)),
			),
			fx.Annotate(
				NewShareService,
				fx.As(new(ShareServiceInterface)),
			),
		),
	)
}
