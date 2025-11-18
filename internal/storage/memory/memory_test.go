package memory

import (
	"errors"
	"testing"
	//"strings"
	//"fmt"
	"time"

	"github.com/dubbersthehoser/mayble/internal/storage"
)


func TestStoreBook(t *testing.T) {
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
		_, err := store.CreateBook(test.Title, test.Author, test.Genre, test.Ratting)
		if err != nil {
			t.Fatalf("case %d, unexpected error: %s", i, err)
		}
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






