package core

import (
	"testing"
	"time"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/memdb"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

func newTestCore() (*Core, error) {
	myStore := memdb.NewMemStorage()
	return New(myStore)
}

func TestCore(t *testing.T) {
	myCore, err := newTestCore()
	if err != nil {
		t.Fatal(err)
	}

	maxAmount := 10

	t.Log("-- Creating --")

	for i:=0; i<maxAmount; i++ {
		bookLoan := storage.NewBookLoan()
		bookLoan.Title = fmt.Sprintf("title_%d", i)
		bookLoan.Author = fmt.Sprintf("author_%d", i)
		bookLoan.Genre = fmt.Sprintf("genre_%d", i)
		bookLoan.Ratting = (i % 6)
		bookLoan.Loan.Name = fmt.Sprintf("borrower_%d", i)
		bookLoan.Loan.Date = time.Now()

		err = myCore.CreateBookLoan(bookLoan)
		if err != nil {
			t.Fatal(err)
		}
	}

	list, err := myCore.ListBookLoans(ByID, ASC)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != maxAmount {
		t.Fatalf("want length '%d', got '%d'", maxAmount, len(list))
	}

	t.Log("-- Updating --")

	for _, bookLoan := range list {
		i := bookLoan.ID
		bookLoan.Title = fmt.Sprintf("title_%d_update", i)
		bookLoan.Author = fmt.Sprintf("author_%d_update", i)
		bookLoan.Genre = fmt.Sprintf("genre_%d_update", i)
		bookLoan.Ratting = 0
		bookLoan.Loan.Name = fmt.Sprintf("borrower_%d_update", i)
		bookLoan.Loan.Date.Add(time.Hour * 24)

		err = myCore.UpdateBookLoan(&bookLoan)
		if err != nil {
			t.Logf("bookLloan: %d\n", bookLoan.ID)
			t.Fatal(err)
		}
	}

	list, err = myCore.ListBookLoans(ByID, ASC)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != maxAmount {
		t.Fatalf("want length '%d', got '%d'", maxAmount, len(list))
	}
	
	t.Log("-- Deleting --")

	for _, bookLoan := range list {
		err = myCore.DeleteBookLoan(&bookLoan)
		if err != nil {
			t.Fatal(err)
		}
	}

	list, err = myCore.ListBookLoans(ByID, ASC)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("want length '%d', got '%d'", 0, len(list))
	}
	t.Log("Done")

}

func TestSave(t *testing.T) {
	myStore := memdb.NewMemStorage()
	core, err := New(myStore)
	if err != nil {
		t.Fatal(err)
	}
	bookAmount := 0
	for i:= 0; i<bookAmount; i++ {
		bookLoan := storage.BookLoan{
			Book: storage.Book{
				ID: storage.ZeroID,
				Title: fmt.Sprintf("title_%d", i),
				Author: fmt.Sprintf("author_%d", i),
				Genre: fmt.Sprintf("genre_%d", i),
				Ratting: i % 6,
			},
			Loan: &storage.Loan{
				ID: storage.ZeroID,
				Name: fmt.Sprintf("name_%d", i),
				Date: time.Now().Add(time.Hour * time.Duration(24 * (i+1))),
			},
		}
		err := core.CreateBookLoan(&bookLoan)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = core.Save()
	if err != nil{
		t.Fatal(err)
	}

	bookLoanList, err := myStore.GetAllBookLoans()
	if err != nil {
		t.Fatal(err)
	}
	if len(bookLoanList) != bookAmount {
		t.Errorf("bookLoanList length should be '%d', got '%d'", len(bookLoanList), bookAmount)
	}
}


func TestListBookLoan(t *testing.T) {
	core, err := newTestCore()
	if err != nil {
		t.Fatal(err)
	}
	bookAmount := 10
	for i:=0; i<bookAmount; i++ {
		bookLoan := storage.BookLoan{
			Book: storage.Book{
				ID: storage.ZeroID,
				Title: fmt.Sprintf("title_%d", i),
				Author: fmt.Sprintf("author_%d", i),
				Genre: fmt.Sprintf("genre_%d", i),
				Ratting: i % 6,
			},
			Loan: &storage.Loan{
				ID: storage.ZeroID,
				Name: fmt.Sprintf("name_%d", i),
				Date: time.Now().Add(time.Hour * time.Duration(24 * (i+1))),
			},
		}
		err := core.CreateBookLoan(&bookLoan)
		if err != nil {
			t.Fatal(err)
		}
	}

	bookList, err := core.ListBookLoans(ByTitle, ASC)
	if err != nil {
		t.Fatal(err)
	}
	
	if len(bookList) != bookAmount {
		t.Fatalf("list size got '%d', got '%d'", len(bookList), bookAmount)
	}
}

func compareBookLoan(a, b storage.BookLoan) error {
	if b.Title != a.Title {
		return fmt.Errorf("want title '%s', got '%s'", a.Title, b.Title)
	}
	if b.Author != a.Author {
		return fmt.Errorf("want title '%s', got '%s'", a.Author, b.Author)
	}
	if b.Genre != a.Genre {
		return fmt.Errorf("want title '%s', got '%s'", a.Genre, b.Genre)
	}
	if b.Ratting != a.Ratting {
		return fmt.Errorf("want title '%d', got '%d'", a.Ratting, b.Ratting)
	}
	if b.Loan.Name != a.Loan.Name {
		return fmt.Errorf("want title '%s', got '%s'", a.Loan.Name, b.Loan.Name)
	}
	if !a.Loan.Date.Equal(b.Loan.Date) {
		return fmt.Errorf("want title '%s', got '%s'", a.Loan.Date, b.Loan.Date)
	}
	return nil
}
