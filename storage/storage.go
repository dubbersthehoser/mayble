package storage

import (
	"errors"
	"time"
	"fmt"

	"github.com/dubbersthehoser/mayble/data"
)

//const MaxBooks int64 = 1000 // the max of books that can be stored.

var (
	ErrInvalidValue  = errors.New("storage: invalid value")   // return if one of the values is invalid.
	ErrEntryExists   = errors.New("storage: entry exists")    // return when creating a new entry and the id exists.
	ErrEntryNotFound = errors.New("storage: entry not found") // return when id could not be found in storage.
	ErrStorageFull   = errors.New("storage: hit storage cap") // return when the BookCap is hit.
)

/********************************
	Storage Interface
*********************************/

type Storage interface {
	// GetAllBookLoans returns a list of stored book loans.
	GetAllBookLoans() ([]data.BookLoan, error)

	// GetBookLoanByID returns stored book by its id.
	GetBookLoanByID(id int64) (data.BookLoan, error)

	// CreateBookLoan adds book loan to storage.
	// returns ErrEntryExists when book id is in storage. Use data.ZeroID as id or data.NewBookLoan().
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

