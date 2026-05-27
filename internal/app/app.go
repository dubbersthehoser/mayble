package app


import (
	"os"
	"errors"

	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/database"
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/csv"
)

type Service struct {
	cfg *config.Config

	db *database.Database

	listeners []func()
}

func NewService(cfg *config.Config) *Service {
	as := &Service{
		cfg: cfg,
		db: nil,

		listeners: make([]func(), 0),
	}
	return as
}

func (as *Service) noDatabase() bool {
	return as.db == nil
}

func (as *Service) CloseDB() error {
	if as.noDatabase() {
		return nil
	}
	return as.db.Conn.Close()
}

func (as *Service) CreateBook(b *models.BookEntry) (int64, error) {
	if as.noDatabase() {
		return 0, nil
	}
	id, err := as.db.CreateBook(b)
	if err == nil {
		as.notify()
	}
	return id, err
}

func (as *Service) UpdateBook(b *models.BookEntry) error {
	if as.noDatabase() {
		return nil
	}
	err := as.db.UpdateBook(b)
	if err == nil {
		as.notify()
	}
	return err
}

func (as *Service) DeleteBook(id int64) error {
	if as.noDatabase() {
		return nil
	}
	err := as.db.DeleteBook(id)
	if err == nil {
		as.notify()
	}
	return err
}

func (as *Service) GetUniqueGenres() ([]string, error) {
	if as.noDatabase() {
		return []string{}, nil
	}
	return as.db.GetUniqueGenres()
}

func (as *Service) GetAllBooks() ([]models.BookEntry, error) {
	if as.noDatabase() {
		return []models.BookEntry{}, nil
	}
	return as.db.GetAllBooks()
}

func (as *Service) GetBookByID(id int64) (models.BookEntry, error) {
	if as.noDatabase() {
		return models.BookEntry{}, errors.New("nil database")
	}
	return as.db.GetBookByID(id)
}

func (as *Service) ExportFile(path string) error {
	if as.noDatabase() {
		return nil
	}
	books, err := as.db.GetAllBooks()
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

func (as *Service) ImportFile(path string) error {
	if as.noDatabase() {
		return nil
	}
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
	as.notify()
	return nil
}

func (as *Service) LoadDatabase() error {
	return as.OpenDatabase(as.cfg.DBFile)
}

func (as *Service) OpenDatabase(path string) error {
	db, err := database.Open(path)
	if err != nil {
		return err
	}

	if as.db == nil {
		as.db = db
	} else {
		err = as.db.Conn.Close()
		if err != nil {
			return err
		}
		as.db.Conn = db.Conn
		as.db.Queries = db.Queries
	}
	as.notify()
	as.cfg.DBFile = path
	return nil
}

func (as *Service) AddListener(fn func()) {
	as.listeners = append(as.listeners, fn)
}

func (as *Service) notify() {
	for _, fn := range as.listeners {
		fn()
	}
}
