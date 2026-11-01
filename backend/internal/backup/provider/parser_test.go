package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProviders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty string", "", nil},
		{"whitespace string", "   ", nil},
		{"single provider", "git", []string{"git"}},
		{"multiple providers", "git:gdrive", []string{"git", "gdrive"}},
		{"with spaces and casing", "  Git : GDrive : git ", []string{"git", "gdrive"}},
		{"with empty segments", "git::gdrive:", []string{"git", "gdrive"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseProviders(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
