package repository

import (
	"time"
	"strconv"
	"fmt"
	"errors"
	"io"
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
		Title:  b.title,
		Author: b.author,
		Genre:  b.genre,
	}

	if b.loaned != "" && b.borrower != "" {
		book.Variant |= Loaned
		date, err := time.Parse(time.DateOnly, b.loaned)
		if err != nil {
			return nil, err
		}
		book.Loaned = date
		book.Borrower = b.borrower
	}

	if b.completed != "" && b.rating != "" && b.rating != "0" {
		book.Variant |= Read
		date, err := time.Parse(time.DateOnly, b.completed)
		if err != nil {
			return nil, err
		}

		book.Read = date
		book.Rating, err = strconv.Atoi(b.rating)
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

func BookEntryFields() []string {
	return []string{
		"Title",
		"Author",
		"Genre",
		"Rating",
		"Read",
		"Borrower",
		"Loaned",
	}
}

type CSVHandler interface {
	ImportFile(path string) error
	ExportFile(path string) error
}


type BookRetriever interface {
	GetAllBooks(Variant) ([]BookEntry, error)
	GetBookByID(id int64) (BookEntry, error)
}

type BookStore interface {
	BookCreator
	BookUpdator
	BookDeletor
}

type GenreRetriever interface {
	GetUniqueGenres() ([]string, error)
}

type BookCreator interface {
	CreateBook(*BookEntry) (int64, error)
}

type BookUpdator interface {
	UpdateBook(*BookEntry) error
}

type BookDeletor interface {
	DeleteBook(id int64) error
}

type BookImporter interface {
	BookImport(io.Reader) error
}

type BookExporter interface {
	BookExport(io.Writer, []BookEntry) error
}
