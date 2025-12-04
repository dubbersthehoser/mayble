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

	saveItem := NewToolbarSave(f.broker)
	saveItem.Disable()

	menuItem := NewToolbarMenu(f.emiter)
	createItem := NewToolbarCreate(f.emiter)
	updateItem := NewToolbarUpdate(f.emiter)
	deleteItem := NewToolbarDelete(f.emiter)

	undoItem := NewToolbarUndo(f.broker)
	redoItem := NewToolbarRedo(f.broker)

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

func NewToolbarSave(b *emiter.Broker) *widget.ToolbarAction {
	ts := &widget.ToolbarAction{
		Icon: theme.DocumentSaveIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: EventSave,
			})
		},
	}
	l := emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case EventSaveDisable:
				ts.Disable()

			case EventSaveEnable:
				ts.Enable()
			}
		},
	}
	b.Subscribe(&l, EventSaveDisable, EventSaveEnable)
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

func NewToolbarUndo(b *emiter.Broker) *widget.ToolbarAction {
	tu := &widget.ToolbarAction{
		Icon: theme.ContentUndoIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: EventUndo,
			})
		},
	}

	l := emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case EventUndoEmpty:
				tu.Disable()
			case EventUndoReady:
				tu.Enable()
			}
		}
	}
	b.Subscribe(&l, EventUndoEmpty, EventUndoReady)
	println("undo listen with id: ", l.id)

}
func NewToolbarRedo(b *emiter.Broker) *widget.ToolbarAction {
	tr := &widget.ToolbarAction{
		Icon: theme.ContentRedoIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: EventRedo,
			})
		},
	}
	l := emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case EventUndoEmpty:
				tu.Disable()
			case EventUndoReady:
				tu.Enable()
			}
		}
	}
	b.Subscribe(&l, EventRedoEmpty, EventRedoReady)
	println("redo listen with id: ", l.id)
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



