package hook

import (
	"encoding/json"
	"net/http"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var data structs.HookData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	res := make(map[string]string)

	if data.Message.Text == "" {
		res["status"] = "OK"
		res["message"] = "Ignore hook with chat content not found"
		Response(w, res, http.StatusOK)
		return
	}

	telegram := telegram.New(&http.Client{}, telegram.BASE_URL)

	bot := bot.NewBot(telegram, data)

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
