package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"

	_"github.com/dubbersthehoser/mayble/internal/gui/controller"
)

func (f *FunkView) TopBar() fyne.CanvasObject {


	// Save Button
	saveBtn := widget.NewButton(
		"Save",
		func() {
			f.emiter.Emit(OnSave)
			//f.Save()
		},
	)
	saveBtn.Importance = widget.HighImportance

	f.emiter.On(OnSave, 
		func() {
			saveBtn.Importance = widget.MediumImportance
			saveBtn.Refresh()
		},
	)

	// TODO Need logic for search.

	// SearchBy Select
	selectSearchBy := widget.NewSelect(
		[]string{"Title", "Author", "Genre"},
		func(s string) {
			fmt.Println("search by not implemeted")
		},
	)
	selectSearchBy.PlaceHolder = "Search By"

	// Search Entry
	searchEnt := widget.NewEntry()
	searchEnt.PlaceHolder = "Search"
	searchEnt.OnChanged = func(s string) {
		fmt.Println("search not implemeted")
	}



	// Add Book Button
	bookBtn := widget.NewButton(
		"New",
		func() {
			fmt.Println("Create button pressed")
		},
	)

	// Delete Book Button
	OnDelete := func() {
		fmt.Println("Delete button pressed")
	}
	deleteBtn := widget.NewButton("Delete", OnDelete)
	deleteBtn.Disable()

	// Update Book Button
	OnUpdate := func() {
		fmt.Println("Update button pressed")
	}
	updateBtn := widget.NewButton("Edit", OnUpdate)
	updateBtn.Disable()


	// On Book Selection
	OnBookSelected := func() {
		updateBtn.Enable()
		//deleteBtn.Enable()
	}
	OnBookUnselected := func() {
		updateBtn.Disable()
		//deleteBtn.Disable()
	}
	
	f.emiter.On(OnSelected, OnBookSelected)
	f.emiter.On(OnUnselected, OnBookUnselected)
	
	// Box
	boxes := []fyne.CanvasObject{
		saveBtn,
		bookBtn,
		updateBtn,
		selectSearchBy,
		//container.New(layout.NewStackLayout(), searchEnt),
		//selectOrderBy,
	}

	o := container.New(layout.NewGridLayout(len(boxes)), boxes...)
	right := widget.NewSeparator()
	right.Hide()
	return container.New(layout.NewBorderLayout(nil, nil, o, right), o, searchEnt, right)
}
