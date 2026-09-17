package database

import (
	"os"
	"path/filepath"
	"testing"
	"context"

	"github.com/dubbersthehoser/mayble/internal/models"
)

func TestOpenMem(t *testing.T) {

	db, err := OpenMem()
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}

	book := models.BookEntry{
		Book: models.Book{
			Title:  "title",
			Author: "author",
			Genre:  "genre",
		},
	}

	id, err := db.CreateBookWithContext(context.Background(), &book)
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}

	actual, err := db.GetBookByID(id)
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}

	book.ID = id

	if actual != book {
		t.Fatalf("expect\n%#v\n  got\n%#v", book, actual)
	}

}

func TestOpen(t *testing.T) {

	dir := os.TempDir()

	path := filepath.Join(dir, "test.db")

	db, err := Open(path)
	if err == nil {
		t.Fatalf("expected error when opening %s which dose not exists", path)
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}
	file.Close()
	file = nil

	db, err = Open(path)
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}
	defer os.Remove(path)

	_, err = os.Lstat(path + ".bak")
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}
	defer os.Remove(path + ".bak")

	book := models.BookEntry{
		Book: models.Book{
			Title:  "title",
			Author: "author",
			Genre:  "genre",
		},
	}

	id, err := db.CreateBookWithContext(context.Background(), &book)
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}

	actual, err := db.GetBookByID(id)
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}

	book.ID = id

	if actual != book {
		t.Fatalf("expect\n%#v\n  got\n%#v", book, actual)
	}

}
