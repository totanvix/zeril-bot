package weather

import (
	"strings"
	"zeril-bot/utils/telegram"
)

func SendForecastOfWeather(chatId int, text string) {
	arr := strings.Fields(text)
	args := arr[1:]
	if len(args) == 0 {
		SendSuggestForecast(chatId, args)
	}
}

func SendSuggestForecast(chatId int, args []string) {
	var buttons []telegram.ButtonCallback
	var btn1, btn2, btn3 telegram.ButtonCallback

	btn1.Text = "Hồ Chí Minh"
	btn1.CallbackData = "/weather ho chi minh"

	btn2.Text = "Hà Nội"
	btn2.CallbackData = "/weather ha noi"

	btn3.Text = "Lâm Đồng"
	btn3.CallbackData = "/weather lam dong"

	buttons = append(buttons, btn1)
	buttons = append(buttons, btn2)
	buttons = append(buttons, btn3)

	if len(args) == 0 {
		telegram.SendMessageWithReplyMarkup(chatId, "Sử dụng cú pháp <code>/weather &lt;tên tỉnh thành phố&gt;</code> hoặc chọn các gợi ý bên dưới để xem thời tiết", buttons)
		return
	}
}
