package database

import (
	"errors"
)

type Service struct {
	db *Database
}

func NewService(db *Database) *Service {
	return &Service{
		db: db,
	}
}

// SetDB closes previous database and sets db.
func (s *Service) SetDB(db *Database) error {
	if db == nil {
		return errors.New("setdb: nil database")
	}
	if err := s.db.Conn.Close(); err != nil {
		return err
	}
	s.db = db
	return nil
}


func (s *Service) DB() *Database {
	return s.db
}

func (s *Service) Close() error {
	return s.db.Conn.Close()
}
