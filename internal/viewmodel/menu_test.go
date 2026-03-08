package viewmodel

import (
	"testing"
	"slices"
	"strings"
	"bytes"
	"io"
	"time"
	"os"

	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

func TestMenuUI(t *testing.T) {
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

	DBPath := binding.NewString()

	menu := NewMenuVM(b, as, DBPath)
	_ = menu


	t.Run("CSVImportAndExport", func(t *testing.T){
		testCSVImportAndExport(t, menu)
	})

	t.Run("DatabaseCreateAndOpen", func(t *testing.T){
		testDatabaseCreateAndOpen(t, menu)
	})
}

func testDatabaseCreateAndOpen(t *testing.T, menu *MenuVM) {
	file, err := os.CreateTemp("", "*.db")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	dbPath := file.Name()
	file.Close()

	errCount := 0
	
	menu.bus.Subscribe(bus.Handler{
		Name: msgUserInfo,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_info:", e.Data.(string))
		},
	})
	menu.bus.Subscribe(bus.Handler{
		Name: msgUserError,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_error:", e.Data.(string))
			errCount += 1
		},
	})
	menu.bus.Subscribe(bus.Handler{
		Name: msgUserSuccess,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_success:", e.Data.(string))
		},
	})

	t.Run("create", func(t *testing.T){
		test_createDatabase(t, menu, dbPath)
		if errCount != 0 {
			t.Fatalf("there was an error in message bus")
		}
	})
	t.Run("open", func(t *testing.T) {
		test_openDatabase(t, menu, dbPath)
		if errCount != 0 {
			t.Fatalf("there was an error in message bus")
		}
	})
}

func test_createDatabase(t *testing.T, menu *MenuVM, path string) {
	var err error

	menu.CreateDatabase(path, err)

	file, err := os.CreateTemp("", "")
	path = file.Name()
	file.Close()
	menu.CreateDatabase(path, err)
	
	expect := path + ".db"

	_, err = os.Lstat(expect)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func test_openDatabase(t *testing.T, menu *MenuVM, path string) {
	menu.OpenDatabase(path, nil)
	if menu.app.dbs == nil {
		t.Fatal("database service is nil")
	}
}




func testCSVImportAndExport(t *testing.T, menu *MenuVM) {

	menu.bus.Subscribe(bus.Handler{
		Name: msgUserInfo,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_info:", e.Data.(string))
		},
	})
	menu.bus.Subscribe(bus.Handler{
		Name: msgUserError,
		Handler: func(e *bus.Event) {
			t.Log("bus.msg_user_error:", e.Data.(string))
		},
	})
	menu.bus.Subscribe(bus.Handler{
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
	defer db.Conn.Close()

	t.Run("import", func(t *testing.T) {
		testImportCSV(t, menu, bytes.NewBuffer([]byte(csvStr)), books)
	})

	t.Run("export", func(t *testing.T) {
		testExportCSV(t, menu, csvStr)
	})
}

func testImportCSV(t *testing.T, menu *MenuVM, input io.Reader, expect []repo.BookEntry) {

	menu.ImportCSV(io.NopCloser(input), nil)
	actual, err := menu.app.bookRetriever.GetAllBooks(repo.Book)
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

func testExportCSV(t *testing.T, menu *MenuVM, expect string) {
	file, err := os.CreateTemp("", "*.csv")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	menu.ExportCSV(file, file.Name(), nil)

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

	menu.ExportCSV(file, file.Name(), nil)
	
	_, err = os.Lstat(expectPath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}









