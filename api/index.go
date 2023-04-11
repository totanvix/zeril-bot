package api

import (
	"net/http"
	"zeril-bot/api/hook"
	"zeril-bot/api/url"

	"github.com/go-chi/chi/v5"
)

func Handler(wri http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/api/hook", hook.Handler)
	r.Get("/url", url.Handler)
	r.ServeHTTP(wri, req)
}
