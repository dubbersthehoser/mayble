package controller

import (
	"strings"
	"github.com/dubbersthehoser/mayble/internal/storage"
)


type Field int
const (
	ByTitle Field = iota
	ByAuthor
	ByGenre
	ByBorrower
)

type Searcher struct {
	list []storage.BookLoan
	by  Field
	selection []int
	selected  int // selected in selection
}

func NewSearcher() *Searcher{
	s := &Searcher{}
	return s
}

func (s *Searcher) SetByField(by Field) {
	s.by = by
}

func (s *Searcher) Search(pattern string) {
	
	prefixFinds := make([]int, 0)
	subFinds := make([]int, 0)
	by := s.by

	for i, bookLoan := range s.list {
		var s string
		switch by {
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
		if strings.HasPrefix(s, pattern) {
			prefixFinds = append(prefixFinds, i)
		}
		if strings.Contains(s, pattern) {
			subFinds = append(subFinds, i)
		}
	}

	s.selection = append(prefixFinds, subFinds...)
	s.selected = 0
	
}

func (s *Searcher) Refresh(l []storage.BookLoan) {
	if l == nil {
		panic("nil value slice")
	}
	s.list = l
}

