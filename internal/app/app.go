package app


import (
	"os"

	"github.com/dubbersthehoser/mayble/internal/config"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/database"
	"github.com/dubbersthehoser/mayble/internal/csv"
)

type Service struct {
	cfg *config.Config
	db *database.Database
}

func NewService(cfg *config.Config, db *database.Database) *Service {
	as := &Service{
		cfg: cfg,
		db: db,
	}
	return as
}

func (as *Service) OpenDB(s string) error {
	ndb, err := database.Open(s)
	if err != nil {
		return err
	}
	err = as.db.Conn.Close()
	if err != nil {
		return err
	}
	as.db.Conn = ndb.Conn
	as.db.Queries = ndb.Queries
	return nil
}

func (as *Service) CloseDB() error {
	return as.db.Conn.Close()
}

func (as *Service) CreateBook(b *repo.BookEntry) (int64, error) {
	return as.db.CreateBook(b)
}

func (as *Service) UpdateBook(b *repo.BookEntry) error {
	return as.db.UpdateBook(b)
}

func (as *Service) DeleteBook(id int64) error {
	return as.db.DeleteBook(id)
}

func (as *Service) GetUniqueGenres() ([]string, error) {
	return as.db.GetUniqueGenres()
}

func (as *Service) GetAllBooks(v repo.Variant) ([]repo.BookEntry, error) {
	return as.db.GetAllBooks(v)
}

func (as *Service) GetBookByID(id int64) (repo.BookEntry, error) {
	return as.db.GetBookByID(id)
}

func (as *Service) ImportFile(path string) error {
	
	books, err := as.db.GetAllBooks(repo.Book)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = csv.Export(file, books)
	if err != nil {
		return err
	}
	return nil
}

func (as *Service) ExportFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	books, err := csv.Import(file)
	if err != nil {
		return err
	}
	for _, book := range books {
		_, err = as.db.CreateBook(&book)
		if err != nil {
			return err
		}
	}
	return nil
}
