package viewmodel

import (
	"strings"
	"io"
	"os"
	"fmt"

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

func (c *MenuVM) ImportCSV(r io.ReadCloser, err error) {
	if err != nil {
		displayError(c.bus, err)
		return
	}
	if r == nil {
		return
	}
	defer r.Close()
	books, err := csv.Import(r)
	if err != nil {
		displayError(c.bus, err)
	}

	for _, book := range books {
		_, err = c.app.bookCreator.CreateBook(&book)
		if err != nil {
			displayError(c.bus, err)
			return
		}
	}
	c.bus.Notify(bus.Event{
		Name: msgDataChanged,
		Data: nil,
	})
	c.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Imported!",
	})
}

func (c *MenuVM) ExportCSV(w io.WriteCloser, filepath string, err error) {
	if err != nil {
		displayError(c.bus, err)
		return
	}
	if w == nil {
		return
	}

	if !strings.HasSuffix(filepath, ".csv") {
		filepath += ".csv"
		w.Close()
		w, err = os.Create(filepath)
		if err != nil {
			displayError(c.bus, err)
			return
		}
	}
	defer w.Close()

	books, err := c.app.bookRetriever.GetAllBooks(repo.Book)
	if err != nil {
		displayError(c.bus, err)
		return
	}

	err = csv.Export(w, books)
	if err != nil {
		displayError(c.bus, err)
		return
	}

	c.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Exported!",
	})
}



func (c *MenuVM) OpenDatabase(path string, err error) {
	if path == "" {
		return
	}
	if err != nil {
		displayError(c.bus, err)
		return
	}

	db, err := database.Open(path)
	if err != nil {
		displayError(c.bus, err)
		return
	}

	err = c.app.changeDB(db)
	if err != nil {
		displayError(c.bus, err)
		return
	}
	c.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
	c.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: fmt.Sprintf("opened: '%s'", path),
	})

	_ = c.DBFile.Set(path)
}

func (c *MenuVM) CreateDatabase(path string, err error) {
	if err != nil {
		displayError(c.bus, err)
		return
	}

	if path == "" {
		return
	}
	_ = c.DBFile.Set(path)

	if !strings.HasSuffix(path, ".db") &&
	   !strings.HasSuffix(path, ".sqlite") &&
	   !strings.HasSuffix(path, ".sqlite3") {
		path += ".db"
	}

	db, err := database.Open(path)

	if err != nil {
		displayError(c.bus, err)
		return
	}

	err = c.app.changeDB(db)
	if err != nil {
		displayError(c.bus, err)
		return
	}

	c.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: fmt.Sprintf("created: '%s'", path),
	})

	c.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
}

func displayError(b *bus.Bus, err error) {
	b.Notify(bus.Event{
		Name: msgUserError,
		Data: err,
	})
}
