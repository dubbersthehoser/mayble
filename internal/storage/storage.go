package storage

import (
	"errors"
	"time"
)

var (
	ErrInvalidValue  = errors.New("storage: invalid value")   // return when entry has invalid data.
	ErrEntryExists   = errors.New("storage: entry exists")    // return when creating a new entry with an id that exists.
	ErrEntryNotFound = errors.New("storage: entry not found") // return when id could not be found in storage.
	ErrStorageFull   = errors.New("storage: hit storage cap") // return when hit a storeage cap.
)

const ZeroID int64 = -127

// IsZeroID check if id is a zero id.
func IsZeroID(id int64) bool {
	return id < 0
}

type Book struct {
	ID      int64
	Title   string
	Author  string
	Genre   string
	Ratting int
}

type BookStore interface {
	// CreateBook with id if id is positive, if negative will be made and return.
	CreateBook(id int64, title, author, genre string, ratting int) (int64, error)
	UpdateBook(*Book) error 
	DeleteBook(*Book) error
	GetBooks() ([]Book, error)
	GetBookByID(ID int64) (*Book, error)
}

type Loan struct {
	ID       int64
	Borrower string
	Date     time.Time
}

type LoanStore interface {
	CreateLoan(BookID int64, borrower string, date time.Time) error
	UpdateLoan(*Loan) error
	DeleteLoan(*Loan) error
	GetLoan(BookID int64) (*Loan, error)
}

type BookLoanStore interface {
	BookStore
	LoanStore
	Close() error
}
