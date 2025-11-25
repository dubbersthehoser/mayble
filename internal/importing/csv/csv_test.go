package csv

import (
	"testing"
	"bytes"
	"time"
	"github.com/dubbersthehoser/mayble/internal/app"
)

func TestExportBooks(t *testing.T) {
	tests := []struct{
		input []app.BookLoan
		expect string
	}{
		{ // case 0
			input: []app.BookLoan{
				app.BookLoan{
					Title: "Test Title",
					Author: "Test Author",
					Genre: "Test Genre",
					Ratting: 5,
					
				},
			},
			expect: "Test Title,Test Author,Test Genre,5,,\n",
		},
		{ // case 1
			input: []app.BookLoan{
				app.BookLoan{
					Title: "Test Title",
					Author: "Test Author",
					Genre: "Test Genre",
					Ratting: 4,
					IsOnLoan: true,
					Borrower: "Brian",
					Date: time.Date(2020, time.Month(9), 21, 0, 0, 0, 0, time.Local),
				},
			},
			expect: "Test Title,Test Author,Test Genre,4,Brian,2020-09-21\n",
		},
		{ // case 2
			input: []app.BookLoan{
				app.BookLoan{
					Title: "Test Title",
					Author: "Test Author",
					Genre: "Test Genre",
					Ratting: 4,
					IsOnLoan: true,
					Borrower: "Brian",
					Date: time.Date(2020, time.Month(9), 21, 0, 0, 0, 0, time.Local),
				},
				app.BookLoan{
					Title: "My Title",
					Author: "My Author",
					Genre: "My Genre",
					Ratting: 2,
				},
			},
			expect: "Test Title,Test Author,Test Genre,4,Brian,2020-09-21\nMy Title,My Author,My Genre,2,,\n",
		},
		{ // case 3
			input: []app.BookLoan{
				app.BookLoan{
					Title: "Title_1",
					Author: "Author_1",
					Genre: "Genre_1",
					Ratting: 0,
					IsOnLoan: false,
					Borrower: "Borrower_Null",
					Date: time.Date(2222, time.Month(2), 2, 0, 0, 0, 0, time.UTC),
				},
			},
			expect: "Title_1,Author_1,Genre_1,0,,\n",
		},
	}


	for i, _case := range tests {
		b := make([]byte, 0)
		buf := bytes.NewBuffer(b)
		exporter := BookLoanCSV{}
		err := exporter.ExportBooks(buf, _case.input)
		if err != nil {
			t.Fatalf("case %d, failed to export: %s", i, err)
		}

		if buf.String() != _case.expect {
			t.Fatalf("case %d, expect:\n'%s',\ngot:\n'%s'", i, _case.expect, buf.String())
		}

	}
}


func TestImportBooks(t *testing.T) {
	input := "The Title,The Author,The Genre,3,,\nMy Title,My Author,My Genre,5,John,2021-01-02\nThe Title,The Author,The Genre,2,,"
	expects := []app.BookLoan{
		app.BookLoan{
			Title: "The Title",
			Author: "The Author",
			Genre: "The Genre",
			Ratting: 3,
			IsOnLoan: false,
		},
		app.BookLoan{
			Title: "My Title",
			Author: "My Author",
			Genre: "My Genre",
			Ratting: 5,
			IsOnLoan: true ,
			Borrower: "John",
			Date: time.Date(2021, time.Month(1), 2, 0, 0, 0, 0, time.UTC),
		},
		app.BookLoan{
			Title: "The Title",
			Author: "The Author",
			Genre: "The Genre",
			Ratting: 2,
			IsOnLoan:  false,
			Borrower: "",
			Date: time.Time{},
		},
	}
	buf := bytes.NewBuffer([]byte(input))
	importer := BookLoanCSV{}

	books, err := importer.ImportBooks(buf)
	if err != nil {
		t.Fatalf("failed to inport books: %s", err)
	}

	var (
		expectLength int = len(expects)
		actualLength int = len(books)
	)
	if expectLength != actualLength {
		t.Fatalf("expected length '%d', got '%d'", expectLength, actualLength)
	}
	for i := range books {
		actual := books[i]
		expect := expects[i]
		
		if actual.Title != expect.Title {
			t.Fatalf("entry %d, expect title '%s', got '%s'", i, expect.Title, actual.Title)
		}
		if actual.Author != expect.Author {
			t.Fatalf("entry %d, expect author '%s', got '%s'", i, expect.Author, actual.Author)
		}
		if actual.Genre != expect.Genre {
			t.Fatalf("entry %d, expect genre '%s', got '%s'", i, expect.Genre, actual.Genre)
		}
		if actual.Ratting != expect.Ratting {
			t.Fatalf("entry %d, expect ratting '%d', got '%d'", i, expect.Ratting, actual.Ratting)
		}

		if actual.IsOnLoan != expect.IsOnLoan {
			t.Fatalf("entry %d, expect is on loan '%t', got '%t'", i, expect.IsOnLoan, actual.IsOnLoan)
		}

		if actual.Borrower != expect.Borrower {
			t.Fatalf("entry %d, expect borrower '%s', got '%s'", i, expect.Borrower, actual.Borrower)
		}
		if !actual.Date.Equal(expect.Date) {
			t.Fatalf("entry %d, expect date '%v', got '%v'", i, expect.Date, actual.Date)
		}

	}
}
