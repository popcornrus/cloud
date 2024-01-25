package handler

import (
	"go.uber.org/fx"
)

type Handlers struct {
	File   *FileHandler
	Folder *FolderHandler
}

func NewHandlers(
	file *FileHandler,
	folder *FolderHandler,
) *Handlers {
	return &Handlers{
		File:   file,
		Folder: folder,
	}
}

func NewHandler() fx.Option {
	return fx.Module(
		"handler",
		fx.Options(),
		fx.Provide(
			NewFileHandler,
			NewFolderHandler,
			NewHandlers,
		),
	)
}
