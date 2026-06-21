package viewmodel

import (
	"slices"
	"cmp"
	"fmt"
	
	"github.com/dubbersthehoser/mayble/internal/search"
	"github.com/dubbersthehoser/mayble/internal/models"
)

type Searching struct {
	column int

	cellSearch  search.CellSearch
	tableSearch search.TableSearch
}

func (s *Searching) GetOptions() []string {
	return []string{
		"All",
		models.BookEntryFields()[models.IdxTitle],
		models.BookEntryFields()[models.IdxAuthor],
		models.BookEntryFields()[models.IdxGenre],
		models.BookEntryFields()[models.IdxBorrower],
	}
}

func (s *Searching) SetBy(c string) {
	if c == "All" {
		s.column = -1
		return
	} 

	s.column = slices.Index(models.BookEntryFields(), c)
}

func (s *Searching) search(data [][]string, search string) (int, int, bool){
	if s.column == -1 {
		return s.searchAll(data, search)
	} else {
		return s.searchColumn(data, search)
	}
}

func (s *Searching) searchColumn(data [][]string, search string) (int, int, bool) {
	dataCol := make([]string, 0)
	for _, row := range data {
		dataCol = append(dataCol, row[s.column])
	}
	type result struct{
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
			row: row,
			score: score,
		}
		results = append(results, r)
	}

	if len(results) == 0 {
		return 0, 0, false
	}

	slices.SortFunc(results, func(a, b result) int {
		return cmp.Compare(a.score, b.score)
	})
	r := results[len(results)-1]
	return r.row, s.column, true
}

func (s *Searching) searchAll(data [][]string, search string) (int, int, bool) {
	if search == "" {
		return 0, 0, false
	}
	type result struct{
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
			row: row,
			col: col,
			score: score,
		}
		results = append(results, r)
	}
	if len(results) == 0 {
		return 0,0, false
	}

	slices.SortFunc(results, func(a, b result) int {
		r := cmp.Compare(a.score, b.score)
		if r == 0 {
			return cmp.Compare(a.row, b.row)
		}
		return r
	})
	for _, r := range results {
		fmt.Printf("%d:%d %d '%s' -> '%s'\n",r.row, r.col, r.score, search, data[r.row][r.col])
	}
	r := results[len(results)-1]
	return r.row, r.col, true
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

