package listing

import (
	"slices"

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
	switch by {
	case "Title":
		return ByTitle, nil
	case "Author":
		return core.ByAuthor, nil
	case "Genre":
		return ByGenre, nil
	case "Ratting":
		return ByRatting, nil
	case "Borrower":
		return ByBorrower, nil
	case "Date":
		return ByDate, nil
	default:
		return ByNothing, errors.New("invalid string order by value")
	}
}

type Ordering int

const (
	ASC Order = iota
	DEC
)

func OrderBookLoans(s []data.BookLoan, by OrderBy, order Ordering) []data.BookLoan {
	compare := func(x, y s.BookLoan) int {
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
				a := strings.ToLower(x.Loan.Name)
				b := strings.ToLower(y.Loan.Name)
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
	return bookLoans
}

func BookLoanToListed(bookLoan *data.BookLoan) *BookLoan {
	view := BookLoanListed{
		BookId:  bookLoan.Book.ID,
		Title:   bookLoan.Title,
		Author:  bookLoan.Author,
		Genre:   bookLoan.Genre,
		Ratting: rattingToString(bookLoan.Ratting),
		Borrower: "n/a",
		Date:     "n/a",
	}
	if bookLoan.IsOnLoan() {
		view.Borrower = bookLoan.Loan.Name
		view.Date  = dateToString(&bookLoan.Loan.Date)
	}
	return &view
}

