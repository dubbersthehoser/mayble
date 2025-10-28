package csv

import (
	"testing"
	"bytes"
	"time"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

func TestExportBooks(t *testing.T) {
	type testCase struct {
		books []storage.BookLoan
		csv   string
	}
	cases := []testCase {
		testCase{
			books: []storage.BookLoan{
				storage.BookLoan{
					Book: storage.Book{
						Title: "Test Title",
						Author: "Test Author",
						Genre: "Test Genre",
						Ratting: 5,
					},
				},

			},
			csv: "Test Title,Test Author,Test Genre,5,,\n",
		},
		testCase{
			books: []storage.BookLoan{
				storage.BookLoan{
					Book: storage.Book{
						Title: "Test Title",
						Author: "Test Author",
						Genre: "Test Genre",
						Ratting: 4,
					},
					Loan: &storage.Loan{
						Name: "Brian",
						Date: time.Date(2020, time.Month(9), 21, 0, 0, 0, 0, time.Local),
					},
				},


			},
			csv: "Test Title,Test Author,Test Genre,4,Brian,2020-09-21\n",
		},
		testCase{
			books: []storage.BookLoan{
				storage.BookLoan{
					Book: storage.Book{
						Title: "Test Title",
						Author: "Test Author",
						Genre: "Test Genre",
						Ratting: 4,
					},
					Loan: &storage.Loan{
						Name: "Brian",
						Date: time.Date(2020, time.Month(9), 21, 0, 0, 0, 0, time.Local),
					},
				},
				storage.BookLoan{
					Book: storage.Book{
						Title: "My Title",
						Author: "My Author",
						Genre: "My Genre",
						Ratting: 2,
					},
				},


			},
			csv: "Test Title,Test Author,Test Genre,4,Brian,2020-09-21\nMy Title,My Author,My Genre,2,,\n",
		},
	}


	for i, _case := range cases {
		b := make([]byte, 0)
		buf := bytes.NewBuffer(b)
		exporter := BookLoanCSV{}
		err := exporter.ExportBooks(buf, _case.books)
		if err != nil {
			t.Fatalf("case %d, failed to export: %s", i, err)
		}

		if buf.String() != _case.csv {
			t.Fatalf("case %d, expect:\n'%s',\ngot:\n'%s'", i, _case.csv, buf.String())
		}

	}
}


func TestImportBooks(t *testing.T) {
	input := "The Title,The Author,The Genre,3,,\nMy Title,My Author,My Genre,5,John,2021-01-02\n"
	expects := []storage.BookLoan{
		storage.BookLoan{
			Book: storage.Book{
				Title: "The Title",
				Author: "The Author",
				Genre: "The Genre",
				Ratting: 3,
			},
			Loan: &storage.Loan{
				Name: "",
				Date: time.Time{},
			},
		},
		storage.BookLoan{
			Book: storage.Book{
				Title: "My Title",
				Author: "My Author",
				Genre: "My Genre",
				Ratting: 5,
			},
			Loan: &storage.Loan{
				Name: "John",
				Date: time.Date(2021, time.Month(1), 2, 0, 0, 0, 0, time.UTC),
			},
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

		if actual.Loan.Name != expect.Loan.Name {
			t.Fatalf("entry %d, expect ratting '%s', got '%s'", i, expect.Loan.Name, actual.Loan.Name)
		}
		if !actual.Loan.Date.Equal(expect.Loan.Date) {
			t.Fatalf("entry %d, expect ratting '%v', got '%v'", i, expect.Loan.Date, actual.Loan.Date)
		}

	}
}








