package structs

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type HookData struct {
	UpdateId      int     `json:"update_id"`
	Message       Message `json:"message"`
	CallbackQuery struct {
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

type TelegramStatus struct {
	Ok bool `json:"ok"`
}