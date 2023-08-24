package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) setHTTPRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	return r
}
