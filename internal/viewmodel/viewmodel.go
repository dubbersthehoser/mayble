package viewmodel

import (
	"time"
	"slices"

	"fyne.io/fyne/v2/data/binding"
)

const dateFormat = "02/01/2006"

func formatDate(t *time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(dateFormat)
}

func parseDate(t string) (*time.Time, error) {
	ret, err := time.Parse(dateFormat, t)
	return &ret, err
}

const capRating = 6

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

func Ratings() []string {
	r := make([]string, capRating)
	for i := range capRating {
		r[i] = formatRating(i)
	}
	return r
}

type listener struct {
	listeners []binding.DataListener
}

func (t *listener) notify() {
	for _, listener := range t.listeners {
		listener.DataChanged()
	}
}

func (t *listener) AddListener(l binding.DataListener) {
	if t.listeners == nil {
		t.listeners = make([]binding.DataListener, 0)
	}
	t.listeners = append(t.listeners, l)
}

func (t *listener) RemoveListener(l binding.DataListener) {
	if t.listeners == nil {
		return
	}
	t.listeners = slices.DeleteFunc(t.listeners, func(ll binding.DataListener) bool {
		return l == ll
	})
}
