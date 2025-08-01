package gui

import (
	"fmt"
	//"time"
	//"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	_"fyne.io/fyne/v2/container"
	_"fyne.io/fyne/v2/layout"
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
				fmt.Println("Add Loan")
			}
		}, 
		u.Window,
	)
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

	onLoanCheck.OnChanged = func(checked bool) {
		if checked {
			loan := u.VM.NewLoan()
			u.NewOnLoanDialog(loan).Show()
			book.SetLoan(loan)
			onLoanCheck.Checked = true
			onLoanCheck.Refresh()
		} else {
			dialog.ShowConfirm("Remove Loan", "Are you sure?", 
				func(ok bool){
					if ok {
						book.UnsetLoan()
					} else {
						onLoanCheck.Checked = true
						onLoanCheck.Refresh()
						
					}
				},
				u.Window,
			)
		}
	}

	f := []*widget.FormItem{
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

	d := dialog.NewForm("New Book", "Add", "Cancel", f,
		func (ok bool) {
			if ok {
				fmt.Println("Yes")
				book.SetTitle(titleEntry.Text)
				book.SetAuthor(authorEntry.Text)
				book.SetGenre(genreEntry.Text)
				book.SetRatting(rattingSelect.Selected)
				u.VM.AddBook(book)
			} else {
				fmt.Println("No")
			}
		}, 
		u.Window,
	)
	d.Resize(GetDialogSize(d.MinSize()))
	return d
}
