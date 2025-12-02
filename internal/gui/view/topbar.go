package view

import (
	//"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"

	//"github.com/dubbersthehoser/mayble/internal/searching"
	"github.com/dubbersthehoser/mayble/internal/emiter"
)


func (f *FunkView) TopBar() fyne.CanvasObject {

	saveItem := NewToolbarSave(f.emiter)
	saveItem.Disable()

	menuItem := NewToolbarMenu(f.emiter)
	createItem := NewToolbarCreate(f.emiter)
	updateItem := NewToolbarUpdate(f.emiter)
	deleteItem := NewToolbarDelete(f.emiter)

	undoItem := NewToolbarUndo(f.emiter)
	redoItem := NewToolbarRedo(f.emiter)
	undoItem.Disable()
	redoItem.Disable()

	nextItem := NewToolbarNext(f.emiter)
	prevItem := NewToolbarPrev(f.emiter)

	selectSearchBy := NewSearchBySelect(f.emiter)

	searchEnt := NewSearchEntry(f.emiter)

	items := []widget.ToolbarItem{
		menuItem,
		saveItem,
		widget.NewToolbarSeparator(),
		createItem,
		updateItem,
		deleteItem,
		widget.NewToolbarSeparator(),
		undoItem,
		redoItem,
		widget.NewToolbarSeparator(),
		nextItem,
		prevItem,
	}
	toolBar := widget.NewToolbar(items...)
	
	boxes := []fyne.CanvasObject{
		toolBar,
		searchEnt,
		selectSearchBy,
	}
	return container.New(layout.NewHBoxLayout(), boxes...)
}


type SearchEntry struct {
	widget.Entry
	emiter *emiter.Emiter
}
func NewSearchEntry(e *emiter.Emiter) *SearchEntry {
	se := &SearchEntry{}
	se.ExtendBaseWidget(se)
	se.emiter = e
	se.Entry.PlaceHolder = "Search"
	se.Entry.OnChanged = se.Search
	se.Entry.OnSubmitted = se.Submit
	se.emiter.OnEvent(OnSetSearchBy, func(_ any) {se.Reset()})
	se.emiter.OnEvent(OnSetOrdering, func(_ any) {se.Reset()})
	return se
}
func (se *SearchEntry) Search(s string) {
	se.emiter.Emit(OnSetSearchPattern, s)
	se.emiter.Emit(OnSearch, nil)
}
func (se *SearchEntry) Submit(s string) {
	se.emiter.Emit(OnSelectNext, nil)
}
func (se *SearchEntry) Reset() {
	se.Entry.SetText("")
}

func (se *SearchEntry) MinSize() fyne.Size {
	size := se.Entry.MinSize()
	size.Width += 350.0 // set horizontal size for vBox
	return size
}


type SearchBySelect struct {
	widget.Select
	emiter *emiter.Emiter
}
func NewSearchBySelect(e *emiter.Emiter) *SearchBySelect {
	sb := &SearchBySelect{}
	sb.ExtendBaseWidget(sb)
	sb.emiter = e
	sb.SetOptions([]string{"Title", "Author", "Genre", "Borrower"})
	sb.PlaceHolder = "Search By"
	sb.OnChanged = sb.OnSelected
	return sb
}
func (sb *SearchBySelect) OnSelected(s string) {
	sb.emiter.Emit(OnSetSearchBy, s)
}

/*****************************
        Toolbar Items
******************************/

func NewToolbarSave(e *emiter.Emiter) *widget.ToolbarAction {
	ts := &widget.ToolbarAction{
		Icon: theme.DocumentSaveIcon(),
		OnActivated: func() {
			e.Emit(OnSave, nil)
		},
	}
	e.OnEvent(OnModification, func(_ any) {
		ts.Enable()
	})
	e.OnEvent(OnSave, func(_ any){
		ts.Disable()
	})
	return ts
}

func NewToolbarMenu(e *emiter.Emiter) *widget.ToolbarAction {
	tm := &widget.ToolbarAction{
		Icon: theme.MenuIcon(),
		OnActivated: func() {
			e.Emit(OnMenuOpen, nil)
		},
	}
	return tm

}

func NewToolbarCreate(e *emiter.Emiter) *widget.ToolbarAction {
	tc := &widget.ToolbarAction{
		Icon: theme.ContentAddIcon(),
		OnActivated: func() {
			e.Emit(OnCreate, nil)
		},
	}
	return tc
}

func NewToolbarUpdate(e *emiter.Emiter) *widget.ToolbarAction {
	tu := &widget.ToolbarAction{
		Icon: theme.DocumentCreateIcon(),
		OnActivated: func() {
			e.Emit(OnUpdate, nil)
		},
	}
	e.OnEvent(OnSelected, func(_ any) {
		tu.Enable()
	})
	e.OnEvent(OnUnselected, func(_ any) {
		tu.Disable()
	})
	return tu
}

func NewToolbarDelete(e *emiter.Emiter) *widget.ToolbarAction {
	td := &widget.ToolbarAction{
		Icon: theme.DeleteIcon(),
		OnActivated: func() {
			e.Emit(OnDelete, nil)
		}, 
	}
	e.OnEvent(OnSelected, func(_ any) {
		td.Enable()
	})
	e.OnEvent(OnUnselected, func(_ any) {
		td.Disable()
	})
	return td
}

func NewToolbarUndo(e *emiter.Emiter) *widget.ToolbarAction {
	tu := &widget.ToolbarAction{
		Icon: theme.ContentUndoIcon(),
		OnActivated: func() {
			e.Emit(OnUndo, nil)
		},
	}

	e.OnEvent(OnUndoEmpty, func(_ any){
		tu.Disable()
	})

	e.OnEvent(OnUndoReady, func(_ any) {
		tu.Enable()
	})
	return tu
}
func NewToolbarRedo(e *emiter.Emiter) *widget.ToolbarAction {
	tr := &widget.ToolbarAction{
		Icon: theme.ContentRedoIcon(),
		OnActivated: func() {
			e.Emit(OnRedo, nil)
		},
	}
	e.OnEvent(OnRedoEmpty, func(_ any){
		tr.Disable()
	})
	e.OnEvent(OnRedoReady, func(_ any) {
		tr.Enable()
	})
	return tr
}

func NewToolbarNext(e *emiter.Emiter) *widget.ToolbarAction {
	tn := &widget.ToolbarAction{
		Icon: theme.MoveDownIcon() ,
		OnActivated: func() {
			e.Emit(OnSelectNext, nil)
		},
	}
	return tn
}
func NewToolbarPrev(e *emiter.Emiter) *widget.ToolbarAction {
	tp := &widget.ToolbarAction{
		Icon: theme.MoveUpIcon(),
		OnActivated: func() {
			e.Emit(OnSelectPrev, nil)
		},
	}
	return tp
}



