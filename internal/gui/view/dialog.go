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

func BasicStringValidator(s string) error {
	if len(s) == 0 {
		return errors.New("empty string")
	}
	return nil
}

type EntryControl struct {
	Label      string
	Binding   binding.String
	Validator func(string) error
}
func (t *EntryControl) GetValidator() func(string) error {
	baseValidator := func(s string) error {
		err := BasicStringValidator(s)
		if err != nil {
			errors.New(t.Label + "is an" + err)
		}
		return nil
	}
	if t.Validator == nil {
		return baseValidator
	}
	return t.Validator
}

type SelectControl struct {
	Label     string
	Binding   binding.String
	Selection func() []string
}
func (s *SelectControl) GetSelection() []string {
	return s.Selection()
}

type EntrySelectControl struct {
	EntryControl
	Selection func() []string
}
func (e *EntrySelectControl) GetSelection() []string {
	return e.Selection()
}

func NewSelectControl(selection func() []string) *EntrySelectControl {
	
}


func GetDialogSize(size fyne.Size) fyne.Size {
	height := size.Height
	s := fyne.NewSize(DialogWidth, 0)
	s.Height = height
	return s
}

type BookFormControl struct {
	Title    binding.String
	Author   binding.String
	Genre    binding.String
	Ratting  binding.String
	LoanName binding.String
	LoanDate *time.Time

	TitleEntry    widget.Entry
	AuthorEntry   widget.Entry
	GenreEntry    widget.SelectEntry
	RattingSelect widget.Select

	OnLoanCheck   widget.Check
	LoanNameEntry widget.Entry
	LoanDateEntry widget.DateEntry
}

func (b *BookFormControl) OnDateChanged(t *time.Time) {  // for widget.DateEntry
	*b.LoanDate = *t
}

func (u *UI) GetBookFormDialog(form *BookFormControl) *dialog.CustomDialog {

	rattings := GetRattingLabels()

	form.TitleEntry := widget.NewEntry()
	form.AuthorEntry := widget.NewEntry()
	form.GenreEntry := widget.NewSelectEntry(u.VM.UniqueGenres())
	form.RattingSelect := widget.NewSelect(rattings, nil)
	form.LoanNameEntry := widget.NewEntry()
	from.LoanDateEntry := widget.NewDateEntry()

	titleEntry.Bind(form.Title)
	authorEntry.Bind(form.Author)
	genreEntry.Bind(form.Genre)
	rattingSelect.Bind(form.Ratting)
	loanNameEntry.Bind(form.LoanName)
	loanDateEntry.OnChanged = form.OnDateChanged

	authorEntry.Validator = nil
	titleEntry.Validator = nil
	loanNameEntry.Validator = nil

	items := []*widget.FormItem{
		widget.NewFormItem("Title*", titleEntry),
		widget.NewFormItem("Author*", authorEntry),
		widget.NewFormItem("Genre", genreEntry),
		widget.NewFormItem("Ratting", rattingSelect),
		widget.NewFormItem("On Loan", onLoanCheck),
		widget.NewFormItem("Name*", loanNameEntry),
		widget.NewFormItem("Date*", loanDateEntry),
	}

	form.OnLoanCheck.OnChanged = func(isChecked bool) {
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
	obj := container.New(layout.NewVBoxLayout(), form)
	d := dialog.NewCustomWithoutButtons(label, obj, u.Window)
	return d
}

func (u *UI) BookDialog(label string, book *EmitBookEntry) dialog.Dialog {

	ctrl := BookFormContol{}
	if book != nil {
		ctrl.Title = book.Title
		ctrl.Author = book.Author
		ctrl.Genre = book.Genre
		ctrl.Ratting = book.Ratting

		if book.IsOnLoan {
			
		}
	}
	d := u.GetBookFormDialog(&ctrl, book)

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
		if err := submitValidation(); err != nil {
			dialog.ShowError(err, u.Window)
			return 
		}
		book := &EventBookEntry{
			Title: titleEntry.Text,
			Author: authorEntry.Text,
			Genre: genreEntry.Text,
			Ratting: GenreEntry.Text,
			IsOnLoan: onLoanCheck.Checked,
		}
		u.Emiter.Emit(BookCreated, book)
		d.Dismiss()
	}

	OnCancel := func() {
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
