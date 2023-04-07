package kqxs

import (
	"fmt"
	"strings"
	"zeril-bot/utils/bot"

	"github.com/mmcdole/gofeed"
)

func Send(chatId int, text string) {
	text = strings.TrimSpace(text)
	arr := strings.Fields(text)
	args := arr[1:]

	if len(args) == 0 {
		SendSuggest(chatId, args)
		return
	}

	zone := text[6:]

	fmt.Println(zone)

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(fmt.Sprintf("https://xosothienphu.com/ket-qua-xo-so-%s.rss", zone))
	fmt.Println(feed.Items[0].Description)

	message := strings.Replace(feed.Items[0].Description, "Giải", "\nGiải", -1)
	message = strings.Replace(message, "[", "\n\n[", -1)

	bot.SendMessage(chatId, feed.Items[0].Title+message)
}

func SendSuggest(chatId int, args []string) {
	var buttons []bot.ButtonCallback
	var btn1, btn2, btn3 bot.ButtonCallback

	btn1.Text = "Miền Nam"
	btn1.CallbackData = "/kqxs mien-nam-xsmn"

	btn2.Text = "Miền Bắc"
	btn2.CallbackData = "/kqxs mien-bac-xsmb"

	btn3.Text = "Miền Trung"
	btn3.CallbackData = "/kqxs mien-trung-xsmt"

	buttons = append(buttons, btn1)
	buttons = append(buttons, btn2)
	buttons = append(buttons, btn3)

	if len(args) == 0 {
		bot.SendMessageWithReplyMarkup(chatId, "Hãy chọn khu vực muốn xem kết quả xổ số", buttons)
		return
	}
}
