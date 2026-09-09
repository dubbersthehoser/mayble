package app


import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/dubbersthehoser/mayble/internal/models"
)

type BookCompare func(a, b models.BookEntry) int

func CompareBookEntry(index int, ascending bool) (BookCompare, error)  {
	if !(models.IdxID <= index && models.IdxBorrower >= index) {
		return nil, fmt.Errorf("compare_books %d: invalid index", index)
	}
	return func(a, b models.BookEntry) int {
		// keep all the non-active values to the bottom of list.
		switch index {
		case models.IdxRating, models.IdxCompletedAt:
			if !a.IsCompleted && !b.IsCompleted {
				return 0
			}
			if !a.IsCompleted {
				return 1
			}
			if !b.IsCompleted {
				return -1
			}
		case models.IdxBorrower, models.IdxLoanedAt:
			if !a.IsLoaned && !b.IsLoaned {
				return 0
			}
			if !a.IsLoaned {
				return 1
			}
			if !b.IsLoaned {
				return -1
			}
		}

		r := -1
		switch index {
		case models.IdxTitle:
			r = cmp.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
		case models.IdxAuthor:
			r = cmp.Compare(strings.ToLower(a.Author), strings.ToLower(b.Author))
		case models.IdxGenre:
			r = cmp.Compare(strings.ToLower(a.Genre), strings.ToLower(b.Genre))
		case models.IdxBorrower:
			r = cmp.Compare(strings.ToLower(a.Borrower), strings.ToLower(b.Borrower))
		case models.IdxLoanedAt:
			r = a.Loaned.LoanedAt.Compare(b.LoanedAt)
		case models.IdxRating:
			r = cmp.Compare(a.Rating, b.Rating)
		case models.IdxCompletedAt:
			r = a.CompletedAt.Compare(b.CompletedAt)
		}
		if ascending {
			return r
		} else {
			return r * -1
		}
	}, nil
}


// SortBooks sort slice of book entries.
func SortBooks(books []models.BookEntry, index int, ascending bool) error {
	comp, err := CompareBookEntry(index, ascending)
	if err != nil {
		return err
	}
	slices.SortFunc(books, comp)
	return nil
}

func SortIndexsThroughBooks(idxs []int, books []models.BookEntry, index int, ascending bool) error {
	if len(idxs) != len(books) {
		return errors.New("length missmatch of indexs and books")
	}
	comp, err := CompareBookEntry(index, ascending)
	if err != nil {
		return err
	}
	slices.SortFunc(idxs, func(a, b int) int {
		return comp(books[idxs[a]], books[idxs[b]])
	})
	return nil
}


