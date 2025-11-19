package app

import (
	"time"
	"testing"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/storage/memory"
)

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
	if actual.Ratting != expect.Ratting {
		return fmt.Errorf("expect ratting %d, got %d", expect.Ratting, actual.Ratting)
	}
	if actual.Borrower != expect.Borrower {
		return fmt.Errorf("expect borrower %s, got %s", expect.Borrower, actual.Borrower)
	}
	if !actual.Date.Equal(expect.Date) {
		return fmt.Errorf("expect date %v, got %v", expect.Date, actual.Date)
	}
	return nil
}

func TestUpdateBookLoan(t *testing.T) {
	store := memory.NewStorage()
	date := time.Now()
	tests := []struct{
		updated BookLoan
		created BookLoan
	}{
		{
			updated: BookLoan{
				Title: "title_1",
				Author: "author_1",
				Genre: "genre_1",
				IsOnLoan: true,
				Ratting: 2,
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
	}

	for i, test := range tests {
		create := test.created
		id, err := createBookLoan(store, &create)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		test.updated.ID = id
	}

	for i, test := range tests {
		update := test.update
		err := updateBookLoan(store, &update)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}

		actual, err := getBookLoanByID(update.ID)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		err = compBookLoan(&update, actual)
		if err != nil {
			t.Fatalf("case %d, %s", i, err)
		}
	}
	
}


