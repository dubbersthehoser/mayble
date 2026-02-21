package database

import (
	"os"
	"database/sql"

	"github.com/dubbersthehoser/mayble/internal/sqlite"
	"github.com/dubbersthehoser/mayble/internal/sqlite/database"
	"github.com/dubbersthehoser/mayble/internal/status"
)

const version int64 = 5

type Database struct {
	Conn      *sql.DB
	Queries *database.Queries
}

// OpenMem create a memory base database.
func OpenMem() (*Database, error) {

	const op status.Op = "database.OpenMem"

	db := &Database{}
	conn, err := sqlite.OpenDB("")
	if err != nil {
		return nil, status.E(op, status.Unexpected, status.LevelError, err)
	}
	db.Conn = conn
	db.Queries = sqlite.GetQueries(db.Conn)
	return db, nil
}

// Open open database from path.
func Open(path string) (*Database, error) {
	
	const op status.Op = "database.Open"

	db := &Database{}
	
	conn, err := sqlite.OpenDB(path)
	if err != nil {
		return nil, status.E(op, status.Unexpected, status.LevelError, err)
	}
	db.Conn = conn
	db.Queries = sqlite.GetQueries(db.Conn)

	if checkIsV1(db.Conn) {
		err := migrate(path, db.Conn)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func dbCopy(pathFrom, pathTo string) error {

	const op status.Op = "database.dbCopy"

	to, err := os.Open(pathTo)
	if err != nil {
		return status.E(op, status.Unexpected, status.LevelError, err)
	}
	from, err := os.Open(pathFrom)
	if err != nil {
		return status.E(op, status.Unexpected, status.LevelError, err)
	}
	_, err = from.WriteTo(to)
	if err != nil {
		return status.E(op, status.Unexpected, status.LevelError, err)
	}
	return nil
}

// migrate up database and create backup.
func migrate(path string, conn *sql.DB)  error {

	const op status.Op = "database.migrage"

	err := dbCopy(path, path + ".bak")
	if err != nil {
		return status.E(op, status.Unexpected, status.LevelError, err)
	}
	return sqlite.MigrateUpTo(conn,  version)
}

// checkIsV1 check if db is a mayble 1.0.0 database.
func checkIsV1(conn *sql.DB) bool {
	v := sqlite.GetVersion(conn)
	return v == 3
}


