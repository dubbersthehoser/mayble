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

func (a *Application) BookQuery(q *repo.Query) ([]repo.BookEntry, error) {
	return nil, nil
}
