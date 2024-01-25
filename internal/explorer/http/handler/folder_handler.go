package handler

import (
	"cloud/internal/explorer/service"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type (
	FolderHandler struct {
		log       *slog.Logger
		validator *validator.Validate

		fs service.FolderServiceInterface
	}

	FolderHandlerInterface interface {
	}
)

func NewFolderHandler(
	log *slog.Logger,
	fs service.FolderServiceInterface,
) *FolderHandler {
	return &FolderHandler{
		log:       log,
		validator: validator.New(),
		fs:        fs,
	}
}
