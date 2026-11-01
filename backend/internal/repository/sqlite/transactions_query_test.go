package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeSortOrder(t *testing.T) {
	tests := []struct {
		name   string
		sortBy string
		order  string
		want   string
	}{
		{name: "valid id desc", sortBy: "id", order: "DESC", want: "id DESC"},
		{name: "valid amount asc", sortBy: "amount", order: "ASC", want: "amount ASC"},
		{name: "default order desc", sortBy: "id", order: "", want: "id DESC"},
		{name: "default sort date_time", sortBy: "", order: "ASC", want: "date_time ASC"},
		{name: "both defaults", sortBy: "", order: "", want: "date_time DESC"},
		{name: "direction case insensitive asc", sortBy: "date_time", order: "asc", want: "date_time ASC"},
		{name: "direction case insensitive mixed", sortBy: "date_time", order: "dEsC", want: "date_time DESC"},
		{name: "sort injection ignored", sortBy: "date_time; DROP TABLE transactions--", order: "DESC", want: "date_time DESC"},
		{name: "order injection ignored", sortBy: "date_time", order: "DESC; --", want: "date_time DESC"},
		{name: "unknown column falls back", sortBy: "1=1", order: "ASC", want: "date_time ASC"},
		{name: "unknown direction falls back", sortBy: "date_time", order: "DROP TABLE", want: "date_time DESC"},
		{name: "or injection ignored", sortBy: "date_time OR 1=1", order: "DESC", want: "date_time DESC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeSortOrder(tt.sortBy, tt.order)
			require.NotEmpty(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}
