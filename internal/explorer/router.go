package explorer

import (
	"cloud/internal/explorer/http/handler"
	md "cloud/internal/explorer/http/middleware"
	"fmt"
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

	public := r
	private := r

	public.Route("/public/v1/explorer", func(r chi.Router) {
		r.Route("/files", func(r chi.Router) {
			r.Get("/{uuid}", handlers.File.Show)
			r.Get("/{uuid}/preview", handlers.File.Preview)
		})

		r.Route("/share", func(r chi.Router) {
			r.Get("/{uuid}", handlers.Share.Show)
			r.Get("/{uuid}/{pin}", handlers.Share.Show)
			r.Get("/{uuid}/download", handlers.File.Download)
		})
	})

	fmt.Printf("\n\n")
	log.Info("Public Routes:")
	chi.Walk(public, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
	fmt.Printf("\n\n")

	private.Route("/api/v1/explorer", func(r chi.Router) {
		r.Use(md.Auth.New())

		r.Route("/files", func(r chi.Router) {
			r.Get("/", handlers.File.List)
			r.Post("/prepare", handlers.File.Prepare)

			r.Route("/{uuid}", func(r chi.Router) {
				r.Patch("/", handlers.File.Update)
				r.Delete("/", handlers.File.Delete)
				r.Get("/download", handlers.File.Download)
				r.Get("/data", handlers.File.Data)

				r.Group(func(rfs chi.Router) {
					rfs.Use(md.MaxBytesReader.New(1024 * 1024 * 2))
					rfs.Post("/upload", handlers.File.Upload)
				})
			})
		})

		r.Route("/share", func(r chi.Router) {
			r.Post("/create", handlers.Share.Create)
			r.Get("/{file}/data", handlers.Share.Data)
			r.Route("/{uuid}", func(ruf chi.Router) {
				ruf.Put("/", handlers.Share.Update)
				ruf.Delete("/", handlers.Share.Delete)
			})
		})

		r.Route("/folders", func(r chi.Router) {
			r.Route("/", func(ru chi.Router) {
			})
		})
	})

	log.Info("Private Routes:")
	chi.Walk(private, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	fmt.Printf("\n\n")

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
