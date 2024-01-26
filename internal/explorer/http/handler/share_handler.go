package handler

import (
	"cloud/internal/explorer/service"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type (
	ShareHandler struct {
		log       *slog.Logger
		validator *validator.Validate

		fs service.ShareServiceInterface
	}

	ShareHandlerInterface interface {
	}
)

func NewShareHandler(
	log *slog.Logger,
	fs service.ShareServiceInterface,
) *ShareHandler {
	return &ShareHandler{
		log:       log,
		validator: validator.New(),
		fs:        fs,
	}
}
