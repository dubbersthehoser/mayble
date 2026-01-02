package view

import (
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/searching"
	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/gui"
	"github.com/dubbersthehoser/mayble/internal/porting"
	csvPorting "github.com/dubbersthehoser/mayble/internal/porting/csv"
)

type FunkView struct {
	window     fyne.Window
	controller *controller.Controller
	View       fyne.CanvasObject
	emiter     *emiter.Emiter
	broker     *emiter.Broker
}

func NewFunkView(control *controller.Controller, window fyne.Window) (FunkView, error) {
	f := FunkView{
		controller: control,
		emiter: emiter.NewEmiter(),
		window: window,
		broker: control.Broker,
	}
	loadOnEventHandlers(&f)
	f.View = f.Body()
	syncView(&f)
	shortcutAdd(&f)
	return f, nil
}

func (f *FunkView) Body() fyne.CanvasObject {
	topBar := f.TopBar()
	table := f.Table()
	body := container.New(layout.NewBorderLayout(topBar, nil, nil, nil), topBar, table)
	return body
}

func NotifyError(b *emiter.Broker, err error) {
	b.Notify(emiter.Event{
		Name: gui.EventDisplayErr,
		Data: err,
	})
}

func (f *FunkView) refresh() {
	f.View.Refresh()
}


func shortcutAdd(f *FunkView) {

	ctrlF := &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierControl}

	ctrlM := &desktop.CustomShortcut{KeyName: fyne.KeyM, Modifier: fyne.KeyModifierControl}
	ctrlN := &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
	ctrlU := &desktop.CustomShortcut{KeyName: fyne.KeyU, Modifier: fyne.KeyModifierControl}
	ctrlShiftDel := &desktop.CustomShortcut{KeyName: fyne.KeyDelete, Modifier: fyne.KeyModifierShift | fyne.KeyModifierControl}
	ctrlShiftD := &desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: fyne.KeyModifierShift | fyne.KeyModifierControl}


	f.window.Canvas().AddShortcut(ctrlF, func(_ fyne.Shortcut) {
		f.broker.Notify(emiter.Event{
			Name: gui.EventSearchFocus,
		})
	})

	f.window.Canvas().AddShortcut(ctrlN, func(_ fyne.Shortcut) {
		f.broker.Notify(emiter.Event{
			Name: gui.EventEditerOpen,
			Data: gui.EventEntryCreate,
		})
	})
	f.window.Canvas().AddShortcut(ctrlU, func(_ fyne.Shortcut) {
		f.broker.Notify(emiter.Event{
			Name: gui.EventEditerOpen,
			Data: gui.EventEntryUpdate,
		})
	})
	f.window.Canvas().AddShortcut(ctrlShiftDel, func(_ fyne.Shortcut) {
		f.broker.Notify(emiter.Event{
			Name: gui.EventEntryDelete,
		})
	})
	f.window.Canvas().AddShortcut(ctrlShiftD, func(_ fyne.Shortcut) {
		f.broker.Notify(emiter.Event{
			Name: gui.EventEntryDelete,
		})
	})

	f.window.Canvas().AddShortcut(ctrlM, func(_ fyne.Shortcut) {
		f.broker.Notify(emiter.Event{
			Name: gui.EventMenuOpen,
		})
	})

	f.window.Canvas().AddShortcut(&fyne.ShortcutUndo{}, func(_ fyne.Shortcut) { // ctrl + Z
		f.broker.Notify(emiter.Event{
			Name: gui.EventUndo,
		})
	})

	f.window.Canvas().AddShortcut(&fyne.ShortcutRedo{}, func(_ fyne.Shortcut) { // ctrl + Y
		f.broker.Notify(emiter.Event{
			Name: gui.EventRedo,
		})
	})
}

func syncView(f *FunkView) {
	if f.controller.App.UndoIsEmpty() {
		f.broker.Notify(emiter.Event{
			Name: gui.EventUndoEmpty,
		})
	}
	if f.controller.App.RedoIsEmpty() {
		f.broker.Notify(emiter.Event{
			Name: gui.EventRedoEmpty,
		})
	}
	if !f.controller.List.HasSelected() {
		f.broker.Notify(emiter.Event{
			Name: gui.EventEntryUnselected,
		})
	}
	if !f.controller.Searcher.HasSelection() {
		f.broker.Notify(emiter.Event{
			Name: gui.EventSelectionNone,
		})
	}

	f.broker.Notify(emiter.Event{
		Name: gui.EventListOrderBy,
		Data: listing.ByTitle,
	})

	f.broker.Notify(emiter.Event{
		Name: gui.EventSearchBy,
		Data: searching.ByTitle,
	})
}


func loadOnEventHandlers(f *FunkView) {

	f.broker.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {

			case gui.EventRedo:
				err := f.controller.App.Redo()
				if err != nil {
					f.broker.Notify(emiter.Event{
						Name: gui.EventDisplayErr,
						Data: err,
					})
					return
				}
				f.broker.Notify(emiter.Event{
					Name: gui.EventDocumentModified,
				})

			case gui.EventUndo:
				err := f.controller.App.Undo()
				if err != nil {
					f.broker.Notify(emiter.Event{
						Name: gui.EventDisplayErr,
						Data: err,
					})
					return
				}
				f.broker.Notify(emiter.Event{
					Name: gui.EventDocumentModified,
				})

			case gui.EventDocumentModified:
				if f.controller.App.UndoIsEmpty() {
					f.broker.Notify(emiter.Event{
						Name: gui.EventUndoEmpty,
					})
				} else {
					f.broker.Notify(emiter.Event{
						Name: gui.EventUndoReady,
					})
				}

				
				if f.controller.App.RedoIsEmpty() {
					f.broker.Notify(emiter.Event{
						Name: gui.EventRedoEmpty,
					})
				} else {
					f.broker.Notify(emiter.Event{
						Name: gui.EventRedoReady,
					})
				}


			case gui.EventEditerOpen:
				subEvent := e.Data.(string)
				var builder *controller.BookLoanBuilder
				switch subEvent {
				case gui.EventEntryCreate:
					builder = controller.NewBookLoanBuilder()

				case gui.EventEntryUpdate:
					if !f.controller.List.HasSelected() {
						return
					}
					book := f.controller.List.Selected()
					builder = controller.NewBuilderWithBookLoan(book)
				}
				if builder == nil {
					panic("unexpected: builder is nil")
				}
				uniqueGenres := f.controller.List.UniqueGenres()
				ShowEditor(f.window, f.broker, builder, uniqueGenres)
			
			case gui.EventEntryDelete:
				if !f.controller.List.HasSelected() {
					return
				}
				book := f.controller.List.Selected()
				builder := controller.NewBuilderWithBookLoan(book)
				builder.Type = controller.Deleting
				f.broker.Notify(emiter.Event{
					Name: gui.EventEntrySubmit,
					Data: builder,
				})
				f.broker.Notify(emiter.Event{
					Name: gui.EventEntryUnselected,
				})

			case gui.EventDocumentImport:
				io := e.Data.(porting.NamedReadCloser)
				defer io.Close()

				porter, err := porting.GetBookLoanPorterByName(io.Name())
				if err != nil {
					NotifyError(f.broker, err)
					return
				}

				books, err := porter.ImportBookLoans(io)
				if err != nil {
					NotifyError(f.broker, err)
					return
				}

				err = f.controller.App.ImportBookLoans(books)
				if err != nil {
					NotifyError(f.broker, err)
					return
				}

				f.broker.Notify(emiter.Event{
					Name: gui.EventDocumentModified,
				})


			case gui.EventDocumentExport:
				io := e.Data.(porting.NamedWriteCloser)
				defer io.Close()

				porter, err := porting.GetBookLoanPorterByName(io.Name())
				if err != nil {
					NotifyError(f.broker, err)
					return
				}

				books, err := f.controller.App.GetBookLoans()
				if err != nil {
					NotifyError(f.broker, err)
					return
				}

				err = porter.ExportBookLoans(io, books)
				if err != nil {
					NotifyError(f.broker, err)
					return
				}

			case gui.EventDocumentExportCSV:
				ioCloser := e.Data.(io.WriteCloser)
				defer ioCloser.Close()

				porter := csvPorting.BookLoanPorter{}

				books, err := f.controller.App.GetBookLoans()
				if err != nil {
					NotifyError(f.broker, err)
					return
				}
				err = porter.ExportBookLoans(ioCloser, books)
				if err != nil {
					NotifyError(f.broker, err)
					return 
				}


			case gui.EventDisplayErr:
				err := e.Data.(error)
				dialog.ShowError(err, f.window)
				
			case gui.EventMenuOpen:
				ShowMenu(f)

			case gui.EventDocumentNew:
				path := e.Data.(string)
				err := f.controller.Config.SetDBFile(path)
				if err != nil {
					NotifyError(f.broker, err)
					return
				}
				
				err = f.controller.Reset()
				if err != nil {
					NotifyError(f.broker, err)
					return
				}
				syncView(f)

			case gui.EventElementUnfocus:
				f.window.Canvas().Unfocus()
			}
		},
	}, 
		gui.EventRedo,
		gui.EventUndo,
		gui.EventEditerOpen,
		gui.EventDisplayErr,
		gui.EventEntryDelete,
		gui.EventMenuOpen,
		gui.EventDocumentImport,
		gui.EventDocumentExport,
		gui.EventDocumentModified,
		gui.EventDocumentNew,
		gui.EventElementUnfocus,
	)
}

