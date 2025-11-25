package searching

import (
	"testing"

	"github.com/dubbersthehoser/mayble/internal/app"
)


func TestRangeRing(t *testing.T) {
	MaxSize := 12
	ring := NewRangeRing(MaxSize)

	for i:=0; i<MaxSize; i++ {
		if i != ring.Selected() {
			t.Fatalf("expect %d, got %d", i, ring.Selected())
		}
		ring.Next()
	}

	if ring.Selected() != 0 {
		t.Fatalf("expect %d, got %d", 0, ring.Selected())
	} 

	for i:=(MaxSize - 1); i>=0; i-- {
		ring.Prev()
		if i != ring.Selected() {
			t.Fatalf("expect %d, got %d", i, ring.Selected())
		}
	}

	if ring.Selected() != 0 {
		t.Fatalf("expect %d, got %d", 0, ring.Selected())
	}
}

func TestSelectionRing(t *testing.T) {
	selection := []int{
		4,
		34,
		2,
		9,
		99,
		124,
	}
	ring := NewSelectionRing(selection)

	for i:=0; i<len(selection); i++ {
		idx := ring.Selected()
		if idx != selection[i] {
			t.Fatalf("expect %d, got %d", selection[i], idx)
		}
		ring.Next()
	}

	if ring.Selected() != selection[0] {
		t.Fatalf("expect %d, got %d", selection[0], ring.Selected())
	}

	for i:=len(selection)-1; i>=0; i-- {
		ring.Prev()
		idx := ring.Selected()
		if idx != selection[i] {
			t.Fatalf("expect %d, got %d", selection[i], idx)
		}
	}

	if ring.Selected() != selection[0] {
		t.Fatalf("expect %d, got %d", selection[0], ring.Selected())
	}
}


func TestStringToField(t *testing.T) {
	tests := []struct{
		input string
		expect Field
	}{
		{
			input: "tiTle",
			expect: ByTitle,
		},
		{
			input: "author",
			expect: ByAuthor,
		},
		{
			input: "Genre",
			expect: ByGenre,
		},
		{
			input: "BORROWER",
			expect: ByBorrower,
		},
	}

	for i, test := range tests {
		expect := test.expect
		actual, err := StringToField(test.input)
	
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		if actual != test.expect {
			t.Fatalf("case %d, expect %#v, got %#v", i, expect, actual)
		}
	}

	errTests := []string{
		"",
		"tile",
		"borower",
		"enre",
		"athor",
		"uthor",
	}

	for i, input := range errTests {
		_, err := StringToField(input)
		if err == nil {
			t.Fatalf("case %d, expect error with '%s'", i, input)
		}
	}
}

func TestSearchBookLoans(t *testing.T) {
	type Input struct{
		by      Field
		pattern string
		list    []app.BookLoan
	}
	tests := []struct{
		input  Input
		expect []int
	}{
		{ // 0
			input: Input{
				by: ByTitle,
				pattern: "cat",
				list: []app.BookLoan{
					app.BookLoan{
						Title: "cat in the hat",
					},
					app.BookLoan{
						Title: "dog in the house",
					},
					app.BookLoan{
						Title: "the cat's house",
					},
					app.BookLoan{
						Title: "category six",
					},
				},
			},
			expect: []int{
				0,
				2,
				3,
			},
		},
		{ // 1
			input: Input{
				by: ByAuthor,
				pattern: "cat",
				list: []app.BookLoan{
					app.BookLoan{
						Author: "cat in the hat",
					},
					app.BookLoan{
						Author: "dog in the house",
					},
					app.BookLoan{
						Author: "cat in the house",
					},
					app.BookLoan{
						Author: "category six",
					},
				},
			},
			expect: []int{
				0,
				2,
				3,
			},
		},
		{ // 2
			input: Input{
				by: ByGenre,
				pattern: "cat",
				list: []app.BookLoan{
					app.BookLoan{
						Genre: "cat in the hat",
					},
					app.BookLoan{
						Genre: "dog in the house",
					},
					app.BookLoan{
						Genre: "cat in the house",
					},
					app.BookLoan{
						Genre: "category six",
					},
				},
			},
			expect: []int{
				0,
				2,
				3,
			},
		},
		{ // 3
			input: Input{
				by: ByBorrower,
				pattern: "cat",
				list: []app.BookLoan{
					app.BookLoan{
						IsOnLoan: true,
						Borrower: "cat in the hat",
					},
					app.BookLoan{
						IsOnLoan: true,
						Borrower: "dog in the house",
					},
					app.BookLoan{
						IsOnLoan: true,
						Borrower: "cat in the house",
					},
					app.BookLoan{
						Borrower: "category six",
					},
				},
			},
			expect: []int{
				0,
				2,
			},
		},
	}

	for i, test := range tests{
		expect := test.expect
		input := test.input
		actual := SearchBookLoans(input.list, input.by, input.pattern)

		if len(expect) != len(actual) {
			t.Fatalf("case %d, expect length %d, got %d", i, len(expect), len(actual))
		}

		for idx, eResult := range expect {
			notFound := true
			for _, aResult := range actual {
				if aResult == eResult {
					notFound = false
					break
				}
			}
			if notFound {
				t.Fatalf("case %d, index %d, value not found %d", i, idx, eResult)
			}
		}
	}
}




