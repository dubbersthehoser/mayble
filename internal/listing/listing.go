package listing

import (
	"slices"
	"errors"
	"strings"
	"time"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/data"
)

type BookLoan struct {
	ID       int64
	Title    string
	Author   string
	Genre    string
	Ratting  string
	Borrower string
	Date     string
}

type OrderBy string

const (
	ByTitle OrderBy = "Title"
	ByAuthor        = "Author"
	ByGenre         = "Genre"
	ByRatting       = "Ratting"
	ByBorrower      = "Borrower"
	ByDate          = "Date"
	ByID            = "ID"
	ByNothing       = ""
)

func SortByList() []string {
	return []string{"Title", "Author", "Genre", "Ratting", "Borrower", "Date"}
}

func StringToOrderBy(s string) (OrderBy, error) {
	s = strings.ToLower(s)
	switch s {
	case "title":
		return ByTitle, nil
	case "author":
		return ByAuthor, nil
	case "genre":
		return ByGenre, nil
	case "ratting":
		return ByRatting, nil
	case "borrower":
		return ByBorrower, nil
	case "date":
		return ByDate, nil
	default:
		return ByNothing, errors.New("invalid string order by value")
	}
}
func MustStringToOrderBy(s string) OrderBy {
	o, err := StringToOrderBy(s)
	if err != nil {
		panic(err)
	}
	return o
}

type Ordering int

const (
	ASC Ordering = iota
	DEC
)

func OrderBookLoans(s []data.BookLoan, by OrderBy, order Ordering) []data.BookLoan {
	compare := func(x, y data.BookLoan) int {
		const (
			GreaterX int = 1
			Equal    int = 0
			LesserX  int = -1
		)
		result := Equal
		switch by {
		case ByID:
			switch {
			case x.Ratting == y.Ratting:
				result = Equal
			case x.Ratting > y.Ratting:
				result = GreaterX
			case x.Ratting < y.Ratting:
				result = LesserX
			}
		case ByTitle:
			a := strings.ToLower(x.Title)
			b := strings.ToLower(y.Title)
			result = strings.Compare(a, b)
		case ByAuthor:
			a := strings.ToLower(x.Author)
			b := strings.ToLower(y.Author)
			result = strings.Compare(a, b)
		case ByGenre:
			a := strings.ToLower(x.Genre)
			b := strings.ToLower(y.Genre)
			result = strings.Compare(a, b)
		case ByRatting:
			switch {
			case x.Ratting == y.Ratting:
				result = Equal
			case x.Ratting > y.Ratting:
				result = GreaterX
			case x.Ratting < y.Ratting:
				result = LesserX
			}
		case ByBorrower, ByDate:
			if x.Loan == nil && y.Loan == nil {
				result = Equal
			} else if x.Loan == nil {
				result = LesserX
			} else if y.Loan == nil {
				result = GreaterX
			} else if by == ByBorrower {
				a := strings.ToLower(x.Loan.Borrower)
				b := strings.ToLower(y.Loan.Borrower)
				result = strings.Compare(a, b)
			} else if by == ByDate {
				result = x.Loan.Date.Compare(y.Loan.Date)
			}
		case ByNothing:
			result = Equal
		}
		if order == DEC {
			result = result * -1
		}
		return result
	}
	slices.SortFunc(s, compare)
	return s
}


func GetRattingStrings() []string {
	return []string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}
}

func RattingToInt(ratting string) (int, error) {
	for i, str := range GetRattingStrings() {
		if str == ratting {
			return i, nil
		}
	}
	return -1, errors.New("listing: invalid ratting string")
}
func MustRattingToInt(r string) int {
	i, err := RattingToInt(r)
	if err != nil {
		panic(err)
	}
	return i
}

func RattingToString(r int) (string, error) {
	strings := GetRattingStrings()
	if r >= len(strings) || r < 0 {
		return "", errors.New("ratting is out of range")
	}
	return strings[r], nil
}
func MustRattingToString(r int) string {
	s, err := RattingToString(r)
	if err != nil {
		panic(err)
	}
	return s
}

func DateToString(date *time.Time) string {
	return fmt.Sprintf("%d/%d/%d", date.Day(), date.Month(), date.Year())
}


func BookLoanToListed(bookLoan *data.BookLoan) *BookLoan {
	view := BookLoan{
		ID:      bookLoan.Book.ID,
		Title:   bookLoan.Title,
		Author:  bookLoan.Author,
		Genre:   bookLoan.Genre,
		Ratting: MustRattingToString(bookLoan.Ratting),
		Borrower: "n/a",
		Date:     "n/a",
	}
	if bookLoan.IsOnLoan() {
		view.Borrower = bookLoan.Loan.Borrower
		view.Date  = DateToString(&bookLoan.Loan.Date)
	}
	return &view
}

