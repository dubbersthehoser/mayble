package csv

import (
	"bytes"
	"slices"
	"fmt"
	"strings"
	"testing"
	"time"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

func compareBookEntry(e, a *repo.BookEntry) error {

	if e.Variant != a.Variant {
		return fmt.Errorf("expect variant '%s', got '%s'", e.Variant, a.Variant)
	}

	if e.Title != a.Title {
		return fmt.Errorf("expect title '%s', got '%s'", e.Title, a.Title)
	}

	if e.Author != a.Author {
		return fmt.Errorf("expect author '%s', got '%s'", e.Author, a.Author)
	}

	if e.Genre != a.Genre {
		return fmt.Errorf("expect genre '%s', got '%s'", e.Genre, a.Genre)
	}

	if e.Read != a.Read {
		return fmt.Errorf("expect read '%s', got '%s'", e.Read, a.Read)
	}

	if e.Rating != a.Rating {
		return fmt.Errorf("expect rating %d, got %d", e.Rating, a.Rating)
	}

	if e.Loaned != a.Loaned {
		return fmt.Errorf("expect loaned '%s', got '%s'", e.Loaned, a.Loaned)
	}

	if e.Borrower != a.Borrower {
		return fmt.Errorf("expect borrower '%s', got '%s'", e.Borrower, a.Borrower)
	}
	return nil

}

func TestImportAndExport(t *testing.T) {
	csvStr := strings.TrimSpace(`
TITLE,AUTHOR,GENRE,RATING,READ,BORROWER,LOANED
Title,Author,Genre,,,,
Title,Author,Genre,2021-02-19,3,,
Title,Author,Genre,,,2021-02-19,Lane
Title,Author,Genre,2021-02-19,3,2021-02-19,Lane
`)
	books := []repo.BookEntry{
		{
			Variant: repo.Book,
			Title:   "Title",
			Author:  "Author",
			Genre:   "Genre",
		},
		{
			Variant: repo.Book | repo.Read,
			Title:   "Title",
			Author:  "Author",
			Genre:   "Genre",
			Read:    time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating:  3,
		},
		{
			Variant:  repo.Book | repo.Loaned,
			Title:    "Title",
			Author:   "Author",
			Genre:    "Genre",
			Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
		},
		{
			Variant:  repo.Book | repo.Loaned | repo.Read,
			Title:    "Title",
			Author:   "Author",
			Genre:    "Genre",
			Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
			Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating:   3,
		},
	}

	t.Run("Import", func(t *testing.T) {
		expect := books
		input := bytes.NewReader([]byte(csvStr))
		actual, err := Import(input)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(expect)+1 != len(actual) {
			t.Fatalf("length of expect missmatch to actual: len(expect)=%d len(actual)=%d", len(expect), len(actual))
		}

		for i := range expect {
			e := &expect[i]
			a := &actual[i]
			if err := compareBookEntry(e, a); err != nil {
				t.Fatalf("[%d] %s", i, err)
			}
		}
	})

	t.Run("Export", func(t *testing.T) {
		expect := csvStr
		input := books
		actualBuf := bytes.NewBuffer([]byte{})
		err := Export(actualBuf, input)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		actual := strings.TrimSpace(actualBuf.String())

		if expect != actual {
			t.Fatalf("expect\n'%s'\ngot\n'%s'", expect, actual)
		}
	})
}


func Test_schemaHeaders(t *testing.T) {
	expect := []string{
		"TITLE", "AUTHOR", "GENRE", "RATING", "READ", "BORROWER", "LOANED",
	}
	actual := schemaHeaders()

	if !slices.Equal(expect, actual) {
		t.Fatalf("expect\n  %#v\ngot  %#v", expect, actual)
	}
}


func Test_mapSchema(t *testing.T) {
	
	tests := []struct{
		name   string
		input  []string
		expect []int
		ok     bool
	}{
		{
			name: "base case",
			input: []string{
				"TITLE", "AUTHOR", "GENRE", "RATING", "READ", "BORROWER", "LOANED",
			},
			expect: []int{
				0,1,2,3,4,5,6,
			},
			ok: true, 
		},
		{
			name: "swaped READ and RATING",
			input: []string{
				"TITLE", "AUTHOR", "GENRE", "READ", "RATING", "BORROWER", "LOANED",
			},
			expect: []int{
				0,1,2,4,3,5,6,
			},
			ok: true, 
		},
		{
			name: "lower case header",
			input: []string{
				"title", "AUTHOR", "GENRE", "READ", "RATING", "BORROWER", "LOANED",
			},
			expect: nil,
			ok: false, 
		},
		{
			name: "missing header",
			input: []string{
				"AUTHOR", "GENRE", "READ", "RATING", "BORROWER", "LOANED",
			},
			expect: nil,
			ok: false, 
		},
	}


	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			actual, ok := mapSchema(c.input)
			if c.ok != ok {
				t.Fatalf("expect %t, got %t", c.ok, ok)
			}

			if !slices.Equal(c.expect, actual) {
				t.Fatalf("expect\n  %#v\ngot\n  %#v", c.expect, actual)
			}
		})
	}
}
