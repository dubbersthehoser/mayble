package sqlite

import (
	"fmt"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
)

func fmtError(err error) error {
	return fmt.Errorf("sqlite: %w", err)
}

// OpenDB data connection to file
func OpenDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite3", path)
}

// GetQueries wapper for database queries.
func GetQueries(db *sql.DB) *database.Queries {
	return database.New(db)
}

// MigrateUp using Schema with current connection.
func MigrateUpTo(db *sql.DB, schemas string, version int64) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmtError(err)
	}
	if err := goose.UpTo(db, schemas, version); err != nil {
		return fmtError(err)
	}
	return nil
}
