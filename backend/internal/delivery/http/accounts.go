package httpdelivery

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type AccountsHandler struct {
	repo domain.AccountRepository
}

func NewAccountsHandler(repo domain.AccountRepository) *AccountsHandler {
	return &AccountsHandler{repo: repo}
}

func (h *AccountsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/summary", h.Summary)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *AccountsHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *AccountsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateAccountInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	acc, err := h.repo.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": acc.ID})
}

func (h *AccountsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	acc, err := h.repo.Get(r.Context(), id)
	if err != nil || acc == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateAccountInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Currency != nil {
		count, _ := h.repo.CountTransactions(r.Context(), id)
		if count > 0 {
			http.Error(w, "cannot change currency: account has transactions", http.StatusBadRequest)
			return
		}
	}
	_, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *AccountsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *AccountsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	sum, err := h.repo.Summary(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	m := map[string]float64{}
	for _, s := range sum {
		m[s.Currency] = s.Balance
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}
