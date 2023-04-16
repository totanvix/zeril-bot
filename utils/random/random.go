package random

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"
	"zeril-bot/utils/channel"
)

func RandomElements(chatId int, text string) {
	arr := strings.Fields(text)

	if len(arr[1:]) == 0 {
		channel.SendMessage(chatId, "Sử dụng cú pháp <code>/random A, B, C</code> để chọn phần tử ngẫu nhiên")
		return
	}

	f := func(c rune) bool {
		return unicode.IsPunct(c) == unicode.IsPunct(',')
	}

	els := strings.FieldsFunc(text[8:], f)

	el := els[rand.Intn(len(els))]

	channel.SendMessage(chatId, fmt.Sprintf("Phần từ được chọn sau khi random: %v", strings.TrimSpace(el)))
}
