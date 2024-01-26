package middleware

import (
	"go.uber.org/fx"
)

type Middleware struct {
	Auth           *AuthMiddleware
	MaxBytesReader *MaxBytesReaderMiddleware
}

func NewMiddlewares(
	auth *AuthMiddleware,
	maxBytesReader *MaxBytesReaderMiddleware,
) *Middleware {
	return &Middleware{
		Auth:           auth,
		MaxBytesReader: maxBytesReader,
	}
}

func NewMiddleware() fx.Option {
	return fx.Module(
		"middleware",
		fx.Provide(
			NewAuthMiddleware,
			NewMaxBytesReaderMiddleware,
			NewMiddlewares,
		),
	)
}
