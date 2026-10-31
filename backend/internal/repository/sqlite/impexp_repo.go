package sqlite

import (
	"kiroku/backend/internal/db"
)

type ImpexpRepo struct{ db *db.DB }

func NewImpexpRepo(database *db.DB) *ImpexpRepo { return &ImpexpRepo{db: database} }