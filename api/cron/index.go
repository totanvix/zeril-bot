package cron

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/request"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]string)
	if token := r.Header.Get("Authorization"); token != "Bearer "+os.Getenv("CRON_SECRET") {
		res["status"] = "ERROR"
		res["message"] = "Token invalid"
		request.Response(w, res, http.StatusOK)
		return
	}

	telegram := telegram.New(&http.Client{}, telegram.BASE_URL)

	chatId, _ := strconv.Atoi(os.Getenv("BOT_OWNER_CHAT_ID"))

	data := structs.HookData{

		Message: structs.Message{
			Text: "/quote",
			From: structs.From{
				ID:        chatId,
				FirstName: "Cron",
			}, Chat: structs.Chat{
				ID:        chatId,
				FirstName: "Cron",
			},
		},
	}

	fmt.Println(data)
	bot := bot.NewBot(telegram, data)
	err := bot.ResolveHook()

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
