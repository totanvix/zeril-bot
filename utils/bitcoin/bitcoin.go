package bitcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"zeril-bot/utils/channel"
	"zeril-bot/utils/structs"

	"github.com/leekchan/accounting"
)

func SendBitcoinPrice(chatId int) {
	acUsd := accounting.Accounting{Symbol: "$", Precision: 2}
	acVnd := accounting.Accounting{Symbol: "", Precision: 0, Thousand: "."}

	btc := getBitcoinPrice()
	p, _ := strconv.ParseFloat(btc.Price, 64)
	usd := acUsd.FormatMoney(p)

	v := exchangeUsdToVnd(p)
	vnd := acVnd.FormatMoney(v) + " đ"

	message := fmt.Sprintf("1 Bitcoin = %s (<b>%s</b>)", usd, vnd)

	channel.SendMessage(chatId, message)
}

func getBitcoinPrice() structs.Btc {
	res, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT")

	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Panic(err)
	}

	var data structs.Btc

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Panic(err)
	}

	return data
}

func exchangeUsdToVnd(p float64) float64 {
	price := fmt.Sprintf("%.2f", p)
	res, err := http.Get("https://api.exchangerate.host/convert?from=USD&to=VND&amount=" + price)

	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Panic(err)
	}

	var data structs.Exchange

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Panic(err)
	}

	return data.Result
}
