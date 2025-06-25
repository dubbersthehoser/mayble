package app

import (
	
	"database/sql"
	"embed"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	
	"github.com/dubbersthehoser/mayble/internal/database"
)

var SchemaDir string
var SchemaFS embed.FS

type Database struct {
	Queries *database.Queries
	DB      *sql.DB
}

func OpenDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	ndb := &Database{
		DB: db,
		Queries: database.New(db),
		}
	return ndb, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func DatabaseMigrateUp(db *Database) error {
	goose.WithVerbose(false)
	goose.SetBaseFS(SchemaFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	if err := goose.Up(db.DB, SchemaDir); err != nil {
		return err
	}
	return nil
}












