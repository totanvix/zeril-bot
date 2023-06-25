package trip

import (
	"encoding/json"
	"net/http"
	"zeril-bot/utils/redis"
)

type data struct {
	Url string `json:"url"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == http.MethodPost {
		var d data
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		redis.Set("trip", d.Url, 0)
		Response(w, d.Url)

	} else {
		url := redis.Get("trip")
		Response(w, url.Val())
	}
}

func Response(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := make(map[string]string)
	resp["url"] = url

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
	return
}
