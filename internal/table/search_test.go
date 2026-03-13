package table

import (
	"testing"
)


func TestEditDist(t *testing.T) {
	tests := []struct{
		name,
		input1,
		input2 string
		expect int
	}{
		{
			name: "no edits",
			input1: "doggie",
			input2: "doggie",
			expect: 0,
		},
		{
			name: "delete edits",
			input1: "doie",
			input2: "doggie",
			expect: 2,
		},
		{
			name: "insert edits",
			input1: "good doggie",
			input2: "doggie",
			expect: 5,
		},
		{
			name: "change edits",
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
	tests := []struct{
		name,
		text,
		search string
		expect int
	}{
		{
			name: "exact match",
			text: "exact match",
			search: "exact match",
			expect: 10000,
		},
		{
			name: "sub match bondry",
			text: "sub match",
			search: "match",
			expect: 6000,
		},
		{
			name: "sub match non-bondry",
			text: "sub match",
			search: "tch",
			expect: 5000,
		},
		{
			name: "fuzzy",
			text: "sub match",
			search: "seb metch",
			expect: 800,
		},
		{
			name: "no match",
			text: "not matching",
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

func TestSearch(t *testing.T) {
	books := [][]string{
			{
				"A Title",
				"A Author",
				"A Genre",
			},
			{
				"B Title",
				"B Author",
				"B Genre",
			},
			{
				"C Title",
				"C Author",
				"C Genre",
			},
			{
				"D Title",
				"D Author",
				"D Genre",
			},
	}

	table := NewTable("SearchTest", []string{"Title", "Author", "Genre"})
	for i, book := range books {
		table.AppendRow(int64(i+1), book)
	}
	results := Search(table, "D Title", "")
	if len(results) == 0 {
		t.Fatal("expected length > 0")
	}

	if results[0].ID != 4 {
		t.Fatalf("expect id of 4, got %d", results[0].ID)
	}
}
