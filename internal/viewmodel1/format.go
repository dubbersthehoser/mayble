package viewmodel

import (
	"time"
)


const dateFormat = "02/01/2006"

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

