package app

import (
	"database/sql"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type Application struct {
	Config  *config.Config
	DB      *sql.DB
	Queries database.Queries
}

func (a *Application) CreateBook(book *repo.BookEntry) error {
	return nil
}

func (a *Application) UpdateBook(book *repo.BookEntry) error {
	return nil
}

func (a *Application) DeleteBook(book *repo.BookEntry) error {
	return nil
}

func (a *Application) GetAllBooks(v repo.Variant) ([]repo.BookEntry, error) {
	es := []repo.BookEntry{
		{
			ID: 0,
			Variant: repo.BookLoaned | repo.BookRead,
			Title: "Harry Potter",
			Author: "J.K. Rolling",
			Genre: "Fantacy",
		},
		{
			ID: 1,
			Variant: repo.BookLoaned | repo.BookRead,
			Title: "Lord of the Rings",
			Author: "J.R.R Tolkien",
			Genre: "Fantacy",
		},
		{
			ID: 2,
			Variant: repo.BookLoaned | repo.BookRead,
			Title: "The Foundation",
			Author: "Asimov",
			Genre: "Sci-fi",
		},
		{
			ID: 3,
			Variant: repo.BookLoaned | repo.BookRead,
			Title: "The Elements of Style",
			Author: "William Strunk jr.",
			Genre: "Writing",
		},
	}
	return es, nil
}

func (a *Application) GetBookByID(id int64) (repo.BookEntry, error) {
	return repo.BookEntry{}, nil
}

func (a *Application) GetUniqueGenres() ([]string, error) {
	return []string{
		"Cat",
		"Dog",
		"Bird",
	}, nil
}
