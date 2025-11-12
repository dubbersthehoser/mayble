package view

import (
	//"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/dialog"
	
	"github.com/dubbersthehoser/mayble/internal/gui/controller"

)

const (
	OnSave string  = "ON_SAVE"
	OnExport       = "ON_EXPORT"

	OnCreate       = "ON_CREATE"
	OnUpdate       = "ON_UPDATE"
	OnDelete       = "ON_DELETE"

	OnUndo         = "ON_UNDO"
	OnRedo         = "ON_REDO"

	OnSort         = "ON_SORT"
	OnSearch       = "ON_SEARCH"
	OnSearchBy     = "ON_SEARCH_BY"

	OnSelectNext   = "ON_SELECT_NEXT"
	OnSelectPrev   = "ON_SELECT_PREV"

	OnSelected     = "ON_SELECTED"
	OnUnselected   = "ON_UNSELECTED"

	OnModification = "ON_MODIFICATION"

	OnMenuOpen     = "ON_MENU_OPEN"
	OnMenuClose    = "ON_MENU_CLOSE"
)

type emiter struct {
	items map[string][]func()
}

func (e *emiter) init() {
	if e.items == nil {
		e.items = make(map[string][]func())
	}
}

func (e *emiter) On(key string, do func()) {
	e.init()
	list, ok := e.items[key]
	if !ok {
		list = make([]func(), 0)
	}
	list = append(list, do)
	e.items[key] = list
}
func (e *emiter) Emit(key string) {
	e.init()
	list, ok := e.items[key]
	if !ok {
		return
	}
	for _, fn := range list {
		fn()
	}
}

/***********************
	FunkView
************************/

// NOTE It's called FunkView because I was frustrated and I needed some amusement.

type FunkView struct {
	window     fyne.Window
	controller *controller.Controller
	View       fyne.CanvasObject
	emiter     emiter
}

func NewFunkView(control *controller.Controller, window fyne.Window) (FunkView, error) {
	f := FunkView{
		controller: control,
		emiter: emiter{},
		window: window,
	}

	f.loadEvents()

	f.View = f.Body()
	f.emiter.Emit(OnSort)
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

func (f *FunkView) loadEvents() {
	f.emiter.On(OnUpdate, f.EventUpdate)
	f.emiter.On(OnCreate, f.EventCreate)
	f.emiter.On(OnModification, f.EventModification)
	f.emiter.On(OnDelete, f.EventDelete)
	f.emiter.On(OnRedo, f.EventRedo)
	f.emiter.On(OnUndo, f.EventUndo)
	f.emiter.On(OnSort, f.EventSort)
	f.emiter.On(OnSearch, f.EventSearch)
	f.emiter.On(OnSearchBy, f.EventSearchBy)
	f.emiter.On(OnMenuOpen, f.EventMenuOpen)
	f.emiter.On(OnSelectNext, f.EventSelectNext)
	f.emiter.On(OnSelectPrev, f.EventSelectPrev)
}

func (f *FunkView) EventSelectNext() {
	f.controller.BookList.SelectNext()
}

func (f *FunkView) EventSelectPrev() {
	f.controller.BookList.SelectPrev()
}

func (f *FunkView) EventSearchBy() {
	f.refresh()
}
func (f *FunkView) EventSearch() {
	f.controller.BookList.Search()
	f.refresh()
}

func (f *FunkView) EventMenuOpen() {
	f.ShowMenu()
}



func (f *FunkView) EventModification() {
	err := f.controller.BookList.Update()
	if err != nil {
		f.displayError(err)
		return
	}
	f.refresh()
}

func (f *FunkView) EventSort() {
	err := f.controller.BookList.Update()
	if err != nil {
		f.displayError(err)
		return
	}
	f.refresh()
}



func (f *FunkView) EventRedo() {
	err := f.controller.Core.Redo()
	if err != nil {
		f.displayError(err)
		return
	}
	f.emiter.Emit(OnModification)
}

func (f *FunkView) EventUndo() {
	err := f.controller.Core.Undo()
	if err != nil {
		f.displayError(err)
		return
	}
	f.emiter.Emit(OnModification)
}

func (f *FunkView) EventDelete() {
	bookLoan, err := f.controller.BookList.Selected()
	if err != nil {
		f.displayError(err)
		return
	}
	builder := controller.NewBuilderWithBookLoan(bookLoan)
	builder.Type = controller.Deleting
	f.controller.BookEditor.Submit(builder)
	f.emiter.Emit(OnModification)
}

func (f *FunkView) EventUpdate() {
	bookLoan, err := f.controller.BookList.Selected()
	if err != nil {
		f.displayError(err)
		return 
	}
	builder := controller.NewBuilderWithBookLoan(bookLoan)
	f.ShowEdit(builder)
}

func (f *FunkView) EventCreate() {
	builder := controller.NewBookLoanBuilder()
	f.ShowEdit(builder)
}
