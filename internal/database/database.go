package database

import (
	"database/sql"
	"os"
	"fmt"

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

// Open open database from path.
func Open(path string) (*Database, error) {

	db := &Database{}

	conn, err := sqlite.OpenDB(path)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	db.Conn = conn
	db.Queries = sqlite.GetQueries(db.Conn)

	if checkIsV1(db.Conn) {
		err := migrate(path, db.Conn)
		if err != nil {
			return nil, fmt.Errorf("datbase: migrate: %w", err)
		}
	}

	return db, nil
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

// migrate up database and create backup.
func migrate(path string, conn *sql.DB) error {
	err := dbCopy(path, path+".bak")
	if err != nil {
		return err
	}
	return sqlite.MigrateUpTo(conn, version)
}

// checkIsV1 check if db is a mayble 1.0.0 database or lower.
func checkIsV1(conn *sql.DB) bool {
	v := sqlite.GetVersion(conn)
	return v < version
}
