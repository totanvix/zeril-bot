package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"zeril-bot/api/hook"
	"zeril-bot/api/url"
	"zeril-bot/utils/structs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	chiMiddle "github.com/go-chi/chi/v5/middleware"
)

func Handler(wri http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	r.Use(chiMiddle.Logger)
	r.Use(PreRequest)
	r.Use(Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Post("/api/hook", hook.Handler)
	r.Get("/url", url.Handler)
	r.ServeHTTP(wri, req)
}

func PreRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var data structs.HookData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			log.Panic(err.Error())
		}

		if data.Message.Text == "" && data.CallbackQuery.Data == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			resp := make(map[string]string)
			resp["status"] = "ERROR"
			resp["message"] = "No message found"

			jsonResp, _ := json.Marshal(resp)
			w.Write(jsonResp)
			return
		}

		if data.Message.Chat.Type == "group" {
			data.Message.Chat.FirstName = data.Message.Chat.Title
		}

		ctx := context.WithValue(r.Context(), "data", data)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			resp := make(map[string]string)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if rvr := recover(); rvr != nil {
				resp["status"] = "ERROR"
				resp["message"] = rvr.(string)

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

			} else {
				resp["status"] = "OK"
			}
			jsonResp, _ := json.Marshal(resp)
			w.Write(jsonResp)
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
