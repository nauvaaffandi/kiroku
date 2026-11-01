package httpdelivery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"kiroku/backend/internal/auth"
	"kiroku/backend/internal/config"
	"kiroku/backend/internal/domain"
)

type mockPushRepo struct {
	subscriptions map[string]*domain.PushSubscription
}

func newMockPushRepo() *mockPushRepo {
	return &mockPushRepo{
		subscriptions: map[string]*domain.PushSubscription{},
	}
}

func (m *mockPushRepo) Create(ctx context.Context, input domain.CreatePushSubscriptionInput) (*domain.PushSubscription, error) {
	if _, exists := m.subscriptions[input.Endpoint]; exists {
		return nil, nil // duplicate
	}
	id := "push_123"
	now := time.Now()
	sub := &domain.PushSubscription{
		ID:          id,
		UserID:      input.UserID,
		Endpoint:    input.Endpoint,
		P256DH:      input.Keys.P256DH,
		Auth:        input.Keys.Auth,
		UserAgent:   input.UserAgent,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	m.subscriptions[input.Endpoint] = sub
	return sub, nil
}

func (m *mockPushRepo) Get(ctx context.Context, id string) (*domain.PushSubscription, error) {
	for _, sub := range m.subscriptions {
		if sub.ID == id {
			return sub, nil
		}
	}
	return nil, nil
}

func (m *mockPushRepo) GetByEndpoint(ctx context.Context, endpoint, userID string) (*domain.PushSubscription, error) {
	sub, ok := m.subscriptions[endpoint]
	if !ok || sub.UserID != userID {
		return nil, nil
	}
	return sub, nil
}

func (m *mockPushRepo) ListByUser(ctx context.Context, userID string) ([]*domain.PushSubscription, error) {
	var list []*domain.PushSubscription
	for _, sub := range m.subscriptions {
		if sub.UserID == userID {
			list = append(list, sub)
		}
	}
	if list == nil {
		list = []*domain.PushSubscription{}
	}
	return list, nil
}

func (m *mockPushRepo) Delete(ctx context.Context, id string) error {
	for ep, sub := range m.subscriptions {
		if sub.ID == id {
			delete(m.subscriptions, ep)
			return nil
		}
	}
	return nil
}

func (m *mockPushRepo) DeleteByEndpoint(ctx context.Context, endpoint, userID string) error {
	sub, ok := m.subscriptions[endpoint]
	if !ok || sub.UserID != userID {
		return nil
	}
	delete(m.subscriptions, endpoint)
	return nil
}

func (m *mockPushRepo) DeleteByUser(ctx context.Context, userID string) error {
	for ep, sub := range m.subscriptions {
		if sub.UserID == userID {
			delete(m.subscriptions, ep)
		}
	}
	return nil
}

func (m *mockPushRepo) Update(ctx context.Context, id string, input domain.UpdatePushSubscriptionInput) (*domain.PushSubscription, error) {
	return nil, nil
}

type noopPushDelivery struct{}

func (d *noopPushDelivery) Send(sub *domain.PushSubscription, payload []byte) error {
	return nil
}

func TestPushSubscribe(t *testing.T) {
	repo := newMockPushRepo()
	cfg := &config.Config{}
	handler := NewPushHandler(repo, cfg)

	reqBody := `{"endpoint":"https://example.com/push/abc","keys":{"p256dh":"base64p256dh","auth":"base64auth"},"user_agent":"test-agent"}`
	req := httptest.NewRequest("POST", "/push/subscribe", strings.NewReader(reqBody))
	req = auth.WithUserID(req, "user_1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Subscribe(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "subscribed" {
		t.Errorf("expected status 'subscribed', got %v", resp["status"])
	}

	// Test duplicate subscription
	w2 := httptest.NewRecorder()
	reqBody2 := `{"endpoint":"https://example.com/push/abc","keys":{"p256dh":"base64p256dh","auth":"base64auth"},"user_agent":"test-agent"}`
	req2 := httptest.NewRequest("POST", "/push/subscribe", strings.NewReader(reqBody2))
	req2 = auth.WithUserID(req2, "user_1")
	req2.Header.Set("Content-Type", "application/json")
	handler.Subscribe(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("expected status 200 for duplicate, got %d", w2.Code)
	}
}

func TestPushUnsubscribe(t *testing.T) {
	repo := newMockPushRepo()
	// Pre-create a subscription
	repo.Create(context.Background(), domain.CreatePushSubscriptionInput{
		UserID:    "user_1",
		Endpoint:  "https://example.com/push/abc",
		Keys:      domain.Keys{P256DH: "p256dh", Auth: "auth"},
		UserAgent: "test-agent",
	})

	cfg := &config.Config{}
	handler := NewPushHandler(repo, cfg)

	reqBody := `{"endpoint":"https://example.com/push/abc"}`
	req := httptest.NewRequest("POST", "/push/unsubscribe", strings.NewReader(reqBody))
	req = auth.WithUserID(req, "user_1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Unsubscribe(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "unsubscribed" {
		t.Errorf("expected status 'unsubscribed', got %v", resp["status"])
	}
}

func TestPushStatus(t *testing.T) {
	repo := newMockPushRepo()
	cfg := &config.Config{}
	handler := NewPushHandler(repo, cfg)

	// No subscription yet
	req := httptest.NewRequest("GET", "/push/status", nil)
	req = auth.WithUserID(req, "user_1")
	w := httptest.NewRecorder()

	handler.Status(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["subscribed"] != false {
		t.Errorf("expected subscribed=false, got %v", resp["subscribed"])
	}

	// Create a subscription
	repo.Create(context.Background(), domain.CreatePushSubscriptionInput{
		UserID:    "user_1",
		Endpoint:  "https://example.com/push/abc",
		Keys:      domain.Keys{P256DH: "p256dh", Auth: "auth"},
		UserAgent: "test-agent",
	})

	w2 := httptest.NewRecorder()
	handler.Status(w2, auth.WithUserID(req, "user_1"))
	json.Unmarshal(w2.Body.Bytes(), &resp)
	if resp["subscribed"] != true {
		t.Errorf("expected subscribed=true, got %v", resp["subscribed"])
	}
	if resp["endpoint"] != "https://example.com/push/abc" {
		t.Errorf("expected endpoint match, got %v", resp["endpoint"])
	}
}

func TestPushSend(t *testing.T) {
	repo := newMockPushRepo()
	repo.Create(context.Background(), domain.CreatePushSubscriptionInput{
		UserID:    "user_1",
		Endpoint:  "https://example.com/push/abc",
		Keys:      domain.Keys{P256DH: "p256dh", Auth: "auth"},
		UserAgent: "test-agent",
	})

	cfg := &config.Config{}
	handler := NewPushHandler(repo, cfg)
	handler.delivery = &noopPushDelivery{}

	reqBody := `{"title":"Test","body":"Hello push","url":"/transactions"}`
	req := httptest.NewRequest("POST", "/push/send", strings.NewReader(reqBody))
	req = auth.WithUserID(req, "user_1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Send(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "sent" {
		t.Errorf("expected status 'sent', got %v", resp["status"])
	}
}

func TestPushSendNoSubscription(t *testing.T) {
	repo := newMockPushRepo()
	cfg := &config.Config{}
	handler := NewPushHandler(repo, cfg)

	reqBody := `{"title":"Test"}`
	req := httptest.NewRequest("POST", "/push/send", strings.NewReader(reqBody))
	req = auth.WithUserID(req, "user_1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Send(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestPushSubscribeInvalidEndpoint(t *testing.T) {
	repo := newMockPushRepo()
	cfg := &config.Config{}
	handler := NewPushHandler(repo, cfg)

	reqBody := `{"endpoint":"http://insecure.com/push","keys":{"p256dh":"abc","auth":"xyz"}}`
	req := httptest.NewRequest("POST", "/push/subscribe", strings.NewReader(reqBody))
	req = auth.WithUserID(req, "user_1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Subscribe(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for non-https endpoint, got %d", w.Code)
	}
}

func TestPushUnsubscribeIDOR(t *testing.T) {
	repo := newMockPushRepo()
	repo.Create(context.Background(), domain.CreatePushSubscriptionInput{
		UserID:    "user_2",
		Endpoint:  "https://example.com/push/abc",
		Keys:      domain.Keys{P256DH: "p256dh", Auth: "auth"},
		UserAgent: "test-agent",
	})

	cfg := &config.Config{}
	handler := NewPushHandler(repo, cfg)

	// user_1 tries to unsubscribe user_2's endpoint
	reqBody := `{"endpoint":"https://example.com/push/abc"}`
	req := httptest.NewRequest("POST", "/push/unsubscribe", strings.NewReader(reqBody))
	req = auth.WithUserID(req, "user_2") // wait, test user_1 trying to delete user_2's sub
	req = auth.WithUserID(req, "user_1")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Unsubscribe(w, req)

	// Verify user_2's subscription still exists in repo
	sub, _ := repo.GetByEndpoint(context.Background(), "https://example.com/push/abc", "user_2")
	if sub == nil {
		t.Errorf("expected user_2 subscription to remain intact after user_1 attempted to unsubscribe it")
	}
}
