package sqlite

import (
	"testing"
	"os"

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



