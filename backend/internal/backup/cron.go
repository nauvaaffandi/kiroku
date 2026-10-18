package backup

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseCronField(field string, min, max int) (uint64, error) {
	var mask uint64
	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			return 0, fmt.Errorf("empty cron field")
		}
		step := 1
		rangePart := part
		if i := strings.Index(part, "/"); i >= 0 {
			rangePart = part[:i]
			s, err := strconv.Atoi(part[i+1:])
			if err != nil || s <= 0 {
				return 0, fmt.Errorf("invalid step %q", part)
			}
			step = s
		}
		lo, hi := min, max
		switch {
		case rangePart == "*":
		case strings.Contains(rangePart, "-"):
			bits := strings.SplitN(rangePart, "-", 2)
			a, err1 := strconv.Atoi(bits[0])
			b, err2 := strconv.Atoi(bits[1])
			if err1 != nil || err2 != nil {
				return 0, fmt.Errorf("invalid range %q", part)
			}
			lo, hi = a, b
		default:
			a, err := strconv.Atoi(rangePart)
			if err != nil {
				return 0, fmt.Errorf("invalid value %q", part)
			}
			lo, hi = a, a
		}
		if lo < min || hi > max || lo > hi {
			return 0, fmt.Errorf("out of range %q", part)
		}
		for v := lo; v <= hi; v += step {
			if v == 7 && max == 7 {
				mask |= 1 << 0
			} else {
				mask |= 1 << uint(v)
			}
		}
	}
	return mask, nil
}

type cronSpec struct {
	minute, hour, dom, month, dow uint64
}

func parseCron(schedule string) (cronSpec, error) {
	fields := strings.Fields(schedule)
	if len(fields) != 5 {
		return cronSpec{}, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}
	var spec cronSpec
	var err error
	if spec.minute, err = parseCronField(fields[0], 0, 59); err != nil {
		return spec, err
	}
	if spec.hour, err = parseCronField(fields[1], 0, 23); err != nil {
		return spec, err
	}
	if spec.dom, err = parseCronField(fields[2], 1, 31); err != nil {
		return spec, err
	}
	if spec.month, err = parseCronField(fields[3], 1, 12); err != nil {
		return spec, err
	}
	if spec.dow, err = parseCronField(fields[4], 0, 7); err != nil {
		return spec, err
	}
	return spec, nil
}

func nextCron(schedule string, from time.Time) (time.Time, error) {
	spec, err := parseCron(schedule)
	if err != nil {
		return time.Time{}, err
	}
	t := from.Truncate(time.Minute).Add(time.Minute)
	limit := from.AddDate(1, 0, 1)
	for t.Before(limit) {
		if spec.minute&(1<<uint(t.Minute())) == 0 {
			t = t.Add(time.Minute)
			continue
		}
		if spec.hour&(1<<uint(t.Hour())) == 0 {
			t = t.Add(time.Minute)
			continue
		}
		if spec.month&(1<<uint(int(t.Month()))) == 0 {
			t = t.Add(time.Minute)
			continue
		}
		if spec.dom&(1<<uint(t.Day())) == 0 {
			t = t.Add(time.Minute)
			continue
		}
		if spec.dow&(1<<uint(int(t.Weekday()))) == 0 {
			t = t.Add(time.Minute)
			continue
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("no next run within one year for %q", schedule)
}
