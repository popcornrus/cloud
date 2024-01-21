package handler

import (
	"go.uber.org/fx"
)

type Handlers struct {
	File *FileHandler
}

func NewHandlers(
	fh *FileHandler,
) *Handlers {
	return &Handlers{
		File: fh,
	}
}

func NewHandler() fx.Option {
	return fx.Module(
		"handler",
		fx.Options(),
		fx.Provide(
			NewFileHandler,
			NewHandlers,
		),
	)
}
