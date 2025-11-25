package searching

import (
	"errors"
	"strings"

	"github.com/dubbersthehoser/mayble/internal/app"
)

type Ring interface {
	Selected() int
	Next() int
	Prev() int
}

type RangeRing struct {
	max int
	selected int
}
func NewRangeRing(max int) *RangeRing {
	return &RangeRing{
		max: max,
	}
}
func (rr *RangeRing) Selected() int {
	return rr.selected
}
func (rr *RangeRing) Next() int {
	idx := (rr.selected+1) % rr.max
	rr.selected = idx
	return idx
}
func (rr *RangeRing) Prev() int {
	index := rr.selected-1
	if index < 0 {
		index = rr.max + index
	}
	rr.selected = index
	return index
}


type SelectionRing struct {
	selected int
	selection []int
}

func NewSelectionRing(selection []int) *SelectionRing {
	return &SelectionRing{
		selection: selection,
	}
}

func (sr *SelectionRing) Selected() int {
	index := sr.selected
	return sr.selection[index]
}

func (sr *SelectionRing) Next() int {
	index := (sr.selected+1) % len(sr.selection)
	sr.selected = index
	return sr.selection[index]
}

func (sr *SelectionRing) Prev() int {
	index := sr.selected-1
	if index < 0 {
		index = len(sr.selection)+index
	}
	sr.selected = index
	return sr.selection[index]
}


type Field int
const (
	ByNothing Field = iota - 1
	ByTitle
	ByAuthor
	ByGenre
	ByBorrower
)

func StringToField(s string) (Field, error) {
	s = strings.ToLower(s)
	switch s {
	case "title":
		return ByTitle, nil
	case "author":
		return ByAuthor, nil
	case "genre":
		return ByGenre, nil
	case "borrower":
		return ByBorrower, nil
	default:
		return ByNothing, errors.New("invalid string field value")
	}
}

func MustStringToField(s string) Field {
	f, err := StringToField(s)
	if err != nil {
		panic(err)
	}
	return f
}

func SearchBookLoans(l []app.BookLoan, f Field, pattern string) []int {
	finds := make([]int, 0)
	for i, bookLoan := range l {
		var s string
		switch f {
		case ByTitle:
			s = bookLoan.Title
		case ByAuthor:
			s = bookLoan.Author
		case ByGenre:
			s = bookLoan.Genre
		case ByBorrower:
			if !bookLoan.IsOnLoan {
				continue
			} 
			s = bookLoan.Borrower
		default:
			panic("searching: invalid search field")
		}
		s = strings.ToLower(s)
		pattern = strings.ToLower(pattern)
		if strings.Contains(s, pattern) {
			finds = append(finds, i)
		}
	}
	return finds
}


