package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/channel"
	"zeril-bot/utils/structs"

	"github.com/go-chi/chi/v5/middleware"
)

func PreRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/trip" {
			next.ServeHTTP(w, r)
			return
		}

		var data structs.HookData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			log.Println("Request is not from Telegram")
			next.ServeHTTP(w, r)
			return
		}

		channel.Create()

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

		ctx := context.WithValue(r.Context(), "data", data)

		if data.CallbackQuery.Data != "" {
			bot.SetChatFrom(data.CallbackQuery.From)
			bot.SetChatType(data.CallbackQuery.Message.Chat.Type)
		} else {
			bot.SetChatFrom(data.Message.From)
			bot.SetChatType(data.Message.Chat.Type)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				resp := make(map[string]string)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				resp["status"] = "ERROR"
				resp["message"] = rvr.(string)
				jsonResp, _ := json.Marshal(resp)

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

				w.Write(jsonResp)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func Handle404NotFound() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)

		res := make(map[string]string)
		res["status"] = "ERROR"
		res["error"] = "http.error"
		res["message"] = "route does not exist"

		mRes, _ := json.Marshal(res)
		w.Write(mRes)
	}
}

func Handle405MethodNotAllowed() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)

		res := make(map[string]string)
		res["status"] = "ERROR"
		res["error"] = "http.error"
		res["message"] = "method is not valid"

		mRes, _ := json.Marshal(res)
		w.Write(mRes)
	}
}
