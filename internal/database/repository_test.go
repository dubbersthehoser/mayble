package database

import (
	"testing"
	"time"

	"github.com/dubbersthehoser/mayble/internal/models"
)

func TestDatabase(t *testing.T) {

	db, err := OpenMem()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	t.Run("create_book", func(t *testing.T) {
		testDatabaseCreateBook(db, t)
	})

	t.Run("get_unique_genres", func(t *testing.T) {
		genres, err := db.GetUniqueGenres()
		if err != nil {
			t.Fatalf("unexpected error %s", err)
		}
		if len(genres) != 1 {
			t.Fatalf("expect 1, got %d", len(genres))
		}
		expect := "A_Genre"
		if expect != genres[0] {
			t.Fatalf("expect '%s', got '%s'", expect, genres[0])
		}

	})

	t.Run("update_book", func(t *testing.T) {
		testDatabaseUpdateBook(db, t)
	})

	t.Run("delete_book", func(t *testing.T) {
		testDatabaseDeleteBook(db, t)
	})

}

func testDatabaseDeleteBook(db *Database, t *testing.T) {

	books, err := db.GetAllBooks()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(books) != 4 {
		t.Fatalf("expect length 4, got %d", len(books))
	}

	for _, book := range books {
		err := db.DeleteBook(book.ID)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	actual, err := db.GetAllBooks()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(actual) != 0 {
		t.Fatalf("expect length 0, got %d", len(actual))
	}

}

func testDatabaseUpdateBook(db *Database, t *testing.T) {
	tests := []struct {
		name    string
		input   models.BookEntry
		willErr bool
	}{
		{
			name: "a simple book",
			input: models.BookEntry{
				ID: 1,
				Book: models.Book{
					Title:  "a_title_update",
					Author: "a_author_update",
					Genre:  "a_genre_update",
				},
			},

			willErr: false,
		},
		{
			name: "a loaned book",
			input: models.BookEntry{
				ID:       2,
				IsLoaned: true,
				Book: models.Book{
					Title:  "a_title_update",
					Author: "a_author_update",
					Genre:  "a_genre_update",
				},
				Loaned: models.Loaned{
					LoanedAt: time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC),
					Borrower: "Bob",
				},
			},
			willErr: false,
		},
		{
			name: "a read book",
			input: models.BookEntry{
				ID:          3,
				IsCompleted: true,
				Book: models.Book{
					Title:  "a_title_update",
					Author: "a_author_update",
					Genre:  "a_genre_update",
				},
				Completed: models.Completed{
					CompletedAt: time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC),
					Rating:      1,
				},
			},
			willErr: false,
		},
		{
			name: "remove loaned and read from book",
			input: models.BookEntry{
				ID: 4,
				Book: models.Book{
					Title:  "a_title_update",
					Author: "a_author_update",
					Genre:  "a_genre_update",
				},
			},
			willErr: false,
		},
		{
			name: "add loaned and read back book",
			input: models.BookEntry{
				ID:          4,
				IsLoaned:    true,
				IsCompleted: true,

				Book: models.Book{
					Title:  "a_title_update",
					Author: "a_author_update",
					Genre:  "a_genre_update",
				},
				Completed: models.Completed{
					CompletedAt: time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC),
					Rating:      1,
				},
				Loaned: models.Loaned{
					LoanedAt: time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC),
					Borrower: "Bob",
				},
			},
			willErr: false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			err := db.UpdateBook(&c.input)
			if c.willErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			abook, err := db.GetBookByID(c.input.ID)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			ebook := c.input
			if abook != ebook {
				t.Fatalf("expect\n%v\n  got\n%v", ebook, abook)
			}
		})

	}
}

func testDatabaseCreateBook(db *Database, t *testing.T) {

	tests := []struct {
		name    string
		input   models.BookEntry
		expect  int64
		willErr bool
	}{
		{
			name: "a simple book",
			input: models.BookEntry{
				Book: models.Book{
					Title:  "A_Title",
					Author: "A_Author",
					Genre:  "A_Genre",
				},
			},
			expect:  1,
			willErr: false,
		},
		{
			name: "a loaned book",
			input: models.BookEntry{
				IsLoaned: true,
				Book: models.Book{
					Title:  "A_Title",
					Author: "A_Author",
					Genre:  "A_Genre",
				},
				Loaned: models.Loaned{
					LoanedAt: time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC),
					Borrower: "Lane",
				},
			},
			expect:  2,
			willErr: false,
		},
		{
			name: "a read book",
			input: models.BookEntry{
				IsCompleted: true,
				Book: models.Book{
					Title:  "A_Title",
					Author: "A_Author",
					Genre:  "A_Genre",
				},
				Completed: models.Completed{
					CompletedAt: time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC),
					Rating:      5,
				},
			},
			expect:  3,
			willErr: false,
		},
		{
			name: "a loaned and read book",
			input: models.BookEntry{
				IsCompleted: true,
				IsLoaned:    true,
				Book: models.Book{
					Title:  "A_Title",
					Author: "A_Author",
					Genre:  "A_Genre",
				},
				Loaned: models.Loaned{
					LoanedAt: time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC),
					Borrower: "Lane",
				},
				Completed: models.Completed{
					CompletedAt: time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC),
					Rating:      5,
				},
			},
			expect:  4,
			willErr: false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			id, err := db.CreateBook(&c.input)
			if c.willErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if id != c.expect {
				t.Fatalf("expect %d, got %d", c.expect, id)
			}

			abook, err := db.GetBookByID(id)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			ebook := c.input
			ebook.ID = id
			if abook != ebook {
				t.Fatalf("expect\n%v\n  got\n%v", ebook, abook)
			}
		})
	}

}
