package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Handler(wri http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/api/hook", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hook"))
	})
	r.Get("/url", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("url"))
	})
	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	r.ServeHTTP(wri, req)
}
