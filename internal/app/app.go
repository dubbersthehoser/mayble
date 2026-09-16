package app

import (
	"os"
	"context"
	"sync"

	"github.com/dubbersthehoser/mayble/internal/command"
	"github.com/dubbersthehoser/mayble/internal/event"
	"github.com/dubbersthehoser/mayble/internal/csv"
	"github.com/dubbersthehoser/mayble/internal/database"
	"github.com/dubbersthehoser/mayble/internal/worker"
	"github.com/dubbersthehoser/mayble/internal/models"
)

type Service struct {
	path string
	mu   sync.RWMutex
	db   *database.Database
}

func NewService(w *worker.Worker, eb *event.EventBus, cb *command.CommandBus) *Service {
	as := &Service{
		db:  nil,
	}
	as.setupCommands(w, eb, cb)
	return as
}

func (as *Service) Path() string {
	return as.path
}

func (as *Service) CloseDB() error {
	if !hasDatabase(as) {
		return nil
	}
	return as.db.Conn.Close()
}

func (as *Service) createBook(b *models.BookEntry) (int64, error) {
	as.mu.Lock()
	defer as.mu.Unlock()
	if !hasDatabase(as) {
		return 0, nil
	}
	id, err := as.db.CreateBookWithContext(context.Background(), b)
	return id, err
}

func (as *Service) updateBook(b *models.BookEntry) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	if !hasDatabase(as) {
		return nil
	}
	err := as.db.UpdateBook(b)
	return err
}

func (as *Service) deleteBook(id int64) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	if !hasDatabase(as) {
		return nil
	}
	err := as.db.DeleteBook(id)
	return err
}

func (as *Service) getAllBooksWithContext(ctx context.Context) ([]models.BookEntry, error) {
	if !hasDatabase(as) {
		return []models.BookEntry{}, nil
	}
	return as.db.GetAllBooksWithContext(ctx)
}

func (as *Service) exportFileWithContext(ctx context.Context, path string) error {
	if !hasDatabase(as) {
		return nil
	}
	books, err := as.db.GetAllBooksWithContext(ctx)
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

func (as *Service) importFileWithContext(ctx context.Context, path string) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	if !hasDatabase(as) {
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
		_, err = as.db.CreateBookWithContext(ctx, &book)
		if err != nil {
			return err
		}
	}
	return nil
}

func (as *Service) createDatabase(path string) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	db, err := database.Create(path)
	if err != nil {
		return err
	}
	if err := swap(as, db); err != nil {
		return err
	}
	as.path = path
	return nil
}

func (as *Service) openDatabase(path string) error {
	as.mu.Lock()
	defer as.mu.Unlock()
	db, err := database.Open(path)
	if err != nil {
		return err
	}
	if err := swap(as, db); err != nil {
		return err
	}
	as.path = path
	return nil
}

func hasDatabase(as *Service) bool {
	return as.db != nil
}

// swap database. Should only be called in [Service.openDatabase], and [Service.CreateDatabase].
func swap(as *Service, db *database.Database) error {
	if as.db == nil {
		as.db = db
	} else {
		err := as.db.Conn.Close()
		if err != nil {
			return err
		}
		as.db.Conn = db.Conn
		as.db.Queries = db.Queries
	}
	return nil
}

func (as *Service) setupCommands(w *worker.Worker, eb *event.EventBus, cb *command.CommandBus) {

	//
	// Opening and Creating Database.
	//
	cb.Register(command.OpenDatabase{}, func(v command.Command) error {
		e := v.(command.OpenDatabase)
		err := as.openDatabase(e.Path)
		if err != nil {
			eb.Notify(event.OpenedDatabase{
				Path: e.Path,
				Failed: true,
				Message: err.Error(),
				Err: err,
			})
		} else {
			eb.Notify(event.OpenedDatabase{
				Path: e.Path,
				Failed: false,
			})
			w.Jobs <- NewJobTakeSnapshot(w, as)
		}
		return nil
	})
	cb.Register(command.CreateDatabase{}, func(v command.Command) error {
		e := v.(command.CreateDatabase)
		err := as.createDatabase(e.Path)
		if err != nil {
			eb.Notify(event.CreatedDatabase{
				Path: e.Path,
				Failed: true,
				Message: err.Error(),
				Err: err,
			})
		} else {
			eb.Notify(event.CreatedDatabase{
				Path: e.Path,
				Failed: false,
			})
			w.Jobs <- NewJobTakeSnapshot(w, as)
		}
		return nil
	})

	//
	// Create, Update, and Delete Book Entry.
	//
	cb.Register(command.CreateBookEntry{}, func(v command.Command) error{
		e := v.(command.CreateBookEntry)
		_, err := as.createBook(&e.Book)
		if err != nil {
			eb.Notify(event.CreatedBookEntry{
				Message: err.Error(),
				Failed: true,
			})
		} else {
			eb.Notify(event.CreatedBookEntry{
				Failed: false,
			})
			w.Jobs <- NewJobTakeSnapshot(w, as)
		}
		return nil
	})
	cb.Register(command.UpdateBookEntry{}, func(v command.Command) error{
		e := v.(command.UpdateBookEntry)
		err := as.updateBook(&e.Book)
		if err != nil {
			eb.Notify(event.UpdatedBookEntry{
				Failed: true,
				Message: err.Error(),
			})
		} else {
			eb.Notify(event.UpdatedBookEntry{
				Failed: false,
			})
			w.Jobs <- NewJobTakeSnapshot(w, as)
		}
		return nil
	})
	cb.Register(command.DeleteBookEntry{}, func(v command.Command) error {
		e := v.(command.DeleteBookEntry)
		err := as.deleteBook(e.BookID)
		if err != nil {
			eb.Notify(event.DeletedBookEntry{
				Message: err.Error(),
				Failed: true,
			})
		} else {
			eb.Notify(event.DeletedBookEntry{
				Failed: false,
			})
			w.Jobs <- NewJobTakeSnapshot(w, as)
		}
		return nil
	})

	//
	// File Import, and Export.
	//
	cb.Register(command.ImportFile{}, func(v command.Command) error {
		e := v.(command.ImportFile)
		w.Jobs <- NewJobImportFile(w, as, e.Path)
		return nil
	})
	cb.Register(command.ExportFile{}, func(v command.Command) error {
		e := v.(command.ExportFile)
		w.Jobs <- NewJobExportFile(w, as, e.Path)
		return nil
	})
}
