package lottery

import (
	"fmt"
	"strings"
	"zeril-bot/utils/structs"

	"github.com/mmcdole/gofeed"
)

func GetDataLottery(zone string) (string, error) {
	switch zone {
	case "mien-nam-xsmn", "mien-bac-xsmb", "mien-trung-xsmt":
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(fmt.Sprintf("https://xosothienphu.com/ket-qua-xo-so-%s.rss", zone))

		if err != nil {
			return "", nil
		}

		message := strings.Replace(feed.Items[0].Description, "Giải", "\nGiải", -1)
		message = strings.Replace(message, "[", "\n\n[", -1)

		if zone == "mien-bac-xsmb" {
			message = strings.Replace(message, "ĐB:", "\n\nĐB:", -1)
		}
		return feed.Items[0].Title + message, nil
	default:
		return "Tôi không hiểu câu lệnh của bạn !!!", nil
	}
}

func GetSuggest(data structs.DataTele) (string, []structs.ButtonCallback) {
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

	message := "Hãy chọn khu vực muốn xem kết quả xổ số"

	return message, buttons
}
