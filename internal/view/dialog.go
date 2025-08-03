package gui

import (
	"fmt"
	//"time"
	//"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"

	"github.com/dubbersthehoser/mayble/internal/controler"
)

const DialogWidth float32 = 400

func GetDialogSize(size fyne.Size) fyne.Size {
	height := size.Height
	s := fyne.NewSize(DialogWidth, 0)
	s.Height = height
	return s
}

func (u *UI) NewOnLoanDialog(loan *controler.LoanVM) dialog.Dialog {

	dateEntry := widget.NewDateEntry()
	nameEntry := widget.NewEntry()

	nameData := binding.NewString()

	nameEntry.Bind(nameData)

	date := loan.Date()
	dateEntry.Date = &date

	nameEntry.Validator = loan.NameValidator

	d := dialog.NewForm("Loaned Book", "Add", "Cancel", 
		[]*widget.FormItem{
			widget.NewFormItem("Name*", nameEntry),
			widget.NewFormItem("Date*", dateEntry),
		}, 
		func(b bool){
			if b {
				loan.SetDate(*dateEntry.Date)
				loan.SetName(nameEntry.Text)
			}
		}, 
		u.Window,
	)
	d.Resize(GetDialogSize(d.MinSize()))
	return d
}

func (u *UI) UpdateBookDialog(book *controler.BookVM) dialog.Dialog {

	rattings := controler.GetRattingStrings()

	titleEntry := widget.NewEntry()
	authorEntry := widget.NewSelectEntry(u.VM.UniqueAuthors())
	genreEntry := widget.NewSelectEntry(u.VM.UniqueGenres())
	rattingSelect := widget.NewSelect(rattings, nil)

	titleData := binding.NewString()
	authorData := binding.NewString()
	genreData := binding.NewString()

	titleData.Set(book.Title())
	authorData.Set(book.Author())
	genreData.Set(book.Genre())

	titleEntry.Bind(titleData)
	authorEntry.Bind(authorData)
	genreEntry.Bind(genreData)

	authorEntry.Validator = book.AuthorValidator
	titleEntry.Validator = book.TitleValidator

	rattingSelect.PlaceHolder = rattings[0]
	rattingSelect.Selected = rattings[0]

	items := []*widget.FormItem{
		widget.NewFormItem("Title*", titleEntry),
		widget.NewFormItem("Author*", authorEntry),
		widget.NewFormItem("Genre", genreEntry),
		widget.NewFormItem( "Ratting", rattingSelect),
	}

	form := widget.NewForm(items...)

	obj := container.New(layout.NewVBoxLayout(), form)

	d := dialog.NewCustomWithoutButtons("Edit Book", obj, u.Window)

	// Buttons
	OnSubmit := func() { // Todo
		fmt.Println("hello world")
		u.Emiter.Emit(UpdatedBookToList, book)
		d.Dismiss()
	}

	OnCancel := func() { // Todo
		fmt.Println("good buy world")
		d.Dismiss()
	}

	OnDelete := func() { // Todo
		fmt.Println("removed ):")
		u.Emiter.Emit(RemovedBookFromList, book)
		d.Dismiss()
	}

	SubmitBtn := widget.NewButton("Submit", OnSubmit)
	CancelBtn := widget.NewButton("Cancel", OnCancel)
	DeleteBtn := widget.NewButton("Delete", OnDelete)

	DeleteBtn.Importance = widget.DangerImportance
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
		DeleteBtn,
	}
	d.SetButtons(btns)
	d.Resize(GetDialogSize(d.MinSize()))
	return d
	
}

func (u *UI) NewBookDialog(book *controler.BookVM) dialog.Dialog {

	rattings := controler.GetRattingStrings()

	titleEntry := widget.NewEntry()
	authorEntry := widget.NewSelectEntry(u.VM.UniqueAuthors())
	genreEntry := widget.NewSelectEntry(u.VM.UniqueGenres())
	rattingSelect := widget.NewSelect(rattings, nil)

	titleData := binding.NewString()
	authorData := binding.NewString()

	titleEntry.Bind(titleData)
	authorEntry.Bind(authorData)

	authorEntry.Validator = book.AuthorValidator
	titleEntry.Validator = book.TitleValidator

	rattingSelect.PlaceHolder = rattings[0]
	rattingSelect.Selected = rattings[0]

	onLoanCheck := widget.NewCheck(
		"", 
		nil,
	)

	OnLoanChecked := func() {
		u.NewOnLoanDialog(book.Loan()).Show()
		onLoanCheck.Checked = true
		onLoanCheck.Refresh()
	}
	OnLoanUnchecked := func() {
		onRemove := func(doRemove bool) {
			if doRemove {
				book.UnsetLoan()
			} else {
				onLoanCheck.Checked = true
				onLoanCheck.Refresh()
				
			}
		}
		dialog.ShowConfirm("Remove Loan", "Are you sure?", onRemove, u.Window)
	}

	onLoanCheck.OnChanged = func(isChecked bool) {
		if isChecked {
			OnLoanChecked()
		} else {
			OnLoanUnchecked()
		}
	}

	items := []*widget.FormItem{
		widget.NewFormItem(
			"Title*", 
			titleEntry,
		),
		widget.NewFormItem(
			"Author*", 
			authorEntry,
		),
		widget.NewFormItem(
			"Genre", 
			genreEntry,
		),
		widget.NewFormItem(
			"Ratting", 
			rattingSelect,
		),
		widget.NewFormItem(
			"On Loan",
			onLoanCheck,
		),
	}

	OnConfirm := func(confirmed bool) {
		if confirmed {
			book.SetTitle(titleEntry.Text)
			book.SetAuthor(authorEntry.Text)
			book.SetGenre(genreEntry.Text)
			book.SetRatting(rattingSelect.Selected)
			u.VM.AddBook(book)
			u.Emiter.Emit(AddedNewBookToList, book)
		}
	}
	d := dialog.NewForm("New Book", "Add", "Cancel", items, OnConfirm, u.Window)
	d.Resize(GetDialogSize(d.MinSize()))
	return d
}
