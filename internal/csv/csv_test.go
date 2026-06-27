package csv

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dubbersthehoser/mayble/internal/models"
)

func compareBookEntries(e, a *models.BookEntry) error {
	if e != a {
		return fmt.Errorf("expect:\n  %#v\ngot:  %#v\n", e, a)
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
	books := []models.BookEntry{
		{
			Book: models.Book{
				Title:   "Title",
				Author:  "Author",
				Genre:   "Genre",
			},
		},
		{
			Book: models.Book{
				Title:   "Title",
				Author:  "Author",
				Genre:   "Genre",
			},
			Completed: models.Completed{
				CompletedAt:    time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:  3,
			},
			IsCompleted: true,
		},
		{
			Book: models.Book{
				Title:    "Title",
				Author:   "Author",
				Genre:    "Genre",
			},
			IsLoaned: true,
			Loaned: models.Loaned{
				LoanedAt:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
			},
		},
		{
			Book: models.Book{
				Title:    "Title",
				Author:   "Author",
				Genre:    "Genre",
			},
			IsLoaned: true,
			IsCompleted: true,

			Loaned: models.Loaned{
				LoanedAt:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
			},
			Completed: models.Completed{
				CompletedAt:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:   3,
			},
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
			if err := compareBookEntries(e, a); err != nil {
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
