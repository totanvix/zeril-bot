package quote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"
)

func SendAQuote(chatId int) {
	quote := getAQuote()
	quoteFormat := fmt.Sprintf("&quot;%s&quot; - <b>%s</b>", quote.Quote, quote.Author)

	telegram.SendMessage(chatId, quoteFormat)
}

func getAQuote() structs.QuoteData {
	res, err := http.Get("https://zenquotes.io/api/random")

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var data []structs.QuoteData
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Fatalln(err)
	}

	return data[0]
}
