package search

import (
	"testing"
)

func TestSearchTable(t *testing.T) {
	table := [][]string{
		{
			"1984", "10/12/2022",
			"jake", "11/01/2022",
		},
	}

	tests := []struct{
		name     string
		search   string
		row, col int
	}{
		{
			name: "hunt for 1984 bug",
			search: "1",
			row: 0, col: 0,
		},
	}

	ts := TableSearch{}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts.Set(table, test.search)
			for !ts.Next() {
				
			}
		})
	}
}

func TestEditDist(t *testing.T) {
	tests := []struct {
		name,
		input1,
		input2 string
		expect int
	}{
		{
			name:   "no edits",
			input1: "doggie",
			input2: "doggie",
			expect: 0,
		},
		{
			name:   "delete edits",
			input1: "doie",
			input2: "doggie",
			expect: 2,
		},
		{
			name:   "insert edits",
			input1: "good doggie",
			input2: "doggie",
			expect: 5,
		},
		{
			name:   "change edits",
			input1: "doggoo",
			input2: "doggie",
			expect: 2,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			actual := EditDist(c.input1, c.input2)
			if actual != c.expect {
				t.Fatalf("expect %d, got %d", c.expect, actual)
			}
		})
	}
}

func Test_searchCompare(t *testing.T) {
	tests := []struct {
		name,
		text,
		search string
		expect int
	}{
		{
			name:   "exact match",
			text:   "exact match",
			search: "exact match",
			expect: 10000,
		},
		{
			name:   "sub match bondry",
			text:   "sub match",
			search: "match",
			expect: 6000,
		},
		{
			name:   "sub match non-bondry",
			text:   "sub match",
			search: "tch",
			expect: 5000,
		},
		{
			name:   "fuzzy",
			text:   "sub match",
			search: "seb metch",
			expect: 800,
		},
		{
			name:   "no match",
			text:   "not matching",
			search: "submatch",
			expect: -1,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			actual := searchCompare(c.text, c.search)
			if actual != c.expect {
				t.Fatalf("expect %d, got %d", c.expect, actual)
			}
		})
	}
}
