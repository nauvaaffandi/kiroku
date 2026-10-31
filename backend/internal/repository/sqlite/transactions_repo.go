package sqlite

import (
	"kiroku/backend/internal/db"
)

type TransactionsRepo struct{ db *db.DB }

func NewTransactionsRepo(database *db.DB) *TransactionsRepo { return &TransactionsRepo{db: database} }