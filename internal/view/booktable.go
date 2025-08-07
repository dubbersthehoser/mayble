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
			if button.Text == labelNormal {
				u.Emiter.Emit(SelectedBookSorting, label)
			}
			if button.Text != labelDesc {
				button.SetText(labelDesc)
				u.VM.SetSorting(label, DESC)
			} else {
				button.SetText(labelAsc)
				u.VM.SetSorting(label, ASC)
			}
			button.Refresh()
		}
		OnNewOrdring := func(data any) {
			button.SetText(labelNormal)
		}
		u.Emiter.On(SelectedBookSorting, OnNewOrdring)
		return button
	}

	fields := []fyne.CanvasObject{
		newSortByButtonAndLabel("Title"),
		newSortByButtonAndLabel("Author"),
		newSortByButtonAndLabel("Genre"),
		newSortByButtonAndLabel("Ratting"),
		newSortByButtonAndLabel("Loan To"),
		newSortByButtonAndLabel("Date"),
	}

	// setting the the title feild to be sorted
	fields[0].(*widget.Button).OnTapped() /* COULD CAUSE UNWANTED STATE CHANGE */

	Heading := container.New(layout.NewGridLayout(len(fields)), fields...)

	OnListLength := func() int {
		return u.VM.ListedCount()
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
		book := u.VM.GetBook(index)
		o.(*fyne.Container).Objects[0].(*widget.Label).Bind(book.Label.Title)
		o.(*fyne.Container).Objects[1].(*widget.Label).Bind(book.Label.Author)
		o.(*fyne.Container).Objects[2].(*widget.Label).Bind(book.Label.Genre)
		o.(*fyne.Container).Objects[3].(*widget.Label).Bind(book.Label.Ratting)
		o.(*fyne.Container).Objects[4].(*widget.Label).Bind(book.Label.LoanName)
		o.(*fyne.Container).Objects[5].(*widget.Label).Bind(book.Label.LoanDate)
	} 

	// The List of Books
	List := widget.NewList(OnListLength, OnCanvasCreation, OnCanvasInit)
	List.HideSeparators = false
	List.OnSelected = func(index int) {
		u.VM.SetSelectedBook(index)
		u.Emiter.Emit(BookSelected, index)
	}
	List.OnUnselected = func(index int) {
		u.Emiter.Emit(BookUnselected, index)
	}

	OnNewBookToTabel := func(book any) {
		u.VM.Sort()
		List.Refresh()
	}
	OnRemoveBookFromTabel := func(book any) {
		u.VM.Sort()
		List.Refresh()
	}
	OnUpdatedBookToTabel := func(book any) {
		u.VM.Sort()
		List.Refresh()
	}

	u.Emiter.On(AddedNewBookToList, OnNewBookToTabel)
	u.Emiter.On(RemovedBookFromList, OnRemoveBookFromTabel)
	u.Emiter.On(UpdatedBookToList, OnUpdatedBookToTabel)

	table := container.New(layout.NewBorderLayout(Heading, nil, nil, nil), Heading, List)
	return table
}
