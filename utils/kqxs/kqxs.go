package kqxs

import (
	"fmt"
	"strings"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"

	"github.com/mmcdole/gofeed"
)

func Send(data structs.DataTele) error {
	text := data.RawMessage
	text = strings.TrimSpace(text)
	arr := strings.Fields(text)
	args := arr[1:]

	if len(args) != 1 {
		return SendSuggest(data)
	}

	zone := text[6:]

	switch zone {
	case "mien-nam-xsmn", "mien-bac-xsmb", "mien-trung-xsmt":
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(fmt.Sprintf("https://xosothienphu.com/ket-qua-xo-so-%s.rss", zone))
		fmt.Println(feed.Items[0].Description)

		message := strings.Replace(feed.Items[0].Description, "Giải", "\nGiải", -1)
		message = strings.Replace(message, "[", "\n\n[", -1)

		if zone == "mien-bac-xsmb" {
			message = strings.Replace(message, "ĐB:", "\n\nĐB:", -1)
		}

		data.ReplyMessage = feed.Items[0].Title + message
		return telegram.SendMessage(data)
	default:
		data.ReplyMessage = "Tôi không hiểu câu lệnh của bạn !!!"
		return telegram.SendMessage(data)
	}
}

func SendSuggest(data structs.DataTele) error {
	var buttons []structs.ButtonCallback
	var btn1, btn2, btn3 structs.ButtonCallback

	btn1.Text = "Miền Nam"
	btn1.CallbackData = "/kqxs mien-nam-xsmn"

	btn2.Text = "Miền Bắc"
	btn2.CallbackData = "/kqxs mien-bac-xsmb"

	btn3.Text = "Miền Trung"
	btn3.CallbackData = "/kqxs mien-trung-xsmt"

	buttons = append(buttons, btn1)
	buttons = append(buttons, btn2)
	buttons = append(buttons, btn3)

	data.ReplyMessage = "Hãy chọn khu vực muốn xem kết quả xổ số"
	return telegram.SendMessageWithReplyMarkup(data, buttons)

}
