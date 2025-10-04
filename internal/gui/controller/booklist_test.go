package controller

import (
	"testing"
	"time"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/memdb"
)

func TestToBookLoanListed(t *testing.T) {
	type testCase struct {
		input  storage.BookLoan
		expect BookLoanListed
	}
	dateNow := time.Now()
	cases := []testCase{
		testCase{
			input: storage.BookLoan{
				Book: storage.Book{
					Title: "This Title",
					Author: "This Author",
					Genre: "This Genre",
					Ratting: 1,
				},
				Loan: &storage.Loan{
					Name: "This Borrower",
					Date: dateNow,
				},
			},
			expect: BookLoanListed{
				Title: "This Title",
				Author: "This Author",
				Genre: "This Genre",
				Ratting: rattingToString(1),
				Borrower: "This Borrower",
				Date: dateToString(&dateNow),
			},
		},
		testCase{
			input: storage.BookLoan{
				Book: storage.Book{
					Title: "",
					Author: "",
					Genre: "",
					Ratting: 0,
				},
			},
			expect: BookLoanListed{
				Title: "",
				Author: "",
				Genre: "",
				Ratting: rattingToString(0),
				Borrower: "n/a",
				Date: "n/a",
			},
		},
		testCase{
			input: storage.BookLoan{
				Book: storage.Book{
					Title: "My Title",
					Author: "My Author",
					Genre: "My Genre",
					Ratting: 4,
				},
			},
			expect: BookLoanListed{
				Title: "My Title",
				Author: "My Author",
				Genre: "My Genre",
				Ratting: rattingToString(4),
				Borrower: "n/a",
				Date: "n/a",
			},
		},
	}

	for i, _case := range cases{
		expect := _case.expect
		actual := toBookLoanListed(&_case.input)
		if expect.Title != actual.Title {
			t.Errorf("case %d, expect title '%s', got '%s'", i, expect.Title, actual.Title)
		}
		if expect.Author != actual.Author {
			t.Errorf("case %d, expect author '%s', got '%s'", i, expect.Author, actual.Author)
		}
		if expect.Genre != actual.Genre {
			t.Errorf("case %d, expect genre '%s', got '%s'", i, expect.Genre, actual.Genre)
		}
		if expect.Ratting != actual.Ratting {
			t.Errorf("case %d, expect ratting '%s', got '%s'", i, expect.Ratting, actual.Ratting)
		}
		if expect.Borrower != actual.Borrower {
			t.Errorf("case %d, expect borrower '%s', got '%s'", i, expect.Borrower, actual.Borrower)
		}
		if expect.Date != actual.Date {
			t.Errorf("case %d, expect date '%s', got '%s'", i, expect.Date, actual.Date)
		}
	}
}


func TestBookList(t *testing.T) {
	expectedListed := make([]BookLoanListed, 0)
	bookAmount := 5
	store := memdb.NewMemStorage()
	for i:=0; i<bookAmount; i++ {
		bookLoan := storage.BookLoan{
			Book: storage.Book{
				ID: storage.ZeroID,
				Title: fmt.Sprintf("title_%d", i),
				Author: fmt.Sprintf("author_%d", i),
				Genre: fmt.Sprintf("genre_%d", i),
				Ratting: i % 6,
			},
			Loan: &storage.Loan{
				ID: storage.ZeroID,
				Name: fmt.Sprintf("name_%d", i),
				Date: time.Now().Add(time.Hour * time.Duration(24 * (i+1))),
			},

		}
		expectedListed = append(expectedListed, *toBookLoanListed(&bookLoan))
		err := store.CreateBookLoan(&bookLoan)
		if err != nil {
			t.Fatal(err)
		}
	}
	core, err := core.New(store)
	if err != nil {
		t.Fatal(err)
	}

	bookList := NewBookList(core)
	err = bookList.Update()
	if err != nil {
		t.Fatal(err)
	}

	if len(bookList.list) != bookAmount {
		t.Fatalf("list size got '%d', got '%d'", len(bookList.list), bookAmount)
	}

	_ = expectedListed

	_, err = bookList.Get(-1)
	if err == nil {
		t.Fatalf("Get(%d) expected error, got %s", -1, err)
	}

	_, err = bookList.Get(bookAmount+1)
	if err == nil {
		t.Fatalf("Get(%d) expected error, got %s", bookAmount+1, err)
	}

	for i:=0; i<bookAmount; i++ {
		actual, err := bookList.Get(i)
		if err != nil {
			t.Fatalf("Get(%d) errored: %s", i, err)
		}
		expect := expectedListed[i]
		if expect.Title != actual.Title {
			t.Errorf("Get(%d) expect title '%s', got '%s'", i, expect.Title, actual.Title)
		}
		if expect.Author != actual.Author {
			t.Errorf("Get(%d) expect author '%s', got '%s'", i, expect.Author, actual.Author)
		}
		if expect.Genre != actual.Genre {
			t.Errorf("Get(%d) expect genre '%s', got '%s'", i, expect.Genre, actual.Genre)
		}
		if expect.Ratting != actual.Ratting {
			t.Errorf("Get(%d) expect ratting '%s', got '%s'", i, expect.Ratting, actual.Ratting)
		}
		if expect.Borrower != actual.Borrower {
			t.Errorf("Get(%d) expect borrower '%s', got '%s'", i, expect.Borrower, actual.Borrower)
		}
		if expect.Date != actual.Date {
			t.Errorf("Get(%d) expect date '%s', got '%s'", i, expect.Date, actual.Date)
		}

		err = bookList.Select(i)
		if err != nil {
			t.Fatalf("Select(%d) errored: %s", i, err)
		}
		
		actual, err = bookList.Selected()
		if err != nil {
			t.Fatalf("Selected() errored: %s", err)
		}
		if actual == nil {
			t.Fatalf("Selected() return nil: %v", actual)
		}
		if expect.Title != actual.Title {
			t.Errorf("Select(%d) expect title '%s', got '%s'", i, expect.Title, actual.Title)
		}
		if expect.Author != actual.Author {
			t.Errorf("Select(%d) expect author '%s', got '%s'", i, expect.Author, actual.Author)
		}
		if expect.Genre != actual.Genre {
			t.Errorf("Select(%d) expect genre '%s', got '%s'", i, expect.Genre, actual.Genre)
		}
		if expect.Ratting != actual.Ratting {
			t.Errorf("Select(%d) expect ratting '%s', got '%s'", i, expect.Ratting, actual.Ratting)
		}
		if expect.Borrower != actual.Borrower {
			t.Errorf("Select(%d) expect borrower '%s', got '%s'", i, expect.Borrower, actual.Borrower)
		}
		if expect.Date != actual.Date {
			t.Errorf("Select(%d) expect date '%s', got '%s'", i, expect.Date, actual.Date)
		}
	}
}
