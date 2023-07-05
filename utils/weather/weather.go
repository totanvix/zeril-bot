package weather

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"zeril-bot/utils/structs"
)

const API_URL = "https://api.openweathermap.org"

var APP_ID = os.Getenv("OPEN_WEATHER_MAP_APP_ID")

func GetSuggestForecast(data structs.DataTele) (string, []structs.ButtonCallback) {
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

	message := "Sử dụng cú pháp <code>/weather tên thành phố</code> hoặc chọn các gợi ý bên dưới để xem thời tiết"

	return message, buttons
}

func GetWeather(cityName string) (structs.WeatherData, error) {
	uri := API_URL + "/data/2.5/weather"
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return structs.WeatherData{}, err
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
		return structs.WeatherData{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return structs.WeatherData{}, err
	}

	var data structs.WeatherData

	if res.StatusCode != 200 {
		log.Println(string(body))
		return data, errors.New("Không tìm thấy thông tin thời tiết")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return structs.WeatherData{}, err
	}

	log.Println("GetWeather OK")

	return data, nil
}
