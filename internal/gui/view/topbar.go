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

	"github.com/dubbersthehoser/mayble/internal/searching"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/gui"

)


func (f *FunkView) TopBar() fyne.CanvasObject {

	saveItem := NewToolbarSave(f.broker)
	saveItem.Disable()

	menuItem := NewToolbarMenu(f.emiter)
	createItem := NewToolbarCreate(f.broker)
	updateItem := NewToolbarUpdate(f.broker)
	deleteItem := NewToolbarDelete(f.broker)

	undoItem := NewToolbarUndo(f.broker)
	redoItem := NewToolbarRedo(f.broker)

	nextItem := NewToolbarNext(f.broker)
	prevItem := NewToolbarPrev(f.broker)

	selectSearchBy := NewSearchBySelect(f.broker)

	searchEnt := NewSearchEntry(f.broker)

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
	broker *emiter.Broker
}
func NewSearchEntry(b *emiter.Broker) *SearchEntry {
	se := &SearchEntry{}
	se.ExtendBaseWidget(se)

	se.broker = b
	se.Entry.PlaceHolder = "Search"
	se.Entry.OnChanged = se.Search
	se.Entry.OnSubmitted = se.Submit

	b.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {

			case gui.EventSearchBy:
				se.Entry.SetText("") 

			case gui.EventListOrdering:
				se.Entry.SetText("")

			default:
				panic("event not found: " + e.Name)
			}
		},
	}, 
		gui.EventSearchBy, 
		gui.EventListOrdering,
	)

	return se
}
func (se *SearchEntry) Search(s string) {
	se.broker.Notify(emiter.Event{
		Name: gui.EventSearchPattern,
		Data: s,
	})
}
func (se *SearchEntry) Submit(s string) {
	se.broker.Notify(emiter.Event{
		Name: gui.EventSelectNext,
	})
}



func (se *SearchEntry) MinSize() fyne.Size {
	size := se.Entry.MinSize()
	size.Width += 350.0 // set horizontal size for VBox
	return size
}


type SearchBySelect struct {
	widget.Select
	broker *emiter.Broker
}
func NewSearchBySelect(b *emiter.Broker) *SearchBySelect {
	sb := &SearchBySelect{}
	sb.ExtendBaseWidget(sb)
	sb.broker = b
	sb.SetOptions([]string{"Title", "Author", "Genre", "Borrower"})
	sb.PlaceHolder = "Search By"
	sb.OnChanged = sb.OnSelected
	return sb
}
func (sb *SearchBySelect) OnSelected(s string) {
	by := searching.MustStringToField(s)
	sb.broker.Notify(emiter.Event{
		Name: gui.EventSearchBy,
		Data: by,
	})
}


/*****************************
        Toolbar Items
******************************/

func NewToolbarSave(b *emiter.Broker) *widget.ToolbarAction {
	ts := &widget.ToolbarAction{
		Icon: theme.DocumentSaveIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventSave,
			})
		},
	}
	b.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case gui.EventSave:
				ts.Disable()

			case gui.EventEntryCreate,
				gui.EventEntryDelete,
				gui.EventEntryUpdate:

				ts.Enable()
			}
		},
	},
		gui.EventSave,
		gui.EventEntryCreate, 
		gui.EventEntryDelete, 
		gui.EventEntryUpdate,
	)
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

func NewToolbarCreate(b *emiter.Broker) *widget.ToolbarAction {
	tc := &widget.ToolbarAction{
		Icon: theme.ContentAddIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventEditOpen,
				Data: gui.EventEntryCreate,
			})
		},
	}
	return tc
}

func NewToolbarUpdate(b *emiter.Broker) *widget.ToolbarAction {
	tu := &widget.ToolbarAction{
		Icon: theme.DocumentCreateIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventEditOpen,
				Data: gui.EventEntryUpdate,
			})
		},
	}
	b.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case gui.EventEntrySelected:
				tu.Enable()
			case gui.EventEntryUnselected:
				tu.Disable()
			}
		},
	},
		gui.EventEntrySelected,
		gui.EventEntryUnselected,
	)
	return tu
}

func NewToolbarDelete(b *emiter.Broker) *widget.ToolbarAction {
	td := &widget.ToolbarAction{
		Icon: theme.DeleteIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventEntryDelete,
			})
		}, 
	}
	b.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case gui.EventEntrySelected:
				td.Enable()

			case gui.EventEntryUnselected:
				td.Disable()

			default:
				panic("event not found: " + e.Name)
			}
		},
	},
		gui.EventEntrySelected,
		gui.EventEntryUnselected,
	)
	return td
}

func NewToolbarUndo(b *emiter.Broker) *widget.ToolbarAction {
	tu := &widget.ToolbarAction{
		Icon: theme.ContentUndoIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventUndo,
			})
		},
	}

	l := emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case gui.EventUndoEmpty:
				tu.Disable()

			case gui.EventUndoReady:
				tu.Enable()

			default:
				panic("event not found: " + e.Name)
			}
		},
	}
	id := b.Subscribe(&l, gui.EventUndoEmpty, gui.EventUndoReady)
	println("undo listen with id: ", id)
	return tu

}

func NewToolbarRedo(b *emiter.Broker) *widget.ToolbarAction {
	tr := &widget.ToolbarAction{
		Icon: theme.ContentRedoIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventRedo,
			})
		},
	}
	l := emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {
			case gui.EventRedoEmpty:
				tr.Disable()

			case gui.EventRedoReady:
				tr.Enable()
			default:
				panic("event not found: " + e.Name)
			}
		},
	}
	id := b.Subscribe(&l, gui.EventRedoEmpty, gui.EventRedoReady)
	println("redo listen with id: ", id)
	return tr
}

func NewToolbarNext(b *emiter.Broker) *widget.ToolbarAction {
	tn := &widget.ToolbarAction{
		Icon: theme.MoveDownIcon() ,
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventSelectNext,
			})
		},
	}
	return tn
}
func NewToolbarPrev(b *emiter.Broker) *widget.ToolbarAction {
	tp := &widget.ToolbarAction{
		Icon: theme.MoveUpIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventSelectPrev,
			})
		},
	}
	return tp
}



