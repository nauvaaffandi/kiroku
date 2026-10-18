package telegram

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"kiroku/backend/internal/config"
	"kiroku/backend/internal/db"
	"kiroku/backend/internal/models"
)

type TelegramService struct {
	db            *db.DB
	cfg           *config.Config
	botToken      string
	botUsername   string
	webhookSecret string
	webhookURL    string
	httpClient    *http.Client
}

func NewTelegramService(database *db.DB, cfg *config.Config) *TelegramService {
	botToken := cfg.TelegramBotToken
	botUsername := cfg.TelegramBotUsername
	if botUsername == "" {
		botUsername = "nauva_kiroku_bot"
	}

	return &TelegramService{
		db:            database,
		cfg:           cfg,
		botToken:      botToken,
		botUsername:   botUsername,
		webhookSecret: cfg.TelegramWebhookSecret,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *TelegramService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/webhook", s.HandleWebhook)
	return r
}

func (s *TelegramService) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if s.botToken == "" {
		log.Warn().Msg("telegram: bot token not configured")
		w.WriteHeader(http.StatusOK)
		return
	}

	if s.webhookSecret != "" {
		secretHeader := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if secretHeader != s.webhookSecret {
			log.Warn().Msg("telegram: unauthorized webhook request")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	var update models.WebhookUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Error().Err(err).Msg("telegram: decode webhook")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if update.Message != nil && update.Message.Text != "" {
		s.handleCommand(update.Message)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *TelegramService) handleCommand(msg *models.Message) {
	text := strings.TrimSpace(msg.Text)
	chatID := fmt.Sprintf("%d", msg.Chat.ID)
	telegramUserID := ""
	if msg.From != nil {
		telegramUserID = fmt.Sprintf("%d", msg.From.ID)
	}

	switch {
	case strings.HasPrefix(text, "/start"):
		s.handleStart(chatID, telegramUserID, text)
	case strings.HasPrefix(text, "/register"):
		s.handleRegister(chatID, telegramUserID, text)
	case strings.HasPrefix(text, "/changepin"):
		s.handleChangePin(chatID, telegramUserID)
	case strings.HasPrefix(text, "/forgotpin"):
		s.handleForgotPin(chatID, telegramUserID)
	default:
		s.sendMessage(chatID, "Unknown command. Available: /register <token>, /changepin, /forgotpin")
	}
}

func (s *TelegramService) handleStart(chatID, telegramUserID, text string) {
	parts := strings.Fields(text)
	if len(parts) >= 2 && strings.HasPrefix(parts[1], "register_") {
		token := strings.TrimPrefix(parts[1], "register_")
		s.handleRegister(chatID, telegramUserID, "/register "+token)
		return
	}
	s.sendMessage(chatID, "Welcome to Kiroku Bot! Available commands: /register <token>, /changepin, /forgotpin")
}

func (s *TelegramService) handleRegister(chatID, telegramUserID, text string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		s.sendMessage(chatID, "Usage: /register <token>\n\nGet a token from Kiroku web settings.")
		return
	}

	token := parts[1]
	var userID string
	err := s.db.QueryRow(`
		SELECT user_id FROM pin_change_tokens 
		WHERE token = ? AND type = 'register' AND expires_at > datetime('now')
	`, token).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			s.sendMessage(chatID, "❌ Invalid or expired registration token. Generate a new one from Kiroku web settings.")
		} else {
			log.Error().Err(err).Msg("telegram: register token lookup")
			s.sendMessage(chatID, "❌ Error validating token")
		}
		return
	}

	_, err = s.db.Exec(`
		INSERT OR REPLACE INTO telegram_bindings (id, user_id, telegram_chat_id, telegram_user_id)
		VALUES (?, ?, ?, ?)
	`, uuid.New().String(), userID, chatID, telegramUserID)
	if err != nil {
		log.Error().Err(err).Msg("telegram: save binding")
		s.sendMessage(chatID, "❌ Failed to save binding")
		return
	}

	_, err = s.db.Exec(`DELETE FROM pin_change_tokens WHERE token = ?`, token)
	if err != nil {
		log.Warn().Err(err).Msg("telegram: cleanup register token")
	}

	s.sendMessage(chatID, "✅ Successfully registered! Your Telegram is now linked to Kiroku.\n\nYou can now use /changepin and /forgotpin commands.")
}

func (s *TelegramService) handleChangePin(chatID, telegramUserID string) {
	var userID string
	err := s.db.QueryRow(`
		SELECT user_id FROM telegram_bindings WHERE telegram_chat_id = ?
	`, chatID).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			s.sendMessage(chatID, "❌ This Telegram chat is not linked to any Kiroku account. Use /register <token> first.")
		} else {
			log.Error().Err(err).Msg("telegram: changepin lookup")
			s.sendMessage(chatID, "❌ Error checking binding")
		}
		return
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(5 * time.Minute)
	_, err = s.db.Exec(`
		INSERT INTO pin_change_tokens (id, user_id, token, type, expires_at)
		VALUES (?, ?, ?, 'changepassword', ?)
	`, uuid.New().String(), userID, token, expiresAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error().Err(err).Msg("telegram: create changepassword token")
		s.sendMessage(chatID, "❌ Failed to generate token")
		return
	}

	msg := fmt.Sprintf(
		"🔑 PIN Change Token (valid 5 minutes):\n\n`%s`\n\nPaste this token in Kiroku web settings along with your new PIN.",
		token,
	)
	s.sendMessage(chatID, msg)
}

func (s *TelegramService) handleForgotPin(chatID, telegramUserID string) {
	var userID string
	err := s.db.QueryRow(`
		SELECT user_id FROM telegram_bindings WHERE telegram_chat_id = ?
	`, chatID).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			s.sendMessage(chatID, "❌ This Telegram chat is not linked to any Kiroku account. Use /register <token> first.")
		} else {
			log.Error().Err(err).Msg("telegram: forgotpin lookup")
			s.sendMessage(chatID, "❌ Error checking binding")
		}
		return
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(10 * time.Minute)
	_, err = s.db.Exec(`
		INSERT INTO pin_change_tokens (id, user_id, token, type, expires_at)
		VALUES (?, ?, ?, 'forgotpassword', ?)
	`, uuid.New().String(), userID, token, expiresAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Error().Err(err).Msg("telegram: create forgotpassword token")
		s.sendMessage(chatID, "❌ Failed to generate token")
		return
	}

	msg := fmt.Sprintf(
		"🔐 PIN Reset Token (valid 10 minutes):\n\n`%s`\n\nUse this token in Kiroku web to reset your PIN.",
		token,
	)
	s.sendMessage(chatID, msg)
}

func (s *TelegramService) sendMessage(chatID, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.botToken)

	payload := map[string]string{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var tgResp models.TelegramResponse
	json.NewDecoder(resp.Body).Decode(&tgResp)

	if !tgResp.Ok {
		log.Warn().Str("chat_id", chatID).Msg("telegram: send failed")
	}

	return nil
}

func (s *TelegramService) SetWebhook(webhookURL string) error {
	if s.botToken == "" {
		return fmt.Errorf("bot token not configured")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", s.botToken)
	payload := map[string]interface{}{
		"url":                  webhookURL,
		"drop_pending_updates": true,
	}
	if s.webhookSecret != "" {
		payload["secret_token"] = s.webhookSecret
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal setWebhook payload: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("create setWebhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("set webhook request: %w", err)
	}
	defer resp.Body.Close()

	var tgResp models.TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&tgResp); err != nil {
		return fmt.Errorf("decode setWebhook response: %w", err)
	}

	if !tgResp.Ok {
		return fmt.Errorf("telegram setWebhook failed: %s", string(tgResp.Result))
	}

	s.webhookURL = webhookURL
	log.Info().Str("webhook_url", webhookURL).Msg("telegram: webhook set")
	return nil
}

func (s *TelegramService) SendMessage(chatID, text string) error {
	return s.sendMessage(chatID, text)
}

func (s *TelegramService) GetBotUsername() string {
	return s.botUsername
}
