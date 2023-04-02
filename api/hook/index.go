package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"zeril-bot/utils"
	"zeril-bot/utils/bitcoin"
	"zeril-bot/utils/lunar"
	"zeril-bot/utils/qr"
	"zeril-bot/utils/quote"
	"zeril-bot/utils/telegram"
	"zeril-bot/utils/weather"
)

func Router(w http.ResponseWriter, r *http.Request) {

	var data utils.HookData
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		var callback utils.CallbackData
		err := json.NewDecoder(r.Body).Decode(&callback)
		if err != nil {
			log.Fatalln(err.Error())
		}

		resolveCallback(callback)
		return
	}

	resolveCommand(data)
}

func resolveCommand(data utils.HookData) {
	name := data.Message.Chat.FirstName
	username := data.Message.Chat.Username
	chatId := data.Message.Chat.ID
	text := data.Message.Text
	arr := strings.Fields(text)

	telegram.SetTypingAction(chatId)

	log.Println(fmt.Sprintf("Yêu cầu từ bạn %s(%s): %s", name, username, text))

	command := arr[0]

	switch command {
	case "/start":
		sendStartMessage(chatId, name)
	case "/help":
		sendHelpMessage(chatId)
	case "/quote":
		quote.SendAQuote(chatId)
	case "/lunar":
		lunar.SendLunarDateNow(chatId)
	case "/weather":
		weather.SendForecastOfWeather(chatId, text)
	case "/bitcoin":
		bitcoin.SendBitcoinPrice(chatId)
	case "/qr":
		qr.SendQRImage(chatId, text)
	default:
		telegram.SendMessage(chatId, "Tôi không hiểu câu lệnh của bạn !!!")
	}
}

func resolveCallback(callback utils.CallbackData) {
	name := callback.CallbackQuery.Chat.FirstName
	username := callback.CallbackQuery.Chat.Username
	chatId := callback.CallbackQuery.Chat.ID
	text := callback.CallbackQuery.Text
	data := callback.CallbackQuery.Data

	telegram.SetTypingAction(chatId)

	log.Println(fmt.Sprintf("Yêu cầu từ bạn %s(%s): %s, callback data: %s", name, username, text, data))

	// command := arr[0]
	// args := arr[1:]

	// switch command {

	// }
}

func sendStartMessage(chatId int, name string) {
	message := fmt.Sprintf("Xin chào %s \n\nGõ <code>/help</code> để xem danh sách các lệnh mà bot hỗ trợ nhé.\n\nBạn cũng có thể truy cập nhanh các chức năng bằng cách nhấn nút Menu bên dưới.", name)
	telegram.SendMessage(chatId, message)
}

func sendHelpMessage(chatId int) {
	telegram.SendMessage(chatId, "<code>/help</code> - Danh sách câu lệnh được hỗ trợ\n\n<code>/quote</code> - Xem trích dẫn hay ngẫu nhiên\n\n<code>/lunar</code> - Xem ngày âm lịch hôm nay\n\n<code>/bitcoin</code> - Xem giá Bitcoin mới nhất\n\n<code>/qr</code> - Tạo mã QR\n\n<code>/weather</code> - Xem tình hình thời tiết các tỉnh")
}
