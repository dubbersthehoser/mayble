package gui

import (
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
	NewBookSortingSelected string = "NewSortingSelected"
)

func (u *UI) NewBookTableComp() fyne.CanvasObject {

	// Top Labels
	newSortByButtonAndLabel := func(label string) fyne.CanvasObject {
		labelAsc := "↑ " + label
		labelDesc := "↓ " + label
		labelNormal := "- " + label

		button := widget.NewButton(
			labelNormal, 
			nil,
		)
		button.Importance = widget.LowImportance
		button.OnTapped = func() {
			u.VM.ChangeOrderBy(label)
			if button.Text == labelNormal {
				u.Emiter.Emit(NewBookSortingSelected, label)
			}
			if button.Text != labelDesc {
				button.SetText(labelDesc)
				u.VM.SetOrderingDESC()
			} else {
				button.SetText(labelAsc)
				u.VM.SetOrderingASC()
			}
			button.Refresh()
		}
		OnNewOrdring := func(data any) {
			button.SetText(labelNormal)
		}
		u.Emiter.On(NewBookSortingSelected, OnNewOrdring)
		return button
	}

	fields := []fyne.CanvasObject{
		newSortByButtonAndLabel("Title"),
		newSortByButtonAndLabel("Author"),
		newSortByButtonAndLabel("Genre"),
		newSortByButtonAndLabel("Ratting"),
		newSortByButtonAndLabel("On Loan"),
	}

	fields[0].(*widget.Button).OnTapped() /* COULD CAUSE UNWANTED STATE CHANGE */

	Heading := container.New(layout.NewGridLayout(len(fields)), fields...)

	OnListLength := func() int {
		return len(u.VM.BookList)
	}

	OnCanvasCreation := func() fyne.CanvasObject {
		titleLabel := widget.NewLabel("")
		authorLabel := widget.NewLabel("")
		genreLabel := widget.NewLabel("")
		rattingLabel := widget.NewLabel("")

		titleLabel.Wrapping = fyne.TextWrapWord
		authorLabel.Wrapping = fyne.TextWrapWord
		genreLabel.Wrapping = fyne.TextWrapWord
		rattingLabel.Wrapping = fyne.TextWrapWord

		titleLabel.Truncation = fyne.TextTruncateEllipsis
		authorLabel.Truncation = fyne.TextTruncateEllipsis
		genreLabel.Truncation = fyne.TextTruncateEllipsis
		rattingLabel.Truncation = fyne.TextTruncateEllipsis

		titleLabel.Selectable = true
		authorLabel.Selectable = true
		genreLabel.Selectable = true
		rattingLabel.Selectable = true

		fields := []fyne.CanvasObject{
			titleLabel,
			authorLabel,
			genreLabel,
			rattingLabel,
		}
		entry := container.New(layout.NewGridLayout(len(fields)), fields...)
		return entry
	}   

	OnCanvasInit := func(id int, o fyne.CanvasObject) {
		index := id
		o.(*fyne.Container).Objects[0].(*widget.Label).SetText(u.VM.BookList[index].Title())
		o.(*fyne.Container).Objects[1].(*widget.Label).SetText(u.VM.BookList[index].Author())
		o.(*fyne.Container).Objects[2].(*widget.Label).SetText(u.VM.BookList[index].Genre())
		o.(*fyne.Container).Objects[3].(*widget.Label).SetText(u.VM.BookList[index].Ratting())
	} 

	// The List of Books
	List := widget.NewList(OnListLength, OnCanvasCreation, OnCanvasInit)
	List.HideSeparators = false
	List.OnSelected = func(id int) {
		u.VM.SetBookSelected(id)
	}
	List.OnUnselected = func(id int) {} // todo

	table := container.New(layout.NewBorderLayout(Heading, nil, nil, nil), Heading, List)
	return table
}
