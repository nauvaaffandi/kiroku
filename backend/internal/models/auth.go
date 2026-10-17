package models

type BootstrapReq struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	TOTPSecret string `json:"totp_secret,omitempty"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TOTPCode string `json:"totp_code,omitempty"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type PasswordChangeRequestReq struct {
	NewPassword string `json:"new_password"`
}

type VerifyPasswordChangeReq struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type RequestPasswordResetReq struct {
	Email string `json:"email"`
}

type VerifyPasswordResetReq struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
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
