package backup

import (
	"strconv"
	"time"
)

func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func nullFloatOrNil(s string) interface{} {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return v
}

func timeOrNow(s string) string {
	if s == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	return s
}
