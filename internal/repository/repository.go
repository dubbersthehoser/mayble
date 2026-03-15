package repository

import (
	"time"
)

type Variant int

const (
	Book   Variant = 0
	Loaned Variant = 1 << iota
	Read
)

func (v Variant) String() string {
	switch v {
	case (Loaned | Read):
		return "book|loaned|Read"
	case (Loaned):
		return "book|loaned"
	case (Read):
		return "book|read"
	case 0:
		return "book"
	default:
		panic("variant not found")
	}
}

type BookEntry struct {
	Variant Variant
	ID      int64

	Title  string
	Author string
	Genre  string

	Rating int
	Read   time.Time

	Borrower string
	Loaned   time.Time
}

const (
	IdxTitle int = iota
	IdxAuthor
	IdxGenre

	IdxRead
	IdxRating

	IdxLoaned
	IdxBorrower
)

func BookEntryFields() []string {
	return []string{
		"Title",
		"Author",
		"Genre",
		"Read",
		"Rating",
		"Loaned",
		"Borrower",
	}
}

type BookRetriever interface {
	GetAllBooks(Variant) ([]BookEntry, error)
	GetBookByID(id int64) (BookEntry, error)
}

type GenreRetriever interface {
	GetUniqueGenres() ([]string, error)
}

type BookCreator interface {
	CreateBook(b *BookEntry) (int64, error)
}

type BookUpdator interface {
	UpdateBook(b *BookEntry) error
}

type BookDeletor interface {
	DeleteBook(b *BookEntry) error
}
