package httpdelivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"kiroku/backend/internal/domain"
)

type TransactionsHandler struct {
	repo domain.TransactionRepository
}

func NewTransactionsHandler(repo domain.TransactionRepository) *TransactionsHandler {
	return &TransactionsHandler{repo: repo}
}

func (h *TransactionsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Post("/batch", h.Batch)
	r.Get("/search", h.Search)
	r.Get("/autocomplete", h.Autocomplete)
	r.Get("/summary", h.Summary)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

func (h *TransactionsHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := domain.TransactionFilter{
		SearchQuery: q.Get("q"),
		Sort:        q.Get("sort"),
		Order:       q.Get("order"),
		Page:        page,
		Limit:       limit,
	}
	if aid := q.Get("account_id"); aid != "" {
		filter.AccountID = &aid
	}
	if cid := q.Get("category_id"); cid != "" {
		filter.CategoryID = &cid
	}
	if typ := q.Get("type"); typ != "" {
		filter.Types = []string{typ}
	}
	if stat := q.Get("status"); stat != "" {
		filter.Status = &stat
	}
	if df := q.Get("date_from"); df != "" {
		if t, err := time.Parse(time.RFC3339, df); err == nil {
			filter.DateFrom = &t
		}
	}
	if dt := q.Get("date_to"); dt != "" {
		if t, err := time.Parse(time.RFC3339, dt); err == nil {
			filter.DateTo = &t
		}
	}

	list, err := h.repo.List(r.Context(), filter)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *TransactionsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTransactionInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	txn, err := h.repo.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": txn.ID})
}

func (h *TransactionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	txn, err := h.repo.Get(r.Context(), id)
	if err != nil || txn == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txn)
}

func (h *TransactionsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req domain.UpdateTransactionInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *TransactionsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *TransactionsHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	list, err := h.repo.Search(r.Context(), q, 50)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *TransactionsHandler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	list, err := h.repo.Autocomplete(r.Context(), q, 10)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *TransactionsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := domain.TransactionFilter{}
	if cid := q.Get("category_id"); cid != "" {
		filter.CategoryID = &cid
	}
	if aid := q.Get("account_id"); aid != "" {
		filter.AccountID = &aid
	}
	if typ := q.Get("type"); typ != "" {
		filter.Types = []string{typ}
	}
	sum, err := h.repo.Summary(r.Context(), filter)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sum)
}

func (h *TransactionsHandler) Batch(w http.ResponseWriter, r *http.Request) {
	var req domain.BatchInput
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	affected, err := h.repo.Batch(r.Context(), req)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"affected": affected, "errors": []string{}})
}
