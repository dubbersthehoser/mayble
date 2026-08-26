package viewmodel

import (
	"fmt"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"
	"github.com/dubbersthehoser/mayble/internal/viewmodel/table"
)

type Window struct {
	cfg *config.Config

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
}

func NewWindow(cfg *config.Config) *Window {

	srv := app.NewService(cfg)

	tbl := table.NewTable(cfg, srv.GetAllBooks)

	w := &Window{
		cfg:          cfg,
		Body:         &Body{},
		StatusLine:   newStatusLine(),
		DBPath:       newDBPath(cfg),
		Table:        tbl,
		UniqueGenres: newUniqueGenres(srv),
		NoData:       &NoDataBody{},
		ShowError:    &ShowError{},
	}

	//
	// Set Up Handlers
	//

	w.Form = newBookForm(
		func() {
			book, err := w.Form.GetBookEntry()
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Waring:", err)
				return
			}

			row := w.Table.Selected.Get().Row
			id, ok := w.Table.Sheet.RowToID(row)
			if !ok {
				log.Printf("Error: row '%d' not found in ids", row)
				return
			}
			book.ID = id
			if err := srv.UpdateBook(book); err != nil {
				log.Println("Error:", err)
				w.StatusLine.sendError(err.Error())
				return
			}
			w.StatusLine.sendSuccess("Updated!")
			w.Form.Reset()
			w.Body.Set(BodyTable)
		},

		func() {
			book, err := w.Form.GetBookEntry()
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Waring:", err)
				return
			}

			if _, err := srv.CreateBook(book); err != nil {
				log.Println("Error:", err)
				w.StatusLine.sendError(err.Error())
				return
			}
			w.StatusLine.sendSuccess("Created!")
			w.Form.Reset()
		},
	)

	w.Controls = &TableControl{
		OnUnselect: func() {
			w.Table.Selected.Unselect()
		},
		OnEdit: func() {
			row := w.Table.Selected.Get().Row
			id, _ := w.Table.Sheet.RowToID(row)
			book, err := srv.GetBookByID(id)
			if err != nil {
				log.Println("Error:", err)
				w.StatusLine.sendError(err.Error())
				return
			}
			w.Form.Set(&book)
			w.Body.Set(BodyBookEdit)
		},
		OnCreate: func() {
			w.Body.Set(BodyBookCreate)
		},
		OnDelete: func() {
			row := w.Table.Selected.Get().Row
			id, _ := w.Table.Sheet.RowToID(row)
			err := srv.DeleteBook(id)
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
			}
		},
	}

	w.FileManage = &FileManage{
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
			if err := srv.CreateDatabase(path); err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				w.NoData.SetDataErr(w.DBPath.Get(), err)
				return
			}
			w.DBPath.Set(path)
			w.StatusLine.sendInfo(fmt.Sprintf("created: %s", w.DBPath.Get()))
			w.Body.Set(BodyTable)
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
			if err := srv.OpenDatabase(path); err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				w.NoData.SetDataErr(w.DBPath.Get(), err)
				return
			}
			w.StatusLine.sendInfo(fmt.Sprintf("opened: %s", w.DBPath.Get()))
			w.Body.Set(BodyTable)
		},

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

			if err := srv.ImportFile(path); err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				w.ShowError.Show(err)
				return
			}
			w.StatusLine.sendSuccess("Imported: " + path)
			w.Table.Sheet.Load()
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

			if err := srv.ExportFile(path); err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				return
			}
			w.StatusLine.sendSuccess("Exported: " + path)
		},
	}

	// Allow table sheet to load data from database when there is a change.
	// Be it from regular CRUD operations, and the opening, and creation of a database file.
	srv.AddListener(func() {
		w.Table.Sheet.Load()
	})

	// start with table view at start up.
	w.Body.Set(BodyTable)

	// only when body can change to NoData is with NoData calls.
	w.NoData.AddListener(func() {
		// The only place this function should be called durring run time.
		w.Body.Set(BodyNoData)
	})

	// try and load database at first start.
	if err := srv.OpenDatabase(w.cfg.DBFile); err != nil {
		log.Println("Warning:", err)
		if w.cfg.DBFile == "" {
			w.NoData.SetNoDB()
		} else {
			w.StatusLine.sendError(err.Error())
			w.NoData.SetDataErr(w.DBPath.Get(), err)
		}
	}

	if err := tbl.Sheet.Load(); err != nil {
		log.Println("error:", err)
	}

	return w
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

// Note: I didn't want to be too depended on Fyne, so I wrap the file open and create functions for their file dialogs.

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
