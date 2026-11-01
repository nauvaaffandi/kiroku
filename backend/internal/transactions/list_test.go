package transactions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SortOrderSuite struct {
	suite.Suite
}

func TestSortOrderSuite(t *testing.T) {
	suite.Run(t, new(SortOrderSuite))
}

func (s *SortOrderSuite) TestValidColumns() {
	tests := []struct {
		name   string
		sortBy string
		order  string
		want   string
	}{
		{"default empty", "", "", "date_time DESC"},
		{"valid id asc", "id", "asc", "id ASC"},
		{"valid amount desc", "amount", "DESC", "amount DESC"},
		{"valid created_at", "created_at", "desc", "created_at DESC"},
		{"case-insensitive order", "date_time", "Desc", "date_time DESC"},
		{"all valid columns", "updated_at", "ASC", "updated_at ASC"},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := sanitizeSortOrder(tc.sortBy, tc.order)
			assert.Equal(s.T(), tc.want, got)
		})
	}
}

func (s *SortOrderSuite) TestInvalidColumnFallsBack() {
	tests := []struct {
		name   string
		sortBy string
		order  string
		want   string
	}{
		{"sql injection drop table", "1; DROP TABLE transactions;--", "DESC", "date_time DESC"},
		{"column not in whitelist", "password", "ASC", "date_time ASC"},
		{"empty column with valid order", "", "ASC", "date_time ASC"},
		{"union select", "id UNION SELECT 1,2,3", "DESC", "date_time DESC"},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := sanitizeSortOrder(tc.sortBy, tc.order)
			assert.Equal(s.T(), tc.want, got)
		})
	}
}

func (s *SortOrderSuite) TestInvalidDirectionFallsBack() {
	tests := []struct {
		name   string
		sortBy string
		order  string
		want   string
	}{
		{"sql injection in order", "date_time", "ASC; DROP TABLE users;--", "date_time DESC"},
		{"invalid direction string", "amount", "ascending", "amount DESC"},
		{"empty order defaults to DESC", "id", "", "id DESC"},
		{"mixed case invalid", "currency", "Up", "currency DESC"},
	}
	for _, tc := range tests {
		s.Run(tc.name, func() {
			got := sanitizeSortOrder(tc.sortBy, tc.order)
			assert.Equal(s.T(), tc.want, got)
		})
	}
}

func (s *SortOrderSuite) TestBothInvalid() {
	got := sanitizeSortOrder("evil", "worse")
	assert.Equal(s.T(), "date_time DESC", got)
}

func (s *SortOrderSuite) TestNoSQLInjectionInOutput() {
	// Output must never contain semicolons, spaces beyond column+direction, or quotes.
	malicious := []string{
		"1; DROP TABLE transactions;--",
		"id; DELETE FROM transactions",
		"amount ASC; SELECT * FROM users",
	}
	for _, m := range malicious {
		got := sanitizeSortOrder(m, "ASC")
		assert.NotContains(s.T(), got, ";", "output should not contain semicolon for input: %s", m)
		assert.NotContains(s.T(), got, "DROP", "output should not contain DROP for input: %s", m)
		assert.NotContains(s.T(), got, "SELECT", "output should not contain SELECT for input: %s", m)
	}
}