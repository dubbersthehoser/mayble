package repository

import (
	"time"
)

type Book struct {
	id int64
	Title  string
	Author string
	Genre  string
}

func (b *Book) ID() int64 {
	return b.id
}

func (b *Book) Type() string {
	return "Book"
}


type BookLoan struct {
	Book

	Borrower string
	Loaned   *time.Time
}

func (bl *BookLoan) ID() int64 {
	return bl.id
}

func (bl *BookLoan) Type() string {
	return "BookLoan"
}


type BookRead struct {
	Book

	Rating int
	Completed *time.Time
}

func (br *BookRead) ID() int64 {
	return br.id
}

func (br *BookRead) Type() string {
	return "BookRead"
}
