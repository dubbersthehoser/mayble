package viewmodel

import (
	"strings"
	"io"
	"os"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/csv"
	"github.com/dubbersthehoser/mayble/internal/database"
	repo "github.com/dubbersthehoser/mayble/internal/repository"

)

type MenuVM struct {
	DBFile binding.String
	vms *vmService
	
}

func NewMenuVM(vms *vmService, dbFile binding.String) *MenuVM {
	m := &MenuVM{
		DBFile: dbFile,
		vms: vms,
	}
	return m
} 

func (c *MenuVM) ImportCSV(r fyne.URIReadCloser, err error) {
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}

	if r == nil {
		return
	}

	books, err := csv.Import(r)
	if err != nil {
		displayError(c.vms.bus, err)
	}

	for _, book := range books {
		err = c.vms.app.bookCreator.CreateBook(&book)
		if err != nil {
			displayError(c.vms.bus, err)
			return
		}
	}
	c.vms.bus.Notify(bus.Event{
		Name: msgDataChanged,
		Data: nil,
	})
	c.vms.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Imported!",
	})
}

func (c *MenuVM) ExportCSV(wURI fyne.URIWriteCloser, err error) {

	
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}

	if wURI == nil {
		return
	}

	filepath := wURI.URI().Path()

	var w io.WriteCloser = wURI

	if !strings.HasSuffix(filepath, ".csv") {
		filepath += ".csv"
		wURI.Close()
		w, err = os.Create(filepath)
		if err != nil {
			displayError(c.vms.bus, err)
			return
		}
	}
	defer w.Close()

	books, err := c.vms.app.bookRetriever.GetAllBooks(repo.Loaned|repo.Read)
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}

	err = csv.Export(w, books)
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}

	c.vms.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Exported!",
	})
}

func (c *MenuVM) OpenDatabase(r fyne.URIReadCloser, err error) {
	if r == nil {
		return
	}
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}
	r.Close()

	filepath := r.URI().Path()

	db, err := database.Open(filepath)
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}

	err = c.vms.app.dbs.SetDB(db)
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}
	c.vms.app.setDB(db)
	_ = c.DBFile.Set(filepath)
	c.vms.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
	c.vms.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: fmt.Sprintf("opened: '%s'", filepath),
	})
}

func (c *MenuVM) CreateDatabase(w fyne.URIWriteCloser, err error) {
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}

	if w == nil {
		return
	}

	w.Close()

	filepath := w.URI().Path()

	if !strings.HasSuffix(filepath, ".db") &&
	   !strings.HasSuffix(filepath, ".sqlite") &&
	   !strings.HasSuffix(filepath, ".sqlite3") {
		filepath += ".db"
	}

	db, err := database.Open(filepath)

	if err != nil {
		displayError(c.vms.bus, err)
		return
	}


	err = c.vms.app.dbs.SetDB(db)
	if err != nil {
		displayError(c.vms.bus, err)
		return
	}

	c.vms.app.setDB(db)

	_ = c.DBFile.Set(filepath)

	c.vms.bus.Notify(bus.Event{
		Name: msgUserInfo,
		Data: fmt.Sprintf("created: '%s'", filepath),
	})

	c.vms.bus.Notify(bus.Event{
		Name: msgDataChanged,
	})
}

func displayError(b *bus.Bus, err error) {
	b.Notify(bus.Event{
		Name: msgUserError,
		Data: err,
	})
}
