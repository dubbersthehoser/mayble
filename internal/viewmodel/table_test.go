package viewmodel

import (
	"testing"
	"slices"
	"time"
	"fmt"
	"math/rand"

	"github.com/dubbersthehoser/mayble/internal/database"
	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/config"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)


func TestTableVM(t *testing.T) {
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
	table := NewTableVM(b, as)


	books := []repo.BookEntry{
		{
			ID: 1,
			Variant: repo.Book,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
		},
		{
			ID: 2,
			Variant: repo.Book | repo.Read,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating: 3,
		},
		{
			ID: 3,
			Variant: repo.Book | repo.Loaned,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
		},
		{
			ID: 4,
			Variant: repo.Book | repo.Loaned | repo.Read,
			Title: "Title",
			Author: "Author",
			Genre: "Genre",
			Loaned: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Borrower: "Lane",
			Read: time.Date(2021, 2, 19, 0, 0, 0, 0, time.UTC),
			Rating: 3,
		},
	}
	for _, book := range books {
		_, err := db.CreateBook(&book)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	t.Run("StoreColumnWidth", func(t *testing.T){
		for i := range table.Headers() {
			actual := table.GetColumnWidth(i)
			if actual != MinColWidth {
				t.Fatalf("expect %f, got %f", MinColWidth, actual)
			}
		}
		for i := range table.Headers() {
			table.StoreColumnWidth(i, MinColWidth-1)
			actual := table.GetColumnWidth(i)
			if actual != MinColWidth {
				t.Fatalf("expect %f, got %f", MinColWidth, actual)
			}
		}
		for i := range table.Headers() {
			expect := MinColWidth + 100.0
			table.StoreColumnWidth(i, expect)
			actual := table.GetColumnWidth(i)
			if actual != expect {
				t.Fatalf("expect %f, got %f", expect, actual)
			}
		}
	})

	t.Run("load, reload, and Sort", func(t *testing.T) {
		err := table.load()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		err = table.reload()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		err = table.Sort()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		testTableVMGet(t, table, books)
	})

	t.Run("Select", func(t *testing.T) {
		row, col := 1, 1
		table.Select(row, col)
		expect := books[row]
		actual, err := table.Selector().getBook()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if actual == nil {
			t.Fatal("unexpected nil")
		}
		if expect != *actual {
			t.Fatalf("expect\n%v\n  got\n%v", expect, *actual)
		}
	})

	t.Run("search", func(t *testing.T) {
		options := table.SearchOptions()
		expect := []string{
			"All",
			"Title",
			"Author",
			"Genre",
			"Borrower",
		}
		if !slices.Equal(options, expect) {
			t.Fatalf("expect\n%#v\n  got\n%#v", expect, options)
		}
		search := "A Title"
		_ = table.Search.Text.Set(search)
		table.search()
		if !table.selector.HasSelected() {
			t.Fatalf("expected to have selected")
		}
	})

	t.Run("bus: data changed", func(t *testing.T){
		table.Select(0, 0)
		b.Notify(bus.Event{
			Name: msgDataChanged,
		})
		if table.selector.HasSelected() {
			t.Fatalf("expected to not have selected")
		}
	})

	t.Run("SetHidden", func(t *testing.T) {
		expect := []string{"Title", "Author"}
		table.SetHidden(expect)
		actual := cfg.GetUITable().ColumnsHidden
		if !slices.Equal(expect, actual) {
			t.Fatalf("expect\n%#v\n  got\n%#v", expect, actual )
		}

		expect = table.Hidden()
		if !slices.Equal(expect, actual) {
			t.Fatalf("expect\n%#v\n  got\n%#v", expect, actual )
		}

		col := len(entryHeaders())-1
		if ok := table.IsItemHidden(0, col); !ok {
			t.Fatalf("expect %t, got %t", true, ok)
		}
		if ok := table.IsHeaderHidden(col); !ok {
			t.Fatalf("expect %t, got %t", false, ok)
		}
		col = 0
		if ok := table.IsItemHidden(0, col); ok {
			t.Fatalf("expect %t, got %t", true, ok)
		}
		if ok := table.IsHeaderHidden(col); ok {
			t.Fatalf("expect %t, got %t", false, ok)
		}

		table.SetHidden([]string{})
		if len(table.Hidden()) != 0 {
			t.Fatalf("expect length %d, got %d", 0, len(table.Hidden()))
		}
	})
}

func testTableVMGet(t *testing.T, table *TableVM, books []repo.BookEntry) {
	row, col := table.Size()
	if row != len(books) {
		t.Fatalf("expect %d, got %d", len(books), row)
	}
	for col != len(entryHeaders()) {
		t.Fatalf("expect %d, got %d", len(entryHeaders()), col)
	}
	for row := range books {
		t.Run(fmt.Sprintf("book#%d.variant='%v'", row, books[row].Variant), func(t *testing.T) {
			id, err := table.GetID(row, 0)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if id != books[row].ID {
				t.Fatalf("expect %d, got %d", books[row].ID, id)
			}
			for col := range entryHeaders() {
				actual := table.Get(row, col)
				if (col == 3 || col == 4) && books[row].Variant & repo.Read == 0 ||
				   (col == 5 || col == 6) && books[row].Variant & repo.Loaned == 0 {
					if actual != stubValue {
						t.Fatalf("expect '%s', got '%s'", stubValue, actual)
					}
					continue
				}
				if actual == stubValue {
					t.Fatalf("unwanted value '%s', at column %d", actual, col)
				}
			}
		})
	}
	actual := table.Get(0, len(books))
	if actual != stubValue {
		t.Fatalf("expect '%s', got '%s'", "N/A", actual)
	}
}

func Test_sortBooks(t *testing.T) {
	expect := []repo.BookEntry{
		{
			Title: "A",
			Author: "A",
			Genre: "A",
			Read: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			Rating: 1,
			Loaned: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			Borrower: "A",
		},
		{
			Title: "B",
			Author: "B",
			Genre: "B",
			Read: time.Date(2, 2, 2, 0, 0, 0, 0, time.UTC),
			Rating: 2,
			Loaned: time.Date(2, 2, 2, 0, 0, 0, 0, time.UTC),
			Borrower: "B",
		},
		{
			Title: "C",
			Author: "C",
			Genre: "C",
			Loaned: time.Date(3, 3, 3, 0, 0, 0, 0, time.UTC),
			Borrower: "C",
			Read: time.Date(3, 3, 3, 0, 0, 0, 0, time.UTC),
			Rating: 3,
		},
		{
			Title: "D",
			Author: "D",
			Genre: "D",
			Loaned: time.Date(4, 4, 4, 0, 0, 0, 0, time.UTC),
			Borrower: "D",
			Read: time.Date(4, 4, 4, 0, 0, 0, 0, time.UTC),
			Rating: 4,
		},
	}

	actual := slices.Clone(expect)
	shuffle := func() {
		rand.Shuffle(len(actual), func(i, j int){
			actual[i], actual[j] = actual[j], actual[i]
		})
	}

	for _, header := range entryHeaders() {
		shuffle()
		slices.Reverse(expect)
		sortBooks(actual, header, true)
		if !slices.Equal(expect, actual) {
			t.Fatal("failed to sort")
		}

		shuffle()
		slices.Reverse(expect)
		sortBooks(actual, header, false)
		if !slices.Equal(expect, actual) {
			t.Fatal("failed to sort")
		}
	}

	err := sortBooks(expect, "invalid", true)
	if err == nil {
		t.Fatal("expected error")
	}
}


func Test_hiddenOptionsToHeaders(t *testing.T) {
	tests := []struct{
		input  []string
		expect []string
	}{
		{
			input: []string{"Loaned"},
			expect: []string{"Loaned", "Borrower"},
		},
		{
			input: []string{"Read"},
			expect: []string{"Read", "Rating"},
		},
		{
			input: []string{"Read", "Loaned"},
			expect: []string{"Read", "Rating", "Loaned", "Borrower"},
		},
		{
			input: []string{"Read", "Loaned", "extra", "extra"},
			expect: []string{"Read", "Rating", "Loaned", "Borrower", "extra", "extra"},
		},
	}

	for i, c := range tests {
		actual := hiddenOptionsToHeaders(c.input)
		slices.Sort(actual)
		slices.Sort(c.expect)
		if n := slices.Compare(c.expect, actual); n != 0 {
			t.Fatalf("[%d] expect\n\t%#v\ngot\n\t%#v", i, c.expect, actual)
		}
	}
}


func Test_hiddenHeadersToOptions(t *testing.T) {
	
	tests := []struct{
		input  []string
		expect []string
	}{
		{
			input: []string{"Loaned", "Borrower"},
			expect: []string{"Loaned"},
		},
		{
			input: []string{"Read", "Rating"},
			expect: []string{"Read"},
		},
		{
			input: []string{"Read", "Rating", "Loaned", "Borrower"},
			expect: []string{"Loaned", "Read"},
		},
		{
			input: []string{"Read", "Rating", "Loaned", "Borrower", "Title", "Author"},
			expect: []string{"Title", "Author", "Loaned", "Read"},
		},
	}
	for i, c := range tests {
		actual := hiddenHeadersToOptions(c.input)
		slices.Sort(actual)
		slices.Sort(c.expect)
		if n := slices.Compare(c.expect, actual); n != 0 {
			t.Fatalf("[%d] expect\n\t%#v\ngot\n\t%#v", i, c.expect, actual)
		}
	}
}




func TestTableControllerVM(t *testing.T) {
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
	controllers := NewTableControllersVM(b, as)
	_ = controllers
}


