package structs

type Chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	FirstName string `json:"first_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Title     string `json:"title,omitempty"`
}

type From struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
	From From   `json:"from"`
}

type HookData struct {
	UpdateId      int     `json:"update_id"`
	Message       Message `json:"message"`
	CallbackQuery struct {
		From    From    `json:"from"`
		Message Message `json:"message"`
		Data    string  `json:"data"`
	} `json:"callback_query"`
}

type Btc struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type Exchange struct {
	Result float64 `json:"result"`
}

type QuoteData struct {
	Quote  string `json:"q"`
	Author string `json:"a"`
}

type Status struct {
	Ok bool `json:"ok"`
}

type WeatherData struct {
	Name    string `json:"name"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float32 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
}

type BotCommands struct {
	Status
	Result []struct {
		Command     string `json:"command"`
		Description string `json:"description"`
	}
}

type SendMessage struct {
	ChatId  int
	Message string
}

type SendPhoto struct {
	ChatId int
	Path   string
}

type SendMessageWithReplyMarkup struct {
	ChatId    int
	Message   string
	ReplyMark []ButtonCallback
}

type BodyReplyMarkup struct {
	ReplyMarkup struct {
		InlineKeyboard [][]ButtonCallback `json:"inline_keyboard"`
	} `json:"reply_markup"`
}

type ButtonCallback struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}
