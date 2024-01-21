package middleware

import (
	"go.uber.org/fx"
)

type Middleware struct {
	Auth *AuthMiddleware
}

func NewMiddlewares(
	auth *AuthMiddleware,
) *Middleware {
	return &Middleware{
		Auth: auth,
	}
}

func NewMiddleware() fx.Option {
	return fx.Module(
		"middleware",
		fx.Provide(
			NewAuthMiddleware,
			NewMiddlewares,
		),
	)
}
