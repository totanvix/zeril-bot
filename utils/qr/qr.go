package qr

import (
	"log"
	"zeril-bot/utils/telegram"

	"github.com/skip2/go-qrcode"
)

func SendQRImage(chatId int, args []string) {
	if len(args) > 1 {
		telegram.SendMessage(chatId, "Không thể tạo mã QR, vì sai cú pháp\nMẫu tạo QR: <code>/qr https://www.google.com/</code>")
		return
	}

	path := "/tmp/qr.png"
	err := qrcode.WriteFile(args[0], qrcode.Medium, 256, path)
	if err != nil {
		log.Fatalln(err)
	}

	telegram.SendAPhoto(chatId, path)
}
