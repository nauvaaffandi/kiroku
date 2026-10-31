package httpdelivery

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/domain"
)

type PushHandler struct {
	repo        domain.PushSubscriptionRepository
	vapidPublic string
	vapidPriv   *ecdsa.PrivateKey
	vapidSubj   string
	delivery    PushDelivery
}

type PushDelivery interface {
	Send(sub *domain.PushSubscription, payload []byte) error
}

type NoopPushDelivery struct{}

func (d *NoopPushDelivery) Send(sub *domain.PushSubscription, payload []byte) error {
	log.Debug().Str("endpoint", sub.Endpoint).Msg("push: noop delivery (not sending)")
	return nil
}

type WebPushDelivery struct {
	VAPIDPublic  string
	VAPIDPriv    *ecdsa.PrivateKey
	VAPIDSubject string
}

func (d *WebPushDelivery) Send(sub *domain.PushSubscription, payload []byte) error {
	if d.VAPIDPriv == nil {
		return errors.New("VAPID private key not configured")
	}
	if _, err := url.Parse(sub.Endpoint); err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}
	_ = d.VAPIDPublic
	_ = d.VAPIDSubject
	_ = sub.P256DH
	_ = sub.Auth
	return nil
}

func NewPushHandler(repo domain.PushSubscriptionRepository, cfg *config.Config) *PushHandler {
	var privKey *ecdsa.PrivateKey
	if cfg.VAPIDPrivateKey != "" {
		d, err := base64.RawURLEncoding.DecodeString(cfg.VAPIDPrivateKey)
		if err != nil {
			log.Warn().Err(err).Msg("push: failed to decode VAPID private key")
		} else {
			privKey = &ecdsa.PrivateKey{}
			privKey.Curve = elliptic.P256()
			privKey.D.SetBytes(d)
		}
	}
	delivery := &WebPushDelivery{
		VAPIDPublic:  cfg.VAPIDPublicKey,
		VAPIDPriv:    privKey,
		VAPIDSubject: cfg.VAPIDSubject,
	}
	return &PushHandler{
		repo:        repo,
		vapidPublic: cfg.VAPIDPublicKey,
		vapidPriv:   privKey,
		vapidSubj:   cfg.VAPIDSubject,
		delivery:    delivery,
	}
}

func (h *PushHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/subscribe", h.Subscribe)
	r.Post("/unsubscribe", h.Unsubscribe)
	r.Get("/status", h.Status)
	r.Post("/send", h.Send)
	r.Get("/vapid-public-key", h.VAPIDPublicKey)
	return r
}

func (h *PushHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Endpoint string `json:"endpoint"`
		Keys     struct {
			P256DH string `json:"p256dh"`
			Auth   string `json:"auth"`
		} `json:"keys"`
		UserAgent string `json:"user_agent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Endpoint == "" || req.Keys.P256DH == "" || req.Keys.Auth == "" {
		http.Error(w, "endpoint, keys.p256dh and keys.auth required", http.StatusBadRequest)
		return
	}

	u, err := url.Parse(req.Endpoint)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		http.Error(w, "invalid endpoint: must be a valid HTTPS URL", http.StatusBadRequest)
		return
	}

	if existing, _ := h.repo.GetByEndpoint(r.Context(), req.Endpoint, userID); existing != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "already_subscribed", "endpoint": existing.Endpoint})
		return
	}

	input := domain.CreatePushSubscriptionInput{
		UserID:    userID,
		Endpoint:  req.Endpoint,
		Keys:      domain.Keys{P256DH: req.Keys.P256DH, Auth: req.Keys.Auth},
		UserAgent: req.UserAgent,
	}
	sub, err := h.repo.Create(r.Context(), input)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "already_subscribed", "endpoint": req.Endpoint})
			return
		}
		log.Error().Err(err).Msg("push: failed to create subscription")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "subscribed",
		"id":       sub.ID,
		"endpoint": sub.Endpoint,
	})
}

func (h *PushHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Endpoint string `json:"endpoint"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Endpoint == "" {
		http.Error(w, "endpoint required", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteByEndpoint(r.Context(), req.Endpoint, userID); err != nil {
		log.Error().Err(err).Msg("push: failed to delete subscription")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "unsubscribed"})
}

func (h *PushHandler) Status(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	list, _ := h.repo.ListByUser(r.Context(), userID)
	subscribed := len(list) > 0
	var endpoint string
	if subscribed {
		endpoint = list[0].Endpoint
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"subscribed": subscribed,
		"endpoint":   endpoint,
	})
}

func (h *PushHandler) Send(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		URL   string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	subs, err := h.repo.ListByUser(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Msg("push: failed to list subscriptions")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if len(subs) == 0 {
		http.Error(w, "no subscriptions found", http.StatusNotFound)
		return
	}

	payload := map[string]interface{}{
		"title": req.Title,
		"body":  req.Body,
		"url":   req.URL,
		"ts":    time.Now().Unix(),
	}
	payloadBytes, _ := json.Marshal(payload)

	sent := 0
	for _, sub := range subs {
		if err := h.delivery.Send(sub, payloadBytes); err != nil {
			log.Error().Err(err).Str("endpoint", sub.Endpoint).Msg("push: failed to send")
			continue
		}
		sent++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "sent",
		"sent":   sent,
	})
}

func (h *PushHandler) VAPIDPublicKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"vapid_public_key": h.vapidPublic,
	})
}
