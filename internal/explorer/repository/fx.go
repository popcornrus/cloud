package repository

import (
	"go.uber.org/fx"
)

func NewRepository() fx.Option {
	return fx.Module(
		"repository",
		fx.Provide(
			fx.Annotate(
				NewFileRepository,
				fx.As(new(FileRepositoryInterface)),
			),
			fx.Annotate(
				NewFolderRepository,
				fx.As(new(FolderRepositoryInterface)),
			),
			fx.Annotate(
				NewShareRepository,
				fx.As(new(ShareRepositoryInterface)),
			),
		),
	)
}
