package viewmodel

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2"

	"github.com/dubbersthehoser/mayble/internal/bus"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type MenuVM struct {
	DBFile binding.String
	dbOpener    DatabaseOpener
	fileHandler repo.CSVHandler
	bus    *bus.Bus
}

func NewMenuVM(b *bus.Bus, fh repo.CSVHandler, dbOpener DatabaseOpener, dbFile binding.String) *MenuVM {
	m := &MenuVM{
		DBFile: dbFile,
		bus:    b,
		dbOpener:    dbOpener,
		fileHandler: fh,
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

	if err := r.Close(); err != nil {
		displayError(c.bus, err)
		return
	}
	
	if err := c.fileHandler.ImportFile(r.URI().Path()); err != nil {
		displayError(c.bus, err)
		return
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

func (c *MenuVM) ExportCSV(w fyne.URIWriteCloser, err error) {
	if err != nil {
		displayError(c.bus, err)
		return
	}
	if w == nil {
		return
	}

	filepath := w.URI().Path()

	if err := w.Close(); err != nil {
		displayError(c.bus, err)
		return
	}

	if !strings.HasSuffix(filepath, ".csv") {
		filepath += ".csv"
	}

	if err := c.fileHandler.ExportFile(filepath); err != nil {
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

	err = c.dbOpener.OpenDB(path)
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

	err = c.dbOpener.OpenDB(path)
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
		Data: err.Error(),
	})
}
