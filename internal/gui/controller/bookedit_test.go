package controller

import (
	"testing"
	"time"

	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/memdb"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

func TestBookLoanBuilder(t *testing.T) {
	
	type testCase struct {
		expect storage.BookLoan
		inputs map[string]string
	}
	
	cases := []testCase{
		testCase{
			expect: storage.BookLoan{
				Book: storage.Book{
					Title: "My Title",
					Author: "My Author",
					Genre: "My Genre",
					Ratting: 0,
				},
				Loan: nil,
			},
			inputs: map[string]string {
				"title": "My Title",
				"author": "My Author",
				"genre": "My Genre",
				"ratting": rattingToString(0),
			},
		},
		testCase{
			expect: storage.BookLoan{
				Book: storage.Book{
					Title: "My Second Title",
					Author: "My Second Author",
					Genre: "My Second Genre",
					Ratting: 4,
				},
				Loan: &storage.Loan{
					Name: "My Borrower",
					Date: time.Date(2020, time.Month(12), 1, 0, 0, 0, 0, time.Local),
				},
			},
			inputs: map[string]string {
				"title": "My Second Title",
				"author": "My Second Author",
				"genre": "My Second Genre",
				"ratting": rattingToString(4),
				"borrower": "My Borrower",
				"date": "2020-12-01",
			},
		},
	}
	
	for i, _case := range cases {
		builder := NewBookLoanBuilder(Creating)
		for key, value := range _case.inputs {
			switch key {
			case "title":
				builder.SetTitle(value)
			case "author":
				builder.SetAuthor(value)
			case "genre":
				builder.SetGenre(value)
			case "ratting":
				builder.SetRattingAsString(value)
			case "borrower":
				builder.SetBorrower(value)
			case "date":
				builder.SetDateAsString(value)
			default:
				panic("invalid case input value")
			}
		}

		actual := builder.Build()
		expect := _case.expect

		if actual.Title != expect.Title {
			t.Errorf("case: %d, with title '%s', got '%s' ", i, expect.Title, actual.Title)
		}
		if actual.Author != expect.Author {
			t.Errorf("case: %d, with author '%s', got '%s' ", i, expect.Author, actual.Author)
		}
		if actual.Genre != expect.Genre {
			t.Errorf("case: %d, with genre '%s', got '%s' ", i, expect.Genre, actual.Genre)
		}
		if actual.Ratting != expect.Ratting {
			t.Errorf("case: %d, with ratting '%d', got '%d' ", i, expect.Ratting, actual.Ratting)

		} else if  actual.Loan != nil  {
			if actual.Loan.Name != expect.Loan.Name {
				t.Errorf("case: %d, with borrower '%s', got '%s' ", i, expect.Loan.Name, actual.Loan.Name)
			}
			if actual.Loan.Date.Equal(expect.Loan.Date) {
				t.Errorf("case: %d, with date '%v', got '%v' ", i, expect.Loan.Name, actual.Loan.Name)
			}
		}
	}

}

func TestBookEditor(t *testing.T) {
	
	store := memdb.NewMemStorage()
	core, err := core.New(store)
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		expect storage.BookLoan
		input  BookLoanBuilder
	}

	cases := []testCase{
		testCase{
			expect: storage.BookLoan{
				Book: storage.Book{
					Title: "My Title",
					Author: "My Author",
					Genre: "My Genre",
					Ratting: 3,
				},
				Loan: &storage.Loan{
					Name: "My Borrower",
					Date: time.Date(1999, time.Month(01), 11, 0, 0, 0, 0, time.Local),
				},
			},
			input: BookLoanBuilder{
				title: "My Title",
				author: "My Author",
				genre: "My Genre",
				ratting: 3,
				isOnLoan: true,
				borrower: "My Borrower",
				date: time.Date(1999, time.Month(01), 11, 0, 0, 0, 0, time.Local),
			},
		},
	}

	bookEditor := NewBookEditor(core)

	for i, _case := range cases {
		expect := _case.expect
		input := _case.input
		input.Type = Creating
		expect.ID = int64(i)
		input.id = expect.ID

		err := bookEditor.Submit(&input)
		if err != nil {
			t.Fatal(err)
		}

		err = core.Save()
		if err != nil {
			t.Fatal(err)
		}

		_, err = store.GetBookLoanByID(expect.ID)

		if err != nil {
			t.Errorf("case %d, errored: %s", i,  err)
		}
	}
}







