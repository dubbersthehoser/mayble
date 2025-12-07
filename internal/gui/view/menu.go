package view

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"

	"github.com/dubbersthehoser/mayble/internal/porting"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/gui"

)


func NewFileDialog(w fyne.Window, handler func(uri fyne.URIReadCloser, err error)) *dialog.FileDialog {
	d := dialog.NewFileOpen(handler, w)
	size := fyne.NewSize(800, 800)
	d.Resize(size)
	return d
}

type CSVImportButton struct {
	widget.Button
	window fyne.Window
	broker *emiter.Broker
}

func NewCSVImportButton(w fyne.Window, b *emiter.Broker) *CSVImportButton {
	cb := &CSVImportButton{
		broker: b,
		window: w,
	}
	cb.ExtendBaseWidget(cb)
	cb.SetText("Import CSV")
	cb.OnTapped = func() {
		d := NewFileDialog(cb.window, cb.Import)
		filter := storage.NewExtensionFileFilter([]string{".csv"})
		d.SetFilter(filter)
		d.Show()
	}
	return cb
}

func (cb *CSVImportButton) Import(uri fyne.URIReadCloser, err error) {
	if err != nil {
		cb.broker.Notify(emiter.Event{
			Name: gui.EventDisplayErr,
			Data: err,
		})
		return 
	}

	if uri == nil { 
		return
	}

	importer, err := porting.GetBookLoanPorterByFilePath(uri.URI().Path())
	if err != nil {
		cb.broker.Notify(emiter.Event{
			Name: gui.EventDisplayErr,
			Data: err,
		})
		return
	}

	books, err := importer.ImportBookLoans(uri)
	if err != nil {
		cb.broker.Notify(emiter.Event{
			Name: gui.EventDisplayErr,
			Data: err,
		})
	}

	cb.broker.Notify(emiter.Event{
		Name: gui.EventDocumentImport,
		Data: books,
	})
}

type CSVExportButton struct {
	
}

func ShowMenu(f *FunkView) {

	csvImportBtn := NewCSVImportButton(f.window, f.broker)
	//csvExportBtn := NewCSVExportButton(f.broker)


	/*
		Export Button
	*/
	//saveImportFile := func(uri fyne.URIWriteCloser, err error) {
	//	if err != nil {
	//		f.displayError(err)
	//		return
	//	}
	//	if uri == nil {
	//		return
	//	}
	//	path := uri.URI().Path()
	//	exporter, err := porting.GetExporterByFilePath(path)
	//	if err != nil {
	//		f.displayError(err)
	//		return
	//	}
	//	books, err := f.controller.App.GetBookLoans()
	//	if err != nil {
	//		f.displayError(err)
	//		return
	//	}
	//	
	//	err = exporter.ExportBooks(uri, books)
	//	if err != nil {
	//		f.displayError(err)
	//		return 
	//	}
	//}

	//exportBtn := widget.NewButton("Export CSV", func() {
	//	d := dialog.NewFileSave(saveImportFile, f.window)
	//	d.SetFileName("book-and-loans.csv")
	//	filter := storage.NewExtensionFileFilter([]string{".csv"})
	//	d.SetFilter(filter)

	//	size := fyne.NewSize(800, 800)
	//	d.Resize(size)
	//	d.Show()
	//})

	obj := container.New(layout.NewVBoxLayout(), csvImportBtn)

	d := dialog.NewCustom("Menu", "Close", obj, f.window)
	size := getDialogSize(d.MinSize())
	d.Resize(size)
	d.Show()
}
