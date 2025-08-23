package storage

import (
	"time"
)

type Book struct {
	ID      int64
	Title   string
	Author  string
	Genre   string
	Ratting int
}

type Loan struct {
	ID     int64
	Name   string
	Date   time.Time
	BookID int64
}

type BookLoan struct {
	Book
	Loan
}

type BookStorage interface {
	GetAllBooks() ([]Book, error)
	CreateBook(*Book) error
	UpdateBook(*Book) error
}

type LoanStorage interface {
	GetAllLoans() ([]Loan, error)
	CreateLoan(*Loan) error
	UpdateLoan(*Loan) error
}

type Storage interface {
	LoanStorage
	BookStorage
	Close() error
}



