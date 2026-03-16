package viewmodel

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

func TestBookForm(t *testing.T) {

	// need to create fyne app for binding to work.
	_ = app.New() // I'm never going to use this lib/framwork ever again.

	form := NewBookForm()

	books := []repo.BookEntry{
		{
			Variant: repo.Book,
			Title:   "Title",
			Author:  "Author",
			Genre:   "Genre",
		},
		{
			Variant: repo.Book | repo.Read,
			Title:   "Title",
			Author:  "Author",
			Genre:   "Genre",
			Read:    time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating:  3,
		},
		{
			Variant:  repo.Book | repo.Loaned,
			Title:    "Title",
			Author:   "Author",
			Genre:    "Genre",
			Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
		},
		{
			Variant:  repo.Book | repo.Loaned | repo.Read,
			Title:    "Title",
			Author:   "Author",
			Genre:    "Genre",
			Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
			Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating:   3,
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
	tests := []struct {
		name    string
		input   repo.BookEntry
		willErr bool
	}{
		{
			name: "complete",
			input: repo.BookEntry{
				Variant:  repo.Book | repo.Loaned | repo.Read,
				Title:    "Title",
				Author:   "Author",
				Genre:    "Genre",
				Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:   3,
			},
			willErr: false,
		},
		{
			name: "without title",
			input: repo.BookEntry{
				Variant:  repo.Book | repo.Loaned | repo.Read,
				Title:    "",
				Author:   "Author",
				Genre:    "Genre",
				Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:   3,
			},
			willErr: true,
		},
		{
			name: "without author",
			input: repo.BookEntry{
				Variant:  repo.Book | repo.Loaned | repo.Read,
				Title:    "Title",
				Author:   "",
				Genre:    "Genre",
				Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:   3,
			},
			willErr: true,
		},
		{
			name: "without genre",
			input: repo.BookEntry{
				Variant:  repo.Book | repo.Loaned | repo.Read,
				Title:    "Title",
				Author:   "Author",
				Genre:    "",
				Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:   3,
			},
			willErr: true,
		},
		{
			name: "without loaned date",
			input: repo.BookEntry{
				Variant:  repo.Book | repo.Loaned | repo.Read,
				Title:    "Title",
				Author:   "Author",
				Genre:    "Genre",
				Loaned:   time.Time{},
				Borrower: "Lane",
				Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:   3,
			},
			willErr: true,
		},
		{
			name: "without loaned borrower",
			input: repo.BookEntry{
				Variant:  repo.Book | repo.Loaned | repo.Read,
				Title:    "Title",
				Author:   "Author",
				Genre:    "Genre",
				Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "",
				Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:   3,
			},
			willErr: true,
		},
		{
			name: "without read date",
			input: repo.BookEntry{
				Variant:  repo.Book | repo.Loaned | repo.Read,
				Title:    "Title",
				Author:   "Author",
				Genre:    "Genre",
				Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read:     time.Time{},
				Rating:   3,
			},
			willErr: true,
		},
		{
			name: "without rating",
			input: repo.BookEntry{
				Variant:  repo.Book | repo.Loaned | repo.Read,
				Title:    "Title",
				Author:   "Author",
				Genre:    "Genre",
				Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Borrower: "Lane",
				Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
				Rating:   0,
			},
			willErr: true,
		},
		{
			name: "just a book",
			input: repo.BookEntry{
				Variant:  repo.Book,
				Title:    "Title",
				Author:   "Author",
				Genre:    "Genre",
				Loaned:   time.Time{},
				Borrower: "",
				Read:     time.Time{},
				Rating:   0,
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
		t.Run(fmt.Sprintf("book: %d", i), func(t *testing.T) {
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
		t.Run(fmt.Sprintf("book: %d", i), func(t *testing.T) {
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

			if book.Variant&repo.Loaned != 0 {
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

			if book.Variant&repo.Read != 0 {
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

func TestSubmissionList(t *testing.T) {
	b := &bus.Bus{}
	list := NewSubmissionList(b, NewBookForm())

	errCount := 0

	b.Subscribe(bus.Handler{
		Name: msgUserError,
		Handler: func(e *bus.Event) {
			errCount += 1
		},
	})

	books := []repo.BookEntry{
		{
			Variant: repo.Book,
			Title:   "Title",
			Author:  "Author",
			Genre:   "Genre",
		},
		{
			Variant: repo.Book | repo.Read,
			Title:   "Title",
			Author:  "Author",
			Genre:   "Genre",
			Read:    time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating:  3,
		},
		{
			Variant:  repo.Book | repo.Loaned,
			Title:    "Title",
			Author:   "Author",
			Genre:    "Genre",
			Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
		},
		{
			Variant:  repo.Book | repo.Loaned | repo.Read,
			Title:    "Title",
			Author:   "Author",
			Genre:    "Genre",
			Loaned:   time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
			Read:     time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating:   3,
		},
	}

	t.Run("add", func(t *testing.T) {
		testSubmissionList_add(t, list, books)
		if errCount != 0 {
			t.Fatalf("unexpected error count")
		}
	})

	t.Run("GetView", func(t *testing.T) {
		testSubmissionListGetView(t, list, books)
		if errCount != 0 {
			t.Fatalf("unexpected error count")
		}
	})

	t.Run("pop", func(t *testing.T) {
		testSubmissionList_pop(t, list, books)
		if errCount != 0 {
			t.Fatalf("unexpected error count")
		}
	})

	t.Run("Edit", func(t *testing.T) {
		for _, book := range books {
			list.append(&book)
		}
		testSubmissionListEdit(t, list, books)
	})

	t.Run("clear", func(t *testing.T) {
		list.Clear()
		if list.Length() != 0 {
			t.Fatalf("expect length 0, got %d", list.Length())
		}
	})
}

func testSubmissionListEdit(t *testing.T, list *SubmissionList, books []repo.BookEntry) {
	lenExpect := list.Length() - 1
	list.Edit(0)
	form := list.form
	if lenExpect != list.Length() {
		t.Fatalf("expect length %d, got %d", lenExpect, list.Length())
	}
	book := form.ToBookEntry()
	if *book != books[0] {
		t.Fatalf("expect\n%v\n  got\n %v", books[0], *book)
	}
}

func testSubmissionList_pop(t *testing.T, list *SubmissionList, books []repo.BookEntry) {
	rbooks := slices.Clone(books)
	slices.Reverse(rbooks)
	for i, book := range rbooks {
		t.Run(fmt.Sprintf("book#%d", i), func(t *testing.T) {
			aBook, err := list.pop()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if aBook == nil {
				t.Fatal("unexpected nil")
			}
			if book != *aBook {
				t.Fatalf("expect\n%v\n  got\n%v", book, aBook)
			}
		})
	}
}

func testSubmissionListGetView(t *testing.T, list *SubmissionList, books []repo.BookEntry) {
	for i, book := range books {
		t.Run(fmt.Sprintf("book#%d", i), func(t *testing.T) {
			actual := list.GetView(i)
			if (book.Variant&repo.Read != 0) && !strings.Contains(actual, "(read)") {
				t.Fatalf("expected to contain string '(read)', in '%s'", actual)
			}
			if (book.Variant&repo.Loaned != 0) && !strings.Contains(actual, "(loaned)") {
				t.Fatalf("expected to contain string '(loaned)', in '%s'", actual)
			}
		})
	}
	actual := list.GetView(list.Length())
	if !strings.Contains(actual, "out of range") {
		t.Fatalf("expected to contain string 'out of range', in '%s'", actual)
	}
	actual = list.GetView(-1)
	if !strings.Contains(actual, "out of range") {
		t.Fatalf("expected to contain string 'out of range', in '%s'", actual)
	}
}

func testSubmissionList_add(t *testing.T, list *SubmissionList, books []repo.BookEntry) {
	form := NewBookForm()

	for i, book := range books {
		t.Run(fmt.Sprintf("book#%d", i), func(t *testing.T) {
			form.Set(&book)
			list.add(form)
		})
	}

	if len(books) != list.Length() {
		t.Fatalf("expect length %d, got %d", len(books), len(list.submissions))
	}
}

func TestCreateBookForm(t *testing.T) {
	b := &bus.Bus{}
	db, err := database.OpenMem()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer db.Conn.Close()
	cfg := &config.Config{}
	as := newAppService(b, cfg, db)
	err = db.Conn.Ping()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	cForm := NewCreateBookForm(b, as)
	_ = cForm

	t.Run("AddSubmission", func(t *testing.T) {
		testCreateBookFormAddsubmission(t, cForm)
	})
	cForm.BookForm.reset()
	t.Run("Submit", func(t *testing.T) {
		testCreateBookFormSubmit(t, cForm)
	})
}

func testCreateBookFormSubmit(t *testing.T, form *CreateBookForm) {
	var ok bool
	id := form.bus.Subscribe(busMsgTestHelper(t, msgUserInfo, func(s string) {
		ok = true
		expect := "No submissions to submit"
		if s != expect {
			t.Fatalf("expect message '%s', got '%s'", expect, s)
		}
	}))

	form.Submit()

	if !ok {
		t.Fatal("message was not signaled")
	}
	form.bus.Unsubscribe(id)

	repo := form.repo
	form.repo = nil

	ok = false
	id = form.bus.Subscribe(busMsgTestHelper(t, msgUserInfo, func(s string) {
		ok = true
		expect := "Not implemented."
		if s != expect {
			t.Fatalf("expect message '%s', got '%s'", expect, s)
		}
	}))

	form.Submit()
	if !ok {
		t.Fatal("message was not signaled")
	}
	form.bus.Unsubscribe(id)
	form.repo = repo
	_ = form.BookForm.Title.Set("title")
	_ = form.BookForm.Author.Set("author")
	_ = form.BookForm.Genre.Set("genre")

	form.AddSubmission()

	ok = false
	id = form.bus.Subscribe(busMsgTestHelper(t, msgUserSuccess, func(s string) {
		ok = true
		expect := "Books Added!"
		if s != expect {
			t.Fatalf("expect message '%s', got '%s'", expect, s)
		}
	}))
	form.Submit()
	if !ok {
		t.Fatal("message was not signaled")
	}

	form.bus.Unsubscribe(id)
}

func testCreateBookFormAddsubmission(t *testing.T, form *CreateBookForm) {

	var ok bool
	id := form.bus.Subscribe(busMsgTestHelper(t, msgUserError, func(s string) {
		ok = true
	}))
	form.AddSubmission()
	if !ok {
		t.Fatal("validation error was not signaled")
	}
	form.bus.Unsubscribe(id)

	_ = form.BookForm.Title.Set("title")
	_ = form.BookForm.Author.Set("author")
	_ = form.BookForm.Genre.Set("genre")

	ok = false
	sid := form.bus.Subscribe(busMsgTestHelper(t, msgUserInfo, func(s string) {
		ok = true
		expect := "Added submission"
		if s != expect {
			t.Fatalf("expect message '%s', got '%s'", expect, s)
		}
	}))
	fid := form.bus.Subscribe(busMsgTestHelper(t, msgUserError, func(s string) {
		t.Fatalf("unexpected validation error message: %s", s)
	}))
	form.AddSubmission()

	if !ok {
		t.Fatalf("validation success was not signaled")
	}
	form.bus.Unsubscribe(sid)
	form.bus.Unsubscribe(fid)
	form.sl.remove(0)
}
