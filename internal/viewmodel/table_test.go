package viewmodel


import (
	"testing"
	"slices"
)

func Test_hiddenOptionsToHeaders(t *testing.T) {
	tests := []struct{
		input  []string
		expect []string
	}{
		{
			input: []string{"Loaned"},
			expect: []string{"Loaned", "Borrower"},
		},
		{
			input: []string{"Read"},
			expect: []string{"Read", "Rating"},
		},
		{
			input: []string{"Read", "Loaned"},
			expect: []string{"Read", "Rating", "Loaned", "Borrower"},
		},
		{
			input: []string{"Read", "Loaned", "extra", "extra"},
			expect: []string{"Read", "Rating", "Loaned", "Borrower", "extra", "extra"},
		},
	}

	for i, c := range tests {
		actual := hiddenOptionsToHeaders(c.input)
		slices.Sort(actual)
		slices.Sort(c.expect)
		if n := slices.Compare(c.expect, actual); n != 0 {
			t.Fatalf("[%d] expect\n\t%#v\ngot\n\t%#v", i, c.expect, actual)
		}
	}
}


func Test_hiddenHeadersToOptions(t *testing.T) {
	
	tests := []struct{
		input  []string
		expect []string

	}{
		{
			input: []string{"Loaned", "Borrower"},
			expect: []string{"Loaned"},
		},
		{
			input: []string{"Read", "Rating"},
			expect: []string{"Read"},
		},
		{
			input: []string{"Read", "Rating", "Loaned", "Borrower"},
			expect: []string{"Loaned", "Read"},
		},
		{
			input: []string{"Read", "Rating", "Loaned", "Borrower", "Title", "Author"},
			expect: []string{"Title", "Author", "Loaned", "Read"},
		},
	}
	for i, c := range tests {
		actual := hiddenHeadersToOptions(c.input)
		slices.Sort(actual)
		slices.Sort(c.expect)
		if n := slices.Compare(c.expect, actual); n != 0 {
			t.Fatalf("[%d] expect\n\t%#v\ngot\n\t%#v", i, c.expect, actual)
		}
	}
}
