package handler

import (
	"go.uber.org/fx"
)

type Handlers struct {
	File   *FileHandler
	Share  *ShareHandler
	Folder *FolderHandler
}

func NewHandlers(
	file *FileHandler,
	share *ShareHandler,
	folder *FolderHandler,
) *Handlers {
	return &Handlers{
		File:   file,
		Share:  share,
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
			NewShareHandler,
			NewHandlers,
		),
	)
}
