package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"zeril-bot/utils/bitcoin"
	"zeril-bot/utils/lunar"
	"zeril-bot/utils/quote"
)

type TelegramData struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date     int    `json:"date"`
		Text     string `json:"text"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities"`
	} `json:"message"`
}

func Router(w http.ResponseWriter, r *http.Request) {
	var data TelegramData
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		fmt.Println(err.Error())
	}

	chatId := data.Message.Chat.ID
	text := data.Message.Text
	command, _ := getCommandAndArgs(text)

	switch command {
	case "/quote":
		quote.SendAQuote(chatId)
	case "/lunar":
		lunar.SendLunarDateNow(chatId)
	case "/bitcoin":
		bitcoin.SendBitcoinPrice(chatId)
	}
}

func getCommandAndArgs(text string) (string, []string) {
	arr := strings.Fields(text)

	return arr[0], arr[1:]
}
