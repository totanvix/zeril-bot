package quote

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"
)

func SendAQuote(data telegram.Data) error {
	quote, err := getAQuote()
	if err != nil {
		return err
	}

	quoteFormat := fmt.Sprintf("&quot;%s&quot; - <b>%s</b>", quote.Quote, quote.Author)
	data.Message = quoteFormat

	return telegram.SendMessage(data)
}

func getAQuote() (*structs.QuoteData, error) {
	res, err := http.Get("https://zenquotes.io/api/random")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data []structs.QuoteData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data[0], nil
}
