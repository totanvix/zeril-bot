package hook

import (
	"encoding/json"
	"net/http"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/structs"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var data structs.HookData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	bot := bot.NewBot(data)
	res := make(map[string]string)

	err = bot.ResolveHook()
	if err != nil {
		res["status"] = "ERROR"
		res["code"] = "internal_error"
		res["message"] = err.Error()
		Response(w, res, http.StatusInternalServerError)
		return
	}

	res["status"] = "OK"
	Response(w, res, http.StatusOK)
}

func Response(w http.ResponseWriter, res map[string]string, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	mRes, _ := json.Marshal(res)
	w.Write(mRes)
}
