package repository

import (
	"time"
)


type Book struct {
	id     int64
	Title  string
	Author string
	Genre  string
}

var _ Resultable = &Book{}

func (b *Book) ID() int64 {
	return b.id
}

func (b *Book) Type() string {
	return string(Main)
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
	return string(Loaned)
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
	return string(Read)
}


type BookBuilder struct {
	id int64
	book struct{
		title string
		author string
		genre  string
	}
}

func (bb *BookBuilder) SetID(i int64) *BookBuilder {
	bb.id = i
	return bb
}

func (bb *BookBuilder) SetTitle(s string) *BookBuilder {
	bb.book.title = s
	return bb
}

func (bb *BookBuilder) SetAuthor(s string) *BookBuilder {
	bb.book.author = s
	return bb
}

func (bb *BookBuilder) SetGenre(s string) *BookBuilder {
	bb.book.genre = s
	return bb
}

func (bb *BookBuilder) Build() (Resultable, error) {
	r := &Book{
		id: bb.id,
		Title: bb.book.title,
		Author: bb.book.author,
		Genre: bb.book.genre,
	}
	return r, nil
}
