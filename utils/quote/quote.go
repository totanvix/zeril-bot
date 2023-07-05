package quote

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"zeril-bot/utils/structs"
)

func GetAQuote() (*structs.QuoteData, error) {
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
