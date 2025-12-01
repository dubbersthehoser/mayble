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
)

func (f *FunkView) TopBar() fyne.CanvasObject {

	// Save
	//------
	OnSaveItem := func() {
		f.emiter.Emit(OnSave, nil)
	}

	saveItem := &widget.ToolbarAction{
		Icon: theme.DocumentSaveIcon(),
		OnActivated: OnSaveItem,
	}

	saveItem.Disable()

	// Events
	EnableSaveItem := func(_ any) {
		saveItem.Enable()
	}
	DisableSaveItem := func(_ any) {
		saveItem.Disable()
	}

	f.emiter.OnEvent(OnModification, EnableSaveItem)
	f.emiter.OnEvent(OnSave, DisableSaveItem)


	// Menu
	//------
	OnMenuItem := func() {
		f.emiter.Emit(OnMenuOpen, nil)
	}
	menuItem := &widget.ToolbarAction{
		Icon: theme.MenuIcon(),
		OnActivated: OnMenuItem,
	}



	// Create 
	//--------
	OnCreateItem := func() {
		f.emiter.Emit(OnCreate, nil)
	}
	createItem := &widget.ToolbarAction{
		Icon: theme.ContentAddIcon(),
		OnActivated: OnCreateItem,
	}


	// Update
	//--------
	OnUpdateItem := func() {
		f.emiter.Emit(OnUpdate, nil)
	}
	updateItem := &widget.ToolbarAction{
		Icon: theme.DocumentCreateIcon(),
		OnActivated: OnUpdateItem,
	}

	// Delete
	//--------
	OnDeleteItem := func() {
		f.emiter.Emit(OnDelete, nil)
	}
	deleteItem := &widget.ToolbarAction{
		Icon: theme.DeleteIcon(),
		OnActivated: OnDeleteItem, 
	}


	// Events
	DisableItemOnMod := func(_ any) {
		updateItem.Disable()
		deleteItem.Disable()
	}
	EnableItemOnMod := func(_ any) {
		updateItem.Enable()
		deleteItem.Enable()
	}
	updateItem.Disable()
	deleteItem.Disable()
	
	f.emiter.OnEvent(OnSelected, EnableItemOnMod)
	f.emiter.OnEvent(OnUnselected, DisableItemOnMod)

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
		f.emiter.Emit(OnRedo, nil)
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

	f.emiter.OnEvent(OnModification, func(_ any) {
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
			f.emiter.Emit(OnSearchBy, s)
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

	searchEnt.OnSubmitted = func(s string) {
		f.emiter.Emit(OnSelectNext, nil)
	}

	resetSearchText := func(_ any) {
		searchEnt.SetText("")
	}
	f.emiter.OnEvent(OnSearchBy, resetSearchText)

	// Next Item
	//-----------
	onNextItem := func() {
		f.emiter.Emit(OnSelectNext, nil)
	}

	nextItemItem := &widget.ToolbarAction{
		Icon: theme.MoveDownIcon() ,
		OnActivated: onNextItem,
	}

	// Previous Item
	// --------------
	onPrevItem := func() {
		f.emiter.Emit(OnSelectPrev, nil)
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
