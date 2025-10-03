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

	t.Log("Creating")

	bookLoan := storage.NewBookLoan()
	bookLoan.Title = "title_1"
	bookLoan.Author = "author_1"
	bookLoan.Genre = "genre_1"
	bookLoan.Ratting = 4
	bookLoan.Loan.Name = "john"
	bookLoan.Loan.Date = time.Now()

	err = myCore.CreateBookLoan(bookLoan)
	if err != nil {
		t.Fatal(err)
	}

	list, err := myCore.ListBookLoans(ByTitle, ASC)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("want length '%d', got '%d'", 1, len(list))
	}

	gotBookLoan := list[0]

	if err := compareBookLoan(*bookLoan, gotBookLoan); err != nil {
		t.Fatal(err)
	}



	t.Log("Updating")

	gotBookLoan.Title = "title_2"
	gotBookLoan.Author = "author_2"
	gotBookLoan.Genre = "genre_2"
	gotBookLoan.Ratting = 2
	gotBookLoan.Loan.Name = "Jack"
	gotBookLoan.Loan.Date.Add(time.Hour * 24)

	err = myCore.UpdateBookLoan(&gotBookLoan)
	if err != nil {
		t.Fatal(err)
	}

	list, err = myCore.ListBookLoans(ByTitle, ASC)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("want length '%d', got '%d'", 1, len(list))
	}
	
	updatedBookLoan := list[0]

	if err := compareBookLoan(gotBookLoan, updatedBookLoan); err != nil {
		t.Fatal(err)
	}

	t.Log("Deleting")

	err = myCore.DeleteBookLoan(&updatedBookLoan)
	if err != nil {
		t.Fatal(err)
	}

	list, err = myCore.ListBookLoans(ByTitle, ASC)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("want length '%d', got '%d'", 0, len(list))
	}

	t.Log("Done")

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
