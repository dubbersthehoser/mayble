package gui

import (
	"fmt"
	"time"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"
)

const DialogWidth float32 = 400

func GetDialogSize(size fyne.Size) fyne.Size {
	height := size.Height
	s := fyne.NewSize(DialogWidth, 0)
	s.Height = height
	return s
}

func (u *UI) UpdateBookDialog(listed *ListedBook) dialog.Dialog {

	label := "Edit Book"
	
	rattings := GetRattingStrings()

	titleEntry := widget.NewEntry()
	authorEntry := widget.NewEntry()
	genreEntry := widget.NewSelectEntry(u.VM.UniqueGenres())
	rattingSelect := widget.NewSelect(rattings, nil)
	loanNameEntry := widget.NewEntry()
	loanDateEntry := widget.NewDateEntry()

	titleEntry.Bind(listed.Entry.Title)
	authorEntry.Bind(listed.Entry.Author)
	genreEntry.Bind(listed.Entry.Genre)
	rattingSelect.Bind(listed.Entry.Ratting)
	loanNameEntry.Bind(listed.Entry.LoanName)
	loanDateEntry.SetDate(listed.Entry.LoanDate)

	authorEntry.Validator = listed.Entry.AuthorValidator
	titleEntry.Validator = listed.Entry.TitleValidator
	loanNameEntry.Validator = nil

	onLoanCheck := widget.NewCheck("", nil)
	onLoanCheck.Checked = false
	onLoanCheck.Refresh()
	loanNameEntry.Disable()
	loanDateEntry.Disable()

	items := []*widget.FormItem{
		widget.NewFormItem("Title*", titleEntry),
		widget.NewFormItem("Author*", authorEntry),
		widget.NewFormItem("Genre", genreEntry),
		widget.NewFormItem("Ratting", rattingSelect),
		widget.NewFormItem("On Loan", onLoanCheck),
		widget.NewFormItem("Name*", loanNameEntry),
		widget.NewFormItem("Date*", loanDateEntry),
	}
	form := widget.NewForm(items...)

	obj := container.New(layout.NewVBoxLayout(), form)
	d := dialog.NewCustomWithoutButtons(label, obj, u.Window)


	onLoanCheck.OnChanged = func(isChecked bool) {
		if isChecked {
			loanNameEntry.Validator = listed.Entry.LoanNameValidator
			loanNameEntry.Enable()
			loanDateEntry.Enable()
		} else {
			loanNameEntry.Validator = func(s string) error {return nil}
			loanNameEntry.Disable()
			loanDateEntry.Disable()
		}
		onLoanCheck.Refresh()
		form.Validate()
		d.Refresh()
	}

	if listed.IsOnLoan() {
		onLoanCheck.Checked = true
		onLoanCheck.OnChanged(true)
	}


	submitValidation := func() error {
		var (
			LoanDateInvalid bool = *loanDateEntry.Date == *new(time.Time)
		)
		if onLoanCheck.Checked && LoanDateInvalid {
			return errors.New("Invalid Loan Date")
		}
		return nil
	}

	
	// Buttons
	OnSubmit := func() {
		if onLoanCheck.Checked == false {
			listed.Entry.UnsetLoan()
		}
		if err := submitValidation(); err != nil {
			dialog.ShowError(err, u.Window)
			return 
		}
		listed.Entry.LoanDate = loanDateEntry.Date
		listed.EntryToData()
		listed.DataToLabel()
		u.Emiter.Emit(UpdatedBookToList, listed)
		d.Dismiss()
	}

	OnCancel := func() {
		fmt.Println("good buy world")
		d.Dismiss()
	}

	OnDelete := func() {
		fmt.Println("removed JK")
		d.Dismiss()
	}

	SubmitBtn := widget.NewButton("Submit", OnSubmit)
	CancelBtn := widget.NewButton("Cancel", OnCancel)
	DeleteBtn := widget.NewButton("Delete", OnDelete)

	SubmitBtn.Importance = widget.HighImportance
	DeleteBtn.Importance = widget.WarningImportance

	OnValidationChanged := func(err error) {
		if err != nil {
			SubmitBtn.Disable()
		} else {
			SubmitBtn.Enable()
		}
	}

	form.SetOnValidationChanged(OnValidationChanged)

	btns := []fyne.CanvasObject{
		SubmitBtn,
		CancelBtn,
		DeleteBtn,

	}

	d.SetButtons(btns)
	d.Resize(GetDialogSize(d.MinSize()))
	return d
}

func (u *UI) NewBookDialog(listed *ListedBook) dialog.Dialog {

	label := "New Book"
	
	rattings := GetRattingStrings()

	titleEntry := widget.NewEntry()
	authorEntry := widget.NewEntry()
	genreEntry := widget.NewSelectEntry(u.VM.UniqueGenres())
	rattingSelect := widget.NewSelect(rattings, nil)
	loanNameEntry := widget.NewEntry()
	loanDateEntry := widget.NewDateEntry()

	titleEntry.Bind(listed.Entry.Title)
	authorEntry.Bind(listed.Entry.Author)
	genreEntry.Bind(listed.Entry.Genre)
	rattingSelect.Bind(listed.Entry.Ratting)
	loanNameEntry.Bind(listed.Entry.LoanName)
	loanDateEntry.Date = listed.Entry.LoanDate

	authorEntry.Validator = listed.Entry.AuthorValidator
	titleEntry.Validator = listed.Entry.TitleValidator
	loanNameEntry.Validator = nil

	onLoanCheck := widget.NewCheck("", nil)
	onLoanCheck.Checked = false
	onLoanCheck.Refresh()
	loanNameEntry.Disable()
	loanDateEntry.Disable()

	items := []*widget.FormItem{
		widget.NewFormItem("Title*", titleEntry),
		widget.NewFormItem("Author*", authorEntry),
		widget.NewFormItem("Genre", genreEntry),
		widget.NewFormItem("Ratting", rattingSelect),
		widget.NewFormItem("On Loan", onLoanCheck),
		widget.NewFormItem("Name*", loanNameEntry),
		widget.NewFormItem("Date*", loanDateEntry),
	}
	form := widget.NewForm(items...)

	obj := container.New(layout.NewVBoxLayout(), form)
	d := dialog.NewCustomWithoutButtons(label, obj, u.Window)

	onLoanCheck.OnChanged = func(isChecked bool) {
		if isChecked {
			loanNameEntry.Validator = listed.Entry.LoanNameValidator
			loanNameEntry.Enable()
			loanDateEntry.Enable()
		} else {
			loanNameEntry.Validator = func(s string) error {return nil}
			loanNameEntry.Disable()
			loanDateEntry.Disable()
		}
		onLoanCheck.Refresh()
		form.Validate()
		d.Refresh()
	}

	submitValidation := func() error {
		var (
			LoanDateInvalid bool = *loanDateEntry.Date == *new(time.Time)
		)
		if onLoanCheck.Checked && LoanDateInvalid {
			return errors.New("Invalid Loan Date")
		}
		return nil
	}

	
	// Buttons
	OnSubmit := func() {
		if onLoanCheck.Checked == false {
			listed.Entry.UnsetLoan()
		}
		if err := submitValidation(); err != nil {
			dialog.ShowError(err, u.Window)
			return 
		}
		listed.Entry.LoanDate = loanDateEntry.Date
		listed.EntryToData()
		listed.DataToLabel()
		u.VM.AddBook(listed)
		u.Emiter.Emit(AddedNewBookToList, listed)
		d.Dismiss()
	}

	OnCancel := func() {
		fmt.Println("good buy world")
		d.Dismiss()
	}

	SubmitBtn := widget.NewButton("Submit", OnSubmit)
	CancelBtn := widget.NewButton("Cancel", OnCancel)

	SubmitBtn.Importance = widget.HighImportance

	OnValidationChanged := func(err error) {
		if err != nil {
			SubmitBtn.Disable()
		} else {
			SubmitBtn.Enable()
		}
	}
	form.SetOnValidationChanged(OnValidationChanged)

	btns := []fyne.CanvasObject{
		SubmitBtn,
		CancelBtn,
	}

	d.SetButtons(btns)
	d.Resize(GetDialogSize(d.MinSize()))
	return d
}
