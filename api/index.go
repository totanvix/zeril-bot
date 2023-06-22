package api

import (
	"net/http"
	"zeril-bot/api/hook"
	"zeril-bot/api/trip"
	"zeril-bot/api/url"
	"zeril-bot/utils/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
)

func Handler(wri http.ResponseWriter, req *http.Request) {

	r := chi.NewRouter()
	r.Use(chiMiddle.Logger)
	r.Use(middleware.PreRequest)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("/api/hook", hook.Handler)
	r.Get("/url", url.Handler)
	r.Get("/trip", trip.Handler)
	r.Post("/trip", trip.Handler)

	r.ServeHTTP(wri, req)
}
