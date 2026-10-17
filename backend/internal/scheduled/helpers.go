package scheduled

import (
	"time"

	"github.com/teambition/rrule-go"
)

func computeNextRun(rruleStr string, after time.Time) *time.Time {
	rule, err := rrule.StrToRRule(rruleStr)
	if err != nil {
		return nil
	}
	next := rule.After(after, true)
	if next.IsZero() {
		return nil
	}
	return &next
}

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
