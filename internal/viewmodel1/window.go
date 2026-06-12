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
	Selected       *EntrySelected
	ColumnSettings *ColumnSettings
	DataTable      *DataTable
	Sorting        *SortingTable
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
	}

	w.Controls = &TableControl{
		OnUnselect: func() {
			w.Selected.Unselect()
		},
		OnEdit: func() {
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

			if err := serv.ImportFile(path); err != nil {
				w.StatusLine.sendError(err.Error())
				log.Println("Error:", err)
				return
			}

			w.StatusLine.sendSuccess("Exported: " + path)
		},
	}

	if err := serv.LoadDatabase(); err != nil {
		w.StatusLine.sendError(err.Error())
		w.Body.Set(BodyNoData)
	}

	w.DBPath.AddListener(func(){
		if err := serv.LoadDatabase(); err != nil {
			w.StatusLine.sendError(err.Error())
			log.Println("Error:", err)
			w.Body.Set(BodyNoData)
		}
		w.StatusLine.sendInfo(fmt.Sprintf("opened:", w.DBPath.Get()))
	})

	w.Sorting.AddListener(func() {
		w.DataTable.load()
	})

	w.Selected = &EntrySelected{}
	return w
}

type DBPath struct {
	cfg *config.Config
	l []func()
}
func newDBPath(cfg *config.Config) *DBPath{
	dbp := &DBPath{
		cfg: cfg,
	}
	return dbp
}
func (p *DBPath) Get() string {
	return p.cfg.DBFile
}
func (p *DBPath) Set(s string) {
	p.cfg.DBFile = s
}
func (p *DBPath) AddListener(fn func()) {
	if p.l == nil {
		p.l = make([]func(), 0)
	}
	p.l = append(p.l, fn)
}

type Body struct {
	value int
	l []func()
}

func (b *Body) Value() int {
	return b.value
}

func (b *Body) Set(v int) {
	b.value = v
	b.notify()
}

func (b *Body) notify() {
	for _, fn := range b.l {
		fn()
	}
}

func (b *Body) AddListener(fn func()) {
	if b.l == nil {
		b.l = make([]func(), 0)
	}
	b.l = append(b.l, fn)
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
