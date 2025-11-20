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

type Book struct {
	ID      int64
	Title   string
	Author  string
	Genre   string
	Ratting int
}

type BookStore interface {
	CreateBook(title, author, genre string, ratting int) (int64, error)
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
	CreateLoan(ID int64, borrower string, date time.Time) error
	UpdateLoan(*Loan) error
	DeleteLoan(*Loan) error
	GetLoan(ID int64) (*Loan, error)
}

type BookLoanStore interface {
	BookStore
	LoanStore
}
