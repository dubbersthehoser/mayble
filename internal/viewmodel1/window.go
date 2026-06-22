package viewmodel

import (
	"log"
	"strings"
	"os"
	"fmt"

	"fyne.io/fyne/v2"

	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/config"
)

const (
	BodyNoData int = iota
	BodyTable
	BodyBookEdit
	BodyBookCreate
)

type Window struct {
	cfg            *config.Config
	service        *app.Service

	Body           *Body
	StatusLine     *StatusLine
	Controls       *TableControl
	FileManage     *FileManage
	DBPath         *DBPath
	UniqueGenres   *UniqueGenres
	Selected       *EntrySelected
	ColumnSettings *ColumnSettings
	DataTable      *DataTable
	Sorting        *SortingTable
	Searching      *Searching
	Search         func(string)
	Form           *BookForm
}

func NewWindow(cfg *config.Config) *Window {
	serv := app.NewService(cfg)
	w := &Window{
		cfg: cfg,
		Body: &Body{},
		StatusLine: newStatusLine(),
		ColumnSettings: newColumnSettings(cfg),
		DBPath: newDBPath(cfg),
		DataTable: newDataTable(cfg, serv),
		Sorting: newSortingTable(cfg),
		Searching: &Searching{},
		Selected: newEntrySelected(),
		UniqueGenres: newUniqueGenres(serv),
	}

	w.Form = newBookForm(
		func() {
			book, err := w.Form.GetBookEntry()
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Waring:", err)
				return
			}

			row, _ := w.Selected.Get()
			id, ok := w.DataTable.rowToID[row]
			if !ok {
				log.Printf("Error: row '%d' not found in ids", row)
				return
			}
			book.ID = id
			if err := serv.UpdateBook(book); err != nil {
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

			if _, err := serv.CreateBook(book); err != nil {
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
			w.Selected.Unselect()
		},
		OnEdit: func() {
			row, _ := w.Selected.Get()
			id := w.DataTable.rowToID[row]
			book, err := serv.GetBookByID(id)
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
			id := w.DataTable.rowToID[w.Selected.row]
			err := serv.DeleteBook(id)
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
			w.DBPath.Set(path)
			w.Body.Set(BodyTable)
		},

		OpenDatabase: func(path string, err error) {
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				return
			}
			if path == "" {
				return
			}
			w.DBPath.Set(path)
			w.Body.Set(BodyTable)
		},

		ImportFile: func(path string, err error) {
			if err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				return
			}
			if path == "" {
				return
			}

			if err := serv.ImportFile(path); err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				return
			}
			w.StatusLine.sendSuccess("Imported: " + path)
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

			if err := serv.ExportFile(path); err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				return
			}

			w.StatusLine.sendSuccess("Exported: " + path)
		},
	}

	w.Search = func(s string) {
		if s == "" {
			w.Selected.Unselect()
			return
		}
		w.Searching.search(w.DataTable.data, s)
	}

	w.Searching.AddListener(func() {
		if w.Searching.Has() {
			w.Selected.Select(
				w.Searching.Selected(),
			)
		} else {
			w.Selected.Unselect()
		}
	})

	w.Body.Set(BodyTable)

	if err := serv.LoadDatabase(); err != nil {
		w.StatusLine.sendError(err.Error())
		log.Println("Warning:", err)
		w.Body.Set(BodyNoData)
	}

	w.DBPath.AddListener(func(){
		if err := serv.LoadDatabase(); err != nil {
			w.StatusLine.sendError(err.Error())
			log.Println("Error:", err)
			w.Body.Set(BodyNoData)
		}
		w.StatusLine.sendInfo(fmt.Sprintf("opened: %s", w.DBPath.Get()))
	})

	w.ColumnSettings.AddListener(func() {
		w.DataTable.load()
		w.Selected.Unselect()
	})

	w.Sorting.AddListener(func() {
		w.DataTable.load()
	})

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
