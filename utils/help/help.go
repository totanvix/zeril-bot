package help

import (
	"fmt"
	"zeril-bot/utils/bot"
	"zeril-bot/utils/channel"
)

func SendStartMessage(chatId int, name string) {
	message := fmt.Sprintf("Xin chào %s \n\nGõ <code>/help</code> để xem danh sách các lệnh mà bot hỗ trợ nhé.\n\nBạn cũng có thể truy cập nhanh các chức năng bằng cách nhấn nút Menu bên dưới.", name)
	channel.SendMessage(chatId, message)
}

func SendHelpMessage(chatId int) {
	messages := ""
	botCommands := bot.GetBotCommands()

	for _, command := range botCommands.Result {
		messages += fmt.Sprintf("<code>/%s</code> - %s\n\n", command.Command, command.Description)
	}

	channel.SendMessage(chatId, messages)
}

func SendGroupId(chatId int, chatType string) {
	if chatType == "group" {
		channel.SendMessage(chatId, fmt.Sprintf("Group ID: <code>%v</code>", chatId))
		return
	}

	channel.SendMessage(chatId, "Không tìm thấy nhóm, bạn cần thêm bot vào nhóm trước khi thực hiện lệnh này !")
}
