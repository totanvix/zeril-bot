package qr

import (
	"strings"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"

	"github.com/skip2/go-qrcode"
)

func SendQRImage(data structs.DataTele) error {
	text := data.RawMessage
	arr := strings.Fields(text)
	args := arr[1:]
	if len(args) == 0 {
		data.ReplyMessage = "Sử dụng cú pháp <code>/qr nội dung</code> để tạo mã QR."
		return telegram.SendMessage(data)
	}

	content := text[4:]
	path := "/tmp/qr.png"
	err := qrcode.WriteFile(content, qrcode.Medium, 256, path)
	if err != nil {
		return err
	}

	return telegram.SendPhoto(data, path)
}
