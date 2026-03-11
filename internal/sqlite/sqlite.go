package sqlite

import (
	"fmt"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
	"github.com/dubbersthehoser/mayble/api"
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

// MigrateUp current connection.
func MigrateUpTo(db *sql.DB, version int64) error {
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(api.SQLiteFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmtError(err)
	}
	if err := goose.UpTo(db, "sqlite/schemas", version); err != nil {
		return fmtError(err)
	}
	return nil
}

// GetVersion 
func GetVersion(db *sql.DB) int64 {
	v, err := goose.GetDBVersion(db)
	if err != nil {
		return -1
	}
	return v

}
