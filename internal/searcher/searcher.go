package searcher

import (
	"github.com/dubbersthehoser/mayble/internal/storage"
)

type Searcher struct {
}

type Field int
const (
	ByTitle Field = iota
	ByAuthor
	ByGenre
	ByBorrower
)

func (s *Searcher) Search(by Field, pattern string) []int64 {
	
}
func (s *Searcher) Refresh() {
	
}

func (s *Searcher) Add(handle int64, bookLoan *storage.BookLoan) {
	
}

