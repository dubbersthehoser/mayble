package display

import (
	"time"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/models"
)

const DateFormat = "02/01/2006"
const MaxRating = 6

func FormatDate(t *time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(DateFormat)
}

func Ratings() []string {
	r := make([]string, MaxRating)
	for i := range MaxRating {
		r[i] = FormatRating(i)
	}
	return r
}

func FormatRating(r int) string {
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

func ParseRating(r string) (int, error) {
	idx := slices.Index(Ratings(), r)
	if idx == -1 {
		return 0, fmt.Errorf("invalid rating string '%s'", r)
	}
	return idx, nil
}

func EntryValues(e *models.BookEntry) []string {

	headers := models.BookEntryFields()
	values := make([]string, len(headers))

	for i, header := range headers {
		switch i {
		case models.IdxID:
			values[i] = fmt.Sprintf("%d", e.ID)
		case models.IdxTitle:
			values[i] = e.Title
		case models.IdxAuthor:
			values[i] = e.Author
		case models.IdxGenre:
			values[i] = e.Genre
		case models.IdxRating:
			values[i] = FormatRating(e.Rating)
		case models.IdxCompletedAt:
			values[i] = FormatDate(&e.CompletedAt)
		case models.IdxBorrower:
			values[i] = e.Borrower
		case models.IdxLoanedAt:
			values[i] = FormatDate(&e.LoanedAt)
		default:
			panic("unknown field:" + header)
		}
	}
	return values
}
