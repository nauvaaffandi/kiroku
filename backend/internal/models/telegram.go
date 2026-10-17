package models

import "encoding/json"

type WebhookUpdate struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	MessageID int      `json:"message_id"`
	Chat      Chat     `json:"chat"`
	From      *User    `json:"from"`
	Text      string   `json:"text"`
	Entities  []Entity `json:"entities"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Entity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type TelegramResponse struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}