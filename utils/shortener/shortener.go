package shortener

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
	"zeril-bot/utils/channel"
	"zeril-bot/utils/redis"

	gonanoid "github.com/matoous/go-nanoid"
)

const rawAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Generate(chatId int, text string) {
	arr := strings.Fields(text)
	args := arr[1:]

	if len(args) != 1 {
		channel.SendMessage(chatId, "Sử dụng cú pháp <code>/shortener https://example.com/</code> để tạo rút gọn liên kết")
		return
	}

	url := text[11:]
	if !isUrl(url) {
		channel.SendMessage(chatId, "URL không đúng định dạng")
		return
	}

	id, err := gonanoid.Generate(rawAlphabet, 6)
	if err != nil {
		log.Panic(err)
	}

	redis.Set(id, url, time.Hour)

	shortUrl := fmt.Sprintf("%s/url/%s", os.Getenv("APP_URL"), id)

	channel.SendMessage(chatId, fmt.Sprintf("URL sau khi rút gọn: <a href='%s'>%s</a>", shortUrl, shortUrl))

}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
