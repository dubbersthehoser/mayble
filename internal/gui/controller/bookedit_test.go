package controller

import (
	"testing"
	"time"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/storage/memory"
)

func TestBookLoanBuilder(t *testing.T) {
	
	type testCase struct {
		expect app.BookLoan
		inputs map[string]string
	}
	
	cases := []testCase{
		testCase{
			expect: app.BookLoan{
				Title: "My Title",
				Author: "My Author",
				Genre: "My Genre",
				Ratting: 0,
			},
			inputs: map[string]string {
				"title": "My Title",
				"author": "My Author",
				"genre": "My Genre",
				"ratting": listing.MustRattingToString(0),
			},
		},
		testCase{
			expect: app.BookLoan{
				Title: "My Second Title",
				Author: "My Second Author",
				Genre: "My Second Genre",
				Ratting: 4,
				IsOnLoan: true,
				Borrower: "My Borrower",
				Date: time.Date(2020, time.Month(12), 1, 0, 0, 0, 0, time.UTC),
			},
			inputs: map[string]string {
				"title": "My Second Title",
				"author": "My Second Author",
				"genre": "My Second Genre",
				"ratting": listing.MustRattingToString(4),
				"borrower": "My Borrower",
				"date": "2020-12-01",
			},
		},
	}
	
	for i, _case := range cases {
		builder := NewBookLoanBuilder()
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

		if !(builder.Date.IsZero() && builder.Borrower == "") {
			builder.IsOnLoan = true
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

		} 
		if actual.IsOnLoan != expect.IsOnLoan {
			t.Errorf("case: %d, expect '%t', got '%t' ", i, expect.IsOnLoan, actual.IsOnLoan)
		}

		if actual.IsOnLoan {
			if actual.Borrower != expect.Borrower {
				t.Errorf("case: %d, expect '%s', got '%s' ", i, expect.Borrower, actual.Borrower)
			}
			if !actual.Date.Equal(expect.Date) {
				t.Errorf("case: %d, date '%v', got '%v' ", i, expect.Date, actual.Date)
			}
		}
	}

}

func TestBookEditor(t *testing.T) {
	
	store := memory.NewStorage()
	a := app.New(store)

	type testCase struct {
		expect app.BookLoan
		input  BookLoanBuilder
	}

	cases := []testCase{
		testCase{
			expect: app.BookLoan{
				Title: "My Title",
				Author: "My Author",
				Genre: "My Genre",
				Ratting: 3,
				IsOnLoan: true,
				Borrower: "My Borrower",
				Date: time.Date(1999, time.Month(01), 11, 0, 0, 0, 0, time.Local),
			},
			input: BookLoanBuilder{
				Title: "My Title",
				Author: "My Author",
				Genre: "My Genre",
				Ratting: 3,
				IsOnLoan: true,
				Borrower: "My Borrower",
				Date: time.Date(1999, time.Month(01), 11, 0, 0, 0, 0, time.Local),
			},
		},
	}

	broker := &emiter.Broker{}

	bookEditor := NewBookEditer(broker, a)

	for i, _case := range cases {
		expect := _case.expect
		input := _case.input
		expect.ID = int64(i)
		input.id = expect.ID

		input.Type = Creating

		err := bookEditor.Submit(&input)
		if err != nil {
			t.Fatal(err)
		}

		_, err = store.GetBookByID(expect.ID)
		if err != nil {
			t.Errorf("case %d, errored: %s", i,  err)
			continue
		}
		_, err = store.GetLoan(expect.ID)
		if err != nil && input.IsOnLoan {
			t.Errorf("case %d, errored: %s", i,  err)
			continue
		}


		input.Type = Updating
		input.SetTitle("New Title")
		
		err = bookEditor.Submit(&input)
		if err != nil {
			t.Fatal(err)
		}

		book, err := store.GetBookByID(expect.ID)
		if err != nil {
			t.Fatalf("case %d, errored: %s", i,  err)
		}
		
		if book.Title != input.Title {
			t.Errorf("case %d, title did not update", i)
			continue
		}

		input.Type = Deleting

		err = bookEditor.Submit(&input)
		if err != nil {
			t.Fatal(err)
		}

		_, err = store.GetBookByID(expect.ID)
		if err == nil {
			t.Errorf("case %d, entry was found after deletion submission", i)
			continue
		}
		if !errors.Is(err, storage.ErrEntryNotFound) {
			t.Errorf("case %d, unexpected error after deletion submission", i)
			continue
		}
		_, err = store.GetLoan(expect.ID)
		if err == nil && input.IsOnLoan {
			t.Errorf("case %d, entry was found after deletion submission", i)
			continue
		}
		if !errors.Is(err, storage.ErrEntryNotFound) {
			t.Errorf("case %d, unexpected error after deletion submission", i)
			continue
		}
	}
}







