package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func NewMenu(w fyne.Window, vm *viewmodel.MenuVM) *fyne.Container {

	csvImportBtn := widget.NewButton("Import", func() {
		d := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			vm.ImportCSV(r, err)
		}, w)
		d.Resize(w.Canvas().Size())
		d.SetTitleText("Import CSV")
		d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		d.Show()
	})
	csvExportBtn := widget.NewButton("Export", func() {
		d := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
			var path string
			if w != nil {
				path = w.URI().Path()
			}
			vm.ExportCSV(w, path, err)
		}, w)
		d.Resize(w.Canvas().Size())
		d.SetTitleText("Export CSV")
		d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		d.Show()
	})

	openDBBtn := widget.NewButton("Open", func() {
		d := dialog.NewFileOpen(func(r fyne.URIReadCloser, err error) {
			var path string
			if r != nil {
				path = r.URI().Path()
				r.Close()
			}
			vm.OpenDatabase(path, err)
		}, w)
		d.Resize(w.Canvas().Size())
		d.SetTitleText("Open Database")
		d.SetFilter(storage.NewExtensionFileFilter([]string{".db", ".sqlite", ".sqlite3"}))
		d.Show()

	})
	saveDBBtn := widget.NewButton("Create", func() {
		d := dialog.NewFileSave(func(w fyne.URIWriteCloser, err error) {
			var path string
			if w != nil {
				path = w.URI().Path()
				w.Close()
			}
			vm.CreateDatabase(path, err)
		}, w)
		d.Resize(w.Canvas().Size())
		d.SetTitleText("Create Database")
		d.SetFilter(storage.NewExtensionFileFilter([]string{".db", ".sqlite", ".sqlite3"}))
		d.Show()
	})

	dbFileLbl := widget.NewLabel("")
	setDBLabel := func() {
		label, _ := vm.DBFile.Get()
		dbFileLbl.SetText(fmt.Sprintf("\"%s\"", label))
	}

	setDBLabel()
	vm.DBFile.AddListener(binding.NewDataListener(func() {
		setDBLabel()
	}))

	newHeaderLabel := func(text string) *widget.Label {
		lbl := widget.NewLabel(text)
		lbl.TextStyle = fyne.TextStyle{
			Bold:      true,
			Underline: true,
		}
		return lbl
	}

	menu := container.NewVBox(
		newHeaderLabel("Database"),
		dbFileLbl,
		container.NewGridWithColumns(3,
			openDBBtn,
			saveDBBtn,
		),
		newHeaderLabel("CSV"),
		container.NewGridWithColumns(3,
			csvImportBtn,
			csvExportBtn,
		),
	)
	return menu
}
