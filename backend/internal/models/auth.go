package models

type BootstrapReq struct {
	Email      string `json:"email"`
	PIN        string `json:"pin"`
	TOTPSecret string `json:"totp_secret,omitempty"`
}

type LoginReq struct {
	Email    string `json:"email"`
	PIN      string `json:"pin"`
	TOTPCode string `json:"totp_code,omitempty"`
}

type ChangePINReq struct {
	CurrentPIN string `json:"current_pin"`
	NewPIN     string `json:"new_pin"`
}

type PinChangeRequestReq struct {
	NewPIN string `json:"new_pin"`
}

type VerifyPinChangeReq struct {
	Token  string `json:"token"`
	NewPIN string `json:"new_pin"`
}

type RequestPinResetReq struct {
	Email string `json:"email"`
}

type VerifyPinResetReq struct {
	Token  string `json:"token"`
	NewPIN string `json:"new_pin"`
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
