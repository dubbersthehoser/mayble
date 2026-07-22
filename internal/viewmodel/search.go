package viewmodel

import (
	"cmp"
	"slices"

	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/search"
)

type Searching struct {
	column int

	cellSearch  search.CellSearch
	tableSearch search.TableSearch

	row    int
	scored [][]int

	l []func()
}

func (s *Searching) GetOptions() []string {
	return []string{
		"All",
		models.BookEntryFields()[models.IdxTitle],
		models.BookEntryFields()[models.IdxAuthor],
		models.BookEntryFields()[models.IdxGenre],
		models.BookEntryFields()[models.IdxBorrower],
		models.BookEntryFields()[models.IdxLoanedAt],
		models.BookEntryFields()[models.IdxCompletedAt],
	}
}

func (s *Searching) SetBy(c string) {
	if c == "All" {
		s.column = -1
		return
	}

	s.column = slices.Index(models.BookEntryFields(), c)
}

func (s *Searching) Selected() (int, int) {
	return s.scored[s.row][0], s.scored[s.row][1]
}

func (s *Searching) Has() bool {
	return len(s.scored) != 0
}

func (s *Searching) Prev() {
	s.row -= 1
	if s.row < 0 {
		s.row = len(s.scored) - 1
	}
	s.notify()
}

func (s *Searching) Next() {
	s.row += 1
	if s.row == len(s.scored) {
		s.row = 0
	}
	s.notify()
}

func (s *Searching) AddListener(fn func()) {
	if s.l == nil {
		s.l = make([]func(), 0)
	}
	s.l = append(s.l, fn)
}

func (s *Searching) notify() {
	for _, fn := range s.l {
		fn()
	}
}

func (s *Searching) search(data [][]string, search string) {
	if s.column == -1 {
		s.searchAll(data, search)
	} else {
		s.searchColumn(data, search)
	}
	s.notify()
}

func (s *Searching) searchColumn(data [][]string, search string) {
	dataCol := make([]string, 0)
	for _, row := range data {
		dataCol = append(dataCol, row[s.column])
	}
	type result struct {
		row, score int
	}
	results := make([]result, 0)
	s.cellSearch.Set(dataCol, search)
	for s.cellSearch.Next() {
		row := s.cellSearch.Pos()
		score := s.cellSearch.Score()
		if score == -1 {
			continue
		}
		r := result{
			row:   row,
			score: score,
		}
		results = append(results, r)
	}

	if len(results) == 0 {
		s.scored = s.scored[:0]
		s.row = 0
		return
	}
	slices.SortFunc(results, func(a, b result) int {
		r := cmp.Compare(a.score, b.score)
		if r == 0 {
			return cmp.Compare(a.row, b.row)
		}
		return cmp.Compare(a.score, b.score)
	})
	s.row = 0
	s.scored = s.scored[:0]
	for _, r := range results {
		row := make([]int, 2)
		row[0] = r.row
		row[1] = s.column
		s.scored = append(s.scored, row)
	}
}

func (s *Searching) searchAll(data [][]string, search string) {
	type result struct {
		row, col, score int
	}
	results := make([]result, 0)
	s.tableSearch.Set(data, search)
	for s.tableSearch.Next() {
		row, col := s.tableSearch.Pos()
		score := s.tableSearch.Score()
		if score == -1 {
			continue
		}
		r := result{
			row:   row,
			col:   col,
			score: score,
		}
		results = append(results, r)
	}
	if len(results) == 0 {
		s.row = 0
		s.scored = s.scored[:0]
	}

	slices.SortFunc(results, func(a, b result) int {
		r := cmp.Compare(a.score, b.score)
		if r == 0 {
			return cmp.Compare(a.row, b.row)
		}
		return r * -1
	})
	s.row = 0
	s.scored = s.scored[:0]
	for _, r := range results {
		row := make([]int, 2)
		row[0] = r.row
		row[1] = r.col
		s.scored = append(s.scored, row)
	}
}

func AllowedSearchOptions(options, headers []string) []string {
	o := make([]string, 0)
	for _, option := range options {
		if slices.Contains(headers, option) || option == "All" {
			o = append(o, option)
		}
	}
	return o
}
