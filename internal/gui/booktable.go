package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"
)

func (u *UIState) NewBookTableComp() fyne.CanvasObject {

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
			u.BookOrderedBy = label
			if button.Text != labelDesc {
				button.SetText(labelDesc)
				u.Emiter.Emit(ChangedOrderByDesc, label)
			} else {
				button.SetText(labelAsc)
				u.Emiter.Emit(ChangedOrderByAsc, label)
			}
			u.Emiter.Emit(ChangedOrderBy, label)
			button.Refresh()
		}
		u.Emiter.On(ChangedOrderBy, func(data any) {
			if data.(string) != label {
				button.SetText(labelNormal)
				button.Refresh()
			}
		})
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

	// The List of Books
	List := widget.NewList(
		func() int { // Length
			return len(u.BookList)
		},
		func() fyne.CanvasObject { // CreateItem

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
		},          
		func(id int, o fyne.CanvasObject) { // UpdateItem
			index := id
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(u.BookList[index])
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText(u.BookList[index])
			o.(*fyne.Container).Objects[2].(*widget.Label).SetText(u.BookList[index])
			o.(*fyne.Container).Objects[3].(*widget.Label).SetText(u.BookList[index])
		})
	List.HideSeparators = false
	List.OnSelected = func(id int) {
		u.BookSelected = id
	}
	List.OnUnselected = func(id int) {} // todo

	table := container.New(layout.NewBorderLayout(Heading, nil, nil, nil), Heading, List)
	return table
}
