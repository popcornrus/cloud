package middleware

import (
	"cloud/external/response"
	"cloud/grpc/users"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"net/http"
)

type AuthMiddleware struct {
	log *slog.Logger
	ch  *cache.Cache
}

func NewAuthMiddleware(
	log *slog.Logger,
	cache *cache.Cache,
) *AuthMiddleware {
	return &AuthMiddleware{
		log: log,
		ch:  cache,
	}
}

func (am *AuthMiddleware) New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "http.middleware.AuthMiddleware.New"

			log := am.log.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			token := r.Header.Get("Authorization")
			if len(token) == 0 {
				log.Warn("token wasn't provided")

				response.Respond(w, response.Response{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized",
				})

				return
			}

			user, err := users.Authorize(log, token)
			if err != nil {
				log.Warn("user wasn't found")

				response.Respond(w, response.Response{
					Status:  http.StatusUnauthorized,
					Message: "Unauthorized",
				})

				return
			}

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
