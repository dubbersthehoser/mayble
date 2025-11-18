package memory

import (
	"errors"
	"testing"
	"time"

	"github.com/dubbersthehoser/mayble/internal/storage"
)


/**********************************
        Testing Store Book 
***********************************/


func TestStoreBookCreate(t *testing.T) {
	store := NewStorage()
	tests := []storage.Book{
		storage.Book{
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			Ratting: 1,
		},
		storage.Book{
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			Ratting: 2,
		},
	}
	for i, test := range tests {
		id, err := store.CreateBook(test.Title, test.Author, test.Genre, test.Ratting)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		if int(id) != i+1 {
			t.Fatalf("case %d, got id %d, want %d", i, id, i+1)
		}
	}
}

func TestStoreBookDelete(t *testing.T) {
	store := NewStorage()
	tests := []storage.Book{
		storage.Book{
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			Ratting: 1,
		},
		storage.Book{
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			Ratting: 2,
		},
	}
	for i, test := range tests {
		id, err := store.CreateBook(test.Title, test.Author, test.Genre, test.Ratting)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		tests[i].ID = id
	}

	for i, test := range tests {
		err := store.DeleteBook(&test)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		err = store.DeleteBook(&test)
		if !errors.Is(err, storage.ErrEntryNotFound) {
			t.Fatalf("case %d, expect error: %s", i, err)
		}
	}
}

func TestStoreUpdate(t *testing.T) {
	store := NewStorage()
	tests := []struct{
		created storage.Book
		updated  storage.Book
	}{
		{
			created: storage.Book{
				Title: "title_1",
				Author: "author_1",
				Genre: "genre_1",
				Ratting: 1,
			},
			updated: storage.Book{
				Title: "new_title_1",
				Author: "new_author_1",
				Genre: "new_genre_1",
				Ratting: 2,
			},
		},
		{
			created: storage.Book{
				Title: "title_2",
				Author: "author_2",
				Genre: "genre_2",
				Ratting: 2,
			},
			updated: storage.Book{
				Title: "new_title_2",
				Author: "new_author_2",
				Genre: "new_genre_2",
				Ratting: 4,
			},
		},
	}

	for i, test := range tests {
		create := test.created
		id, err := store.CreateBook(create.Title, create.Author, create.Genre, create.Ratting)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		tests[i].created.ID = id
		tests[i].updated.ID = id
	}

	for i, test := range tests {
		err := store.UpdateBook(&test.updated)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}

		book, ok := store.Books[test.updated.ID]
		if !ok {
			t.Fatalf("case %d, unexpected could not find book: id=%d", i, test.updated.ID)
		}
		if book.Title != test.updated.Title {
			t.Fatalf("case %d, expect title %s, got %s", i, test.updated.Title, book.Title)
		} 
		if book.Author != test.updated.Author {
			t.Fatalf("case %d, expect author %s, got %s", i, test.updated.Author, book.Author)
		} 
		if book.Genre != test.updated.Genre {
			t.Fatalf("case %d, expect genre %s, got %s", i, test.updated.Genre, book.Genre)
		} 
		if book.Ratting != test.updated.Ratting {
			t.Fatalf("case %d, expect ratting %d, got %d", i, test.updated.Ratting, book.Ratting)
		} 
	}


}

func TestGetBooks(t *testing.T) {
	store := NewStorage()
	tests := []storage.Book{
		storage.Book{
			Title: "title_1",
			Author: "author_1",
			Genre: "genre_1",
			Ratting: 1,
		},
		storage.Book{
			Title: "title_2",
			Author: "author_2",
			Genre: "genre_2",
			Ratting: 2,
		},
	}
	for i, test := range tests {
		id, err := store.CreateBook(test.Title, test.Author, test.Genre, test.Ratting)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
		tests[i].ID = id
	}

	actuals, err := store.GetBooks()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(actuals) != len(tests) {
		t.Fatalf("expect length %d, got %d", len(tests), len(actuals))
	}
}

func TestStoreBookErrors(t *testing.T) {
	store := NewStorage()

	book := storage.Book{
		ID: 10,
		Title: "title",
		Author: "author",
		Genre: "genre",
		Ratting: 3,
	}

	// UpdateBook()
	err := store.UpdateBook(nil)
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatal("update book gave wrong error value")
	}
	err = store.UpdateBook(&book)
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatal("update book gave wrong error value")
	}
	
	// DeleteBook()
	err = store.DeleteBook(nil)
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatal("delete book gave wrong error value")
	}
	err = store.DeleteBook(&book)
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatal("delete book gave wrong error value")
	}

}

/**********************************
        Testing Store Loan 
***********************************/

func TestStoreLoanCreate(t *testing.T) {
	store := NewStorage()
	date := time.Now()
	tests := []storage.Loan{
		storage.Loan{
			ID: 1,
			Borrower: "borrower_1",
			Date:	date,
		},
		storage.Loan{
			ID: 2,
			Borrower: "borrower_2",
			Date:	date.Add(time.Hour * 24),
		},
	}

	for i, test := range tests{
		err := store.CreateLoan(test.ID, test.Borrower, test.Date)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		loan, ok := store.Loans[test.ID]
		if !ok {
			t.Fatalf("case %d, unexpected could not find create loan", i)
		}
		if loan.Borrower != test.Borrower {
			t.Fatalf("case %d, expect %s, got %s", i, test.Borrower, loan.Borrower)
		}
		if !loan.Date.Equal(test.Date) {
			t.Fatalf("case %d, expect '%v', got '%v'", i, test.Date, loan.Date)
		}
	}
}

func TestStoreLoanUpdate(t *testing.T) {
	store := NewStorage()
	date := time.Now()
	tests := []struct{
		created storage.Loan
		updated storage.Loan
	}{
		{
			created: storage.Loan{
				ID: 0,
				Borrower: "borrower_1",
				Date:	date,
			},
			updated: storage.Loan{
				ID: 0,
				Borrower: "new_borrower_1",
				Date: date.Add(time.Hour * (24)),
			},
		},
		{
			created: storage.Loan{
				ID: 1,
				Borrower: "borrower_2",
				Date:	date.Add(time.Hour * 24),
			},
			updated: storage.Loan{
				ID: 1,
				Borrower: "new_borrower_2",
				Date: date.Add(time.Hour * 48),
			},
		},
	}

	for i, test := range tests {
		err := store.CreateLoan(test.created.ID, test.created.Borrower, test.created.Date)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s",i, err)
		}
	}

	for i, test := range tests {
		err := store.UpdateLoan(&test.updated)
		if err != nil {
			t.Fatalf("cate %d, unexpected error: %s", i, err)
		}
		actual, ok := store.Loans[test.created.ID]
		if !ok {
			t.Fatalf("case %d, unexpected could not find updated loan", i)
		}

		if actual.Borrower != test.updated.Borrower {
			t.Fatalf("case %d, expect %s, got %s", i, test.updated.Borrower, actual.Borrower)
		}

		if !test.updated.Date.Equal(actual.Date) {
			t.Fatalf("case %d, expect '%v', got '%v'", i, test.updated.Date, actual.Date)
		}
	}
}

func TestStoreLoanDelete(t *testing.T) {
	store := NewStorage()
	date := time.Now()
	tests := []storage.Loan{
		storage.Loan{
			ID: 1,
			Borrower: "borrower_1",
			Date:	date,
		},
		storage.Loan{
			ID: 2,
			Borrower: "borrower_2",
			Date:	date.Add(time.Hour * 24),
		},
	}

	for i, test := range tests{
		err := store.CreateLoan(test.ID, test.Borrower, test.Date)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	for i, test := range tests{
		err := store.DeleteLoan(&test)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}
	if len(store.Loans) != 0 {
		t.Fatalf("expected length of 0, got %d", len(store.Loans))
	}
}

func TestStoreGetLoan(t *testing.T) {
	store := NewStorage()
	date := time.Now()
	tests := []storage.Loan{
		storage.Loan{
			ID: 1,
			Borrower: "borrower_1",
			Date:	date,
		},
		storage.Loan{
			ID: 2,
			Borrower: "borrower_2",
			Date:	date.Add(time.Hour * 24),
		},
	}
	for i, test := range tests{
		err := store.CreateLoan(test.ID, test.Borrower, test.Date)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
	}

	for i, test := range tests{
		loan, err := store.GetLoan(test.ID)
		if err != nil {
			t.Fatalf("case %d, unexpeced error: %s", i, err)
		}

		if loan.Borrower != test.Borrower {
			t.Fatalf("case %d, expect %s, got %s", i, test.Borrower, loan.Borrower)
		}
		if !loan.Date.Equal(test.Date) {
			t.Fatalf("case %d, expect '%v', got '%v'", i, test.Date, loan.Date)
		}

	}
}

func TestStoreLoanErrors(t *testing.T) {
	store := NewStorage()

	// CreateLoan()
	loan := storage.Loan{
		ID: 10,
		Borrower: "jake",
		Date: time.Now(),
	}
	err := store.CreateLoan(loan.ID, loan.Borrower, loan.Date)
	if err != nil {
		t.Fatal("unexpected error")
	}
	err = store.CreateLoan(loan.ID, loan.Borrower, loan.Date)
	if !errors.Is(err, storage.ErrEntryExists) {
		t.Fatal("CreateLoan gave wrong error value")
	}

	// UpdateLoan()
	err = store.UpdateLoan(nil)
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatal("UpdateLoan gave wrong error value")
		
	}

	loan.ID = 20
	err = store.UpdateLoan(&loan)
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatal("UpdateLoan gave wrong error value")
	}

	// DeleteLoan()
	err = store.DeleteLoan(nil)
	if !errors.Is(err, storage.ErrInvalidValue) {
		t.Fatal("DeleteLoan gave wrong error value")
	}

	err = store.DeleteLoan(&loan)
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatal("DeleteLoan gave wrong error value")
	}


	// GetLoan()
	_, err = store.GetLoan(loan.ID)
	if !errors.Is(err, storage.ErrEntryNotFound) {
		t.Fatal("GetLoan gave wrong error value")
	}
}

func TestGetNewBookID(t *testing.T) {
	store := NewStorage()

	for i := range 10 {
		expect := int64(i) + 1
		id := store.getNewBookID()
		if expect != id {
			t.Fatalf("%d'th call got %d, want %d", i, id, expect)
		}
		store.CreateBook("title", "author", "genre", 3)
	}

}






