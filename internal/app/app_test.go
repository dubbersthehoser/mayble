package app

import (
	"time"
	"testing"
	"fmt"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/storage/memory"
	"github.com/dubbersthehoser/mayble/internal/storage"
)


func TestAppCreateBookLoan(t *testing.T) {
	store := memory.NewStorage()
	app, err := New(store)
	if err != nil {
		t.Fatalf("unexpeced error: %s", err)
	}

	date := time.Now()

	tests := []BookLoan{
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			Ratting: 0,
			IsOnLoan: false,
		},
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			Ratting: 1,
			IsOnLoan: false,
		},
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_3",
			Author: "author_3",
			Genre: "genre_3",
			Ratting: 2,
			IsOnLoan: false,
			Borrower: "borrower_3",
			Date: date,

		},
	}

	for i, test := range tests {
		err := app.CreateBookLoan(&test)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	bookLoans, err := app.GetBookLoans()
	if err != nil {
		t.Fatalf("unexpeced error: %s", err)
	}

	if len(bookLoans) != len(tests) {
		t.Fatalf("expect length %d, got %d", len(tests), len(bookLoans))
	}
}

func TestAppError(t *testing.T) {
	store := memory.NewStorage()
	app, err := New(store)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	err = app.UpdateBookLoan(nil)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatal("error value was not catched")
	}

	err = app.CreateBookLoan(nil)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatal("error value was not catched")
	}

	err = app.DeleteBookLoan(nil)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatal("error value was not catched")
	}

	err = app.ImportBookLoans(nil)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatal("error value was not catched")
	}
}

func TestAppDeleteBookLoan(t *testing.T) {
	store := memory.NewStorage()
	app, err := New(store)
	if err != nil {
		t.Fatalf("unexpeced error: %s", err)
	}

	date := time.Now()
	tests := []BookLoan{
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			Ratting: 0,
			IsOnLoan: false,
		},
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			Ratting: 1,
			IsOnLoan: true,
			Borrower: "borrower_2",
			Date: date,

		},
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_3",
			Author: "author_3",
			Genre: "genre_3",
			Ratting: 2,
			IsOnLoan: false,
		},
	}

	for i, test := range tests {
		err := app.CreateBookLoan(&test)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	bookLoans, err := app.GetBookLoans()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	for i, bookLoan := range bookLoans {
		err := app.DeleteBookLoan(&bookLoan)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}
	
	bookLoans, err = app.GetBookLoans()
	if len(bookLoans) != 0 {
		t.Fatalf("expected length %d, got %d", 0, len(bookLoans))
	}
}

func TestAppUpdateBookLoan(t *testing.T) {
	store := memory.NewStorage()
	app, err := New(store)
	if err != nil {
		t.Fatalf("unexpeced error: %s", err)
	}

	date := time.Now()
	
	tests := []struct{
		create BookLoan
		update BookLoan
	}{
		{
			create: BookLoan{
				ID: storage.ZeroID,
				Title: "title_1",
				Author: "author_1",
				Genre: "genre_1",
				Ratting: 1,
			},
			update: BookLoan{
				Title: "new_title_1",
				Author: "new_author_1",
				Genre: "new_genre_1",
				Ratting: 2,
				IsOnLoan: true,
				Borrower: "borrower_1",
				Date: date,
			},
		},
		{
			create: BookLoan{
				ID: storage.ZeroID,
				Title: "new_title_2",
				Author: "new_author_2",
				Genre: "new_genre_2",
				Ratting: 3,
				IsOnLoan: true,
				Borrower: "borrower_2",
				Date: date,
			},
			update: BookLoan{
				Title: "title_2",
				Author: "author_2",
				Genre: "genre_2",
				Ratting: 2,
			},
		},
		{
			create: BookLoan{
				ID: 21,
				Title: "title_3",
				Author: "author_3",
				Genre: "genre_3",
				Ratting: 4,
				IsOnLoan: true,
				Borrower: "borrower_3",
				Date: date.Add(time.Hour * 24),
			},
			update: BookLoan{
				ID: 31,
				Title: "new_title_3",
				Author: "new_author_3",
				Genre: "new_genre_3",
				Ratting: 5,
				IsOnLoan: true,
				Borrower: "new_borrower_3",
				Date: date.Add(time.Hour * 24),
			},
		},
	}

	for i, test := range tests {
		id, err := createBookLoan(app.memory, &test.create)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		tests[i].update.ID = id
		if storage.IsZeroID(tests[i].create.ID) {
			tests[i].create.ID = id
		}
	}

	for i, test := range tests {
		err := app.UpdateBookLoan(&test.update)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		expect := &test.update
		actual, err := getBookLoanByID(app.memory, test.update.ID)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		err = compBookLoans(expect, actual)
		if err != nil {
			t.Fatalf("case %d, %s", i, err)
		}
	}

	for !app.UndoIsEmpty() {
		app.Undo()
	}
	for i, test := range tests {
		expect := &test.create
		actual, err := getBookLoanByID(app.memory, test.update.ID)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		err = compBookLoans(expect, actual)
		if err != nil {
			t.Fatalf("case %d, %s", i, err)
		}
	}

	for !app.RedoIsEmpty() {
		app.Redo()
	}
	for i, test := range tests {
		expect := &test.update
		actual, err := getBookLoanByID(app.memory, test.update.ID)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		err = compBookLoans(expect, actual)
		if err != nil {
			t.Fatalf("case %d, %s", i, err)
		}
	}

}

func TestAppImportBookLoans(t *testing.T) {
	store := memory.NewStorage()
	app, err := New(store)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	date := time.Now()

	test := []BookLoan{
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			Ratting: 1,
			IsOnLoan: true,
			Borrower: "borrower_1",
			Date: date,
		},
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			Ratting: 2,
		},
	}

	err = app.ImportBookLoans(test)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	bookLoans, err := app.GetBookLoans()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(bookLoans) != len(test) {
		t.Fatalf("expect length %d, got %d", len(test), len(bookLoans))
	}

	app.Undo()

	bookLoans, err = app.GetBookLoans()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(bookLoans) != 0 {
		t.Fatalf("expect length %d, got %d", 0, len(bookLoans))
	}
	app.Redo()
	bookLoans, err = app.GetBookLoans()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(bookLoans) != len(test) {
		t.Fatalf("expect length %d, got %d", len(test), len(bookLoans))
	}
}

func TestAppSave(t *testing.T) {
	store := memory.NewStorage()
	app, err := New(store)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	date := time.Now()

	test := []BookLoan{
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			Ratting: 1,
			IsOnLoan: true,
			Borrower: "borrower_1",
			Date: date,
		},
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			Ratting: 2,
		},
	}

	for _, test := range test {
		err := app.CreateBookLoan(&test)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	expect, err := getAllBookLoans(app.memory)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	app.Save()

	expect, err = getAllBookLoans(app.memory)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	actual, err := getAllBookLoans(app.storage)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(expect) != len(actual) {
		t.Fatalf("expect length %d, got %d", len(expect), len(actual))
	}
}

func TestAppLoad(t *testing.T) {
	store := memory.NewStorage()
	app, err := New(store)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	date := time.Now()

	test := []BookLoan{
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			Ratting: 1,
			IsOnLoan: true,
			Borrower: "borrower_1",
			Date: date,
		},
		BookLoan{
			ID: storage.ZeroID,
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			Ratting: 2,
		},
	}

	// place book loans in to storage,
	// then load it into memory.
	//
	for i, bookLoan := range test {
		id, err := createBookLoan(app.storage, &bookLoan)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		test[i].ID = id
	}

	err = app.load()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	actual, err := app.GetBookLoans()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(actual) != len(test) {
		t.Fatalf("expect length %d, got %d", len(test), len(actual))
	}
}


















func compBookLoans(expect, actual *BookLoan) error {
	if actual.ID != expect.ID {
		println("actual:", actual.ID)
		println("expect:", expect.ID)
		return fmt.Errorf("expect id %d, got %d", expect.ID, actual.ID)
	}
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
			ID: storage.ZeroID,
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			IsOnLoan: true,
			Ratting: 1,
			Borrower: "borrower_1",
			Date: date,
		},
		BookLoan{
			ID: storage.ZeroID,
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
				ID: storage.ZeroID,
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
				ID: storage.ZeroID,
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
				ID: storage.ZeroID,
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
				ID: storage.ZeroID,
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
