package view

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/dialog"
	
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/searching"
	"github.com/dubbersthehoser/mayble/internal/listing"
	"github.com/dubbersthehoser/mayble/internal/gui"
	"github.com/dubbersthehoser/mayble/internal/porting"
)



/***********************
	FunkView
************************/

// NOTE It's called FunkView because I was frustrated and I needed some amusement.

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
			case gui.EventSave:
				f.controller.App.Save()
				f.broker.Notify(emiter.Event{
					Name: gui.EventSaveDisable,
				})

			case gui.EventRedo:
				f.controller.App.Redo()
				if f.controller.App.RedoIsEmpty() {
					f.broker.Notify(emiter.Event{
						Name: gui.EventRedoEmpty,
					})
				}

			case gui.EventUndo:
				f.controller.App.Undo()
				if f.controller.App.UndoIsEmpty() {
					f.broker.Notify(emiter.Event{
						Name: gui.EventUndoEmpty,
					})
				}

			case gui.EventEditerOpen:
				subEvent := e.Data.(string)
				var builder *controller.BookLoanBuilder
				switch subEvent {
				case gui.EventEntryCreate:
					builder = controller.NewBookLoanBuilder()

				case gui.EventEntryUpdate:
					book := f.controller.List.Selected()
					builder = controller.NewBuilderWithBookLoan(book)
				}
				if builder == nil {
					panic("unexpected: builder is nil")
				}
				ShowEditor(f.window, f.broker, builder)
			
			case gui.EventEntryDelete:
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
				


			case gui.EventDisplayErr:
				err := e.Data.(error)
				dialog.ShowError(err, f.window)
				
			case gui.EventMenuOpen:
				ShowMenu(f)

			}

				
		},
	}, 
		gui.EventSave,
		gui.EventRedo,
		gui.EventUndo,
		gui.EventEditerOpen,
		gui.EventDisplayErr,
		gui.EventEntryDelete,
		gui.EventMenuOpen,
		gui.EventDocumentImport,
		gui.EventDocumentExport,
	)
}

