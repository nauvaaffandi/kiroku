package sqlite

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/teambition/rrule-go"
)

func genID() string {
	return uuid.New().String()
}

func joinStrings(slice []string, sep string) string {
	return strings.Join(slice, sep)
}

func randomPastel() string {
	return "#" + fmt.Sprintf("%06x", uuid.New().ClockSequence()%0xFFFFFF)
}

func computeNextRun(rruleStr string) *time.Time {
	rule, err := rrule.StrToRRule(rruleStr)
	if err != nil {
		return nil
	}
	next := rule.After(time.Now(), true)
	if next.IsZero() {
		return nil
	}
	return &next
}

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
