package memdb

import (
	"testing"
	"strings"
	"fmt"
	"time"

	"github.com/dubbersthehoser/mayble/internal/storage"
)


func TestMemStorageBasic(t *testing.T) {
	memStore := NewMemStorage()
	bookAmount := 12
	tests := make([]*storage.BookLoan, bookAmount)
	for i:=0; i<bookAmount; i++ {
		bookLoan := storage.BookLoan{
			Book: storage.Book{
				ID: int64(i),
				Title: fmt.Sprintf("title_%d", i),
				Author: fmt.Sprintf("author_%d", i),
				Genre: fmt.Sprintf("genre_%d", i),
				Ratting: i % 6,
			},
			Loan: &storage.Loan{
				ID: int64(i),
				Name: fmt.Sprintf("name_%d", i),
				Date: time.Now().Add(time.Hour * time.Duration(24 * (i+1))),
			},
		}
		tests[i] = &bookLoan
	}
	
	t.Run("CreateBookLoan", func(t *testing.T) {
		for _, bookLoan := range tests {
			err := memStore.CreateBookLoan(bookLoan)
			if err != nil {
				t.Fatal(err)
			}
		}
	})
	
	t.Run("GetBookLoanByID", func(t *testing.T) {
		for _, bookLoan := range tests {
			_, err := memStore.GetBookLoanByID(bookLoan.ID)
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("GetAllBookLoans", func(t *testing.T) {
		testMap := make(map[int64]bool)

		for _, bookLoan := range tests {
			testMap[bookLoan.ID] = false
		}

		bookLoans, err := memStore.GetAllBookLoans()
		if err != nil {
			t.Fatal(err)
		}

		for _, bookLoan := range bookLoans {
			testMap[bookLoan.ID] = true
		}

		for id, ok := range testMap {
			if !ok {
				t.Errorf("entry.id = %d, was not in list", id)
			}
		}
	})

	t.Run("UpdateBookLoan", func(t *testing.T) {
		
		for _, bookLoan := range tests {
			bookLoan.Title = bookLoan.Title + "_update"
			bookLoan.Author = bookLoan.Author + "_update"
			bookLoan.Genre = bookLoan.Genre + "_update"
			bookLoan.Ratting = 0
		}

		for _, bookLoan := range tests {
			err := memStore.UpdateBookLoan(bookLoan)
			if err != nil {
				t.Fatal(err)
			}
		}

		for _, bookLoan := range tests{
			ret, err := memStore.GetBookLoanByID(bookLoan.ID)
			if err != nil {
				t.Fatal(err)
			}
			switch {
			case !strings.HasSuffix(ret.Title, "_update"):
				t.Error("title missing update suffix")
			case !strings.HasSuffix(ret.Author, "_update"):
				t.Error("author missing update suffix")
			case !strings.HasSuffix(ret.Genre, "_update"):
				t.Error("genre missing update suffix")
			case ret.Ratting != 0:
				t.Error("ratting not update value")
			}
		}
	})

	t.Run("DeleteBookLoan", func(t *testing.T) {
		for _, bookLoan := range tests{
			err := memStore.DeleteBookLoan(bookLoan)
			if err != nil {
				t.Fatalf("%s: %#v\n",err, bookLoan)
			}
		}

		bookLoans, err := memStore.GetAllBookLoans()
		if err != nil {
			t.Fatal(err)
		}

		if len(bookLoans) != 0 {
			t.Error("not all entries were deleted")
		}
	})
}

