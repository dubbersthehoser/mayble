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

	"github.com/dubbersthehoser/mayble/internal/listing"
)

func (f *FunkView) Table() fyne.CanvasObject {

	/*************************
		Table Header
	**************************/

	// the labelResets, newLabelReset, and resetHeaderLabel is for helping with switching button text.
	labelResets := make([]func(), 0)
	newLabelReset := func(button *widget.Button, label string) func() {
		return func() {
			button.SetText("- " + label)
		}
	} 
	resetHeaderLabels := func() {
		for _, fn := range labelResets {
			fn()
		}
	}
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
			if button.Text == labelNormal {
				resetHeaderLabels()
			}
			if button.Text != labelDesc {
				button.SetText(labelDesc)
				f.controller.BookList.SetOrdering(listing.DEC)
			} else {
				button.SetText(labelAsc)
				f.controller.BookList.SetOrdering(listing.ASC)
			}
			button.Refresh()
			by := listing.MustStringToOrderBy(label)
			f.controller.BookList.SetOrderBy(by)
			f.emiter.Emit(OnSort)
		}
		//f.emiter.On(OnModification, func() {
		//	button.SetText(labelNormal)
		//})
		return button
	}
	fields := make([]fyne.CanvasObject, 0)
	labels := listing.SortByList()
	for _, label := range labels {
		obj := newSortByButtonAndLabel(label)
		labelResets = append(labelResets, newLabelReset(obj.(*widget.Button), label))
		fields = append(fields, obj)
	}

	// setting the the title field to be sorted
	//fields[0].(*widget.Button).OnTapped() /* COULD CAUSE UNWANTED STATE CHANGE */

	Heading := container.New(layout.NewGridLayout(len(fields)), fields...)


	/***************************
		List's Methods
	****************************/
	OnListLength := func() int {
		n := f.controller.BookList.Len()
		return n
	}

	OnCanvasCreation := func() fyne.CanvasObject {

		titleLabel := widget.NewLabel("")
		authorLabel := widget.NewLabel("")
		genreLabel := widget.NewLabel("")
		rattingLabel := widget.NewLabel("")
		onLoanName := widget.NewLabel("")
		onLoanDate := widget.NewLabel("")

		titleLabel.Wrapping = fyne.TextWrapWord
		authorLabel.Wrapping = fyne.TextWrapWord
		genreLabel.Wrapping = fyne.TextWrapWord
		rattingLabel.Wrapping = fyne.TextWrapWord
		onLoanName.Wrapping = fyne.TextWrapWord
		onLoanDate.Wrapping = fyne.TextWrapWord
		

		titleLabel.Truncation = fyne.TextTruncateEllipsis
		authorLabel.Truncation = fyne.TextTruncateEllipsis
		genreLabel.Truncation = fyne.TextTruncateEllipsis
		rattingLabel.Truncation = fyne.TextTruncateEllipsis
		onLoanName.Truncation = fyne.TextTruncateEllipsis
		onLoanDate.Truncation = fyne.TextTruncateEllipsis

		titleLabel.Selectable = false
		authorLabel.Selectable = false
		genreLabel.Selectable = false
		rattingLabel.Selectable = false
		onLoanName.Selectable = false
		onLoanDate.Selectable = false

		fields := []fyne.CanvasObject{
			titleLabel,
			authorLabel,
			genreLabel,
			rattingLabel,
			onLoanName,
			onLoanDate,
		}
		entry := container.New(layout.NewGridLayout(len(fields)), fields...)
		return entry
	}   
	OnCanvasInit := func(index int, o fyne.CanvasObject) {
		book, err := f.controller.BookList.Get(index)
		if err != nil {
			f.displayError(err)
			return
		}
		o.(*fyne.Container).Objects[0].(*widget.Label).SetText(book.Title)
		o.(*fyne.Container).Objects[1].(*widget.Label).SetText(book.Author)
		o.(*fyne.Container).Objects[2].(*widget.Label).SetText(book.Genre)
		o.(*fyne.Container).Objects[3].(*widget.Label).SetText(book.Ratting)
		o.(*fyne.Container).Objects[4].(*widget.Label).SetText(book.Borrower)
		o.(*fyne.Container).Objects[5].(*widget.Label).SetText(book.Date)
	} 

	OnSelect := func(index int) {
		fmt.Println("Selected book")
		err := f.controller.BookList.Select(index)
		if err != nil {
			f.displayError(err)
			return
		}
		f.emiter.Emit(OnSelected)
	}
	OnUnselect := func(index int) {
		fmt.Println("Unselected book")
		f.controller.BookList.Unselect()
		f.emiter.Emit(OnUnselected)
	}

	/*******************
		List
	********************/
	List := widget.NewList(OnListLength, OnCanvasCreation, OnCanvasInit)
	List.HideSeparators = false
	List.OnSelected = OnSelect
	List.OnUnselected = OnUnselect

	listOnModification := func() {
		List.UnselectAll()
	}

	listOnSearch := func() {
		fmt.Println("list on search")
		idx := f.controller.BookList.SelectedIndex
		List.Select(widget.ListItemID(idx))
	}
	listOnSelectNext := listOnSearch
	listOnSelectPrev := listOnSearch

	f.emiter.On(OnModification, listOnModification)
	f.emiter.On(OnSearch, listOnSearch)
	f.emiter.On(OnSelectNext, listOnSelectNext)
	f.emiter.On(OnSelectPrev, listOnSelectPrev)


	// Table
	table := container.New(layout.NewBorderLayout(Heading, nil, nil, nil), Heading, List)
	return table
}
