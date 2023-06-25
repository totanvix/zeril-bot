package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"zeril-bot/utils/channel"
	"zeril-bot/utils/structs"
)

var APP_ID = os.Getenv("OPEN_WEATHER_MAP_APP_ID")
var API_URL = "https://api.openweathermap.org"

func SendForecastOfWeather(chatId int, text string) {
	text = strings.TrimSpace(text)
	arr := strings.Fields(text)
	args := arr[1:]

	if len(args) == 0 {
		SendSuggestForecast(chatId, args)
		return
	}

	cityName := text[9:]
	data, err := GetWeather(cityName)
	if err != nil {
		channel.SendMessage(chatId, "Không tìm thấy thông tin thời tiết")
		return
	}

	channel.SendMessage(chatId, fmt.Sprintf("🏙 Thời tiết hiện tại ở <b>%s</b>\n\n🌡 Nhiệt độ: <b>%.2f°C</b>\n\n💧 Độ ẩm: <b>%v&#37;</b>\n\nℹ️ Tổng quan: %s", data.Name, data.Main.Temp, data.Main.Humidity, data.Weather[0].Description))
}

func SendSuggestForecast(chatId int, args []string) {
	var buttons []structs.ButtonCallback
	var btn1, btn2, btn3 structs.ButtonCallback

	btn1.Text = "Hồ Chí Minh"
	btn1.CallbackData = "/weather ho chi minh"

	btn2.Text = "Hà Nội"
	btn2.CallbackData = "/weather ha noi"

	btn3.Text = "Nha Trang"
	btn3.CallbackData = "/weather nha trang"

	buttons = append(buttons, btn1)
	buttons = append(buttons, btn2)
	buttons = append(buttons, btn3)

	if len(args) == 0 {
		channel.SendMessageWithReplyMarkup(chatId, "Sử dụng cú pháp <code>/weather tên thành phố</code> hoặc chọn các gợi ý bên dưới để xem thời tiết", buttons)
		return
	}
}

func GetWeather(cityName string) (structs.WeatherData, error) {
	uri := API_URL + "/data/2.5/weather"
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		log.Panic(err)
	}

	q := req.URL.Query()
	q.Add("appid", APP_ID)
	q.Add("q", cityName)
	q.Add("lang", "vi")
	q.Add("units", "metric")

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Panic(err)
	}

	var data structs.WeatherData

	if res.StatusCode != 200 {
		log.Println(string(body))
		return data, errors.New("Không tìm thấy thông tin thời tiết")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Panic(err)
	}

	log.Println("GetWeather OK")

	return data, nil
}
