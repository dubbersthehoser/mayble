package sqlite

import (
	"testing"
	"os"
	"time"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/storage"
)

func newStore(t *testing.T) (*Storage, error){
	tmp, err := os.CreateTemp("", "sqlite_*.db")
	if err != nil {
		return  nil, err
	}
	tmp.Close()

	t.Cleanup(
		func() {
			_ = os.Remove(tmp.Name())
			t.Logf("removed: %s", tmp.Name())

		},
	)
	return NewStorage(tmp.Name())
	
}

func TestErrors(t *testing.T) {
	store, err := newStore(t)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	_, err = store.GetBookByID(0)
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryNotFound, err)
	}
	_, err = store.GetBookByID(-1)
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatalf("expected error: %s, got %s", storage.ErrInvalidValue, err)
	}


	_, err = store.CreateBook(10, "title_10", "author_10", "genre_10", 0)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	_, err = store.CreateBook(10, "title_10", "author_10", "genre_10", 0)
	if !errors.Is(err, storage.ErrEntryExists) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryExists, err)
	}
	_, err = store.CreateBook(20, "title_20", "author_20", "genre_20", -1)
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = store.CreateBook(20, "title_20", "author_20", "genre_20", 6)
	if err == nil {
		t.Fatal("expected error")
	}


	
	err = store.UpdateBook(&storage.Book{
		ID: 0,
	})
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryNotFound, err)
	}
	err = store.UpdateBook(&storage.Book{
		ID: -1,
	})
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatalf("expected error: %s, got %s", storage.ErrInvalidValue, err)
	}
	err = store.UpdateBook(&storage.Book{
		ID: 10,
		Ratting: -1,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	err = store.UpdateBook(&storage.Book{
		ID: 10,
		Ratting: 6,
	})
	if err == nil {
		t.Fatal("expected error")
	}


	
	err = store.DeleteBook(&storage.Book{
		ID: 0,
	})
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryNotFound, err)
	}
	err = store.DeleteBook(&storage.Book{
		ID: -1,
	})
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatalf("expected error: %s, got %s", storage.ErrInvalidValue, err)
	}


	err = store.CreateLoan(0, "", time.Now())
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryNotFound, err)
	}
	err = store.CreateLoan(-1, "", time.Now())
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatalf("expected error: %s, got %s", storage.ErrInvalidValue, err)
	}
	err = store.CreateLoan(10, "", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	err = store.CreateLoan(10, "", time.Now())
	if !errors.Is(err, storage.ErrEntryExists) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryExists, err)
	}


	err = store.UpdateLoan(&storage.Loan{
		ID: 9,
	})
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryNotFound, err)
	}
	err = store.UpdateLoan(&storage.Loan{
		ID: 0,
	})
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryNotFound, err)
	}
	err = store.UpdateLoan(&storage.Loan{
		ID: -1,
	})
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatalf("expected error: %s, got %s", storage.ErrInvalidValue, err)
	}


	err = store.DeleteLoan(&storage.Loan{
		ID: -1,
	})
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatalf("expected error: %s, got %s", storage.ErrInvalidValue, err)
	}
	err = store.DeleteLoan(&storage.Loan{
		ID: 0,
	})
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatalf("expected error: %s, got %s", storage.ErrEntryNotFound, err)
	}


}


/*
        Book Store
*/


func TestDeleteBook(t *testing.T) {
	store, err := newStore(t)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tests := []storage.Book{
		storage.Book{
			ID: 0,
		},
		storage.Book{
			ID: 20,
		},
		storage.Book{
			ID: 44,
		},
		storage.Book{
			ID: 34,
		},
	}

	for i, test := range tests {
		_, err = store.CreateBook(
			test.ID,
			test.Title,
			test.Author,
			test.Genre,
			test.Ratting,
		)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	books, err := store.GetBooks()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	for i, book := range books {
		err = store.DeleteBook(&book)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	books, err = store.GetBooks()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(books) != 0 {
		t.Fatalf("expect length %d, got %d", 0, len(books))
	}
}

func TestUpdateBook(t *testing.T) {
	store, err := newStore(t)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tests := []struct{
		create storage.Book
		update storage.Book
	}{
		{ // case 0
			create: storage.Book{
				ID: 0,
				Title: "title_0",
				Author: "author_0",
				Genre: "genre_0",
				Ratting: 0,
			},
			update: storage.Book{
				ID: 0,
				Title: "title_0_updated",
				Author: "author_0_updated",
				Genre: "genre_0_updated",
				Ratting: 5,
			},
		},
		{ // case 0
			create: storage.Book{
				ID: 100,
				Title: "title_0",
				Author: "author_0",
				Genre: "genre_0",
				Ratting: 0,
			},
			update: storage.Book{
				ID: 100,
				Title: "",
				Author: "",
				Genre: "",
				Ratting: 0,
			},
		},
	}

	for i, test := range tests {
		book := test.create
		_, err := store.CreateBook(book.ID, book.Title, book.Author, book.Genre, book.Ratting)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	for i, test := range tests {
		expect := test.update
		err := store.UpdateBook(&expect)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}

		actual, err := store.GetBookByID(expect.ID)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}

		if expect.Title != actual.Title {
			t.Fatalf("case %d, expect '%s', got '%s'", i, expect.Title, actual.Title)
		}
		if expect.Author != actual.Author {
			t.Fatalf("case %d, expect '%s', got '%s'", i, expect.Author, actual.Author)
		}
		if expect.Genre != actual.Genre {
			t.Fatalf("case %d, expect '%s', got '%s'", i, expect.Genre, actual.Genre)
		}
		if expect.Ratting != actual.Ratting {
			t.Fatalf("case %d, expect  %d, got %d", i, expect.Ratting, actual.Ratting)
		}
		if expect.ID != actual.ID {
			t.Fatalf("case %d, expect  %d, got %d", i, expect.ID, actual.ID)
		}

	}


}

func TestCreateBook(t *testing.T) {
	store, err := newStore(t)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}


	tests := []struct{
		book storage.Book
		expect int64
	}{
		{ // case 0
			storage.Book{
				ID: storage.ZeroID,
				Title: "title_0",
				Author: "author_0",
				Genre: "genre_0",
				Ratting: 0,
			},
			1,
		},
		{ // case 1
			storage.Book{
				ID: 101,
				Title: "title_0",
				Author: "author_0",
				Genre: "genre_0",
				Ratting: 0,
			},
			101,
		},
	}



	for i, test := range tests {
		actual, err := store.CreateBook(
			test.book.ID,
			test.book.Title,
			test.book.Author,
			test.book.Genre,
			test.book.Ratting,
		)

		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		if actual != test.expect {
			t.Fatalf("case %d, expect %d, got %d", i, actual, test.expect)
		}
	}
}


/*
        Loan Store
*/


func TestCreateLoan(t *testing.T) {
	store, err := newStore(t)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tests := []storage.Loan{
		storage.Loan{
			ID: 0,
			Borrower: "borrower_0",
			Date: time.Date(2012, time.Month(4), 9, 0, 0, 0, 0, time.UTC),
		},
	}

	for i, test := range tests {
		_, err := store.CreateBook(test.ID, "", "", "", 0)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		err = store.CreateLoan(test.ID, test.Borrower, test.Date)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}
}

func TestUpdateLoan(t *testing.T) {
	store, err := newStore(t)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	
	tests := []struct{
		create storage.Loan
		update storage.Loan
	}{
		{ // case 0
			create: storage.Loan{
				ID: 0,
				Borrower: "borrower_0",
				Date: time.Date(2012, time.Month(8), 11, 0, 0, 0, 0, time.UTC),
			},
			update: storage.Loan{
				ID: 0,
				Borrower: "borrower_0_updated",
				Date: time.Date(2022, time.Month(3), 11, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	
	for i, test := range tests {
		create := test.create
		_, err := store.CreateBook(create.ID, "", "", "", 0)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		err = store.CreateLoan(create.ID, create.Borrower, create.Date)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	for i, test := range tests {
		update := test.update
		err := store.UpdateLoan(&update)

		actual, err := store.GetLoan(update.ID)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}

		if actual.Borrower != update.Borrower {
			t.Fatalf("case %d, expect '%s', got '%s'", i, update.Borrower, actual.Borrower)
		}
		if !actual.Date.Equal(update.Date) {
			t.Fatalf("case %d, expect '%#v', got '%#v'", i, update.Date, actual.Date)
		}
	}
}

func TestDeleteLoan(t *testing.T) {
	store, err := newStore(t)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tests := []storage.Loan{
		storage.Loan{
			ID: 0,
		},
		storage.Loan{
			ID: 1,
		},
		storage.Loan{
			ID: 2,
		},
		storage.Loan{
			ID: 3,
		},
	}

	for i, test := range tests {
		_, err := store.CreateBook(test.ID, "", "", "", 0)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		err = store.CreateLoan(test.ID, test.Borrower, test.Date)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	for i, test := range tests {
		err := store.DeleteLoan(&test)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		_, err = store.GetLoan(test.ID)
		if !errors.Is(err, storage.ErrEntryNotFound) {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}
}


