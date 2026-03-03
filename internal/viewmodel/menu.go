package viewmodel

import (
	"strings"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/dubbersthehoser/mayble/internal/bus"
	"github.com/dubbersthehoser/mayble/internal/csv"
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
	displayErr := func(err error) {
		c.vms.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: err,
		})
	}
	if err != nil {
		displayErr(err)
		return
	}

	books, err := csv.Import(r)
	if err != nil {
		displayErr(err)
	}

	for _, book := range books {
		err = c.vms.app.bookCreator.CreateBook(&book)
		if err != nil {
			displayErr(err)
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
	
	displayErr := func(err error) {
		c.vms.bus.Notify(bus.Event{
			Name: msgUserError,
			Data: err,
		})
	}

	if err != nil {
		displayErr(err)
		return
	}

	filepath := wURI.URI().Path()

	var w io.Writer = wURI

	if !strings.HasSuffix(filepath, ".csv") {
		filepath += ".csv"
		wURI.Close()
		w, err = os.Create(filepath)
		if err != nil {
			displayErr(err)
			return
		}
	}

	books, err := c.vms.app.bookRetriever.GetAllBooks(repo.Loaned|repo.Read)
	if err != nil {
		displayErr(err)
		return
	}

	err = csv.Export(w, books)
	if err != nil {
		displayErr(err)
		return
	}

	c.vms.bus.Notify(bus.Event{
		Name: msgUserSuccess,
		Data: "Exported!",
	})
}



func (c *MenuVM) OpenDatabase(r fyne.URIReadCloser, err error) {
	println(err)
}

func (c *MenuVM) CreateDatabase(r fyne.URIWriteCloser, err error) {
	println(err)
}



