package table

import (
	"fmt"
	"testing"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/models"
)

func TestTable(t *testing.T) {

	books := make([]models.BookEntry, 0)

	for i := range 6 {
		book := models.BookEntry{
			ID: int64(i),
			Book: models.Book{
				Title:  fmt.Sprintf("title-%d", i),
				Author: fmt.Sprintf("author-%d", i),
				Genre:  fmt.Sprintf("genre-%d", i),
			},
		}
		books = append(books, book)
	}

	source := func() ([]models.BookEntry, error) {
		return books, nil
	}

	cfg := config.NewConfigWithDefaults("")

	table := NewTable(cfg, source)

	// Test Sheet

	// Load
	err := table.Sheet.Load()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	table.Settings.SetIDHidden(true)
	table.Sorting.SetOrderBy(models.BookEntryFields()[models.IdxGenre])
	table.Sorting.SetAscending(true)
	table.Sorting.Sort()

	{ // Get
		p := Point{
			Row: len(books) - 1,
			Col: 2,
		}
		expect := fmt.Sprintf("genre-%d", len(books)-1)
		actual, err := table.Sheet.Get(p)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if expect != actual {
			t.Fatalf("expect '%s', got '%s'", expect, actual)
		}
	}

}
