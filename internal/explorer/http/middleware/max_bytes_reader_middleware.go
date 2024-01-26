package middleware

import (
	"net/http"
)

type MaxBytesReaderMiddleware struct {
}

func NewMaxBytesReaderMiddleware() *MaxBytesReaderMiddleware {
	return &MaxBytesReaderMiddleware{}
}

func (mbrm *MaxBytesReaderMiddleware) New(maxBytes int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
