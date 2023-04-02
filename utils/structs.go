package utils

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type HookData struct {
	Message struct {
		Text string `json:"text"`
		Chat Chat   `json:"chat"`
	} `json:"message"`
}

type CallbackData struct {
	CallbackQuery struct {
		Chat Chat   `json:"chat"`
		Text string `json:"text"`
		Data string `json:"data"`
	}
}
