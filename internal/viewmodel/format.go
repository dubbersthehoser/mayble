package viewmodel

import (
	"fmt"
	"slices"
	"time"
)

const dateFormat = "02/01/2006"

const capRating = 6

func Ratings() []string {
	r := make([]string, capRating)
	for i := range capRating {
		r[i] = formatRating(i)
	}
	return r
}

func formatDate(t *time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(dateFormat)
}

func formatRating(r int) string {
	switch r {
	case 0:
		return ""
	case 1:
		return "⭐"
	case 2:
		return "⭐⭐"
	case 3:
		return "⭐⭐⭐"
	case 4:
		return "⭐⭐⭐⭐"
	case 5:
		return "⭐⭐⭐⭐⭐"
	default:
		return "ERROR"
	}
}

func parseDate(t string) (*time.Time, error) {
	ret, err := time.Parse(dateFormat, t)
	return &ret, err
}

func parseRating(r string) (int, error) {
	idx := slices.Index(Ratings(), r)
	if idx == -1 {
		return 0, fmt.Errorf("invalid rating string '%s'", r)
	}
	return idx, nil
}
