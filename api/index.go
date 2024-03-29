package api

import (
	"net/http"
	"zeril-bot/api/cron"
	"zeril-bot/api/hook"
	"zeril-bot/api/url"
	"zeril-bot/utils/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
)

func Handler(wri http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	r.Use(chiMiddle.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(middleware.Handle404NotFound())
	r.MethodNotAllowed(middleware.Handle405MethodNotAllowed())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi there !"))
	})
	r.Post("/api/hook", hook.Handler)
	r.Get("/api/cron", cron.Handler)
	r.Get("/url", url.Handler)

	r.ServeHTTP(wri, req)
}
