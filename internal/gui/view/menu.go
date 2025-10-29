package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"

	_"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/gui/controller/porting"
)


func (f *FunkView) ShowMenu() {


	/*
		Import Button
	*/
	openImportFile := func(uri fyne.URIReadCloser, err error) {
		if err != nil {
			f.displayError(err)
			return
		}
		if uri == nil {
			return
		}	
		fmt.Println("opening import:", uri.URI().Path())
		impoter, err := porting.GetImporterByFilePath(uri.URI().Path())
		if err != nil {
			f.displayError(err)
			return
		}
		books, err := impoter.ImportBooks(uri)
		if err != nil {
			f.displayError(err)
			return
		}
		_ = books

	}
	
	importBtn := widget.NewButton("Import CSV", func() {
		d := dialog.NewFileOpen(openImportFile, f.window)
		filter := storage.NewExtensionFileFilter([]string{".csv"})
		d.SetFilter(filter)
		size := fyne.NewSize(800, 800)
		d.Resize(size)
		d.Show()
	})


	/*
		Export Button
	*/
	saveImportFile := func(uri fyne.URIWriteCloser, err error) {
		if err != nil {
			f.displayError(err)
			return
		}
		if uri == nil {
			return
		}
		path := uri.URI().Path()
		fmt.Println("saving import:", path)
		exporter, err := porting.GetExporterByFilePath(path)
		if err != nil {
			f.displayError(err)
			return
		}
		books, err := f.controller.Core.GetAllBookLoans()
		if err != nil {
			f.displayError(err)
			return
		}
		
		err = exporter.ExportBooks(uri, books)
		if err != nil {
			f.displayError(err)
			return 
		}
	}
	exportBtn := widget.NewButton("Export CSV", func() {
		d := dialog.NewFileSave(saveImportFile, f.window)
		d.SetFileName("book-and-loans.csv")
		filter := storage.NewExtensionFileFilter([]string{".csv"})
		d.SetFilter(filter)

		size := fyne.NewSize(800, 800)
		d.Resize(size)
		d.Show()
	})

	obj := container.New(layout.NewVBoxLayout(), importBtn, exportBtn)

	d := dialog.NewCustom("Menu", "Close", obj, f.window)
	d.Show()
}
