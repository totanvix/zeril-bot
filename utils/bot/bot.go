package bot

import (
	"errors"
	"log"
	"strings"
	"time"

	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"
)

type Bot struct {
	HookData structs.HookData
	rCh      chan rChannel
	Telegram *telegram.Telegram
}

type rChannel struct {
	err error
}

const numberOfCh = 2

func NewBot(telegram *telegram.Telegram, hookData structs.HookData) *Bot {
	ch := make(chan rChannel, numberOfCh)
	return &Bot{Telegram: telegram, HookData: hookData, rCh: ch}
}

func (b Bot) ResolveHook() error {
	go b.setTypingAction()

	switch {
	case b.isCommand():
		go b.resolveCommand()
	case b.isCallbackCommand():
		go b.resolveCallbackCommand()
	}

	for i := 0; i < numberOfCh; i++ {
		select {
		case r := <-b.rCh:
			if r.err != nil {
				return r.err
			}
		case <-time.After(10 * time.Second):
			return errors.New("Timeout")
		}
	}

	close(b.rCh)

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
		err = b.sendAQuote(data)
	case "/groupid", "/groupid@zerill_bot":
		err = b.sendGroupId(data)
	case "/lunar", "/lunar@zerill_bot":
		err = b.sendLunarDateNow(data)
	case "/weather", "/weather@zerill_bot":
		err = b.sendForecastOfWeather(data)
	case "/bitcoin", "/bitcoin@zerill_bot":
		err = b.sendBitcoinPrice(data)
	case "/qr", "/qr@zerill_bot":
		err = b.sendQRImage(data)
	case "/random", "/random@zerill_bot":
		err = b.sendSelectedElement(data)
	case "/kqxs", "/kqxs@zerill_bot":
		err = b.sendLottery(data)
	case "/shortener", "/shortener@zerill_bot":
		err = b.generateShortenerURL(data)
	default:
		err = b.invalidCommand(data)
	}

	return err
}

func (b Bot) resolveCallbackCommand() error {
	var err error

	defer func() {
		b.rCh <- rChannel{err: err}
	}()

	data := b.getTelegramData()

	switch data.Command {
	case "/weather":
		err = b.sendForecastOfWeather(data)
	case "/kqxs":
		err = b.sendLottery(data)
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

func (b Bot) getBotCommands() (*structs.BotCommands, error) {
	return telegram.GetBotCommands()
}
