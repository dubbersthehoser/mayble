package storage

import (
	"errors"

	"github.com/dubbersthehoser/mayble/internal/data"
)

var (
	ErrInvalidValue  = errors.New("storage: invalid value")   // return when entry has invalid data.
	ErrEntryExists   = errors.New("storage: entry exists")    // return when creating a new entry with an id that exists.
	ErrEntryNotFound = errors.New("storage: entry not found") // return when id could not be found in storage.
	ErrStorageFull   = errors.New("storage: hit storage cap") // return when hit a storeage cap.
)


type Storage interface {

	// GetAllBookLoans all book loans in store.
	GetAllBookLoans() ([]data.BookLoan, error)

	// GetBookLoanByID returns stored book by its id, and ErrEntryNotFound if not found.
	GetBookLoanByID(id int64) (data.BookLoan, error)

	// CreateBookLoan adds book loan to storage, and returns its new id.
	// return ErrEntryExists when book id is in storage. Use data.ZeroID as id or data.NewBookLoan().
	// return ErrInvalidValue when given nil
	CreateBookLoan(*data.BookLoan) (int64, error)

	// UpdateBookLoan update book loan in storage.
	// return ErrEntryNotFound when book id is not in storage.
	// return ErrInvalidValue when book loan is nil
	UpdateBookLoan(*data.BookLoan) error

	// DeleteBookLoan remove book loan from storage.
	// returns ErrEntryNotFound when book id is not in storage.
	DeleteBookLoan(*data.BookLoan) error

	// Close run clean up code, or close a connection. Can be NOP.
	Close() error
}

