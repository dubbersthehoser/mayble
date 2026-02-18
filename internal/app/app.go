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
			Variant: repo.Book,
			Title: "A Example Title",
			Author: "A Example Author",
			Genre: "A Example Genre",
		},
		{
			Variant: repo.Book,
			Title: "B Example Title",
			Author: "B Example Author",
			Genre: "B Example Genre",
		},
		{
			Variant: repo.Book,
			Title: "C Example Title",
			Author: "C Example Author",
			Genre: "C Example Genre",
		},
		{
			Variant: repo.Book,
			Title: "D Example Title",
			Author: "D Example Author",
			Genre: "D Example Genre",
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
