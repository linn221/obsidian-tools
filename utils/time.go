package utils

import (
	"regexp"
	"strings"
	"time"
)

// return if date equal regardless of the timezone
func DateEqual(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func MinusDay(a, b time.Time) int {
	a = time.Date(a.Year(), a.Month(), a.Day(), 0, 0, 0, 0, time.UTC)
	b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.UTC)

	return int(a.Sub(b).Hours() / 24)
}

func ParseMyDate(s string) time.Time {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ToLower(s)

	months := map[string]string{
		"jan": "jan", "january": "jan",
		"feb": "feb", "february": "feb",
		"mar": "mar", "march": "mar",
		"apr": "apr", "april": "apr",
		"may": "may",
		"jun": "jun", "june": "jun",
		"jul": "jul", "july": "jul",
		"aug": "aug", "august": "aug",
		"sep": "sep", "sept": "sep", "september": "sep",
		"oct": "oct", "october": "oct",
		"nov": "nov", "november": "nov",
		"dec": "dec", "december": "dec",
	}

	for name, short := range months {
		re := regexp.MustCompile(`\b` + name + `\.?\b`)
		s = re.ReplaceAllString(s, short)
	}

	layouts := []struct {
		value   string
		hasYear bool
	}{
		{"2/1/06", true},
		{"2-1-06", true},
		{"2/1/2006", true},
		{"2-1-2006", true},
		{"2/1", false},
		{"2-1", false},
		{"Jan 2 2006", true},
		{"Jan 2", false},
		{"2 Jan 2006", true},
		{"2 Jan", false},
	}

	for _, l := range layouts {
		t, err := time.ParseInLocation(l.value, s, time.Local)
		if err != nil {
			continue
		}

		if !l.hasYear {
			t = time.Date(
				time.Now().Year(),
				t.Month(),
				t.Day(),
				0, 0, 0, 0,
				time.Local,
			)
		}

		return t
	}

	return time.Time{}
}
