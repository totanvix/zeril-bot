package qr

import (
	"log"
	"strings"
	"zeril-bot/utils/telegram"

	"github.com/skip2/go-qrcode"
)

func SendQRImage(chatId int, text string) {
	arr := strings.Fields(text)
	args := arr[1:]
	if len(args) == 0 {
		telegram.SendMessage(chatId, "Sử dụng cú pháp <code>/qr &lt;nội dung&gt;</code> để tạo mã QR.")
		return
	}

	content := text[4:]
	path := "/tmp/qr.png"
	err := qrcode.WriteFile(content, qrcode.Medium, 256, path)
	if err != nil {
		log.Fatalln(err)
	}

	telegram.SendAPhoto(chatId, path)
}
