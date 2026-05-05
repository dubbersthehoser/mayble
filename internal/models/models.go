package model

import (
	"time"
	"errors"
	"strconv"
	"fmt"
)

type Book struct {
	Title  string
	Author string
	Genre  string
}

type Loaned struct {
	Borrower string
	Loaned   time.Time
}

type Completed struct {
	Rating    int
	Completed time.Time
}


type BookEntry struct {
	ID      int64
	Book
	Loaned
	Completed
	IsCompleted bool
	IsLoaned    bool
}

type BookEntryBuilder struct {
	id        int64
	title     string
	author    string
	genre     string
	completed string
	loaned    string
	borrower  string
	rating    string
}

func NewBookEntryBuilder() *BookEntryBuilder {
	return &BookEntryBuilder{
		id: -127,
	}
}

func (b *BookEntryBuilder) Build() (*BookEntry, error) {

	if b.id < 0 {
		return nil, errors.New("id was not set")
	}

	if b.title == "" {
		return nil, errors.New("title was not set")
	}
	if b.author == "" {
		return nil, errors.New("author was not set")
	}
	if b.genre == "" {
		return nil, errors.New("genre was not set")
	}

	book := &BookEntry{
		ID:     b.id,
		Book: Book{ 
			Title:  b.title,
			Author: b.author,
			Genre:  b.genre,
		},
	}

	if b.loaned != "" && b.borrower != "" {
		book.IsLoaned = true
		date, err := time.Parse(time.DateOnly, b.loaned)
		if err != nil {
			return nil, err
		}
		book.Loaned.Loaned = date
		book.Loaned.Borrower = b.borrower
	}

	if b.completed != "" && b.rating != "" && b.rating != "0" {
		book.IsCompleted = true
		date, err := time.Parse(time.DateOnly, b.completed)
		if err != nil {
			return nil, err
		}

		book.Completed.Completed = date
		book.Completed.Rating, err = strconv.Atoi(b.rating)
		if err != nil {
			return nil, fmt.Errorf("invalid rating '%s'", b.rating)
		}
	}
	return book, nil
}

func (b *BookEntryBuilder) SetID(id int64) *BookEntryBuilder {
	b.id = id
	return b
}

func (b *BookEntryBuilder) SetTitle(t string) *BookEntryBuilder {
	b.title = t
	return b
}

func (b *BookEntryBuilder) SetAuthor(a string) *BookEntryBuilder {
	b.author = a
	return b
}

func (b *BookEntryBuilder) SetGenre(g string) *BookEntryBuilder {
	b.genre = g
	return b
}

func (b *BookEntryBuilder) SetLoaned(d string) *BookEntryBuilder {
	b.loaned = d
	return b
}

func (b *BookEntryBuilder) SetCompleted(d string) *BookEntryBuilder {
	b.completed = d
	return b
}

func (b *BookEntryBuilder) SetBorrower(n string) *BookEntryBuilder {
	b.borrower = n
	return b
}

func (b *BookEntryBuilder) SetRating(r string) *BookEntryBuilder {
	b.rating = r
	return b
}


const (
	IdxTitle int = iota
	IdxAuthor
	IdxGenre
	IdxRating
	IdxRead
	IdxLoaned
	IdxBorrower
)

// BookEntryFields returns the names of each field name of BookEntry
func BookEntryFields() []string {
	return []string{
		"Title",
		"Author",
		"Genre",
		"Rating",
		"Read",
		"Loaned",
		"Borrower",
	}
}
