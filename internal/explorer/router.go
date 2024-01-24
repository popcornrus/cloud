package explorer

import (
	"cloud/internal/explorer/http/handler"
	md "cloud/internal/explorer/http/middleware"
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
	//r.Use(logger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(corsHandler())
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with your allowed origin(s)
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Content-Range", "Range"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1/explorer", func(ri chi.Router) {
		ri.Get("/files/{uuid}", handlers.File.Show)
		ri.Get("/files/{uuid}/preview", handlers.File.Preview)

		ri.Route("/files", func(ru chi.Router) {
			ru.Use(md.Auth.New())
			ru.Get("/", handlers.File.List)
			ru.Post("/prepare", handlers.File.Prepare)

			ru.Route("/{uuid}", func(ruf chi.Router) {
				ruf.Patch("/", handlers.File.Update)
				ruf.Delete("/", handlers.File.Delete)
				ruf.Get("/download", handlers.File.Download)
				ruf.Get("/data", handlers.File.Data)
				ruf.Post("/upload", handlers.File.Upload)
			})
		})
	})

	return r
}

func corsHandler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Range, Range")
				w.WriteHeader(http.StatusOK)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
