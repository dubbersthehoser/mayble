package storage

import (
	"errors"
	"time"
	"fmt"

	"github.com/dubbersthehoser/mayble/data"
)

//const MaxBooks int64 = 1000 // the max of books that can be stored.

const ZeroID int64 = -127 

var (
	ErrInvalidValue  = errors.New("storage: invalid value")   // return if one of the values is invalid.
	ErrEntryExists   = errors.New("storage: entry exists")    // return when creating a new entry and the id exists.
	ErrEntryNotFound = errors.New("storage: entry not found") // return when id could not be found in storage.
	ErrStorageFull   = errors.New("storage: hit storage cap") // return when the BookCap is hit.
)


/***********************
	Validation
************************/

func ValidateLoan(loan Loan) error {
	var (
		NameIsZero      bool = loan.Name == ""
		DateIsZero      bool = loan.Date.Equal(time.Time{})
		IDIsZero        bool = loan.ID == ZeroID
		IDIsInvalid     bool = !ValidID(loan.ID)
	)
	switch {
	case NameIsZero:
		return fmt.Errorf("loan name is zero value")
	case DateIsZero:
		return fmt.Errorf("loan date is zero value")
	case IDIsZero:
		return fmt.Errorf("loan id is zero value")
	case IDIsInvalid:
		return fmt.Errorf("loan id is invalid: %d", loan.ID)
	default:
		return nil
	}
}

func ValidateBook(book data.Book) error {
	var (
		TitleIsZero      bool = book.Title == ""
		AuthorIsZero     bool = book.Author == ""
		GenreIsZero      bool = book.Genre == ""
		RattingIsInvalid bool = !ValidRatting(book.Ratting)
		IDIsZero         bool = book.ID == ZeroID
		IDIsInvalid      bool = !ValidID(book.ID)
	)
	switch {
	case TitleIsZero:
		return fmt.Errorf("book title is zero value")
	case AuthorIsZero:
		return fmt.Errorf("book author is zero value")
	case GenreIsZero:
		return fmt.Errorf("book genre is zero value")
	case RattingIsInvalid:
		return fmt.Errorf("book ratting is invalid: %d", book.Ratting)
	case IDIsZero:
		return fmt.Errorf("book ratting is zero value")
	case IDIsInvalid:
		return fmt.Errorf("book id is invalid: %d", book.ID)
	default:
		return nil
	}
}

func ValidID(id int64) bool {
	return id >= 0
}

func ValidRatting(ratting int) bool {
	return ratting >= 0 && ratting < 6
}



/********************************
	Storage Interface
*********************************/

type Storage interface {
	// GetAllBookLoans returns a list of stored book loans.
	GetAllBookLoans() ([]data.BookLoan, error)

	// GetBookLoanByID returns stored book by its id.
	GetBookLoanByID(id int64) (data.BookLoan, error)

	// CreateBookLoan adds book loan to storage.
	// returns ErrEntryExists when book id is in storage. Use ZeroID for id or NewBookLoan().
	CreateBookLoan(*data.BookLoan) (int64, error)

	// UpdateBookLoan update book loan in storage.
	// returns ErrEntryNotFound when book id is not in storage.
	UpdateBookLoan(*data.BookLoan) error

	// DeleteBookLoan remove book loan from storage.
	// returns ErrEntryNotFound when book id is not in storage.
	DeleteBookLoan(*data.BookLoan) error

	// Close whatever implementation. Can be nop.
	Close() error
}

