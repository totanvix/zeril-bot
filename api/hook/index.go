package hook

import (
	"encoding/json"
	"net/http"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/request"
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

	if data.CallbackQuery.Data == "" && data.Message.Text == "" {
		res["status"] = "ERROR"
		res["message"] = "Ignore hook with chat content not found"
		request.Response(w, res, http.StatusOK)
		return
	}

	telegram := telegram.New(&http.Client{}, telegram.BASE_URL)

	bot := bot.NewBot(telegram, data)

	err = bot.ResolveHook()

	if err != nil {
		res["status"] = "ERROR"
		res["code"] = "internal_error"
		res["message"] = err.Error()
		request.Response(w, res, http.StatusInternalServerError)
		return
	}

	res["status"] = "OK"
	request.Response(w, res, http.StatusOK)
}
