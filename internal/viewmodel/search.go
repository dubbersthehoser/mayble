package viewmodel

import (
	"strings"
	"unicode"
	"slices"
	"cmp"

	"fyne.io/fyne/v2/data/binding"
)



// EditDist an Levenshtein distance function.
// Returns the total number edits to make s and t match.
func EditDist(s, t string) int {
	
	if len(s) == 0 {
		return len(t)
	}
	if len(t) == 0 {
		return len(s)
	}
	
	height := len(t)+1
	width := len(s)+1

	topbuf := make([]int, width)
	buffer := make([]int, width)

	for i := range width {
		topbuf[i] = i
	}
	
	for y:=1; y<height; y++ {
		buffer[0] = y
		for x:=1; x<width; x++ {
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
// Score goes from 0-n where 0 is the lowest and n is the highest.
// No match returns -1.
func searchCompare(text, search string) int {

	text = strings.ToLower(text)
	search = strings.ToLower(search)

	const (
		ExactMatch    int = 10000
		NoMatch       int = -1

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
	if distance*100 > (maxLength*fuzzyTheshold) {
		return NoMatch
	}
	score := Fuzzy + (distance*fuzzyStep)
	if score < 0 {
		return NoMatch
	} else {
		return score
	}
}

type cellSearchResult struct {
	id       int64
	row, col int
	score    int
}

type TableSearch struct {
	Text   binding.String
	Option binding.String
	table  *table
}

func NewTableSearch(t *table) *TableSearch {
	return &TableSearch{
		Text:   binding.NewString(),
		Option: binding.NewString(),
		table:  t,
	}
}

func (ts *TableSearch) Options() []string {
	return []string{
		"All",
		"Title",
		"Author",
		"Genre",
		"Borrower",
	}
}

func (ts *TableSearch) search(search string) []cellSearchResult {
	if search == "" {
		return []cellSearchResult{}
	}
	result := []cellSearchResult{}
	option, _ := ts.Option.Get()
	ts.table.walkVisableValues(func(row, col int, c *dataCell){
		if search == "" {
			return
		}
		if option != ts.Options()[0] && c.header != option {
			return
		}
		score := searchCompare(c.view, search)
		if score == -1 {
			return
		}
		r := cellSearchResult{
			score: score,
			row: row,
			col: col,
			id: c.id,
		}
		result = append(result, r)
	})
	slices.SortFunc(result, func(a, b cellSearchResult) int {
		return cmp.Compare(a.score, b.score) * -1
	})
	return result
}



