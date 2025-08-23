package sqlitedb

import (
	"os"
	"time"
	"testing"
	"path/filepath"


	"github.com/dubbersthehoser/mayble/internal/storage"
)


func TestDatabase(t *testing.T) {
	
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current working directory: %s", err)
	}

	projectRoot := filepath.Join(cwd, "..", "..")
	
	schemaPath := filepath.Join(projectRoot, "sql/schemas/")
	db := NewDatabase(nil, schemaPath)
	t.Log("creating temp file...")
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "*.sqlite")
	if err != nil {
		t.Fatal(err)
	}
	DBFile := tempFile.Name()
	t.Logf("created: %s", DBFile)
	tempFile.Close()
	defer func() {
		os.Remove(DBFile)
	}()

	if err := db.Open(DBFile); err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	t.Log("migrating up database...")
	err = db.MigrateUp()
	if err != nil {
		t.Fatal(err)
	}

	books := []*storage.Book{
		&storage.Book{
			Title: "EXAPLE_title_1",
			Author: "EXAPLE_author_1",
			Genre: "EXAPLE_genre_1",
			Ratting: 0,
		},
		&storage.Book{
			Title: "EXAPLE_title_2",
			Author: "EXAPLE_author_2",
			Genre: "EXAPLE_genre_2",
			Ratting: 1,
		},
		&storage.Book{
			Title: "EXAPLE_title_3",
			Author: "EXAPLE_author_3",
			Genre: "EXAPLE_genre_3",
			Ratting: 2,
		},
		&storage.Book{
			Title: "EXAPLE_title_4",
			Author: "EXAPLE_author_4",
			Genre: "EXAPLE_genre_4",
			Ratting: 3,
		},
		&storage.Book{
			Title: "EXAPLE_title_5",
			Author: "EXAPLE_author_5",
			Genre: "EXAPLE_genre_5",
			Ratting: 4,
		},
		&storage.Book{
			Title: "EXAPLE_title_6",
			Author: "EXAPLE_author_6",
			Genre: "EXAPLE_genre_6",
			Ratting: 5,
		},
	}

	invalidBooks := []*storage.Book{
		&storage.Book{
			Title: "this dog",
			Author: "should",
			Genre: "error",
			Ratting: 6,
		},
		&storage.Book{
			Title: "this dog",
			Author: "should",
			Genre: "error",
			Ratting: -1,
		},
		&storage.Book{
			Title: "this dog",
			Author: "should",
			Genre: "error",
			Ratting: 10,
		},
		&storage.Book{
			Title: "this dog",
			Author: "should",
			Genre: "error",
			Ratting: 1000,
		},
		&storage.Book{
			Title: "this dog",
			Author: "should",
			Genre: "error",
			Ratting: -137,
		},
		&storage.Book{
			Title: "this dog",
			Author: "should",
			Genre: "error",
			Ratting: 69,
		},
	}

	t.Run("CreateValidBooks", func(t *testing.T) {
		for i, book := range books {
			abook, err := db.createBookWithRet(book)
			if err != nil {
				t.Error(err)
				continue
			}
			if abook.Title != book.Title {
				t.Errorf("books[%d] got title='%s', want title='%s'", i, abook.Title, book.Title)
			}
			if abook.Author != book.Author {
				t.Errorf("books[%d] got author='%s', want author='%s'", i, abook.Author, book.Author)
			}
			if abook.Genre != book.Genre {
				t.Errorf("books[%d] got genre='%s', want genre='%s'", i, abook.Genre, book.Genre)
			}
			if abook.Ratting != int64(book.Ratting) {
				t.Errorf("books[%d] got ratting=%d, want ratting=%d", i, abook.Ratting, book.Ratting)
			}
			book.ID = abook.ID
		}
	})

	t.Run("CreateInvalidBooks", func(t *testing.T) {
		for _, invalid := range invalidBooks {
			_, err = db.createBookWithRet(invalid)
			if err == nil {
				t.Errorf("gave invalid ratting, did not error: ratting=%d", invalid.Ratting)
			}
		}
	})

	t.Run("UpdateValidBooks", func(t *testing.T) {
		TitleTo := "update_title_example_0"
		AuthorTo := "update_author_example_1"
		GenreTo := "update_genre_example_2"
		RattingTo := 0

		for i, book := range books {
			book.Title = TitleTo
			book.Author = AuthorTo
			book.Genre = GenreTo
			book.Ratting = RattingTo
			abook, err := db.updateBookWithRet(book)
			if err != nil {
				t.Error(err)
				continue
			}
			if abook.Title != book.Title {
				t.Errorf("books[%d] got title='%s', want title='%s'", i, abook.Title, book.Title)
			}
			if abook.Author != book.Author {
				t.Errorf("books[%d] got author='%s', want author='%s'", i, abook.Author, book.Author)
			}
			if abook.Genre != book.Genre {
				t.Errorf("books[%d] got genre='%s', want genre='%s'", i, abook.Genre, book.Genre)
			}
			if abook.Ratting != int64(book.Ratting) {
				t.Errorf("books[%d] got ratting=%d, want ratting=%d", i, abook.Ratting, book.Ratting)
			}
		}
	})

	t.Run("UpdateInvalidBooks", func(t *testing.T) {
		for _, invalid := range invalidBooks {
			_, err = db.updateBookWithRet(invalid)
			if err == nil {
				t.Errorf("gave invalid ratting, did not error: ratting=%d", invalid.Ratting)
			}
		}
	})

	timeNow := time.Now()

	loans := []*storage.Loan{
		&storage.Loan{
			Name: "Loan One",
			Date: timeNow.Add(time.Hour),
			BookID: books[0].ID,
		},
		&storage.Loan{
			Name: "Loan Two",
			Date: timeNow.Add(time.Hour * (24 * 1)),
			BookID: books[1].ID,
		},
		&storage.Loan{
			Name: "Loan Three",
			Date: timeNow.Add(time.Hour * (24 * 2)),
			BookID: books[2].ID,
		},
		&storage.Loan{
			Name: "Loan Four",
			Date: time.Now().Add(time.Hour * (24 * 3)),
			BookID: books[3].ID,
		},
		&storage.Loan{
			Name: "Loan Five",
			Date: time.Now().Add(time.Hour * (24 * 4)),
			BookID: books[4].ID,
		},
		&storage.Loan{
			Name: "Loan Six",
			Date: time.Now().Add(time.Hour * (24 * 5)),
			BookID: books[5].ID,
		},
	}

	t.Run("CreateLoan", func(t *testing.T) {
		for i, loan := range loans {
			aloan, err := db.createLoanWithRet(loan)
			if err != nil {
				t.Error(err)
				continue
			}
			if aloan.Name != loan.Name {
				t.Errorf("loan[%d] got name='%s', want name='%s'", i, aloan.Name, loan.Name)
			}
			if aloan.Date != loan.Date.Unix() {
				t.Errorf("loan[%d] got date=%d, want date=%d", i, aloan.Date, loan.Date.Unix())
			}
			if aloan.BookID != loan.BookID {
				t.Errorf("loan[%d] got bookid=%d, want bookid=%d", i, aloan.BookID, loan.BookID)
			}
			loans[i].ID = aloan.ID
		}
	})

	t.Run("UpdateLoan", func(t *testing.T) {
		NameTo := "example_name_0"
		DateTo := time.Now().Add(time.Hour * (24 * 12))
		for i, loan := range loans {
			loan.Name = NameTo
			loan.Date = DateTo
			aloan, err := db.updateLoanWithRet(loan)
			if err != nil {
				t.Error(err)
				continue
			}
			if aloan.Name != loan.Name {
				t.Errorf("loan[%d] got name='%s', want name='%s'", i, aloan.Name, loan.Name)
			}
			if aloan.Date != loan.Date.Unix() {
				t.Errorf("loan[%d] got date=%d, want date=%d", i, aloan.Date, loan.Date.Unix())
			}
		}
	})


	t.Run("DeleteBook", func(t *testing.T) {
		list, err := db.GetAllBooks()
		if err != nil {
			t.Error(err)
		}
		listedCount := len(list)
		deletedBook := books[0]

		if err := db.DeleteBook(deletedBook); err != nil {
			t.Error(err)
		}
		list, err = db.GetAllBooks()
		if err != nil {
			t.Error(err)
		}
		if listedCount == len(list) {
			t.Errorf("list size match: %d == %d", listedCount, len(list))
		}
		for _, book := range list {
			if book.ID == deletedBook.ID {
				t.Errorf("deleted book id was found: %d", deletedBook.ID)
			}
		}
		loans, err := db.GetAllLoans()
		if err != nil {
			t.Error(err)
		}
		for _, loan := range loans {
			if loan.BookID == deletedBook.ID {
				t.Errorf("deleted book id was found in loans: %d", deletedBook.ID)
			}
		}
	})

	t.Run("DeleteLoan", func(t *testing.T) {
		list, err := db.GetAllLoans()
		if err != nil {
			t.Error(err)
		}
		listedCount := len(list)
		deletedLoan := loans[1]

		if err := db.DeleteLoan(deletedLoan); err != nil {
			t.Error(err)
		}

		list, err = db.GetAllLoans()
		if err != nil {
			t.Error(err)
		}
		if listedCount == len(list) {
			t.Errorf("list size match: %d == %d", listedCount, len(list))
		}

		for _, loan := range list {
			if loan.ID == deletedLoan.ID {
				t.Errorf("deleted book id was found: %d", deletedLoan.ID)
			}
		}

		books, err := db.GetAllBooks()
		if err != nil {
			t.Error(err)
		}

		found := false
		for _, book := range books {
			if deletedLoan.BookID == book.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("deleted loan's book was not found: %d", deletedLoan.BookID)
		}
	})
}
















