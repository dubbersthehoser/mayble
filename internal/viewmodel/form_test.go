package viewmodel

import (
	"testing"
	"time"
	"fmt"

	"fyne.io/fyne/v2/app"

	repo "github.com/dubbersthehoser/mayble/internal/repository"

)

func TestBookForm(t *testing.T) {
	
	// need to create fyne app for binding to work.
	_ = app.New() // I'm never going to use this lib/framwork ever again.
	
	form := NewBookForm()

	books := []repo.BookEntry{
		{
			Variant: repo.Book,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
		},
		{
			Variant: repo.Book | repo.Read,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating: 3,
		},
		{
			Variant: repo.Book | repo.Loaned,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
		},
		{
			Variant: repo.Book | repo.Loaned | repo.Read,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
			Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating: 3,
		},
	}

	t.Run("Set", func(t *testing.T) {
		testBookFormSet(t, form, books)
	})

	t.Run("Get", func(t *testing.T) {
		testBookFormToBookEntry(t, form, books)
	})

	t.Run("validate", func(t *testing.T) {
		testBookForm_validate(t, form)
	})
}

func testBookForm_validate(t *testing.T, form *BookForm) {
	tests := []struct{
		name string
		input repo.BookEntry
		willErr bool
	}{
		{
			name: "complete",
			input: repo.BookEntry{
				Variant: repo.Book | repo.Loaned | repo.Read,
				Title: "Title",
				Author: "Author",
				Genre: "Genre",
				Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating: 3,
			},
			willErr: false,
		},
		{
			name: "without title",
			input: repo.BookEntry{
				Variant: repo.Book | repo.Loaned | repo.Read,
				Title: "",
				Author: "Author",
				Genre: "Genre",
				Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating: 3,
			},
			willErr: true,
		},
		{
			name: "without author",
			input: repo.BookEntry{
				Variant: repo.Book | repo.Loaned | repo.Read,
				Title: "Title",
				Author: "",
				Genre: "Genre",
				Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating: 3,
			},
			willErr: true,
		},
		{
			name: "without genre",
			input: repo.BookEntry{
				Variant: repo.Book | repo.Loaned | repo.Read,
				Title: "Title",
				Author: "Author",
				Genre: "",
				Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating: 3,
			},
			willErr: true,
		},
		{
			name: "without loaned date",
			input: repo.BookEntry{
				Variant: repo.Book | repo.Loaned | repo.Read,
				Title: "Title",
				Author: "Author",
				Genre: "Genre",
				Loaned: time.Time{},
				Borrower: "Lane",
				Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating: 3,
			},
			willErr: true,
		},
		{
			name: "without loaned borrower",
			input: repo.BookEntry{
				Variant: repo.Book | repo.Loaned | repo.Read,
				Title: "Title",
				Author: "Author",
				Genre: "Genre",
				Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "",
				Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating: 3,
			},
			willErr: true,
		},
		{
			name: "without read date",
			input: repo.BookEntry{
				Variant: repo.Book | repo.Loaned | repo.Read,
				Title: "Title",
				Author: "Author",
				Genre: "Genre",
				Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read: time.Time{},
				Rating: 3,
			},
			willErr: true,
		},
		{
			name: "without rating",
			input: repo.BookEntry{
				Variant: repo.Book | repo.Loaned | repo.Read,
				Title: "Title",
				Author: "Author",
				Genre: "Genre",
				Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating: 0,
			},
			willErr: true,
		},
		{
			name: "just a book",
			input: repo.BookEntry{
				Variant: repo.Book,
				Title: "Title",
				Author: "Author",
				Genre: "Genre",
				Loaned: time.Time{},
				Borrower: "",
				Read: time.Time{},
				Rating: 0,
			},
			willErr: false,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			form.Set(&c.input)
			err := form.validate()
			if err == nil && c.willErr {
				t.Fatal("expected error")
			}
			if !c.willErr {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return
			}
		})
	}
}

func testBookFormToBookEntry(t *testing.T, form *BookForm, books []repo.BookEntry) {
	for i, book := range books {
		t.Run(fmt.Sprintf("book: %d", i), func(t *testing.T){
			form.Set(&book)
			abook := *form.ToBookEntry()
			if abook != book {
				t.Fatalf("expect\n%v\n  got\n%v", book, abook)
			}
		})
	}
}

func testBookFormSet(t *testing.T, form *BookForm, books []repo.BookEntry) {
	
	for i, book := range books {
		t.Run(fmt.Sprintf("book: %d", i), func(t *testing.T){
			form.Set(&book)

			aTitle, _ := form.Title.Get()
			aAuthor, _ := form.Author.Get()
			aGenre, _ := form.Genre.Get()

			if aTitle != book.Title {
				t.Fatalf("expect '%s', got '%s'", book.Title, aTitle)
			}

			if aAuthor != book.Author {
				t.Fatalf("expect '%s', got '%s'", book.Author, aAuthor)
			}

			if aGenre != book.Genre {
				t.Fatalf("expect '%s', got '%s'", book.Genre, aGenre)
			}

			if book.Variant & repo.Loaned != 0 {
				ok, _ := form.IsLoaned.Get()
				if !ok {
					t.Fatalf("expect %t, got %t", true, ok)
				}
				eDate := formatDate(&book.Loaned)
				eBorrower := book.Borrower

				aDate, _ := form.Date.Get()
				aBorrower, _ := form.Borrower.Get()
				if aDate != eDate {
					t.Fatalf("expect loaned %s, got %s", eDate, aDate)
				}
				if aBorrower != eBorrower {
					t.Fatalf("expect %s, got %s", eBorrower, aBorrower)
				}
			}

			if book.Variant & repo.Read != 0 {
				ok, _ := form.IsRead.Get()
				if !ok {
					t.Fatalf("expect %t, got %t", true, ok)
				}
				eDate := formatDate(&book.Read)
				eRating := formatRating(book.Rating)

				aDate, _ := form.Completed.Get()
				aRating, _ := form.Rating.Get()
				if aDate != eDate {
					t.Fatalf("expect read %s, got %s", eDate, aDate)
				}
				if aRating != eRating {
					t.Fatalf("expect %s, got %s", eRating, aRating)
				}
			}
		})
	}
}


