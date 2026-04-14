package app


import (
	"github.com/dubbersthehoser/mayble/internal/config"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
	"github.com/dubbersthehoser/mayble/internal/database"
	"github.com/dubbersthehoser/mayble/internal/bus"
)


type Service struct {
	cfg *config.Config
	dbs *database.Service
	store repo.BookStore

	bookRetriever  repo.BookRetriever
}

func (as *Service) OpenDB(s string) error {
	err := as.dbs.Open(s)
	if err != nil {
		return err
	}
	as.cfg.DBFile = s
	return nil
}

func NewService(bus *bus.Bus, cfg *config.Config, db *database.Database) *Service {
	as := &Service{
		cfg: cfg,
		dbs: database.NewService(db),
	}
	as.bookRetriever = db
	as.store = db
	return as
}

// ??

func (as *Service) CreateBook(b *repo.BookEntry) (int64, error) {
	
}

func (as *Service) UpdateBook(b *repo.BookEntry) error {
}

func (as *Service) DeleteBook(id int64) error {
}

func (as *Service) GetUniqueGenres() ([]string, error) {
}

func (as *Service) GetAllBooks(repo.Variant) ([]repo.BookEntry, error) {
}

func (as *Service) GetBookByID(id int64) (repo.BookEntry, error) {
}
