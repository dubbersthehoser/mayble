package database

type Service struct {
	db *Database
}

func NewService(db *Database) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Open(path string) error {
	ndb, err := Open(path)
	if err != nil {
		return err
	}
	err = s.db.Conn.Close()
	if err != nil {
		return err
	}
	s.db.Conn = ndb.Conn
	s.db.Queries = ndb.Queries
	return nil
}

func (s *Service) DB() *Database {
	return s.db
}

func (s *Service) Close() error {
	return s.db.Conn.Close()
}
