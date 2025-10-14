package view

import (
	_"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	_"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"

	"github.com/dubbersthehoser/mayble/internal/gui/controller"
)



const DialogWidth float32 = 400

func (f *FunkView) BookEdit() (fyne.CanvasObject, error) {


	var builder *controller.BookLoanBuilder

	if f.controller.BookList.IsSelected() {
		bookLoan, err := f.controller.BookList.Selected()
		println(err)
		if err != nil {
			return nil, err
		}
		builder = controller.NewBuilderWithBookLoan(bookLoan)
	} else {
		builder = controller.NewBookLoanBuilder()
	}

	rattings := controller.GetRattingStrings()

	type EntryField struct { // container for entry form items
		Entry fyne.CanvasObject
		Label fyne.CanvasObject
	}

	ValidationInfo := widget.NewLabel("")

	titleField := EntryField{
		Entry: widget.NewEntry(),
		Label: widget.NewLabel("Title"),
	}
	authorField := EntryField{
		Entry: widget.NewEntry(),
		Label: widget.NewLabel("Author"),
	}
	genreField := EntryField{
		Entry: widget.NewSelectEntry([]string{"placeholder_1", "placeholder_2", "placeholder_2"}),
		Label: widget.NewLabel("Genre"),
	}
	rattingField := EntryField{
		Entry: widget.NewSelect(rattings, nil),
		Label: widget.NewLabel("Ratting"),
	}
	onLoanField := EntryField{
		Entry: widget.NewCheck("", nil),
		Label: widget.NewLabel("On Loan"),
	}
	loanNameField := EntryField{
		Entry: widget.NewEntry(),
		Label: widget.NewLabel("Borrower"),
	}
	loanDateField := EntryField{
		Entry: widget.NewDateEntry(),
		Label: widget.NewLabel("Date"),
	}

	{ // Set data in the entries for the current edit status of the form
		var (
			isBeingUpdated  bool  = builder.Type == controller.Updating
			isBeingCreated  bool  = builder.Type == controller.Creating
			isBeingLoaned   bool  = builder.IsOnLoan
		)

		if isBeingCreated {
			rattingField.Entry.(*widget.Select).SetSelectedIndex(0) 
		}
		
		if isBeingUpdated {
			titleField.Entry.(*widget.Entry).Text = builder.Title
			authorField.Entry.(*widget.Entry).Text = builder.Author
			genreField.Entry.(*widget.SelectEntry).Text = builder.Genre
			rattingField.Entry.(*widget.Select).Selected = rattings[builder.Ratting]
		}
		if isBeingLoaned && isBeingUpdated {
			loanNameField.Entry.(*widget.Entry).Text = builder.Borrower
			loanDateField.Entry.(*widget.DateEntry).Date = &builder.Date
		}
	}

	onCancel := func() {}

	onSubmit := func() {
		f.controller.BookEditor.Submit(builder)
		f.Update()
	}

	submitBtn := widget.NewButton("Submit", onSubmit)
	cancelBtn := widget.NewButton("Cancel", onCancel)

	validateData := builder.Validate

	onDataChange := func() {
		err := validateData()
		if err != nil {
			submitBtn.Disable()
			ValidationInfo.SetText(err.Error())
			ValidationInfo.Importance = widget.DangerImportance
			ValidationInfo.Refresh()
		} else {
			submitBtn.Enable()
			ValidationInfo.SetText("")
			ValidationInfo.Importance = widget.MediumImportance
			ValidationInfo.Refresh()
		}
	}

	titleField.Entry.(*widget.Entry).OnChanged = func(s string) {
		builder.SetTitle(s)
		onDataChange()
	}
	authorField.Entry.(*widget.Entry).OnChanged = func(s string) {
		builder.SetAuthor(s)
		onDataChange()
	}
	genreField.Entry.(*widget.SelectEntry).OnChanged = func(s string) {
		builder.SetGenre(s)
		onDataChange()
	}
	rattingField.Entry.(*widget.Select).OnChanged = func(s string) {
		builder.SetRattingAsString(s)
		onDataChange()
	}
	loanNameField.Entry.(*widget.Entry).OnChanged = func(s string) {
		builder.SetBorrower(s)
		onDataChange()
	}
	loanDateField.Entry.(*widget.DateEntry).OnChanged = func(d *time.Time) {
		builder.SetDate(d)
		onDataChange()
	}

	authorField.Entry.(*widget.Entry).Validator = nil
	titleField.Entry.(*widget.Entry).Validator = nil
	loanNameField.Entry.(*widget.Entry).Validator = nil

	onLoanCheck := onLoanField.Entry.(*widget.Check)

	formItems := []fyne.CanvasObject{
		titleField.Label,    titleField.Entry,
		authorField.Label,   authorField.Entry,
		genreField.Label,    genreField.Entry,
		rattingField.Label,  rattingField.Entry,
		onLoanField.Label,   onLoanField.Entry,
		loanNameField.Label, loanNameField.Entry,
		loanDateField.Label, loanDateField.Entry,
	}

	formlayout := container.New(layout.NewFormLayout(), formItems...)

	onLoanCheck.OnChanged = func(isChecked bool) {
		NameEntry := loanNameField.Entry.(*widget.Entry)
		NameLabel := loanNameField.Label.(*widget.Label)
		DateEntry := loanDateField.Entry.(*widget.DateEntry)
		DateLabel := loanDateField.Label.(*widget.Label)
		if isChecked {
			NameEntry.Enable()
			DateEntry.Enable()
			
			NameLabel.Importance = widget.MediumImportance
			NameLabel.Refresh()
			DateLabel.Importance = widget.MediumImportance
			DateLabel.Refresh()
		} else {
			NameEntry.Disable()
			DateEntry.Disable()

			NameLabel.Importance = widget.LowImportance
			NameLabel.Refresh()
			DateLabel.Importance = widget.LowImportance
			DateLabel.Refresh()
		}
		builder.SetIsOnLoan(isChecked)
		onDataChange()
		onLoanCheck.Refresh()
	}
	onLoanCheck.OnChanged(builder.IsOnLoan)

	obj := container.New(layout.NewVBoxLayout(), formlayout, ValidationInfo, submitBtn, cancelBtn)

	return obj, nil
}
