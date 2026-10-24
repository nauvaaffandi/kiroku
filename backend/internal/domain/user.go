package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	PinHash   string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type LoginInput struct {
	Email string
	Pin   string
}
