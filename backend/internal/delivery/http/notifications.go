package httpdelivery

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/domain"
	"kiroku/backend/internal/notification"
)

// NotificationsHandler exposes the notification inbox REST API.
type NotificationsHandler struct {
	service *notification.Service
}

func NewNotificationsHandler(service *notification.Service) *NotificationsHandler {
	return &NotificationsHandler{service: service}
}

func (h *NotificationsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/unread-count", h.UnreadCount)
	r.Patch("/read-all", h.MarkAllRead)
	r.Delete("/read", h.DeleteRead)
	r.Patch("/{id}/read", h.MarkRead)
	r.Delete("/{id}", h.Delete)
	return r
}

// List returns a paginated list of notifications for the authenticated user.
func (h *NotificationsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	unreadOnly := r.URL.Query().Get("unread") == "true" || r.URL.Query().Get("unread") == "1"

	notifs, err := h.service.List(r.Context(), userID, domain.ListFilter{UnreadOnly: unreadOnly, Limit: limit, Offset: offset})
	if err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "ListNotifications").Str("path", r.URL.Path).Msg("notifications/list failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}

// Create persists a new notification for the authenticated user. Intended for other
// backend modules to integrate with, without coupling to the notification implementation.
func (h *NotificationsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req domain.CreateNotificationInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.RecipientID == "" {
		req.RecipientID = userID
	}
	if req.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	notif, err := h.service.Create(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "CreateNotification").Str("path", r.URL.Path).Msg("notifications/create failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notif)
}

// UnreadCount returns the count of unread notifications for the authenticated user.
func (h *NotificationsHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	count, err := h.service.UnreadCount(r.Context(), userID)
	if err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "UnreadNotificationCount").Str("path", r.URL.Path).Msg("notifications/unread-count failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"unread_count": count})
}

// MarkRead marks a single notification read.
func (h *NotificationsHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	notif, err := h.service.MarkRead(r.Context(), id, userID)
	if err != nil || notif == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// MarkAllRead marks every unread notification as read.
func (h *NotificationsHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if _, err := h.service.MarkAllRead(r.Context(), userID); err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "MarkAllNotificationsRead").Str("path", r.URL.Path).Msg("notifications/mark-all-read failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Delete soft-deletes a single notification owned by the authenticated user and cancels pending deliveries.
func (h *NotificationsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Error().Err(err).Str("component", "http").Str("method", "DeleteNotification").Str("path", r.URL.Path).Str("id", id).Msg("notifications/delete failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// DeleteRead removes all read notifications for the authenticated user.
func (h *NotificationsHandler) DeleteRead(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if _, err := h.service.DeleteRead(r.Context(), userID); err != nil {
		log.Error().Err(err).Str("component", "http").Str("method", "DeleteReadNotifications").Str("path", r.URL.Path).Msg("notifications/delete-read failed")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
