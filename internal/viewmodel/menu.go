package viewmodel

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2"

	"github.com/dubbersthehoser/mayble/internal/bus"
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type Menu struct {
	DBFile binding.String
	dbOpener    databaseOpener
	fileHandler repo.CSVHandler
	bus    *bus.Bus
}

func NewMenu(b *bus.Bus, fh repo.CSVHandler, dbOpener databaseOpener, dbFile binding.String) *Menu {
	m := &Menu{
		DBFile: dbFile,
		bus:    b,
		dbOpener:    dbOpener,
		fileHandler: fh,
	}
	return m
}

func (c *Menu) ImportCSV(r fyne.URIReadCloser, err error) {
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
		Name: msgUserSuccess,
		Data: "Imported!",
	})
}

func (c *Menu) ExportCSV(w fyne.URIWriteCloser, err error) {
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

	err = os.Remove(filepath)
	if err != nil {
		c.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: err.Error(),
		})
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

func (c *Menu) OpenDatabase(path string, err error) {
	if path == "" {
		return
	}

	if err != nil {
		displayError(c.bus, err)
		return
	}

	err = c.dbOpener.OpenDatabase(path)
	if err != nil {
		displayError(c.bus, err)
		return
	}

	c.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: fmt.Sprintf("opened: '%s'", path),
	})
	_ = c.DBFile.Set(path)
}

func (c *Menu) CreateDatabase(path string, err error) {
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

	err = c.dbOpener.OpenDatabase(path)
	if err != nil {
		displayError(c.bus, err)
		return
	}

	c.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: fmt.Sprintf("created: '%s'", path),
	})
}

func displayError(b *bus.Bus, err error) {
	b.Notify(bus.Event{
		Name: msgUserError,
		Data: err.Error(),
	})
}
