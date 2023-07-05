package shortener

import (
	"fmt"
	"net/url"
	"os"
	"time"
	"zeril-bot/utils/redis"

	gonanoid "github.com/matoous/go-nanoid"
)

const rawAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Generate(url string) (string, error) {
	if !isUrl(url) {
		return "URL không đúng định dạng", nil
	}

	id, err := gonanoid.Generate(rawAlphabet, 6)
	if err != nil {
		return "", err
	}

	redis.Set(id, url, time.Hour)
	shortUrl := fmt.Sprintf("%s/url/%s", os.Getenv("APP_URL"), id)

	return fmt.Sprintf("URL sau khi rút gọn: <a href='%s'>%s</a>", shortUrl, shortUrl), nil
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
