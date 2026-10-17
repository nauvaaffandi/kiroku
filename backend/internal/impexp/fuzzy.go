package impexp

import (
	"strings"

	"github.com/google/uuid"

	"kiroku/backend/internal/models"
)

func (h *Handler) fuzzyMatchCategory(input string, createMissing bool) models.CategoryMatch {
	input = strings.TrimSpace(input)
	if input == "" {
		return models.CategoryMatch{Input: input}
	}

	var bestMatch models.CategoryMatch
	bestMatch.Input = input
	bestMatch.Distance = 100

	rows, err := h.db.Query("SELECT id, name FROM categories WHERE archived = 0")
	if err != nil {
		return bestMatch
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		dist := levenshtein(input, name)
		if dist < bestMatch.Distance {
			bestMatch = models.CategoryMatch{
				Input:    input,
				Matched:  name,
				MatchID:  id,
				Distance: dist,
			}
		}
	}

	if bestMatch.Distance <= 2 {
		return bestMatch
	}

	if createMissing {
		id := uuid.New().String()
		h.db.Exec("INSERT INTO categories (id, name, color, sort_order) VALUES (?, ?, ?, 0)", id, input, randomPastel())
		return models.CategoryMatch{Input: input, Matched: input, MatchID: id, Distance: 0, Created: true}
	}

	return models.CategoryMatch{Input: input, Distance: bestMatch.Distance}
}

func (h *Handler) fuzzyMatchAccount(input string, createMissing bool) models.AccountMatch {
	input = strings.TrimSpace(input)
	if input == "" {
		return models.AccountMatch{Input: input}
	}

	var bestMatch models.AccountMatch
	bestMatch.Input = input
	bestMatch.Distance = 100

	rows, err := h.db.Query("SELECT id, name FROM accounts WHERE archived = 0")
	if err != nil {
		return bestMatch
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		rows.Scan(&id, &name)
		dist := levenshtein(input, name)
		if dist < bestMatch.Distance {
			bestMatch = models.AccountMatch{
				Input:    input,
				Matched:  name,
				MatchID:  id,
				Distance: dist,
			}
		}
	}

	if bestMatch.Distance <= 2 {
		return bestMatch
	}

	if createMissing {
		id := uuid.New().String()
		h.db.Exec("INSERT INTO accounts (id, name, type, currency, sort_order) VALUES (?, ?, 'cash', 'JPY', 0)", id, input)
		return models.AccountMatch{Input: input, Matched: input, MatchID: id, Distance: 0, Created: true}
	}

	return models.AccountMatch{Input: input, Distance: bestMatch.Distance}
}

func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
		dp[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		dp[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dp[i][j] = min(
				dp[i-1][j]+1,
				dp[i][j-1]+1,
				dp[i-1][j-1]+cost,
			)
		}
	}
	return dp[len(a)][len(b)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
