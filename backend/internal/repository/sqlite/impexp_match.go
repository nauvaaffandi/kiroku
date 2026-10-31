package sqlite

import (
	"strings"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

func (r *ImpexpRepo) fuzzyMatchCategory(input string) domain.CategoryMatch {
	input = strings.TrimSpace(input)
	if input == "" {
		return domain.CategoryMatch{Input: input}
	}
	best := domain.CategoryMatch{Input: input, Distance: 100}
	rows, err := r.db.Query(`SELECT id, name FROM categories WHERE archived = 0`)
	if err != nil {
		return best
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		if rows.Scan(&id, &name) != nil {
			continue
		}
		dist := levenshtein(input, name)
		if dist < best.Distance {
			best = domain.CategoryMatch{Input: input, Matched: name, MatchID: id, Distance: dist}
		}
	}
	if best.Distance <= 2 {
		return best
	}
	id := genID()
	_, _ = r.db.Exec(`INSERT INTO categories (id, name, color, sort_order) VALUES (?, ?, ?, 0)`, id, input, randomPastel())
	return domain.CategoryMatch{Input: input, Matched: input, MatchID: id, Distance: 0}
}

func (r *ImpexpRepo) fuzzyMatchAccount(input string) domain.AccountMatch {
	input = strings.TrimSpace(input)
	if input == "" {
		return domain.AccountMatch{Input: input}
	}
	best := domain.AccountMatch{Input: input, Distance: 100}
	rows, err := r.db.Query(`SELECT id, name FROM accounts WHERE archived = 0`)
	if err != nil {
		return best
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		if rows.Scan(&id, &name) != nil {
			continue
		}
		dist := levenshtein(input, name)
		if dist < best.Distance {
			best = domain.AccountMatch{Input: input, Matched: name, MatchID: id, Distance: dist}
		}
	}
	if best.Distance <= 2 {
		return best
	}
	id := genID()
	_, _ = r.db.Exec(`INSERT INTO accounts (id, name, type, currency, sort_order) VALUES (?, ?, 'cash', 'JPY', 0)`, id, input)
	return domain.AccountMatch{Input: input, Matched: input, MatchID: id, Distance: 0}
}

func resolveAccountByName(database *db.DB, name string) *string {
	if name == "" {
		return nil
	}
	var id string
	err := database.QueryRow(`SELECT id FROM accounts WHERE name = ? AND archived = 0`, name).Scan(&id)
	if err != nil {
		return nil
	}
	return &id
}

func resolveCategoryByName(database *db.DB, name string) *string {
	if name == "" {
		return nil
	}
	var id string
	err := database.QueryRow(`SELECT id FROM categories WHERE name = ? AND archived = 0`, name).Scan(&id)
	if err != nil {
		return nil
	}
	return &id
}