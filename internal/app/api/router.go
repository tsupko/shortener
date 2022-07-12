package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(m *RequestHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			_ = chi.URLParam(r, "id")
			m.handleGetRequest(w, r)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			m.handlePostRequest(w, r)
		})
	})
	return r
}
