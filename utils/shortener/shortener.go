package shortener

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/redis"

	gonanoid "github.com/matoous/go-nanoid"
)

const rawAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Do(chatId int, text string) {
	arr := strings.Fields(text)
	args := arr[1:]
	if len(args) == 0 {
		bot.SendMessage(chatId, "Sử dụng cú pháp <code>/shorturl https://example.com/</code> để tạo rút gọn liên kết")
		return
	}

	url := text[10:]
	id, err := gonanoid.Generate(rawAlphabet, 8)
	if err != nil {
		log.Fatalln(err)
	}

	redis.Set(id, url, time.Hour)

	bot.SendMessage(chatId, fmt.Sprintf("URL sau khi rút gọn: <code>%s/url/%s</code>", os.Getenv("APP_URL"), id))
}
