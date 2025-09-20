package memdb

import (
	"testing"
	"strings"

	"github.com/dubbersthehoser/mayble/internal/storage"
)


func TestMemStorageBasic(t *testing.T) {

	memStore := NewMemStorage()

	tests := []*storage.BookLoan{
		&storage.BookLoan{
			Book: storage.Book{
				Title: "title_1",
				Author: "author_1",
				Genre: "genre_1",
				Ratting: 0,
			},
		},
		&storage.BookLoan{
			Book: storage.Book{
				Title: "title_2",
				Author: "author_2",
				Genre: "genre_2",
				Ratting: 1,
			},
		},
		&storage.BookLoan{
			Book: storage.Book{
				Title: "title_3",
				Author: "author_3",
				Genre: "genre_3",
				Ratting: 2,
			},
		},
		&storage.BookLoan{
			Book: storage.Book{
				Title: "title_4",
				Author: "author_4",
				Genre: "genre_4",
				Ratting: 3,
			},
		},
		&storage.BookLoan{
			Book: storage.Book{
				Title: "title_5",
				Author: "author_5",
				Genre: "genre_5",
				Ratting: 4,
			},
		},
		&storage.BookLoan{
			Book: storage.Book{
				Title: "title_6",
				Author: "author_6",
				Genre: "genre_6",
				Ratting: 5,
			},
		},
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
				t.Fatal(err)
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

