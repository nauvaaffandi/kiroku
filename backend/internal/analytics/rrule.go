package analytics

import (
	"fmt"
	"time"
)

func (h *Handler) expandRRULE(rruleStr string, dateFrom, dateTo time.Time) []time.Time {
	if rruleStr == "" {
		return []time.Time{}
	}

	var dates []time.Time
	parts := splitString(rruleStr, ";")
	freq := ""
	interval := 1
	byDay := []string{}

	for _, part := range parts {
		kv := splitString(part, "=")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "FREQ":
			freq = kv[1]
		case "INTERVAL":
			var val int
			fmt.Sscanf(kv[1], "%d", &val)
			interval = val
		case "BYDAY":
			byDay = splitString(kv[1], ",")
		}
	}

	current := dateFrom
	switch freq {
	case "DAILY":
		for !current.After(dateTo) {
			dates = append(dates, current)
			current = current.AddDate(0, 0, interval)
		}
	case "WEEKLY":
		weekDays := map[string]time.Weekday{
			"MO": time.Monday, "TU": time.Tuesday, "WE": time.Wednesday,
			"TH": time.Thursday, "FR": time.Friday, "SA": time.Saturday, "SU": time.Sunday,
		}
		for !current.After(dateTo) {
			for _, bd := range byDay {
				if wd, ok := weekDays[bd]; ok {
					offset := int(wd - current.Weekday())
					if offset < 0 {
						offset += 7
					}
					nextDate := current.AddDate(0, 0, offset)
					if !nextDate.After(dateTo) {
						dates = append(dates, nextDate)
					}
				}
			}
			current = current.AddDate(0, 0, 7*interval)
		}
	case "MONTHLY":
		for !current.After(dateTo) {
			dates = append(dates, current)
			current = current.AddDate(0, interval, 0)
		}
	case "YEARLY":
		for !current.After(dateTo) {
			dates = append(dates, current)
			current = current.AddDate(interval, 0, 0)
		}
	}
	return dates
}

func (h *Handler) getFxRate(from, to string) float64 {
	if from == to {
		return 1.0
	}
	var rate float64
	err := h.db.QueryRow("SELECT rate FROM fx_rates WHERE base = ? AND quote = ? ORDER BY date DESC LIMIT 1", from, to).Scan(&rate)
	if err != nil {
		return 1.0
	}
	return rate
}
