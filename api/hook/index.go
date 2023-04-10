package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"zeril-bot/utils/bitcoin"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/kqxs"
	"zeril-bot/utils/lunar"
	"zeril-bot/utils/qr"
	"zeril-bot/utils/quote"
	"zeril-bot/utils/random"
	"zeril-bot/utils/shortener"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/weather"
)

func Router(w http.ResponseWriter, r *http.Request) {
	var data structs.HookData
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		log.Fatalln(err.Error())
	}

	if data.Message.Text == "" && data.CallbackQuery.Data == "" {
		log.Println("No message found")
		w.Header().Set("status", "200")
		fmt.Fprintln(w, "No message found")
		return
	}

	if data.Message.Chat.Type == "group" {
		data.Message.Chat.FirstName = data.Message.Chat.Title
	}

	if data.CallbackQuery.Data != "" {
		ResolveCallback(data)
		return
	}

	ResolveCommand(data)
}

func ResolveCommand(data structs.HookData) {
	name := data.Message.Chat.FirstName
	chatId := data.Message.Chat.ID
	text := data.Message.Text
	arr := strings.Fields(text)

	bot.SetTypingAction(chatId)

	log.Println(fmt.Sprintf("Yêu cầu từ bạn %s: %s", name, text))

	command := arr[0]

	switch command {
	case "/start", "/start@zerill_bot":
		bot.SendStartMessage(chatId, name)
	case "/help", "/help@zerill_bot":
		bot.SendHelpMessage(chatId)
	case "/groupid", "/groupid@zerill_bot":
		bot.SendGroupId(chatId, data.Message.Chat.Type)
	case "/quote", "/quote@zerill_bot":
		quote.SendAQuote(chatId)
	case "/lunar", "/lunar@zerill_bot":
		lunar.SendLunarDateNow(chatId)
	case "/weather", "/weather@zerill_bot":
		weather.SendForecastOfWeather(chatId, text)
	case "/bitcoin", "/bitcoin@zerill_bot":
		bitcoin.SendBitcoinPrice(chatId)
	case "/qr", "/qr@zerill_bot":
		qr.SendQRImage(chatId, text)
	case "/random", "/random@zerill_bot":
		random.RandomElements(chatId, text)
	case "/kqxs", "/kqxs@zerill_bot":
		kqxs.Send(chatId, text)
	case "/shorturl", "/shorturl@zerill_bot":
		shortener.Do(chatId, text)
	default:
		bot.SendMessage(chatId, "Tôi không hiểu câu lệnh của bạn !!!")
	}
}

func ResolveCallback(callback structs.HookData) {
	name := callback.CallbackQuery.Message.Chat.FirstName
	chatId := callback.CallbackQuery.Message.Chat.ID
	text := callback.CallbackQuery.Message.Text
	data := callback.CallbackQuery.Data

	bot.SetTypingAction(chatId)

	log.Println(fmt.Sprintf("Yêu cầu từ bạn %s: %s, callback data: %s", name, text, data))

	arr := strings.Fields(data)
	command := arr[0]

	switch command {
	case "/weather":
		weather.SendForecastOfWeather(chatId, data)
	case "/kqxs":
		kqxs.Send(chatId, data)
	}
}
