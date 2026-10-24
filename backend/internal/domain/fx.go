package domain

import "time"

type FXRate struct {
	Base       string    `json:"base"`
	Quote      string    `json:"quote"`
	Rate       float64   `json:"rate"`
	Date       time.Time `json:"date"`
}

type FXQuote struct {
	Base      string    `json:"base"`
	Quote     string    `json:"quote"`
	Rate      float64   `json:"rate"`
	UpdatedAt time.Time `json:"updated_at"`
}
