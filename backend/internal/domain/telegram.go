package domain

type TelegramWebhookUpdate struct {
	UpdateID int `json:"update_id"`
	Message  *struct {
		MessageID int `json:"message_id"`
		From      *struct {
			ID        int64  `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID    int64  `json:"id"`
			Type  string `json:"type"`
		} `json:"chat"`
		Date int64  `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

type TelegramBinding struct {
	ChatID   int64  `json:"chat_id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
