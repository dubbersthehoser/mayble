// Package util contains helper functions for creating importing exporting functions.
//
package util

import (
	"time"
	"strconv"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/app"
)

const BookLoanFieldCount int = 6

const (
	TitleIndex int = 0
	AuthorIndex    = 1
	GenreIndex     = 2
	RattingIndex   = 3
	BorrowerIndex  = 4
	DateIndex      = 5
)

func BookLoanToFields(book app.BookLoan) ([]string, error) {
	fields := make([]string, BookLoanFieldCount)
	fields[TitleIndex] = book.Title
	fields[AuthorIndex] = book.Author
	fields[GenreIndex] = book.Genre

	ratting := strconv.Itoa(book.Ratting)
	fields[RattingIndex] = ratting

	if book.IsOnLoan {
		fields[BorrowerIndex] = book.Borrower
		fields[DateIndex] = book.Date.Format(time.DateOnly)
	}
	return fields, nil
}

func BookLoanFromFields(fields []string) (*app.BookLoan, error) {
	if len(fields) != BookLoanFieldCount {
		return nil, errors.New("invalid number of fields")
	}

	book := &app.BookLoan{ID: app.ZeroID}

	book.Title = fields[TitleIndex]
	book.Author = fields[AuthorIndex]
	book.Genre = fields[GenreIndex]
	ratting, err := strconv.Atoi(fields[RattingIndex])
	if err != nil {
		return nil, errors.New("failed to parse ratting field")
	}
	book.Ratting = ratting

	var (
		HasBorrower bool = fields[BorrowerIndex] != ""
		HasDate     bool = fields[DateIndex] != ""
	)
	if HasBorrower && HasDate {
		book.IsOnLoan = true
	} else {
		book.IsOnLoan = false
	}
	if book.IsOnLoan {
		book.Borrower = fields[BorrowerIndex]
		date, err := time.Parse(time.DateOnly, fields[DateIndex])
		if err != nil {
			return nil, errors.New("failed to parse date field")
		}
		book.Date = date
	}
	return book, nil
}
