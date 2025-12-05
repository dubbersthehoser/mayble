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

func (f *FunkView) displayError(err error) {
	if err  != nil {
		dialog.ShowError(err, f.window)
	}
}

func (f *FunkView) refresh() {
	f.View.Refresh()
}


/*******************
	Events
********************/

const (
	OnShowError string = "ON_SHOW_ERR"
	OnSave         = "ON_SAVE"
	OnExport       = "ON_EXPORT"

	OnCreate       = "ON_CREATE"
	OnUpdate       = "ON_UPDATE"
	OnDelete       = "ON_DELETE"

	OnUndo         = "ON_UNDO"
	OnUndoEmpty    = "ON_UNDO_EMPTY"
	OnUndoReady    = "ON_UNDO_READY"
	OnRedo         = "ON_REDO"
	OnRedoEmpty    = "ON_REDO_EMPTY"
	OnRedoReady    = "ON_REDO_READY"

	OnSearch           = "ON_SEARCH"
	OnSetSearchPattern = "ON_SET_SEARCH_PATTERN"
	OnSetSearchBy      = "ON_SET_SEARCH_BY"

	OnSetOrdering  = "ON_SET_ORDERING"
	OnSetOrderBy   = "ON_SET_ORDER_BY"

	OnSelectNext   = "ON_SELECT_NEXT"
	OnSelectPrev   = "ON_SELECT_PREV"

	OnSelected     = "ON_SELECTED"
	OnUnselected   = "ON_UNSELECTED"

	OnModification = "ON_MODIFICATION"

	OnMenuOpen     = "ON_MENU_OPEN"
	OnMenuClose    = "ON_MENU_CLOSE"
)

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

	f.emiter.OnEvent(OnMenuOpen, handleMenuOpen(f))
	f.emiter.OnEvent(OnShowError, handleShowError(f))

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
				f.ShowEditor(builder)
			}
		},
	}, 
		gui.EventSave,
		gui.EventRedo,
		gui.EventUndo,
		gui.EventEditerOpen,
	)
}


func handleShowError(f *FunkView) func(any) {
	return func(data any) {
		err, ok := data.(error)
		if !ok {
			panic("given invalid data for Error event")
		}
		f.displayError(err)
	}
}


func handleMenuOpen(f *FunkView) func(any) {
	return func(data any) {
		f.ShowMenu()
	}
}

