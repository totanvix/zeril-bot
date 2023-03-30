package quote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"zeril-bot/utils/telegram"
)

type QuoteData struct {
	Quote  string `json:"q"`
	Author string `json:"a"`
}

func SendAQuote(chatId int) {
	quote := getAQuote()
	quoteFormat := fmt.Sprintf("&quot;%s&quot; - <b>%s</b>", quote.Quote, quote.Author)

	telegram.SendMessage(chatId, quoteFormat)
}

func getAQuote() QuoteData {
	res, err := http.Get("https://zenquotes.io/api/random")

	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var data []QuoteData
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Println(err)
	}

	return data[0]
}
