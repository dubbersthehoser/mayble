package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/dubbersthehoser/mayble/internal/sqlite"
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
)

const version int64 = 5

type Database struct {
	Conn    *sql.DB
	Queries *database.Queries
}

// OpenMem create a memory base database.
func OpenMem() (*Database, error) {

	db := &Database{}
	conn, err := sqlite.OpenDB("")
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	db.Conn = conn
	db.Queries = sqlite.GetQueries(db.Conn)

	err = sqlite.MigrateUpTo(conn, version)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	return db, nil
}

// Create create a new database from path.
func Create(path string) (*Database, error) {

	db := &Database{}

	_, err := os.Lstat(path)
	if err == nil {
		return nil, fmt.Errorf("database: create %s: file exists", path)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("database: create %s: %w", path, err)
	}

	conn, err := sqlite.OpenDB(path)
	if err != nil {
		return nil, fmt.Errorf("database: create %s: %w", path, err)
	}

	db.Conn = conn
	db.Queries = sqlite.GetQueries(db.Conn)

	if err := migrate(db.Conn); err != nil {
		return nil, fmt.Errorf("database: migrate %s: %w", path, err)
	}
	return db, nil
}

// Open open database from path.
func Open(path string) (*Database, error) {

	db := &Database{}

	_, err := os.Lstat(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil, fmt.Errorf("database: open %s: file not found", path)
	case err != nil:
		return nil, fmt.Errorf("database: open %s: %w", path, err)
	}

	conn, err := sqlite.OpenDB(path)
	if err != nil {
		return nil, fmt.Errorf("database: open %s: %w", path, err)
	}
	db.Conn = conn
	db.Queries = sqlite.GetQueries(db.Conn)

	if checkIsV1(db.Conn) {
		err := dbBackup(path)
		if err != nil {
			return nil, fmt.Errorf("database: db_backup %s: %w", path, err)
		}
		err = migrate(db.Conn)
		if err != nil {
			return nil, fmt.Errorf("database: migrate %s: %w", path, err)
		}
	}

	return db, nil
}

func dbBackup(path string) error {
	return dbCopy(path, path+".bak")
}

// dbCopy copy file pathFromt to pathTo
func dbCopy(pathFrom, pathTo string) error {
	to, err := os.Create(pathTo)
	if err != nil {
		return err
	}
	from, err := os.Open(pathFrom)
	if err != nil {
		return err
	}
	_, err = from.WriteTo(to)
	if err != nil {
		return err
	}
	return nil
}

// migrate database up to current version.
func migrate(conn *sql.DB) error {
	return sqlite.MigrateUpTo(conn, version)
}

// checkIsV1 check if db is a mayble 1.0.0 database or lower.
func checkIsV1(conn *sql.DB) bool {
	v := sqlite.GetVersion(conn)
	return v < version
}
