package gui

// Main App Header

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"

	_"github.com/dubbersthehoser/mayble/internal/controler"
)

const (
	ChangedOrderBy string = "ChangedOrderBy"
	ChangedSearch         = "ChangedSearch"
	ChangedSearchBy       = "ChangedSearchBy"
	SaveButtonClicked     = "SaveButtonClicked"
)


func (u *UI) NewHeaderComp() fyne.CanvasObject {

	// Save Button
	saveBtn := widget.NewButton(
		"Save",
		func() {
			u.Emiter.Emit(SaveButtonClicked, nil)
		},
	)
	saveBtn.Importance = widget.HighImportance
	u.Emiter.On(
		SaveButtonClicked, 
		func(data any) {
			saveBtn.Importance = widget.MediumImportance
			saveBtn.Refresh()
		},
	)


	// OrderBy Select
	selectOrderBy := widget.NewSelect(
		[]string{"Title", "Author", "Genre", "Ratting"},
		func(s string) {
			u.Emiter.Emit(ChangedOrderBy, s)
		},
	)
	selectOrderBy.PlaceHolder = "Order By"

	// SearchBy Select
	selectSearchBy := widget.NewSelect(
		[]string{"Title", "Author", "Genre"},
		func(s string) {
			u.Emiter.Emit(ChangedSearchBy, s)
		},
	)
	selectSearchBy.PlaceHolder = "Search By"

	// Search Entry
	searchEnt := widget.NewEntry()
	searchEnt.PlaceHolder = "Search"
	searchEnt.OnChanged = func(s string) {
		u.Emiter.Emit(ChangedSearch, s)
	}


	// Add Book Button
	bookBtn := widget.NewButton(
		"New",
		func() {
			book := u.VM.NewBook()
			u.NewBookDialog(book).Show()
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
		book := u.VM.SelectedBook
		u.UpdateBookDialog(book).Show()
	}
	updateBtn := widget.NewButton("Edit", OnUpdate)
	updateBtn.Disable()

	// On Book Selection
	OnBookSelected := func(nothing any) {
		updateBtn.Enable()
		//deleteBtn.Enable()
	}
	OnBookUnselected := func(nothing any) {
		updateBtn.Disable()
		//deleteBtn.Disable()
	}
	
	u.Emiter.On(BookSelected, OnBookSelected)
	u.Emiter.On(BookUnselected, OnBookUnselected)
	
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
