package quote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"zeril-bot/utils/channel"
	"zeril-bot/utils/structs"
)

func SendAQuote(chatId int) {
	quote := getAQuote()
	quoteFormat := fmt.Sprintf("&quot;%s&quot; - <b>%s</b>", quote.Quote, quote.Author)
	channel.SendMessage(chatId, quoteFormat)
}

func getAQuote() structs.QuoteData {
	res, err := http.Get("https://zenquotes.io/api/random")

	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var data []structs.QuoteData
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Panic(err)
	}

	return data[0]
}
