package searching

import (
	"strings"

	"github.com/dubbersthehoser/mayble/internal/data"
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
	return (rr.selected+1) % rr.max
}
func (rr *RangeRing) Prev() int {
	index := rr.selected-1
	if index < 0 {
		index = rr.selected + index
	}
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
	return sr.selection[index]
}

func (sr *SelectionRing) Prev() int {
	index := sr.selected-1
	if index < 0 {
		index = sr.selected+index
	}
	return sr.selection[index]
}


//func Search(s []string, pattern string) []int {}

type Field int
const (
	ByTitle Field = iota
	ByAuthor
	ByGenre
	ByBorrower
)

func SearchBookLoans(l []data.BookLoan, f Field, pattern string) []int {
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
			if bookLoan.Loan == nil {
				continue
			} 
			s = bookLoan.Genre
		default:
			panic("unknown search field")
		}
		s = strings.ToLower(s)
		pattern = strings.ToLower(pattern)
		if strings.Contains(s, pattern) {
			finds = append(finds, i)
		}
	}
	return finds
}


