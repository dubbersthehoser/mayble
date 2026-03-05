package viewmodel

import (
	"testing"
	"slices"
	"strings"
	"bytes"
	"io"
	"time"
	"os"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/database"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)



func Test_DatabaseCreateAndOpen(t *testing.T) {
	b := &bus.Bus{}

	as := &appService{}

	file, err := os.CreateTemp("", "*.db")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	dbPath := file.Name()
	file.Close()

	errCount := 0
	
	b.Subscribe(bus.Handler{
		Name: msgUserInfo,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_info:", e.Data.(string))
		},
	})
	b.Subscribe(bus.Handler{
		Name: msgUserError,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_error:", e.Data.(string))
			errCount += 1
		},
	})
	b.Subscribe(bus.Handler{
		Name: msgUserSuccess,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_success:", e.Data.(string))
		},
	})

	t.Run("create", func(t *testing.T){
		test_createDatabase(t, as, b, dbPath)
		if errCount != 0 {
			t.Fatalf("there was an error in message bus")
		}
	})
	t.Run("open", func(t *testing.T) {
		test_openDatabase(t, as, b, dbPath)
		if errCount != 0 {
			t.Fatalf("there was an error in message bus")
		}
	})
}
func test_createDatabase(t *testing.T, as *appService, b *bus.Bus, path string) {
	var err error
	createDatabase(path, as, b)

	file, err := os.CreateTemp("", "")
	path = file.Name()
	file.Close()
	createDatabase(path, as, b)
	
	expect := path + ".db"

	_, err = os.Lstat(expect)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func test_openDatabase(t *testing.T, as *appService, b *bus.Bus, path string) {
	as.dbs = nil
	openDatabase(path, as, b)
	if as.dbs == nil {
		t.Fatal("database service is nil")
	}
}





func Test_CSVImportAndExport(t *testing.T) {
	b := &bus.Bus{}

	b.Subscribe(bus.Handler{
		Name: msgUserInfo,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_info:", e.Data.(string))
		},
	})
	b.Subscribe(bus.Handler{
		Name: msgUserError,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_error:", e.Data.(string))
		},
	})
	b.Subscribe(bus.Handler{
		Name: msgUserSuccess,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_success:", e.Data.(string))
		},
	})

	csvStr := strings.TrimSpace(`
Title,Author,Genre,,,,
Title,Author,Genre,2021-02-19,3,,
Title,Author,Genre,,,2021-02-19,Lane
Title,Author,Genre,2021-02-19,3,2021-02-19,Lane
`)
	books := []repo.BookEntry{
		{
			Variant: repo.Book,
			ID: 1,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
		},
		{
			Variant: repo.Book | repo.Read,
			ID: 2,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating: 3,
		},
		{
			Variant: repo.Book | repo.Loaned,
			ID: 3,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
		},
		{
			Variant: repo.Book | repo.Loaned | repo.Read,
			ID: 4,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
			Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating: 3,
		},
	}

	db, err := database.OpenMem()

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	t.Run("import", func(t *testing.T) {
		testImportCSV(t, b, db, bytes.NewBuffer([]byte(csvStr)), books)
	})

	t.Run("export", func(t *testing.T) {
		testExportCSV(t, b, db, csvStr)
	})
}

func testImportCSV(t *testing.T, b *bus.Bus, db *database.Database, input io.Reader, expect []repo.BookEntry) {

	importCSV(input, b, db)
	actual, err := db.GetAllBooks(repo.Book)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(expect) != len(actual) {
		t.Fatalf("expect length %d, got %d", len(expect), len(actual))
	}

	for _, book := range expect {
		if !slices.Contains(actual, book) {
			t.Fatalf("book not found:\n  %v", book)
		}
	}
}

func testExportCSV(t *testing.T, b *bus.Bus, db *database.Database, expect string) {
	file, err := os.CreateTemp("", "*.csv")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	exportCSV(file, file.Name(), b, db)

	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	raw, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	actual := strings.TrimSpace(string(raw))
	if actual != expect {
		t.Fatalf("expect\n'%s'\ngot\n'%s'", expect, actual)
	}

	file.Close()
	os.Remove(file.Name())



	file, err = os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	filepath := file.Name()
	expectPath := filepath + ".csv"

	exportCSV(file, file.Name(), b, db)
	
	_, err = os.Lstat(expectPath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}









