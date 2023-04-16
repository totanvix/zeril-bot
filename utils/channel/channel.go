package channel

import (
	"sync"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/structs"
)

var sendMessageChan = make(chan structs.SendMessage)
var sendPhotoChan = make(chan structs.SendPhoto)
var sendMessageWithReplyMarkupChan = make(chan structs.SendMessageWithReplyMarkup)
var Wg sync.WaitGroup

func init() {
	go func() {
		for {
			select {
			case sm := <-sendMessageChan:
				bot.SendMessage(sm.ChatId, sm.Message)
			case sp := <-sendPhotoChan:
				bot.SendPhoto(sp.ChatId, sp.Path)
			case srm := <-sendMessageWithReplyMarkupChan:
				bot.SendMessageWithReplyMarkup(srm.ChatId, srm.Message, srm.ReplyMark)
			}
			Wg.Done()
		}
	}()
}

func SendMessage(chatId int, message string) {
	sendMessageChan <- structs.SendMessage{ChatId: chatId, Message: message}
}

func SendPhoto(chatId int, path string) {
	sendPhotoChan <- structs.SendPhoto{ChatId: chatId, Path: path}
}

func SendMessageWithReplyMarkup(chatId int, message string, replyMark []structs.ButtonCallback) {
	sendMessageWithReplyMarkupChan <- structs.SendMessageWithReplyMarkup{
		ChatId:    chatId,
		Message:   message,
		ReplyMark: replyMark,
	}
}

func GetWg() *sync.WaitGroup {
	return &Wg
}
