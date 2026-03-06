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

// SetDB closes previous db.Conn and sets db.Queries and db.Conn.
func (s *Service) SetDB(db *Database) error {
	if db == nil {
		return errors.New("app_service.setdb: nil database")
	}
	if err := s.db.Conn.Close(); err != nil {
		return err
	}
	s.db.Conn = db.Conn
	s.db.Queries = db.Queries
	return nil
}


func (s *Service) DB() *Database {
	return s.db
}

func (s *Service) Close() error {
	return s.db.Conn.Close()
}
