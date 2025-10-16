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

	_"github.com/dubbersthehoser/mayble/internal/gui/controller"
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

	// Events
	DisableUpdateItem := func() {
		updateItem.Disable()
	}
	EnableUpdateItem := func() {
		updateItem.Enable()
	}
	updateItem.Disable()
	
	f.emiter.On(OnSelected, EnableUpdateItem)
	f.emiter.On(OnUnselected, DisableUpdateItem)


	// Undo Redo
	//-----------
	//
	// TODO How to sync undo and redo buttons with core?
	//
	// Undo
	//------
	OnUndoItem := func() {
		fmt.Println("Undo button pressed")
	}
	undoItem := &widget.ToolbarAction{
		Icon: theme.ContentUndoIcon(),
		OnActivated: OnUndoItem,
	}

	// Redo
	//------
	OnRedoItem := func() {
		fmt.Println("Redo button pressed")
	}
	redoItem := &widget.ToolbarAction{
		Icon: theme.ContentRedoIcon(),
		OnActivated: OnRedoItem,
	}

	// Search
	//---------
	//
	// TODO Add search logic
	//
	// Search By
	//-----------
	selectSearchBy := widget.NewSelect(
		[]string{"Title", "Author", "Genre"},
		func(s string) {
			fmt.Println("search by not implemeted")
		},
	)
	selectSearchBy.PlaceHolder = "Search By"

	// Search Entry
	//--------------
	searchEnt := widget.NewEntry()
	searchEnt.PlaceHolder = "Search"
	searchEnt.OnChanged = func(s string) {
		fmt.Println("search not implemeted")
	}


	// Toolbar Items
	//---------------
	items := []widget.ToolbarItem{
		menuItem,
		widget.NewToolbarSeparator(),
		createItem,
		updateItem,
		saveItem,
		widget.NewToolbarSeparator(),
		undoItem,
		redoItem,
		widget.NewToolbarSeparator(),
	}
	toolBar := widget.NewToolbar(items...)
	

	// Canvas
	//--------
	boxes := []fyne.CanvasObject{
		toolBar,
		selectSearchBy,
	}
	o := container.New(layout.NewGridLayout(len(boxes)), boxes...)
	right := widget.NewSeparator()
	right.Hide()
	return container.New(layout.NewBorderLayout(nil, nil, o, right), o, searchEnt, right)
}
