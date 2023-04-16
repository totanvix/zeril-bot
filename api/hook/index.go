package hook

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"zeril-bot/utils/bitcoin"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/channel"
	_ "zeril-bot/utils/channel"
	"zeril-bot/utils/help"
	"zeril-bot/utils/kqxs"
	"zeril-bot/utils/lunar"
	"zeril-bot/utils/qr"
	"zeril-bot/utils/quote"
	"zeril-bot/utils/random"
	"zeril-bot/utils/shortener"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/weather"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("data").(structs.HookData)

	channel.GetWg().Add(2)

	if data.CallbackQuery.Data != "" {
		resolveCallback(data)
	} else {
		resolveCommand(data)
	}

	channel.GetWg().Wait()
}

func resolveCommand(data structs.HookData) {
	name := data.Message.From.FirstName
	chatId := data.Message.Chat.ID
	text := data.Message.Text
	arr := strings.Fields(text)

	setBotIsTyping(chatId)

	log.Println(fmt.Sprintf("Yêu cầu từ bạn %s: %s", name, text))

	command := arr[0]

	switch command {
	case "/start", "/start@zerill_bot":
		help.SendStartMessage(chatId, name)
	case "/help", "/help@zerill_bot":
		help.SendHelpMessage(chatId)
	case "/groupid", "/groupid@zerill_bot":
		help.SendGroupId(chatId, string(data.Message.Chat.Type))
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
	case "/shortener", "/shortener@zerill_bot":
		shortener.Generate(chatId, text)
	default:
		channel.SendMessage(chatId, "Tôi không hiểu câu lệnh của bạn !!!")
	}
}

func resolveCallback(callback structs.HookData) {
	name := callback.CallbackQuery.Message.From.FirstName
	chatId := callback.CallbackQuery.Message.Chat.ID
	text := callback.CallbackQuery.Message.Text
	data := callback.CallbackQuery.Data

	setBotIsTyping(chatId)

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

func setBotIsTyping(chatId int) {
	go func() {
		bot.SetTypingAction(chatId)
		channel.GetWg().Done()
	}()
}
