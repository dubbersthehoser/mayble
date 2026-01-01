package view

import (
	"path/filepath"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"

	"github.com/dubbersthehoser/mayble/internal/porting"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/gui"
	_"github.com/dubbersthehoser/mayble/internal/config"
)


func ShowMenu(f *FunkView) {

	openFile := NewOpenFileButton(f.window, f.broker)
	newFile := NewNewFileButton(f.window, f.broker)
	
	current := widget.NewLabel(f.controller.Config.DBFile)

	importBtn := NewImportButton(f.window, f.broker)
	exportBtn := NewExportButton(f.window, f.broker)

	openFile.Alignment = widget.ButtonAlignLeading
	newFile.Alignment = widget.ButtonAlignLeading
	importBtn.Alignment = widget.ButtonAlignLeading
	exportBtn.Alignment = widget.ButtonAlignLeading
	
	cancelBtn := NewButtonWithKeyed("Close")

	body := container.New(layout.NewVBoxLayout(), current, openFile, newFile, importBtn, exportBtn, cancelBtn)

	//d := dialog.NewCustom("Database", "Close", body, f.window)
	d := dialog.NewCustomWithoutButtons("Database", body, f.window)
	size := getDialogSize(d.MinSize())
	d.Resize(size)

	cancelBtn.OnTapped = func() {
		d.Dismiss()
	}

	var listnerID int

	listnerID = f.broker.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			d.Dismiss()
			f.broker.Unsubscribe(listnerID, gui.EventDocumentNew)
		},
	},
		gui.EventDocumentNew,
	)
	d.Show()
}

type OpenFileButton struct {
	widget.Button

	window fyne.Window
	broker *emiter.Broker
}
func NewOpenFileButton(w fyne.Window, b *emiter.Broker) *OpenFileButton {
	pf := &OpenFileButton{
		window: w,
		broker: b,
	}
	pf.ExtendBaseWidget(pf)

	pf.Icon = theme.FolderOpenIcon()
	pf.SetText("Open Database")

	pf.OnTapped = func() {
		d := NewFileReadDialog(pf.window, pf.OpenFile)
		d.SetConfirmText("Open")
		d.SetTitleText("Open Database")
		d.SetFileName("mayble.db")
		filter := storage.NewExtensionFileFilter([]string{".db", ".sqlite", ".sqlite3"})
		d.SetFilter(filter)
		d.Show()

	}
	return pf
}
func (pf *OpenFileButton) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn:
		pf.Tapped(nil)
	}
}

func (pf *OpenFileButton) OpenFile(uri fyne.URIReadCloser, err error) {

	if err != nil {
		pf.broker.Notify(emiter.Event{
			Name: gui.EventDisplayErr,
			Data: err,
		})
		return 
	}

	if uri == nil {
		return
	}

	filePath := uri.URI().Path()

	pf.broker.Notify(emiter.Event{
		Name: gui.EventDocumentNew,
		Data: filePath,
	})

} 

type NewFileButton struct {
	widget.Button
	window fyne.Window
	broker *emiter.Broker
}
func NewNewFileButton(w fyne.Window, b *emiter.Broker) *NewFileButton {
	nf := &NewFileButton{
		broker: b,
		window: w,
	}
	nf.ExtendBaseWidget(nf)

	nf.Icon = theme.FolderNewIcon()
	nf.SetText("New Database")

	nf.OnTapped = func() {
		d := NewFileWriteDialog(nf.window, nf.NewFile)
		d.SetConfirmText("Create")
		d.SetTitleText("Create Database")
		d.SetFileName("mayble.db")
		filter := storage.NewExtensionFileFilter([]string{".db", ".sqlite", ".sqlite3"})
		d.SetFilter(filter)
		d.Show()
	}

	return nf
}
func (nf *NewFileButton) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn:
		nf.Tapped(nil)
	}
}

func (nf *NewFileButton) NewFile(uri fyne.URIWriteCloser, err error) {
	if err != nil {
		nf.broker.Notify(emiter.Event{
			Name: gui.EventDisplayErr,
			Data: err,
		})
		return 
	}

	if uri == nil {
		return
	}

	filePath := uri.URI().Path()

	_ = os.Remove(filePath)

	ext := filepath.Ext(filePath)
	switch ext {
	case ".db", ".sqlite", ".sqlite3":
		break
	default:
		filePath = filePath + ".db"
	}

	nf.broker.Notify(emiter.Event{
		Name: gui.EventDocumentNew,
		Data: filePath,
	})
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
	eb.Icon = theme.LogoutIcon()
	eb.SetText("Export CSV")
	eb.OnTapped = func() {
		d := NewFileWriteDialog(eb.window, eb.OpenFile)
		filter := storage.NewExtensionFileFilter([]string{".csv"})
		d.SetFilter(filter)
		d.SetTitleText("Export")
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
	ib.SetText("Import CSV")
	ib.Icon = theme.LoginIcon()
	ib.OnTapped = func() {
		d := NewFileReadDialog(ib.window, ib.OpenFile)
		filter := storage.NewExtensionFileFilter([]string{".csv"})
		d.SetFilter(filter)
		d.SetTitleText("Import")
		d.Show()
	}
	return ib
}
func (ib *ImportButton) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn:
		ib.Tapped(nil)
	}
}

func (eb *ExportButton) TypedKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn:
		eb.Tapped(nil)
	}
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

	filePath := uri.URI().Path()

	err = uri.Close()
	if err != nil {
		eb.broker.Notify(emiter.Event{
			Name: gui.EventDisplayErr,
			Data: err,
		})
	}
	err = os.Remove(filePath)
	if err != nil {
		NotifyError(eb.broker, err)
		return 
	}

	ext := filepath.Ext(filePath)
	switch ext {
	case ".csv":
		break
	default:
		filePath = filePath + ".csv"
	}

	file, err := os.Create(filePath)
	if err != nil {
		NotifyError(eb.broker, err)
		return 
	}

	eb.broker.Notify(emiter.Event{
		Name: gui.EventDocumentExportCSV,
		Data: file,
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
