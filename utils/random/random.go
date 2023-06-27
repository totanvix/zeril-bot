package random

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"
)

func RandomElements(data structs.DataTele) error {
	text := data.RawMessage
	arr := strings.Fields(text)

	if len(arr[1:]) == 0 {
		data.ReplyMessage = "Sử dụng cú pháp <code>/random A, B, C</code> để chọn phần tử ngẫu nhiên"
		return telegram.SendMessage(data)

	}

	f := func(c rune) bool {
		return unicode.IsPunct(c) == unicode.IsPunct(',')
	}

	els := strings.FieldsFunc(text[8:], f)
	el := els[rand.Intn(len(els))]

	data.ReplyMessage = fmt.Sprintf("Phần từ được chọn sau khi random: %v", strings.TrimSpace(el))

	return telegram.SendMessage(data)
}
