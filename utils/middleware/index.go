package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5/middleware"
)

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
