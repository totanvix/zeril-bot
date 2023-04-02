package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"zeril-bot/utils/bitcoin"
	"zeril-bot/utils/lunar"
	"zeril-bot/utils/qr"
	"zeril-bot/utils/quote"
	"zeril-bot/utils/telegram"
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

	name := data.Message.From.FirstName
	username := data.Message.From.Username
	chatId := data.Message.Chat.ID
	text := data.Message.Text
	arr := strings.Fields(text)

	command := arr[0]
	args := arr[1:]

	telegram.SetTypingAction(chatId)

	log.Println(fmt.Sprintf("Yêu cầu từ bạn %s(%s): %s", name, username, text))

	switch command {
	case "/start":
		message := fmt.Sprintf("Xin chào %s \n\nGõ <code>/help</code> để xem danh sách các lệnh mà bot hỗ trợ nhé.\n\nBạn cũng có thể truy cập nhanh các chức năng bằng cách nhấn nút Menu bên dưới.", name)
		telegram.SendMessage(chatId, message)
	case "/help":
		telegram.SendMessage(chatId, "<code>/help</code> - Danh sách câu lệnh được hỗ trợ\n\n<code>/quote</code> - Xem trích dẫn hay ngẫu nhiên\n\n<code>/lunar</code> - Xem ngày âm lịch hôm nay\n\n<code>/bitcoin</code> - Xem giá Bitcoin mới nhất\n\n<code>/qr</code> - Tạo mã QR\n\n<code>/weather</code> - Xem tình hình thời tiết các tỉnh")
	case "/quote":
		quote.SendAQuote(chatId)
	case "/lunar":
		lunar.SendLunarDateNow(chatId)
	case "/bitcoin":
		bitcoin.SendBitcoinPrice(chatId)
	case "/qr":
		if len(args) == 0 {
			telegram.SendMessage(chatId, "Sử dụng cú pháp <code>/qr nội dung viết liền, không khoảng cách</code> để tạo mã QR.")
			return
		}

		qr.SendQRImage(chatId, args)
	default:
		telegram.SendMessage(chatId, "Tôi không hiểu câu lệnh của bạn !!!")
	}
}
