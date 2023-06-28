package bot

import (
	"fmt"
	"log"
	"strings"

	"zeril-bot/utils/bitcoin"
	"zeril-bot/utils/kqxs"
	"zeril-bot/utils/lunar"
	"zeril-bot/utils/qr"
	"zeril-bot/utils/quote"
	"zeril-bot/utils/random"
	"zeril-bot/utils/shortener"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"
	"zeril-bot/utils/weather"
)

type Bot struct {
	HookData structs.HookData
	rCh      chan rChannel
}

type rChannel struct {
	err error
}

func NewBot(hookData structs.HookData) *Bot {
	ch := make(chan rChannel, 2)

	return &Bot{HookData: hookData, rCh: ch}
}

func (b Bot) ResolveHook() error {
	go b.setTypingAction()

	switch {
	case b.isCommand():
		go b.resolveCommand()
	case b.isCallbackCommand():
		go b.resolveCallbackCommand()
	}

	for r := range b.rCh {
		if r.err != nil {
			return r.err
		}
	}

	return nil
}

func (b Bot) setTypingAction() {
	data := b.getTelegramData()
	rawMessage := b.getRawMessage()

	log.Printf("Yêu cầu từ bạn %s: %s", data.FirstName, rawMessage)

	err := telegram.SetTypingAction(data)

	defer func() {
		b.rCh <- rChannel{err: err}
	}()
}

func (b Bot) getTelegramData() structs.DataTele {
	rawMessage := b.getRawMessage()

	return structs.DataTele{
		ChatId:     b.getChatId(),
		ChatType:   b.getChatType(),
		Username:   b.getUsername(),
		FirstName:  b.getFirstName(),
		RawMessage: rawMessage,
		Command:    b.getCommand(),
	}
}

func (b Bot) resolveCommand() error {
	var err error

	defer func() {
		b.rCh <- rChannel{err: err}
	}()

	data := b.getTelegramData()

	switch data.Command {
	case "/start", "/start@zerill_bot":
		err = b.sendStartMessage(data)
	case "/help", "/help@zerill_bot":
		err = b.sendHelpMessage(data)
	case "/quote", "/quote@zerill_bot":
		err = quote.SendAQuote(data)
	case "/groupid", "/groupid@zerill_bot":
		err = b.sendGroupId(data)
	case "/lunar", "/lunar@zerill_bot":
		err = lunar.SendLunarDateNow(data)
	case "/weather", "/weather@zerill_bot":
		err = weather.SendForecastOfWeather(data)
	case "/bitcoin", "/bitcoin@zerill_bot":
		err = bitcoin.SendBitcoinPrice(data)
	case "/qr", "/qr@zerill_bot":
		err = qr.SendQRImage(data)
	case "/random", "/random@zerill_bot":
		err = random.RandomElements(data)
	case "/kqxs", "/kqxs@zerill_bot":
		err = kqxs.Send(data)
	case "/shortener", "/shortener@zerill_bot":
		shortener.Generate(data)
	default:
		err = b.invalidCommand(data)
	}

	return err
}

func (b Bot) resolveCallbackCommand() error {
	var err error

	defer func() {
		b.rCh <- rChannel{err: err}
		close(b.rCh)
	}()

	data := b.getTelegramData()

	switch data.Command {
	case "/weather":
		err = weather.SendForecastOfWeather(data)
	case "/kqxs":
		err = kqxs.Send(data)
	}

	return err
}

func (b Bot) isCommand() bool {
	return b.HookData.CallbackQuery.Data == ""
}

func (b Bot) isCallbackCommand() bool {
	return b.HookData.CallbackQuery.Data != ""
}

func (b Bot) getChatType() string {
	if b.isCallbackCommand() {
		return b.HookData.Message.Chat.Type
	}
	return b.HookData.Message.Chat.Type
}

func (b Bot) getUsername() string {
	if b.isCallbackCommand() {
		return b.HookData.CallbackQuery.From.Username
	}
	return b.HookData.Message.From.Username
}

func (b Bot) getRawMessage() string {
	if b.isCallbackCommand() {
		return b.HookData.CallbackQuery.Data
	}
	return b.HookData.Message.Text
}

func (b Bot) getCommand() string {
	var arr []string

	if b.isCallbackCommand() {
		data := b.HookData.CallbackQuery.Data
		arr = strings.Fields(data)
	} else {
		text := b.HookData.Message.Text
		arr = strings.Fields(text)
	}

	return arr[0]
}

func (b Bot) getChatId() int {
	if b.isCallbackCommand() {
		return b.HookData.CallbackQuery.Message.Chat.ID
	}
	return b.HookData.Message.Chat.ID
}

func (b Bot) getFirstName() string {
	if b.isCallbackCommand() {
		return b.HookData.CallbackQuery.Message.From.FirstName
	}
	return b.HookData.Message.From.FirstName
}

func (b Bot) sendStartMessage(data structs.DataTele) error {
	message := fmt.Sprintf("Xin chào %s \n\nGõ <code>/help</code> để xem danh sách các lệnh mà bot hỗ trợ nhé.\n\nBạn cũng có thể truy cập nhanh các chức năng bằng cách nhấn nút Menu bên dưới.", data.FirstName)
	data.ReplyMessage = message
	return telegram.SendMessage(data)
}

func (b Bot) sendHelpMessage(data structs.DataTele) error {
	messages := ""
	botCommands := b.getBotCommands()

	for _, command := range botCommands.Result {
		messages += fmt.Sprintf("<code>/%s</code> - %s\n\n", command.Command, command.Description)
	}

	data.ReplyMessage = messages

	return telegram.SendMessage(data)
}

func (b Bot) getBotCommands() structs.BotCommands {
	return telegram.GetBotCommands()
}

func (b Bot) sendGroupId(data structs.DataTele) error {
	if data.ChatType == "group" {
		data.ReplyMessage = fmt.Sprintf("Group ID: <code>%v</code>", data.ChatId)
	} else {
		data.ReplyMessage = "Không tìm thấy nhóm, bạn cần thêm bot vào nhóm trước khi thực hiện lệnh này !"
	}

	return telegram.SendMessage(data)
}

func (b Bot) invalidCommand(data structs.DataTele) error {
	data.ReplyMessage = "Tôi không hiểu câu lệnh của bạn !!!"
	return telegram.SendMessage(data)
}
