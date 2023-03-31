package lunar

import (
	"fmt"
	"strconv"
	"time"
	"zeril-bot/utils/telegram"

	"github.com/Lofanmi/chinese-calendar-golang/calendar"
)

func SendLunarDateNow(chatId int) {
	t := time.Now()
	c := calendar.ByTimestamp(t.Unix())

	d := strconv.Itoa(int(c.Lunar.GetDay()))
	m := strconv.Itoa(int(c.Lunar.GetMonth()))
	y := strconv.Itoa(int(c.Lunar.GetYear()))

	message := fmt.Sprintf("Âm lịch hôm nay là %s/%s/%s", d, m, y)

	telegram.SendMessage(chatId, message)
}
