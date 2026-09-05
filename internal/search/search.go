package search

import (
	"strings"
	"unicode"
)

type Point struct {
	Col   int
	Row   int
}

type Traverser interface {
	Next() (string, bool)
	Point() Point
}

type Searcher struct {
	tvr    Traverser
	isDone bool
	score  int
	search string
}

func (s *Searcher) Set(tvr Traverser, search string) *Searcher {
	s.tvr = tvr
	s.isDone = false
	s.search = search
	return s
}

func (s *Searcher) Point() Point {
	return s.tvr.Point()
}

func (s *Searcher) Score() int {
	return s.score
}

func (s *Searcher) IsFinished() bool {
	return s.isDone
}

func (s *Searcher) Next() bool {
	item, ok := s.tvr.Next()
	if ok {
		s.score = searchCompare(item, s.search)
	}
	s.isDone = !ok
	return ok
}

type ColumnTraverse struct {
	data    [][]string
	setCol  int
	currRow int
}

func (ct *ColumnTraverse) Set(data [][]string, col int) *ColumnTraverse {
	ct.data = data
	ct.setCol = col
	ct.currRow = -1
	return ct
}

func (ct *ColumnTraverse) Next() (string, bool) {
	ct.currRow += 1
	if len(ct.data) <= ct.currRow {
		return "", false
	}
	return ct.data[ct.currRow][ct.setCol], true
}

func (ct *ColumnTraverse) Point() Point {
	return Point{Row: ct.currRow, Col: ct.setCol}
}

type TableTraverse struct {
	data    [][]string
	currCol int
	currRow int
}

func (tt *TableTraverse) Set(data [][]string) *TableTraverse {
	tt.data = data
	tt.currCol = -1
	tt.currRow = 0
	return tt
}

func (tt *TableTraverse) Next() (string, bool) {
	tt.currCol += 1
	if len(tt.data[tt.currRow]) <= tt.currCol {
		tt.currCol = 0
		tt.currRow += 1
	}
	if len(tt.data) <= tt.currRow {
		return "", false
	}

	return tt.data[tt.currRow][tt.currCol], true
}

func (tt *TableTraverse) Point() Point {
	return Point{Row: tt.currRow, Col: tt.currCol}
}

// EditDist an Levenshtein distance function.
// Returns the total number edits to make s and t match.
func EditDist(s, t string) int {

	if len(s) == 0 {
		return len(t)
	}
	if len(t) == 0 {
		return len(s)
	}

	height := len(t) + 1
	width := len(s) + 1

	topbuf := make([]int, width)
	buffer := make([]int, width)

	for i := range width {
		topbuf[i] = i
	}

	for y := 1; y < height; y++ {
		buffer[0] = y
		for x := 1; x < width; x++ {
			if t[y-1] != s[x-1] {
				del := 1 + topbuf[x]
				ins := 1 + buffer[x-1]
				cha := 1 + topbuf[x-1]
				buffer[x] = min(del, ins, cha)
			} else {
				buffer[x] = topbuf[x-1]
			}
		}
		buffer, topbuf = topbuf, buffer
	}
	return topbuf[width-1]
}

// searchCompare compare text to a search and get its score.
// The score goes from 0-n where 0 is the lowest and n is the highest, highest being the closes match.
// If there is no match it will return -1.
func searchCompare(text, search string) int {

	text = strings.ToLower(text)
	search = strings.ToLower(search)

	const (
		ExactMatch int = 10000
		NoMatch    int = -1

		SubString     int = 5000 // Base sub-string search score.
		BoundaryBonus int = 1000 // Sub-string bonus for being a prefix of a word.
		PrefixBonus   int = 1000 // When search is the start of the text.

		Fuzzy         int = 1000 // Base fuzzy search score
		fuzzyTheshold int = 40   // Precentage theshold for the length of longest string to the edit distance.
		fuzzyStep     int = -100 // Reduced score per edit distance step.
	)

	if text == search {
		return ExactMatch
	}

	if idx := strings.Index(text, search); idx != -1 {
		score := SubString
		// check whether the search string is at the start of a word in text.
		inBoundary := idx == 0 || !unicode.IsLetter(rune(text[idx-1]))
		isPrefix := strings.HasPrefix(text, search)
		if inBoundary {
			score += BoundaryBonus
		}
		if isPrefix {
			score += PrefixBonus
		}
		return score
	}

	distance := EditDist(text, search)
	maxLength := max(len(text), len(search))
	if distance*100 > (maxLength * fuzzyTheshold) {
		return NoMatch
	}
	score := Fuzzy + (distance * fuzzyStep)
	if score < 0 {
		return NoMatch
	} else {
		return score
	}
}
