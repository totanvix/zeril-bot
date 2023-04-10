package shortener

import (
	"fmt"
	"log"
	"net/url"
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
		bot.SendMessage(chatId, "Sử dụng cú pháp <code>/shortener https://example.com/</code> để tạo rút gọn liên kết")
		return
	}

	url := text[11:]
	fmt.Println(url)

	if !isUrl(url) {
		bot.SendMessage(chatId, "URL không đúng định dạng")
		return
	}

	id, err := gonanoid.Generate(rawAlphabet, 6)
	if err != nil {
		log.Fatalln(err)
	}

	redis.Set(id, url, time.Hour)

	shortUrl := fmt.Sprintf("%s/url/%s", os.Getenv("APP_URL"), id)

	bot.SendMessage(chatId, fmt.Sprintf("URL sau khi rút gọn: <a href='%s'>%s</a>", shortUrl, shortUrl))

}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
