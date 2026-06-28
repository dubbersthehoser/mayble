package app


import (
	"os"
	"errors"
	"slices"
	"fmt"
	"cmp"
	"strings"

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
	as.cfg.DBFile = path
	as.notify()
	return nil
}

// AddListener listen for database changes.
func (as *Service) AddListener(fn func()) {
	as.listeners = append(as.listeners, fn)
}

func (as *Service) notify() {
	for _, fn := range as.listeners {
		fn()
	}
}

// SortBooks sort slice of book entries.
func SortBooks(books []models.BookEntry, index int, ascending bool) error {

	if !(models.IdxID <= index && models.IdxBorrower >= index) {
		return fmt.Errorf("sort_books: invalid index '%d'", index)
	}

	slices.SortFunc(books, func(a, b models.BookEntry) int {

		// keep all the non-active values to the bottom of list.
		switch index {
		case models.IdxRating, models.IdxCompletedAt:
			if !a.IsCompleted && !b.IsCompleted {
				return 0
			}
			if !a.IsCompleted {
				return 1
			}
			if !b.IsCompleted {
				return -1
			}
		case models.IdxBorrower, models.IdxLoanedAt:
			if !a.IsLoaned && !b.IsLoaned {
				return 0
			}
			if !a.IsLoaned {
				return 1
			}
			if !b.IsLoaned {
				return -1
			}
		}
		
		r := -1
		switch index {
		case models.IdxTitle:
			r = cmp.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
		case models.IdxAuthor:
			r = cmp.Compare(strings.ToLower(a.Author), strings.ToLower(b.Author))
		case models.IdxGenre:
			r = cmp.Compare(strings.ToLower(a.Genre), strings.ToLower(b.Genre))
		case models.IdxBorrower:
			r = cmp.Compare(strings.ToLower(a.Borrower), strings.ToLower(b.Borrower))
		case models.IdxLoanedAt:
			r = a.Loaned.LoanedAt.Compare(b.LoanedAt)
		case models.IdxRating:
			r = cmp.Compare(a.Rating, b.Rating)
		case models.IdxCompletedAt:
			r = a.CompletedAt.Compare(b.CompletedAt)
		}
		if ascending {
			return r
		} else {
			return r * -1
		}
	})
	return nil
}
