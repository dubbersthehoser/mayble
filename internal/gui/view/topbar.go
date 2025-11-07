package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"

	"github.com/dubbersthehoser/mayble/internal/gui/controller"
)

func (f *FunkView) TopBar() fyne.CanvasObject {

	// Save
	//------
	OnSaveItem := func() {
		f.emiter.Emit(OnSave)
	}

	saveItem := &widget.ToolbarAction{
		Icon: theme.DocumentSaveIcon(),
		OnActivated: OnSaveItem,
	}

	saveItem.Disable()

	// Events
	EnableSaveItem := func() {
		saveItem.Enable()
	}
	DisableSaveItem := func() {
		saveItem.Disable()
	}

	f.emiter.On(OnModification, EnableSaveItem)
	f.emiter.On(OnSave, DisableSaveItem)


	// Menu
	//------
	OnMenuItem := func() {
		fmt.Println("Menu button pressed")
		f.emiter.Emit(OnMenuOpen)
	}
	menuItem := &widget.ToolbarAction{
		Icon: theme.MenuIcon(),
		OnActivated: OnMenuItem,
	}



	// Create 
	//--------
	OnCreateItem := func() {
		fmt.Println("Create button pressed")
		f.emiter.Emit(OnCreate)
	}
	createItem := &widget.ToolbarAction{
		Icon: theme.ContentAddIcon(),
		OnActivated: OnCreateItem,
	}


	// Update
	//--------
	OnUpdateItem := func() {
		fmt.Println("Update button pressed")
		f.emiter.Emit(OnUpdate)
	}
	updateItem := &widget.ToolbarAction{
		Icon: theme.DocumentCreateIcon(),
		OnActivated: OnUpdateItem,
	}

	// Delete
	//--------
	OnDeleteItem := func() {
		fmt.Println("Delete button pressed")
		f.emiter.Emit(OnDelete)
	}
	deleteItem := &widget.ToolbarAction{
		Icon: theme.DeleteIcon(),
		OnActivated: OnDeleteItem, 
	}


	// Events
	DisableItemOnMod := func() {
		updateItem.Disable()
		deleteItem.Disable()
	}
	EnableItemOnMod := func() {
		updateItem.Enable()
		deleteItem.Enable()
	}
	updateItem.Disable()
	deleteItem.Disable()
	
	f.emiter.On(OnSelected, EnableItemOnMod)
	f.emiter.On(OnUnselected, DisableItemOnMod)

	// Undo Redo
	//-----------
	//
	// Undo
	//------
	OnUndoItem := func() {
		fmt.Println("Undo button pressed")
		f.emiter.Emit(OnUndo)
	}
	undoItem := &widget.ToolbarAction{
		Icon: theme.ContentUndoIcon(),
		OnActivated: OnUndoItem,
	}

	// Redo
	//------
	OnRedoItem := func() {
		fmt.Println("Redo button pressed")
		f.emiter.Emit(OnRedo)
	}
	redoItem := &widget.ToolbarAction{
		Icon: theme.ContentRedoIcon(),
		OnActivated: OnRedoItem,
	}

	undoItem.Disable()
	redoItem.Disable()

	f.emiter.On(OnModification, func() {
		if f.controller.Core.IsUndo() {
			undoItem.Enable()
		} else {
			undoItem.Disable()
		}
		if f.controller.Core.IsRedo() {
			redoItem.Enable()
		} else {
			redoItem.Disable()
		}
	})


	// Search
	//---------
	//
	// Search By
	//-----------
	selectSearchBy := widget.NewSelect(
		[]string{"Title", "Author", "Genre", "Borrower"},
		func(s string) {
			switch s {
			case "Title":
				f.controller.BookList.SetSearchBy(controller.ByTitle)
			case "Author":
				f.controller.BookList.SetSearchBy(controller.ByAuthor)
			case "Genre":
				f.controller.BookList.SetSearchBy(controller.ByGenre)
			case "Borrower":
				f.controller.BookList.SetSearchBy(controller.ByBorrower)
			default:
				panic("search by not found")
			}
		f.emiter.Emit(OnSearchBy)
		},
	)
	selectSearchBy.PlaceHolder = "Search By"

	// Search Entry
	//--------------
	searchEnt := widget.NewEntry()
	searchEnt.PlaceHolder = "Search"
	searchEnt.OnChanged = func(s string) {
		f.controller.BookList.SetSearch(s)
		f.emiter.Emit(OnSearch)
	}

	// Next Item
	//-----------
	onNextItem := func() {
		f.emiter.Emit(OnSelectNext)
	}

	nextItemItem := &widget.ToolbarAction{
		Icon: theme.MoveDownIcon() ,
		OnActivated: onNextItem,
	}

	// Previous Item
	// --------------
	onPrevItem := func() {
		f.emiter.Emit(OnSelectPrev)
	}

	prevItemItem := &widget.ToolbarAction{
		Icon: theme.MoveUpIcon(),
		OnActivated: onPrevItem,
	}


	// Toolbar Items
	//---------------
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
		nextItemItem,
		prevItemItem,
	}
	toolBar := widget.NewToolbar(items...)
	

	// Canvas
	//--------
	boxes := []fyne.CanvasObject{
		toolBar,
		searchEnt,
	}
	o := container.New(layout.NewGridLayout(len(boxes)), boxes...)
	right := widget.NewSeparator()
	right.Hide()
	return container.New(layout.NewBorderLayout(nil, nil, o, nil), o, selectSearchBy)
}
