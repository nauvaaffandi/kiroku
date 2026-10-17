package categories

import (
	"github.com/google/uuid"
)

func join(slice []string, sep string) string {
	result := ""
	for i, s := range slice {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func randomPastel() string {
	// Generate a pastel color: high lightness, low saturation
	// Pastel colors in HSL: h=0-360, s=30-50%, l=70-85%
	// Convert to hex
	// Simple implementation: random hue, fixed saturation/lightness
	return "#" + randomHex(3)
}

func randomHex(n int) string {
	const letters = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[uuid.New().ClockSequence()%16]
	}
	return string(b)
}
