package database

import (
	"testing"
	"path/filepath"
	"os"

	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

func TestOpenMem(t *testing.T) {
	
	db, err := OpenMem()
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}

	book := repo.BookEntry{
		Title: "title",
		Author: "author",
		Genre: "genre",
	}

	id, err := db.CreateBook(&book)
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
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}
	defer os.Remove(path)

	_, err = os.Lstat(path + ".bak")
	if err != nil {
		t.Fatalf("unexpected error: '%s'", err)
	}
	defer os.Remove(path + ".bak")

	book := repo.BookEntry{
		Title: "title",
		Author: "author",
		Genre: "genre",
	}

	id, err := db.CreateBook(&book)
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
