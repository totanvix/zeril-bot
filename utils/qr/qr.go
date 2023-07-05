package qr

import (
	"github.com/skip2/go-qrcode"
)

func CreateImage(content string) (string, error) {
	path := "/tmp/qr.png"
	err := qrcode.WriteFile(content, qrcode.Medium, 256, path)
	if err != nil {
		return "", err
	}

	return path, nil
}
