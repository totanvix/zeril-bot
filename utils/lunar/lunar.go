package lunar

import (
	"fmt"
	"strconv"
	"time"
	"zeril-bot/utils/bot"

	"github.com/Lofanmi/chinese-calendar-golang/calendar"
)

func SendLunarDateNow(chatId int) {
	t := time.Now()
	c := calendar.ByTimestamp(t.Unix())

	d := c.Lunar.GetDay()
	m := strconv.Itoa(int(c.Lunar.GetMonth()))
	y := c.Lunar.GetYear()

	if len(m) == 1 {
		m = "0" + m
	}

	message := fmt.Sprintf("🗓 %v/%v/%v", d, m, y)

	bot.SendMessage(chatId, message)
}
