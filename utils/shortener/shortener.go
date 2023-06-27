package shortener

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
	"zeril-bot/utils/redis"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"

	gonanoid "github.com/matoous/go-nanoid"
)

const rawAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Generate(data structs.DataTele) error {
	text := data.RawMessage
	arr := strings.Fields(text)
	args := arr[1:]
	fmt.Println(args)

	if len(args) == 0 {
		data.ReplyMessage = "Sử dụng cú pháp <code>/shortener https://example.com/</code> để tạo rút gọn liên kết"
		return telegram.SendMessage(data)
	}

	url := text[11:]
	if !isUrl(url) {
		data.ReplyMessage = "URL không đúng định dạng"
		return telegram.SendMessage(data)
	}

	id, err := gonanoid.Generate(rawAlphabet, 6)
	if err != nil {
		log.Panic(err)
	}

	redis.Set(id, url, time.Hour)

	shortUrl := fmt.Sprintf("%s/url/%s", os.Getenv("APP_URL"), id)
	data.ReplyMessage = fmt.Sprintf("URL sau khi rút gọn: <a href='%s'>%s</a>", shortUrl, shortUrl)

	return telegram.SendMessage(data)
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
