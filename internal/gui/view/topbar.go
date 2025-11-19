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
		f.emiter.Emit(OnMenuOpen)
	}
	menuItem := &widget.ToolbarAction{
		Icon: theme.MenuIcon(),
		OnActivated: OnMenuItem,
	}



	// Create 
	//--------
	OnCreateItem := func() {
		f.emiter.Emit(OnCreate)
	}
	createItem := &widget.ToolbarAction{
		Icon: theme.ContentAddIcon(),
		OnActivated: OnCreateItem,
	}


	// Update
	//--------
	OnUpdateItem := func() {
		f.emiter.Emit(OnUpdate)
	}
	updateItem := &widget.ToolbarAction{
		Icon: theme.DocumentCreateIcon(),
		OnActivated: OnUpdateItem,
	}

	// Delete
	//--------
	OnDeleteItem := func() {
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
		f.emiter.Emit(OnUndo, nil)
	}
	undoItem := &widget.ToolbarAction{
		Icon: theme.ContentUndoIcon(),
		OnActivated: OnUndoItem,
	}

	// Redo
	//------
	OnRedoItem := func() {
		f.emiter.Emit(OnRedo)
	}
	redoItem := &widget.ToolbarAction{
		Icon: theme.ContentRedoIcon(),
		OnActivated: OnRedoItem,
	}

	checkUndoBtn := func() {
		if f.controller.App.UndoIsEmpty() {
			undoItem.Enable()
		} else {
			undoItem.Disable()
		}
	}
	checkRedoBtn := func() {
		if f.controller.App.RedoIsEmpty() {
			redoItem.Enable()
		} else {
			redoItem.Disable()
		}
	}

	f.emiter.On(OnModification, func() {
		checkRedoBtn()
		checkUndoBtn()
	})

	checkRedoBtn()
	checkUndoBtn()



	// Search
	//---------
	//
	// Search By
	//-----------
	selectSearchBy := widget.NewSelect(
		[]string{"Title", "Author", "Genre", "Borrower"},
		func(s string) {
			f.emit(OnSearchBy, s)
		},
	)
	selectSearchBy.PlaceHolder = "Search By"

	// Search Entry
	//--------------
	searchEnt := widget.NewEntry()
	searchEnt.PlaceHolder = "Search"
	searchEnt.OnChanged = func(s string) {
		f.emiter.Emit(OnSearch, s)
	}

	searchEnt.OnSubmitted = func (s string) {
		f.emiter.Emit(OnSelectNext)
	}

	f.emiter.On(OnSort, func(_ any) {
		searchEnt.SetText("")
	})

	f.emiter.On(OnSearchBy, func(_ any)) {
		searchEnt.SetText("")
	})

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
