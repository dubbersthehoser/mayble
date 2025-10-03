package controller

import (
	"testing"
	"time"
	"fmt"

	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/storage"
	"github.com/dubbersthehoser/mayble/internal/memdb"
)

func TestBookList(t *testing.T) {
	bookAmount := 5
	store := memdb.NewMemStorage()
	for i:=0; i<bookAmount; i++ {
		bookLoan := storage.BookLoan{
			Book: storage.Book{
				Title: fmt.Sprintf("title_%d", i),
				Author: fmt.Sprintf("author_%d", i),
				Genre: fmt.Sprintf("genre_%d", i),
				Ratting: i % 6,
			},
			Loan: &storage.Loan{
				Name: fmt.Sprintf("name_%d", i),
				Date: time.Now().Add(time.Hour * time.Duration(24 * (i+1))),
			},
		}
		err := store.CreateBookLoan(&bookLoan)
		if err != nil {
			t.Fatal(err)
		}
	}
	core, err := core.New(store)
	if err != nil {
		t.Fatal(err)
	}

	bookList := NewBookList(core)
	err = bookList.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(bookList.list) != bookAmount {
		t.Fatalf("list size got '%d', got '%d'", len(bookList.list), bookAmount)
	}

	for _, bookLoan := range bookList.list {
		fmt.Printf("%v\n", bookLoan)
	}
}
