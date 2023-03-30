package quote

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"zeril-bot/api/telegram"
)

type QuoteData struct {
	Quote string `json:"quote"`
}

func SendAQuote(chatId int) {
	quote := getAQuote()
	telegram.SendMessage(chatId, quote)
}

func getAQuote() string {
	res, err := http.Get("https://api.kanye.rest/")

	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var data QuoteData
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Println(err)
	}

	return data.Quote
}
