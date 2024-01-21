package root

import (
	"cloud/internal/root/http/handler"
	md "cloud/internal/root/http/middleware"
	"cloud/internal/root/http/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"time"
)

func NewRouter(
	log *slog.Logger,
	handlers *handler.Handlers,
	md *md.Middleware,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(logger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with your allowed origin(s)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(ri chi.Router) {
		ri.Route("/users", func(ru chi.Router) {
			ru.Post("/sign-up", handlers.User.SignUp)
			ru.Post("/sign-in", handlers.User.SignIn)

			ru.Route("/me", func(rup chi.Router) {
				rup.Use(md.AuthMiddleware.New())

				rup.Get("/", handlers.User.Get)
				rup.Put("/update", handlers.User.Update)
			})
		})
	})

	return r
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
