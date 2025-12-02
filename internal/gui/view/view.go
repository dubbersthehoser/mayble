package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/dialog"
	
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/searching"
	"github.com/dubbersthehoser/mayble/internal/listing"
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
}

func NewFunkView(control *controller.Controller, window fyne.Window) (FunkView, error) {
	f := FunkView{
		controller: control,
		emiter: emiter.NewEmiter(),
		window: window,
	}

	loadOnEventHandlers(&f)

	f.View = f.Body()
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
	OnRedo         = "ON_REDO"

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

func loadOnEventHandlers(f *FunkView) {
	f.emiter.OnEvent(OnUpdate, handleUpdate(f))
	f.emiter.OnEvent(OnCreate, handleCreate(f))
	f.emiter.OnEvent(OnModification, handleModification(f))
	f.emiter.OnEvent(OnDelete, handleDelete(f))
	f.emiter.OnEvent(OnRedo, handleRedo(f))
	f.emiter.OnEvent(OnUndo, handleUndo(f))
	f.emiter.OnEvent(OnMenuOpen, handleMenuOpen(f))
	f.emiter.OnEvent(OnSelectNext, handleSelectNext(f))
	f.emiter.OnEvent(OnSelectPrev, handleSelectPrev(f))
	f.emiter.OnEvent(OnSelected, handleSelected(f))
	f.emiter.OnEvent(OnSetSearchBy, handleSetSearchBy(f))
	f.emiter.OnEvent(OnSetSearchPattern, handleSetSearchPattern(f))
	f.emiter.OnEvent(OnSearch, handleSearch(f))
	f.emiter.OnEvent(OnSetOrdering, handleSetOrdering(f))
	f.emiter.OnEvent(OnSetOrderBy, handleSetOrderBy(f))
	f.emiter.OnEvent(OnShowError, handleShowError(f))
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

func handleSetOrdering(f *FunkView) func(any) {
	return func(data any) {
		o, ok := data.(listing.Ordering)
		if !ok {
			panic("invalid data to OnOrdering event")
		}
		fmt.Println("view: set ordering: ", o)
		f.controller.List.SetOrdering(o)
		err := f.controller.List.Update()
		if err != nil {
			f.displayError(err)
			return
		}
		f.emiter.Emit(OnUnselected, nil)
		f.refresh()
	}
}

func handleSetOrderBy(f *FunkView) func(any) {
	return func(data any) {
		var by listing.OrderBy
		switch v := data.(type) {
		case string:
			by = listing.MustStringToOrderBy(v)
		case listing.OrderBy:
			by = v
		default:
			panic("invalid data to OnOrderBy event")
		}
		fmt.Println("view: set order by: ", by)
		f.controller.List.SetOrderBy(by)

	}
}

func handleSelected(f *FunkView) func(any) {
	return func(data any) {
		index, ok := data.(int)
		if !ok {
			panic("given invalid data for Selected event")
		}

		err := f.controller.List.Select(index)
		if err != nil {
			f.displayError(err)
		}
	}
}


func handleSelectNext(f *FunkView) func(any) {
	return func(data any) {
		f.controller.List.SelectNext()
	}
}
func handleSelectPrev(f *FunkView) func(any) {
	return func(data any) {
		f.controller.List.SelectPrev()
	}
}

func handleUnselect(f *FunkView) func(any) {
	return func(_ any) {
		f.controller.List.Unselect()
	}
}

func handleSetSearchBy(f *FunkView) func(any) {
	return func(data any) {
		by, ok := data.(string)
		if !ok {
			panic("given invalid data for SearchBy event")
		}
		field := searching.MustStringToField(by)
		f.controller.List.SetSearchBy(field)
	}
}
func handleSetSearchPattern(f *FunkView) func(any) {
	return func(data any) {
		pattern, ok := data.(string)
		if !ok {
			panic("given invalid data for Search event")
		}
		f.controller.List.SetSearchPattern(pattern)
	}
}
func handleSearch(f *FunkView) func(any) {
	return func(_ any) {
		f.controller.List.Search()
	}
}

func handleMenuOpen(f *FunkView) func(any) {
	return func(data any) {
		f.ShowMenu()
	}
}

func handleModification(f *FunkView) func(any) {
	return func(_ any) {
		err := f.controller.List.Update()
		if err != nil {
			f.displayError(err)
			return
		}
		f.refresh()
	}
}

func handleRedo(f *FunkView) func(any) {
	return func(data any) {
		err := f.controller.App.Redo()
		if err != nil {
			f.displayError(err)
			return
		}
		f.emiter.Emit(OnModification, nil)
	}
}

func handleUndo(f *FunkView) func(any) {
	return func(data any) {
		err := f.controller.App.Undo()
		if err != nil {
			f.displayError(err)
			return
		}
		f.emiter.Emit(OnModification, nil)
	}
}

func handleDelete(f *FunkView) func(any) {
	return func(data any) {
		bookLoan, err := f.controller.List.Selected()
		if err != nil {
			f.displayError(err)
			return
		}
		builder := controller.NewBuilderWithBookLoan(bookLoan)
		builder.Type = controller.Deleting
		f.controller.Editor.Submit(builder)
		f.emiter.Emit(OnModification, nil)
	}
}

func handleUpdate(f *FunkView) func(any) {
	return func(data any) {
		bookLoan, err := f.controller.List.Selected()
		if err != nil {
			f.displayError(err)
			return 
		}
		builder := controller.NewBuilderWithBookLoan(bookLoan)
		f.ShowEdit(builder)
	}
}

func handleCreate(f *FunkView) func(any) {
	return func (_ any) {
		builder := controller.NewBookLoanBuilder()
		f.ShowEdit(builder)
	}
}

