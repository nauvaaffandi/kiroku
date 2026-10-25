package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// PIN is a custom type that accepts both JSON string and JSON number
// (frontend <input type="number"> serializes PIN as a JSON number).
type PIN string

func (p *PIN) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*p = PIN(s)
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*p = PIN(strconv.Itoa(int(f)))
		return nil
	}
	return fmt.Errorf("pin must be a string or number, got %s", string(data))
}

type BootstrapReq struct {
	Email      string `json:"email"`
	PIN        PIN    `json:"pin"`
	TOTPSecret string `json:"totp_secret,omitempty"`
}

type LoginReq struct {
	Email    string `json:"email"`
	PIN      PIN    `json:"pin"`
	TOTPCode string `json:"totp_code,omitempty"`
}

type ChangePINReq struct {
	CurrentPIN PIN `json:"current_pin"`
	NewPIN     PIN `json:"new_pin"`
}

type PinChangeRequestReq struct {
	NewPIN PIN `json:"new_pin"`
}

type VerifyPinChangeReq struct {
	Token  string `json:"token"`
	NewPIN PIN    `json:"new_pin"`
}

type RequestPinResetReq struct {
	Email string `json:"email"`
}

type VerifyPinResetReq struct {
	Token  string `json:"new_pin"`
	NewPIN PIN    `json:"new_pin"`
}

type AuthResponse struct {
	Status string `json:"status"`
	UserID string `json:"user_id,omitempty"`
}

type MeResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type StatusResponse struct {
	HasUsers bool `json:"has_users"`
}
