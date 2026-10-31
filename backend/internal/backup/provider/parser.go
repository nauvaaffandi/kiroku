package provider

import (
	"strings"
)

// ParseProviders parses a colon-separated provider string (e.g. "git:gdrive")
// into a normalized, deduplicated slice of provider names.
// Empty string returns an empty slice (local-only mode).
func ParseProviders(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ":")
	seen := make(map[string]bool)
	var result []string
	for _, p := range parts {
		name := strings.ToLower(strings.TrimSpace(p))
		if name != "" && !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}
	return result
}
