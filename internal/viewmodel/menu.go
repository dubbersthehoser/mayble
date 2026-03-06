package viewmodel

import (
	"strings"
	"io"
	"os"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/csv"
	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/database"
	repo "github.com/dubbersthehoser/mayble/internal/repository"

)

type MenuVM struct {
	DBFile binding.String
	app *appService
	bus *bus.Bus
}

func NewMenuVM(b *bus.Bus, app *appService, dbFile binding.String) *MenuVM {
	m := &MenuVM{
		DBFile: dbFile,
		bus: b,
		app: app,
	}
	return m
} 

func (c *MenuVM) ImportCSV(r fyne.URIReadCloser, err error) {
	if err != nil {
		displayError(c.bus, err)
		return
	}
	if r == nil {
		return
	}
	defer r.Close()
	importCSV(r, c.bus, c.app.bookCreator)

}

func (c *MenuVM) ExportCSV(w fyne.URIWriteCloser, err error) {
	if err != nil {
		displayError(c.bus, err)
		return
	}
	if w == nil {
		return
	}
	exportCSV(w, w.URI().Path(), c.bus, c.app.bookRetriever)
}



func (c *MenuVM) OpenDatabase(r fyne.URIReadCloser, err error) {
	if r == nil {
		return
	}
	if err != nil {
		displayError(c.bus, err)
		return
	}
	r.Close()
	openDatabase(r.URI().Path(), c.app, c.bus)
	_ = c.DBFile.Set(r.URI().Path())
}

func (c *MenuVM) CreateDatabase(w fyne.URIWriteCloser, err error) {
	if err != nil {
		displayError(c.bus, err)
		return
	}

	if w == nil {
		return
	}

	w.Close()
	filepath := w.URI().Path()
	_ = c.DBFile.Set(filepath)
}

func displayError(b *bus.Bus, err error) {
	b.Notify(bus.Event{
		Name: msgUserError,
		Data: err,
	})
}



func importCSV(r io.Reader, b *bus.Bus, bc repo.BookCreator) {
	books, err := csv.Import(r)
	if err != nil {
		displayError(b, err)
	}

	for _, book := range books {
		_, err = bc.CreateBook(&book)
		if err != nil {
			displayError(b, err)
			return
		}
	}
	b.Notify(bus.Event{
		Name: msgDataChanged,
		Data: nil,
	})
	b.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Imported!",
	})
}

func exportCSV(w io.WriteCloser, filepath string, b *bus.Bus, br repo.BookRetriever) {
	var err error
	if !strings.HasSuffix(filepath, ".csv") {
		filepath += ".csv"
		w.Close()
		w, err = os.Create(filepath)
		if err != nil {
			displayError(b, err)
			return
		}
	}
	defer w.Close()

	books, err := br.GetAllBooks(repo.Book)
	if err != nil {
		displayError(b, err)
		return
	}

	err = csv.Export(w, books)
	if err != nil {
		displayError(b, err)
		return
	}

	b.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Exported!",
	})
}

func openDatabase(path string, as *appService, b *bus.Bus) {
	db, err := database.Open(path)
	if err != nil {
		displayError(b, err)
		return
	}

	err = as.changeDB(db)
	if err != nil {
		displayError(b, err)
		return
	}
	b.Notify(bus.Event{
		Name: msgDataChanged,
	})
	b.Notify(bus.Event{
		Name: msgUserInfo,
		Data: fmt.Sprintf("opened: '%s'", path),
	})
}

func createDatabase(path string, as *appService, b *bus.Bus) {
	if !strings.HasSuffix(path, ".db") &&
	   !strings.HasSuffix(path, ".sqlite") &&
	   !strings.HasSuffix(path, ".sqlite3") {
		path += ".db"
	}

	db, err := database.Open(path)

	if err != nil {
		displayError(b, err)
		return
	}

	err = as.changeDB(db)
	if err != nil {
		displayError(b, err)
		return
	}

	b.Notify(bus.Event{
		Name: msgUserInfo,
		Data: fmt.Sprintf("created: '%s'", path),
	})

	b.Notify(bus.Event{
		Name: msgDataChanged,
	})
	
}
