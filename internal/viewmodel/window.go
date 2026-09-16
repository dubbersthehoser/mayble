package viewmodel

import (
	"fmt"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/snapshot"
	"github.com/dubbersthehoser/mayble/internal/command"
	"github.com/dubbersthehoser/mayble/internal/event"
	"github.com/dubbersthehoser/mayble/internal/worker"
	"github.com/dubbersthehoser/mayble/internal/models"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/table"
)

type Window struct {
	cfg *config.Config
	eb  *event.EventBus
	cb  *command.CommandBus
	srv *app.Service

	Body         *Body
	StatusLine   *StatusLine
	Controls     *TableControl
	FileManage   *FileManage
	DBPath       *DBPath
	UniqueGenres *UniqueGenres
	Table        *table.Table
	Form         *BookForm
	NoData       *NoDataBody
	ShowError    *ShowError
	Worker       *worker.Worker
}

func (w *Window) HandleWorkerEvent(ev worker.Event) {
	switch ev.Type {
	case worker.Started:
		log.Print(ev.Format())

	case worker.Finished:
		log.Print(ev.Format())
		w.eb.Notify(ev.Data)

	case worker.Failed:
		log.Printf("Warning: %s: %s", ev.Format(), ev.Err)
		w.StatusLine.sendError(ev.Message)

	default:
		log.Printf("Error: %d unknown worker event type", ev.Type)
	}
}

func NewWindow(cfg *config.Config) *Window {

	worker := worker.NewWorker()
	eb := event.NewEventBus()
	cb := command.NewCommandBus()
	srv := app.NewService(worker, eb, cb)

	w := &Window{
		cfg:          cfg,
		eb:           eb,
		cb:           cb,
		srv:          srv,
		Body:         &Body{},
		StatusLine:   newStatusLine(),
		DBPath:       newDBPath(cfg),
		Table:        table.NewTable(cfg, worker, cb, eb),
		UniqueGenres: newUniqueGenres(eb),
		NoData:       newNoDataBody(eb),
		ShowError:    &ShowError{},
		Worker:       worker,
	}

	w.Body.Set(BodyNoData)

	setupStatusLineToEvents(w.StatusLine, eb)
	setupBodyToEvents(w.Body, eb)

	// Set Up Handlers
	w.Form = newBookForm(eb,
		func() { // OnUpdate
			book, err := w.Form.GetBookEntry()
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Waring:", err)
				return
			}

			cell, has := w.Table.Selected.Get()
			if !has {
				err := fmt.Errorf("updating book entry: nothing selected")
				log.Println("Error:", err)
				w.StatusLine.sendError(err.Error())
				return
			}

			book.ID = cell.ID
			cb.Dispatch(command.UpdateBookEntry{Book: *book})
		},

		func() { // OnCreate
			book, err := w.Form.GetBookEntry()
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Waring:", err)
				return
			}
			cb.Dispatch(command.CreateBookEntry{Book: *book})
		},
	)

	w.Controls = &TableControl{
		OnUnselect: func() {
			w.Table.Selected.Set(w.Table.Sheet.Version(), models.Cell{}, false)
		},
		OnEdit: func() {
			cell, has := w.Table.Selected.Get()
			if !has {
				err := fmt.Errorf("edit book entry: nothing selected")
				log.Println("Error:", err)
				w.StatusLine.sendError(err.Error())
				return
			}
			book, err := snapshot.Current.Load().GetBookEntryByID(cell.ID)
			if err != nil {
				err := fmt.Errorf("edit book entry: %w", err)
				log.Println(err)
				w.StatusLine.sendError(err.Error())
				return
			}
			w.Form.Set(book)
			w.Body.Set(BodyBookEdit)
		},
		OnCreate: func() {
			w.Body.Set(BodyBookCreate)
		},

		OnDelete: func() {
			cell, has := w.Table.Selected.Get()
			if !has {
				err := fmt.Errorf("delete book entry: nothing selected")
				log.Println("Error:", err)
				w.StatusLine.sendError(err.Error())
				return
			}
			cb.Dispatch(command.DeleteBookEntry{BookID: cell.ID})
		},
	}

	w.FileManage = &FileManage{

		// Opening and Creating Database
		CreateDatabase: func(path string, err error) {
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				w.ShowError.Show(err)
				return
			}
			if path == "" {
				return
			}
			if !strings.HasSuffix(path, ".db") &&
				!strings.HasSuffix(path, ".sqlite") &&
				!strings.HasSuffix(path, ".sqlite3") {
				path += ".db"
			}
			cb.Dispatch(command.CreateDatabase{Path: path})
		},

		OpenDatabase: func(path string, err error) {
			if err != nil {
				w.StatusLine.sendError(err.Error())
				w.ShowError.Show(err)
				log.Println("Error:", err)
				return
			}
			if path == "" {
				return
			}
			w.DBPath.Set(path)
			cb.Dispatch(command.OpenDatabase{Path: path})
		},

		// Importing and Exporting
		ImportFile: func(path string, err error) {
			if err != nil {
				w.StatusLine.sendError(err.Error())
				w.ShowError.Show(err)
				log.Println("Error:", err)
				return
			}
			if path == "" {
				return
			}

			cb.Dispatch(command.ImportFile{Path: path})
		},

		ExportFile: func(path string, err error) {
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				return
			}
			if path == "" {
				return
			}
			if !strings.HasSuffix(path, ".csv") {
				path += ".csv"
			}

			cb.Dispatch(command.ExportFile{Path: path})
		},
	}
	return w
}


func FirstLoad(w *Window) {
	w.cb.Dispatch(command.OpenDatabase{
		Path: w.cfg.DBFile,
	})
}


type TableControl struct {
	OnCreate   func()
	OnUnselect func()
	OnEdit     func()
	OnDelete   func()
}

type FileManage struct {
	OpenDatabase   func(path string, err error)
	CreateDatabase func(path string, err error)

	ImportFile func(path string, err error)
	ExportFile func(path string, err error)
}

// NOTE: I didn't want to be too depended on Fyne, so I wrap the file open and create functions for their file dialogs.

func WrapFyneFileOpen(fn func(string, error)) func(fyne.URIReadCloser, error) {
	return func(r fyne.URIReadCloser, err error) {
		var path string
		if r != nil {
			if e := r.Close(); e != nil {
				err = e
			}
			path = r.URI().Path()
		}
		fn(path, err)
	}
}

func WrapFyneFileCreate(fn func(string, error)) func(fyne.URIWriteCloser, error) {
	return func(w fyne.URIWriteCloser, err error) {
		var path string
		if w != nil {
			if e := w.Close(); e != nil {
				err = e
			}
			if e := os.Remove(w.URI().Path()); e != nil {
				err = e
			}
			path = w.URI().Path()
		}
		fn(path, err)
	}
}

func setupStatusLineToEvents(sl *StatusLine, eb *event.EventBus) {
	eb.Subscribe(event.CreatedBookEntry{}, func(v event.Event){
		e := v.(event.CreatedBookEntry)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendSuccess("Entry Added!")
		}
	})

	eb.Subscribe(event.UpdatedBookEntry{}, func(v event.Event){
		e := v.(event.UpdatedBookEntry)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendSuccess("Entry Updated!")
		}
	})
	eb.Subscribe(event.DeletedBookEntry{}, func(v event.Event){
		e := v.(event.DeletedBookEntry)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendSuccess("Entry Removed!")
		}
	})
	eb.Subscribe(event.ImportedFile{}, func(v event.Event){
		e := v.(event.ImportedFile)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendInfo(fmt.Sprintf("Imported %s", e.Path))
		}
	})
	eb.Subscribe(event.ExportedFile{}, func(v event.Event) {
		e := v.(event.ExportedFile)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendInfo(fmt.Sprintf("Exported to %s", e.Path))
		}
	})
	eb.Subscribe(event.OpenedDatabase{}, func(v event.Event) {
		e := v.(event.OpenedDatabase)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendInfo(fmt.Sprintf("Opened %s", e.Path))
		}
	})
	eb.Subscribe(event.CreatedDatabase{}, func(v event.Event) {
		e := v.(event.CreatedDatabase)
		if e.Failed {
			sl.sendError(e.Message)
		} else {
			sl.sendInfo(fmt.Sprintf("Created %s", e.Path))
		}
	})
}

func setupBodyToEvents(b *Body, eb *event.EventBus) {
	eb.Subscribe(event.UpdatedBookEntry{}, func(v event.Event) {
		e := v.(event.UpdatedBookEntry)
		// when updating an entry go back to table when completed successfully.
		if !e.Failed {
			b.Set(BodyTable)
		}
	})
	eb.Subscribe(event.OpenedDatabase{}, func(v event.Event) {
		e := v.(event.OpenedDatabase)
		if e.Failed && b.Value() != BodyTable {
			b.Set(BodyNoData)
		} else {
			b.Set(BodyTable)
		}
	})
	eb.Subscribe(event.CreatedDatabase{}, func(v event.Event) {
		e := v.(event.CreatedDatabase)
		if e.Failed && b.Value() != BodyTable {
			b.Set(BodyNoData)
		} else {
			b.Set(BodyTable)
		}
	})
}
