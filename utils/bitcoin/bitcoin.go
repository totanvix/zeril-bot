package bitcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"zeril-bot/utils/structs"

	"github.com/leekchan/accounting"
)

func GetBitcoinPrice() (string, error) {
	acUsd := accounting.Accounting{Symbol: "$", Precision: 2}

	btc, err := getCurrentPrice()
	if err != nil {
		return "", err
	}

	p, err := strconv.ParseFloat(btc.Price, 64)
	if err != nil {
		return "", err
	}

	usd := acUsd.FormatMoney(p)

	return fmt.Sprintf("1 Bitcoin = %s", usd), nil
}

func getCurrentPrice() (*structs.Btc, error) {
	res, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data structs.Btc

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
