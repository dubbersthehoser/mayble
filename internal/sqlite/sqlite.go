package sqlite

import (
	"fmt"
	"io/fs"
	"embed"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	
	"github.com/dubbersthehoser/mayble/internal/sqlitedb/database"
	"github.com/dubbersthehoser/mayble/api"
)


var schemaDir string = "sqlite/schemas"

// Schema contains schemas for migrations. Primary use is for embed file system and testing.
type Schema struct {
	Dir string // directory path to the schema directory.
	FS   fs.FS // root filesystem location of the schema directory.
}

// Database contains queries, Schema, and db connection
type Database struct {
	Queries *database.Queries // sqlc queries
	DB      *sql.DB           // database connection
	Schema   Schema           
}

// NewDatabase create a new database. 
func NewDatabase() *Database {
	goose.SetVerbose(false)
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(api.SQLiteFS)
	return &Database{
		Schema: Schema{
			Dir: schemaDir,
			FS:  api.SQLiteFS,
		},
	}
}

func fmtError(err error) error {
	return fmt.Errorf("sqlite: %w", err)
}

// EnableForeignKeys set sqlite's foreign_key to ON for current connection
func (d *Database) EnableForeignKeys() error {
	_, err := d.DB.Exec("PRAGMA foreign_keys = ON;")
	return err
}

// DisableForeignKeys set sqlite's foreign_key to OFF for current connection
func (d *Database) DisableForeignKeys() error {
	_, err := d.DB.Exec("PRAGMA foreign_keys = OFF;")
	return err
}

// Open create a connection to a sqlite file.
func (d *Database) Open(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return fmtError(err)
	}	
	d.DB = db
	d.Queries =  database.New(db)
	d.EnableForeignKeys()
	return nil
}

// Close connection and set DB to nil.
func (d *Database) Close() error {
	err := d.DB.Close()
	if err != nil {
		return err
	}
	d.DB = nil
	return err
}

// MigrateUp using Schema with current connection.
func (db *Database) MigrateUp() error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmtError(err)
	}
	
	if err := goose.Up(db.DB, db.Schema.Dir); err != nil {
		return fmtError(err)
	}
	return nil
}
