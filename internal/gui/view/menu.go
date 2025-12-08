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

func ShowMenu(f *FunkView) {

	importBtn := NewImportButton(f.window, f.broker)
	exportBtn := NewExportButton(f.window, f.broker)

	obj := container.New(layout.NewVBoxLayout(), importBtn, exportBtn)

	d := dialog.NewCustom("Menu", "Close", obj, f.window)
	size := getDialogSize(d.MinSize())
	d.Resize(size)
	d.Show()
}

// ReadIOPorting a wrapper for porting.NamedReadCloser.
type ReadIOPorting struct {
	uri fyne.URIReadCloser
}
func (r *ReadIOPorting) Read(b []byte) (int, error) {
	return r.uri.Read(b)
}
func (r *ReadIOPorting) Close() error {
	return r.uri.Close()
}
func (r *ReadIOPorting) Name() string {
	return r.uri.URI().Name()
}

// ReadIOPorting a wrapper for porting.NamedWriteCloser.
type WriteIOPorting struct {
	uri fyne.URIWriteCloser
}
func (w *WriteIOPorting) Write(b []byte) (int, error) {
	return w.uri.Write(b)
}
func (w *WriteIOPorting) Close() error {
	return w.uri.Close()
}
func (w *WriteIOPorting) Name() string {
	return w.uri.URI().Name()
}


func NewFileReadDialog(w fyne.Window, handler func(uri fyne.URIReadCloser, err error)) *dialog.FileDialog {
	d := dialog.NewFileOpen(handler, w)
	size := fyne.NewSize(800, 800)
	d.Resize(size)
	return d
}

func NewFileWriteDialog(w fyne.Window, handler func(uri fyne.URIWriteCloser, err error)) *dialog.FileDialog {
	d := dialog.NewFileSave(handler, w)
	size := fyne.NewSize(800, 800)
	d.Resize(size)
	return d
}


type ImportButton struct {
	widget.Button
	window fyne.Window
	broker *emiter.Broker
}

type ExportButton struct {
	widget.Button
	window fyne.Window
	broker *emiter.Broker
}

func NewExportButton(w fyne.Window, b *emiter.Broker) *ExportButton {
	eb := &ExportButton{
		broker: b,
		window: w,
	}
	eb.ExtendBaseWidget(eb)
	eb.SetText("Export")
	eb.OnTapped = func() {
		d := NewFileWriteDialog(eb.window, eb.OpenFile)
		filter := storage.NewExtensionFileFilter([]string{".csv"})
		d.SetFilter(filter)
		d.Show()
	}
	return eb
}

func NewImportButton(w fyne.Window, b *emiter.Broker) *ImportButton {
	ib := &ImportButton{
		broker: b,
		window: w,
	}
	ib.ExtendBaseWidget(ib)
	ib.SetText("Import")
	ib.OnTapped = func() {
		d := NewFileReadDialog(ib.window, ib.OpenFile)
		filter := storage.NewExtensionFileFilter([]string{".csv"})
		d.SetFilter(filter)
		d.Show()
	}
	return ib
}


func (eb *ExportButton) OpenFile(uri fyne.URIWriteCloser, err error) {
	if err != nil {
		eb.broker.Notify(emiter.Event{
			Name: gui.EventDisplayErr,
			Data: err,
		})
		return 
	}

	if uri == nil {
		return
	}
	pIO := &WriteIOPorting{uri: uri}

	eb.broker.Notify(emiter.Event{
		Name: gui.EventDocumentExport,
		Data: porting.NamedWriteCloser(pIO),
	})
}

func (cb *ImportButton) OpenFile(uri fyne.URIReadCloser, err error) {
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
	pIO := &ReadIOPorting{uri: uri}

	cb.broker.Notify(emiter.Event{
		Name: gui.EventDocumentImport,
		Data: porting.NamedReadCloser(pIO),
	})
}

