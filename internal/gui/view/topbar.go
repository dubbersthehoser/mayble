package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"

	"github.com/dubbersthehoser/mayble/internal/searching"
	"github.com/dubbersthehoser/mayble/internal/emiter"
	"github.com/dubbersthehoser/mayble/internal/gui"

)


func (f *FunkView) TopBar() fyne.CanvasObject {

	menuItem := NewToolbarMenu(f.broker)
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
		widget.NewToolbarSeparator(),
		undoItem,
		redoItem,
		widget.NewToolbarSeparator(),
		createItem,
		updateItem,
		deleteItem,
		widget.NewToolbarSeparator(),
		nextItem,
		prevItem,
		searchEnt,
		selectSearchBy,

	}
	toolBar := widget.NewToolbar(items...)
	return toolBar
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

func (se *SearchEntry) ToolbarObject() fyne.CanvasObject {
	return se
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

	b.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			by := e.Data.(searching.Field)
			switch by {
			case searching.ByTitle:
				sb.Selected = sb.Options[0]
			case searching.ByAuthor:
				sb.Selected = sb.Options[1]
			case searching.ByGenre:
				sb.Selected = sb.Options[2]
			case searching.ByBorrower:
				sb.Selected = sb.Options[3]
			}
			sb.Refresh()
		},
	},
		gui.EventSearchBy,
	)
	return sb
}
func (sb *SearchBySelect) OnSelected(s string) {
	by := searching.MustStringToField(s)
	sb.broker.Notify(emiter.Event{
		Name: gui.EventSearchBy,
		Data: by,
	})
}

func (sb *SearchBySelect) ToolbarObject() fyne.CanvasObject {
	return sb
}


/*****************************
        Toolbar Items
******************************/

func NewToolbarMenu(b *emiter.Broker) *widget.ToolbarAction {
	tm := &widget.ToolbarAction{
		Icon: theme.MenuIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventMenuOpen,
			})
		},
	}
	return tm

}

func NewToolbarCreate(b *emiter.Broker) *widget.ToolbarAction {
	tc := &widget.ToolbarAction{
		Icon: theme.ContentAddIcon(),
		OnActivated: func() {
			b.Notify(emiter.Event{
				Name: gui.EventEditerOpen,
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
				Name: gui.EventEditerOpen,
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

	b.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {

			case gui.EventSelectionNone:
				tn.Disable()

			case gui.EventSelection:
				tn.Enable()
			}
		},
	},
		gui.EventSelectionNone,
		gui.EventSelection,
	)

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
	b.Subscribe(&emiter.Listener{
		Handler: func(e *emiter.Event) {
			switch e.Name {

			case gui.EventSelectionNone:
				tp.Disable()

			case gui.EventSelection:
				tp.Enable()
			}
		},
	},
		gui.EventSelectionNone,
		gui.EventSelection,
	)
	return tp
}



