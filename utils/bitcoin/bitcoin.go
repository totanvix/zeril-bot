package bitcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"

	"github.com/leekchan/accounting"
)

func SendBitcoinPrice(data structs.DataTele) error {
	acUsd := accounting.Accounting{Symbol: "$", Precision: 2}
	acVnd := accounting.Accounting{Symbol: "", Precision: 0, Thousand: "."}

	btc, err := getBitcoinPrice()
	if err != nil {
		return err
	}

	p, err := strconv.ParseFloat(btc.Price, 64)
	if err != nil {
		return err
	}

	usd := acUsd.FormatMoney(p)
	v, err := exchangeUsdToVnd(p)
	if err != nil {
		return err
	}

	vnd := acVnd.FormatMoney(*v) + " đ"

	message := fmt.Sprintf("1 Bitcoin = %s (<b>%s</b>)", usd, vnd)

	data.ReplyMessage = message

	return telegram.SendMessage(data)
}

func getBitcoinPrice() (*structs.Btc, error) {
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

func exchangeUsdToVnd(p float64) (*float64, error) {
	price := fmt.Sprintf("%.2f", p)

	res, err := http.Get("https://api.exchangerate.host/convert?from=USD&to=VND&amount=" + price)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data structs.Exchange

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data.Result, nil
}
