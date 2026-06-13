package search

import (
	"strings"
	"unicode"
)

// EditDist an Levenshtein distance function.
//
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

// searchCompare get the compare score for a search.
//
// Score goes from 0-n where 0 is the lowest and n is the highest.
// No match returns -1.
func searchCompare(text, search string) int {

	text = strings.ToLower(text)
	search = strings.ToLower(search)

	const (
		ExactMatch int = 10000
		NoMatch    int = -1

		SubString     int = 5000 // Base sub-string search score.
		BoundaryBonus int = 1000 // Sub-string bonus for being a prefix of a word.

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
		if inBoundary {
			score += BoundaryBonus
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

type CellSearch struct {
	search string
	cells []string
	curr int
	score int
}

func (cs *CellSearch) Pos() int {
	return cs.curr
}

func (cs *CellSearch) Score() int {
	return cs.score
}

func (cs *CellSearch) IsFinished() bool {
	if len(cs.cells) >= cs.curr {
		return true
	}
	return false
}

func (cs *CellSearch) Next() bool {
	if cs.IsFinished() {
		return false
	}
	cs.score = searchCompare(cs.cells[cs.curr], cs.search)
	if cs.score == -1 {
		return false
	}
	cs.curr += 1
	return true
}

func (cs *CellSearch) Set(c []string, search string) {
	cs.search = search
	cs.cells = c
	cs.score = -1
	cs.curr = 0
}

type TableSearch struct {
	search string
	table [][]string
	row, col int
	cellSearch *CellSearch
}

func (ts *TableSearch) Set(table [][]string, search string) {
	ts.search = search
	ts.table = table
	ts.row, ts.col = 0, 0
}

func (ts *TableSearch) IsFinished() bool {
	if len(ts.table) >= ts.row {
		return true
	}
	return false
}

func (ts *TableSearch) Pos() (int, int) {
	return ts.row, ts.cellSearch.curr
}

func (ts *TableSearch) Score() int {
	return ts.cellSearch.score
}

func (ts *TableSearch) Next() bool {
	if ts.IsFinished() {
		return false
	}
	if ts.cellSearch == nil {
		ts.cellSearch = NewCellSearch(ts.table[ts.row], ts.search)
	}
	
	for !ts.cellSearch.Next() {
		if ts.cellSearch.IsFinished() {
			ts.row += 1
			if ts.IsFinished() {
				return false
			}
			ts.cellSearch.Research(ts.table[ts.row], ts.search)
		}
	}
	return true
}
