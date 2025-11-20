package app

import (
	"time"
	"testing"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/storage/memory"
)


func TestApp(t *testing.T) {
	store := memeory.NewStorage()
	app := NewApp(store)



}

















func compBookLoans(expect, actual *BookLoan) error {
	if actual.Title != expect.Title {
		return fmt.Errorf("expect title %s, got %s", expect.Title, actual.Author)
	}
	if actual.Author != expect.Author {
		return fmt.Errorf("expect author %s, got %s", expect.Author, actual.Author)
	}
	if actual.Genre != expect.Genre {
		return fmt.Errorf("expect genre %s, got %s", expect.Genre, actual.Genre)
	}
	if actual.Ratting != expect.Ratting {
		return fmt.Errorf("expect ratting %d, got %d", expect.Ratting, actual.Ratting)
	}
	if actual.IsOnLoan != expect.IsOnLoan {
		return fmt.Errorf("expect is on loan %t, got %t", expect.IsOnLoan, actual.IsOnLoan)
	}
	if actual.Borrower != expect.Borrower && actual.IsOnLoan {
		return fmt.Errorf("expect borrower %s, got %s", expect.Borrower, actual.Borrower)
	}
	if !actual.Date.Equal(expect.Date)  && actual.IsOnLoan {
		return fmt.Errorf("expect date %v, got %v", expect.Date, actual.Date)
	}
	return nil
}

func TestDeleteBookLoan(t *testing.T) {
	store := memory.NewStorage()
	date := time.Now()
	tests := []BookLoan{
		BookLoan{
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			IsOnLoan: true,
			Ratting: 1,
			Borrower: "borrower_1",
			Date: date,
		},
		BookLoan{
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			IsOnLoan: true,
			Ratting: 3,
			Borrower: "borrower_2",
			Date: date.Add(time.Hour * 24),
		},
	}
	for i, test := range tests {
		create := test
		id, err := createBookLoan(store, &create)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		tests[i].ID = id

	}
	bookLoans, err := getAllBookLoans(store)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(bookLoans) != len(tests) {
		t.Fatalf("expected length %d, got %d", len(tests), len(bookLoans))
	}

	for i, test := range tests {
		err := deleteBookLoan(store, &test)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}
	
	bookLoans, err = getAllBookLoans(store)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(bookLoans) != 0 {
		t.Fatalf("expected length %d, got %d", 0, len(bookLoans))
	}

}

func TestUpdateBookLoan(t *testing.T) {
	store := memory.NewStorage()
	date := time.Now()
	tests := []*struct{
		updated BookLoan
		created BookLoan
	}{
		{
			updated: BookLoan{
				Title: "title_1",
				Author: "author_1",
				Genre: "genre_1",
				Ratting: 2,
				IsOnLoan: true,
				Borrower: "borrower_1",
				Date: date,
			},
			created: BookLoan{
				Title: "new_title_1",
				Author: "new_author_1",
				Genre: "new_genre_1",
				Ratting: 3,
				IsOnLoan: true,
				Borrower: "new_borrower_1",
				Date: date.Add(time.Hour * 24),
			},
		},
		{
			updated: BookLoan{
				Title: "title_2",
				Author: "author_2",
				Genre: "genre_2",
				IsOnLoan: true,
				Ratting: 2,
				Borrower: "borrower_2",
				Date: date,
			},
			created: BookLoan{
				Title: "new_title_2",
				Author: "new_author_2",
				Genre: "new_genre_2",
				Ratting: 3,
				IsOnLoan: false,
				Borrower: "new_borrower_2",
				Date: date.Add(time.Hour * 24),
			},
		},
		{
			updated: BookLoan{
				Title: "title_3",
				Author: "author_3",
				Genre: "genre_3",
				IsOnLoan: false,
				Ratting: 3,
				Borrower: "borrower_3",
				Date: date,
			},
			created: BookLoan{
				Title: "new_title_2",
				Author: "new_author_2",
				Genre: "new_genre_2",
				Ratting: 3,
				IsOnLoan: true,
				Borrower: "new_borrower_2",
				Date: date.Add(time.Hour * 24),
			},
		},
		{
			updated: BookLoan{
				Title: "title_4",
				Author: "author_4",
				Genre: "genre_4",
				IsOnLoan: false,
				Ratting: 4,
				Borrower: "borrower_4",
				Date: date,
			},
			created: BookLoan{
				Title: "new_title_4",
				Author: "new_author_4",
				Genre: "new_genre_4",
				Ratting: 4,
				IsOnLoan: false,
				Borrower: "new_borrower_4",
				Date: date.Add(time.Hour * 24),
			},
		},
	}

	for i, test := range tests {
		create := test.created
		id, err := createBookLoan(store, &create)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		tests[i].updated.ID = id
	}

	for i, test := range tests {
		update := test.updated
		err := updateBookLoan(store, &update)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}

		actual, err := getBookLoanByID(store, update.ID)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		err = compBookLoans(&update, actual)
		if err != nil {
			t.Fatalf("case %d, %s", i, err)
		}
	}
}



