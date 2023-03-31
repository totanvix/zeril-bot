package bitcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"zeril-bot/utils/telegram"

	"github.com/leekchan/accounting"
)

type Btc struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type Exchange struct {
	Result float64 `json:"result"`
}

func SendBitcoinPrice(chatId int) {
	acUsd := accounting.Accounting{Symbol: "$", Precision: 2}
	acVnd := accounting.Accounting{Symbol: "", Precision: 0, Thousand: "."}

	btc := getBitcoinPrice()
	p, _ := strconv.ParseFloat(btc.Price, 64)
	usd := acUsd.FormatMoney(p)

	v := exchangeUsdToVnd(p)
	vnd := acVnd.FormatMoney(v) + " đ"

	message := fmt.Sprintf("1 Bitcoin = %s (<b>%s</b>)", usd, vnd)

	telegram.SendMessage(chatId, message)
}

func getBitcoinPrice() Btc {
	res, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var data Btc

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}

	return data
}

func exchangeUsdToVnd(p float64) float64 {
	price := fmt.Sprintf("%.2f", p)
	res, err := http.Get("https://api.exchangerate.host/convert?from=USD&to=VND&amount=" + price)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var data Exchange

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}

	return data.Result
}
