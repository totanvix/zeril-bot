package bot

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"
	"zeril-bot/utils/bitcoin"
	"zeril-bot/utils/lottery"
	"zeril-bot/utils/lunar"
	"zeril-bot/utils/qr"
	"zeril-bot/utils/quote"
	"zeril-bot/utils/shortener"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/weather"
)

func (b Bot) sendStartMessage(data structs.DataTele) error {
	message := fmt.Sprintf("Xin chào %s \n\nGõ <code>/help</code> để xem danh sách các lệnh mà bot hỗ trợ nhé.\n\nBạn cũng có thể truy cập nhanh các chức năng bằng cách nhấn nút Menu bên dưới.", data.FirstName)
	data.ReplyMessage = message
	return b.Telegram.SendMessage(data)
}

func (b Bot) sendHelpMessage(data structs.DataTele) error {
	messages := ""
	botCommands, err := b.getBotCommands()

	if err != nil {
		return err
	}

	for _, command := range botCommands.Result {
		messages += fmt.Sprintf("<code>/%s</code> - %s\n\n", command.Command, command.Description)
	}

	data.ReplyMessage = messages

	return b.Telegram.SendMessage(data)
}

func (b Bot) sendGroupId(data structs.DataTele) error {
	if data.ChatType == "group" {
		data.ReplyMessage = fmt.Sprintf("Group ID: <code>%v</code>", data.ChatId)
	} else {
		data.ReplyMessage = "Không tìm thấy nhóm, bạn cần thêm bot vào nhóm trước khi thực hiện lệnh này !"
	}

	return b.Telegram.SendMessage(data)
}

func (b Bot) invalidCommand(data structs.DataTele) error {
	data.ReplyMessage = "Tôi không hiểu câu lệnh của bạn !!!"
	return b.Telegram.SendMessage(data)
}

func (b Bot) sendAQuote(data structs.DataTele) error {
	quote, err := quote.GetAQuote()
	if err != nil {
		return err
	}

	quoteFormat := fmt.Sprintf("&quot;%s&quot; - <b>%s</b>", quote.Quote, quote.Author)
	data.ReplyMessage = quoteFormat

	return b.Telegram.SendMessage(data)
}

func (b Bot) sendLunarDateNow(data structs.DataTele) error {
	path, err := lunar.DownloadAndCropImage()
	if err != nil {
		return err
	}

	return b.Telegram.SendPhoto(data, path)
}

func (b Bot) sendForecastOfWeather(data structs.DataTele) error {
	text := data.RawMessage
	text = strings.TrimSpace(text)
	arr := strings.Fields(text)
	args := arr[1:]

	if len(args) == 0 {
		message, buttons := weather.GetSuggestForecast(data)
		data.ReplyMessage = message
		return b.Telegram.SendMessageWithReplyMarkup(data, buttons)
	}

	cityName := text[9:]
	wData, err := weather.GetWeather(cityName)
	if err != nil {
		data.ReplyMessage = "Không tìm thấy thông tin thời tiết"
		return b.Telegram.SendMessage(data)
	}

	data.ReplyMessage = fmt.Sprintf("🏙 Thời tiết hiện tại ở <b>%s</b>\n\n🌡 Nhiệt độ: <b>%.2f°C</b>\n\n💧 Độ ẩm: <b>%v&#37;</b>\n\nℹ️ Tổng quan: %s", wData.Name, wData.Main.Temp, wData.Main.Humidity, wData.Weather[0].Description)

	return b.Telegram.SendMessage(data)
}

func (b Bot) sendBitcoinPrice(data structs.DataTele) error {
	message, err := bitcoin.GetBitcoinPrice()
	if err != nil {
		return err
	}

	data.ReplyMessage = message

	return b.Telegram.SendMessage(data)
}

func (b Bot) sendQRImage(data structs.DataTele) error {
	text := data.RawMessage
	arr := strings.Fields(text)
	args := arr[1:]
	if len(args) == 0 {
		data.ReplyMessage = "Sử dụng cú pháp <code>/qr nội dung</code> để tạo mã QR."
		return b.Telegram.SendMessage(data)
	}

	content := text[4:]
	path, err := qr.CreateImage(content)
	if err != nil {
		return err
	}

	return b.Telegram.SendPhoto(data, path)
}

func (b Bot) sendSelectedElement(data structs.DataTele) error {
	text := data.RawMessage
	arr := strings.Fields(text)

	if len(arr[1:]) == 0 {
		data.ReplyMessage = "Sử dụng cú pháp <code>/random A, B, C</code> để chọn phần tử ngẫu nhiên"
		return b.Telegram.SendMessage(data)

	}

	f := func(c rune) bool {
		return unicode.IsPunct(c) == unicode.IsPunct(',')
	}

	els := strings.FieldsFunc(text[8:], f)
	el := els[rand.Intn(len(els))]

	data.ReplyMessage = fmt.Sprintf("Phần từ được chọn sau khi random: %v", strings.TrimSpace(el))

	return b.Telegram.SendMessage(data)
}

func (b Bot) sendLottery(data structs.DataTele) error {
	text := data.RawMessage
	text = strings.TrimSpace(text)
	arr := strings.Fields(text)
	args := arr[1:]

	if len(args) != 1 {
		message, buttons := lottery.GetSuggest(data)
		data.ReplyMessage = message
		return b.Telegram.SendMessageWithReplyMarkup(data, buttons)
	}

	zone := text[6:]

	message, err := lottery.GetDataLottery(zone)
	if err != nil {
		return err
	}

	data.ReplyMessage = message

	return b.Telegram.SendMessage(data)
}

func (b Bot) generateShortenerURL(data structs.DataTele) error {
	text := data.RawMessage
	arr := strings.Fields(text)
	args := arr[1:]

	if len(args) == 0 {
		data.ReplyMessage = "Sử dụng cú pháp <code>/shortener https://example.com/</code> để tạo rút gọn liên kết"
		return b.Telegram.SendMessage(data)
	}

	url := text[11:]
	message, err := shortener.Generate(url)
	if err != nil {
		return err
	}

	data.ReplyMessage = message

	return b.Telegram.SendMessage(data)
}
