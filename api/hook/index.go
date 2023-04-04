package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"zeril-bot/utils/bitcoin"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/lunar"
	"zeril-bot/utils/qr"
	"zeril-bot/utils/quote"
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
		log.Fatalln("No message found")
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
		sendStartMessage(chatId, name)
	case "/help", "/help@zerill_bot":
		sendHelpMessage(chatId)
	case "/groupid", "/groupid@zerill_bot":
		sendGroupId(chatId)
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
	}
}

func sendStartMessage(chatId int, name string) {
	message := fmt.Sprintf("Xin chào %s \n\nGõ <code>/help</code> để xem danh sách các lệnh mà bot hỗ trợ nhé.\n\nBạn cũng có thể truy cập nhanh các chức năng bằng cách nhấn nút Menu bên dưới.", name)
	bot.SendMessage(chatId, message)
}

func sendHelpMessage(chatId int) {
	bot.SendMessage(chatId, "<code>/help</code> - Danh sách câu lệnh được hỗ trợ\n\n<code>/quote</code> - Xem trích dẫn hay ngẫu nhiên\n\n<code>/lunar</code> - Xem ngày âm lịch hôm nay\n\n<code>/bitcoin</code> - Xem giá Bitcoin mới nhất\n\n<code>/qr</code> - Tạo mã QR\n\n<code>/weather</code> - Xem tình hình thời tiết các tỉnh")
}

func sendGroupId(chatId int) {
	bot.SendMessage(chatId, fmt.Sprintf("Group ID: <code>%v</code>", chatId))
}
