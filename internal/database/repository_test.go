package database

import (
	"testing"
	"time"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

func Test_bookBuilder(t *testing.T) {

	t.Run("invalid builds", func(t *testing.T) {
		builder := newBookBuilder()
		builder.SetID(432)
		builder.SetTitle("A Title")
		builder.SetAuthor("A Author")
		builder.SetGenre("A Genre")
		_, err := builder.Build()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		builder.SetID(-127)
		_, err = builder.Build()
		if err == nil {
			t.Fatalf("expected error")
		}

		builder.SetID(324)
		builder.SetTitle("")
		_, err = builder.Build()
		if err == nil {
			t.Fatalf("expected error")
		}
		builder.SetTitle("A Title")
		builder.SetAuthor("")
		_, err = builder.Build()
		if err == nil {
			t.Fatalf("expected error")
		}
		builder.SetAuthor("A Author")
		builder.SetGenre("")
		_, err = builder.Build()
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("build", func(t *testing.T) {
		testBuilderBuild(t)
	})
}

func testBuilderBuild(t *testing.T) {
	type input struct {
		id int64
		title,
		author,
		genre,
		loanDate,
		borrower,
		readDate string
		rating int64
	}
	tests := []struct {
		name   string
		input  input
		expect repo.BookEntry
	}{
		{
			name: "just a book",
			input: input{
				id:     132,
				title:  "A Title",
				author: "A Author",
				genre:  "A Genre",
			},
			expect: repo.BookEntry{
				ID:     132,
				Title:  "A Title",
				Author: "A Author",
				Genre:  "A Genre",
			},
		},
		{
			name: "just loaned book",
			input: input{
				id:       132,
				title:    "A Title",
				author:   "A Author",
				genre:    "A Genre",
				loanDate: "2020-02-03",
				borrower: "Lane",
			},
			expect: repo.BookEntry{
				Variant:  repo.Loaned,
				ID:       132,
				Title:    "A Title",
				Author:   "A Author",
				Genre:    "A Genre",
				Loaned:   time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
			},
		},
		{
			name: "just read book",
			input: input{
				id:       132,
				title:    "A Title",
				author:   "A Author",
				genre:    "A Genre",
				readDate: "2020-02-03",
				rating:   3,
			},
			expect: repo.BookEntry{
				Variant: repo.Read,
				ID:      132,

				Title:  "A Title",
				Author: "A Author",
				Genre:  "A Genre",
				Read:   time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC),
				Rating: 3,
			},
		},
		{
			name: "just book and read loaned",
			input: input{
				id:       132,
				title:    "A Title",
				author:   "A Author",
				genre:    "A Genre",
				readDate: "2020-02-03",
				rating:   3,
				loanDate: "2020-02-03",
				borrower: "Lane",
			},
			expect: repo.BookEntry{
				Variant:  repo.Read | repo.Loaned,
				ID:       132,
				Title:    "A Title",
				Author:   "A Author",
				Genre:    "A Genre",
				Read:     time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC),
				Rating:   3,
				Loaned:   time.Date(2020, time.February, 3, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			builder := newBookBuilder()
			actual, err := builder.SetID(c.input.id).
				SetTitle(c.input.title).
				SetAuthor(c.input.author).
				SetGenre(c.input.genre).
				SetReadDate(c.input.readDate).
				SetRating(c.input.rating).
				SetLoanedDate(c.input.loanDate).
				SetBorrower(c.input.borrower).
				Build()

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if actual == nil {
				t.Fatalf("unexpected nil")
			}

			if *actual != c.expect {
				t.Fatalf("expect\n%v\n  got\n%v", c.expect, *actual)
			}
		})
	}
}

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

	books, err := db.GetAllBooks(0)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(books) != 4 {
		t.Fatalf("expect length 4, got %d", len(books))
	}

	for _, book := range books {
		err := db.DeleteBook(&book)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	actual, err := db.GetAllBooks(0)
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
		input   repo.BookEntry
		willErr bool
	}{
		{
			name: "a simple book",
			input: repo.BookEntry{
				ID:     1,
				Title:  "a_title_update",
				Author: "a_author_update",
				Genre:  "a_genre_update",
			},
			willErr: false,
		},
		{
			name: "a loaned book",
			input: repo.BookEntry{
				ID:       2,
				Variant:  repo.Loaned,
				Title:    "a_title_update",
				Author:   "a_author_update",
				Genre:    "a_genre_update",
				Loaned:   time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC),
				Borrower: "Bob",
			},
			willErr: false,
		},
		{
			name: "a read book",
			input: repo.BookEntry{
				ID:      3,
				Variant: repo.Read,
				Title:   "a_title_update",
				Author:  "a_author_update",
				Genre:   "a_genre_update",
				Read:    time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC),
				Rating:  1,
			},
			willErr: false,
		},
		{
			name: "remove loaned and read from book",
			input: repo.BookEntry{
				ID:      4,
				Variant: repo.Book,
				Title:   "a_title_update",
				Author:  "a_author_update",
				Genre:   "a_genre_update",
			},
			willErr: false,
		},
		{
			name: "add loaned and read back book",
			input: repo.BookEntry{
				ID:       4,
				Variant:  repo.Loaned | repo.Read,
				Title:    "a_title_update",
				Author:   "a_author_update",
				Genre:    "a_genre_update",
				Read:     time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC),
				Rating:   1,
				Loaned:   time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC),
				Borrower: "Bob",
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
		input   repo.BookEntry
		expect  int64
		willErr bool
	}{
		{
			name: "a simple book",
			input: repo.BookEntry{
				Title:  "A_Title",
				Author: "A_Author",
				Genre:  "A_Genre",
			},
			expect:  1,
			willErr: false,
		},
		{
			name: "a loaned book",
			input: repo.BookEntry{
				Variant:  repo.Loaned,
				Title:    "A_Title",
				Author:   "A_Author",
				Genre:    "A_Genre",
				Loaned:   time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
			},
			expect:  2,
			willErr: false,
		},
		{
			name: "a read book",
			input: repo.BookEntry{
				Variant: repo.Read,
				Title:   "A_Title",
				Author:  "A_Author",
				Genre:   "A_Genre",
				Read:    time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC),
				Rating:  5,
			},
			expect:  3,
			willErr: false,
		},
		{
			name: "a loaned and read book",
			input: repo.BookEntry{
				Variant:  repo.Loaned | repo.Read,
				Title:    "A_Title",
				Author:   "A_Author",
				Genre:    "A_Genre",
				Loaned:   time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read:     time.Date(2020, time.February, 2, 0, 0, 0, 0, time.UTC),
				Rating:   5,
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
