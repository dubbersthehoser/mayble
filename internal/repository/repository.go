package repository

import (
	"time"
)

type Action string
const (
	Delete Action = "delete"
	Update Action = "update"
	Create Action = "create"
	Select Action = "select"
)

type Query struct {
	
	Variant Variant

	Action   Action

	BookID    int64
	SortBy    string
	OrderBy   string

	Entry    *BookEntry
}

type Variant int

const (
	Book   Variant = 1 << iota
	Loaned
	Read
)

func (v Variant) String() string {
	switch v {
	case (Book|Loaned|Read):
		return "book|loaned|Read"
	case (Book|Loaned):
		return "book|loaned"
	case (Book|Read):
		return "book|read"
	case (Book):
		return "book"
	case 0:
		return ""
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

type BookQuerier interface {
	BookQuery(q *Query) ([]BookEntry, error)
}

