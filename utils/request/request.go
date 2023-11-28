package request

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, res map[string]string, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	mRes, _ := json.Marshal(res)
	w.Write(mRes)
}
