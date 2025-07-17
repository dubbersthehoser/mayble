package gui

import (
	"fmt"
	"time"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	_"fyne.io/fyne/v2/container"
	_"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"
)

type Loan struct {
	Name string
	Date time.Time
}

type Book struct {
	Title string
	Author string
	Genre string
	Ratting string
	IsOnLoan bool
	OnLoan Loan
}

const DialogWidth float32 = 400

func (u *UIState) NewOnLoanDialog(loan *Loan) dialog.Dialog {
	if loan == nil {
		loan = &Loan{}
	}
	dateEntry := widget.NewDateEntry()
	nameEntry := widget.NewEntry()

	nameEntry.Validator = func(s string) error {
		if len(s) == 0 {
			return errors.New("Must have a Name")
		}
		return nil
	}

	d := dialog.NewForm("Loaned Book", "Add", "Cancel", 
		[]*widget.FormItem{
			widget.NewFormItem("Name*", nameEntry),
			widget.NewFormItem("Date*", dateEntry),
		}, 
		func(b bool){
			if b {
				loan.Date = *dateEntry.Date
				loan.Name = nameEntry.Text
			}
		}, 
		u.Window,
	)
	height := d.MinSize().Height
	size := fyne.NewSize(DialogWidth, 0)
	size.Height = height
	d.Resize(size)
	return d
}

func (u *UIState) NewBookDialog(book *Book) dialog.Dialog {

	if book == nil {
		book = &Book{}
	}

	rattings := GetRattingStrings()

	titleEntry := widget.NewEntry()
	authorEntry := widget.NewSelectEntry(u.UniqueAuthors)
	genreEntry := widget.NewSelectEntry(u.UniqueGenres)
	rattingSelect := widget.NewSelect(rattings, nil)

	titleEntry.Text = book.Title
	authorEntry.Text = book.Author
	genreEntry.Text = book.Genre
	rattingSelect.Selected = book.Ratting


	titleEntry.Validator = func(s string) error {
		if len(s) == 0 {
			return errors.New("Must have a Title")
		}
		return nil
	}

	authorEntry.Validator = func(s string) error {
		if len(s) == 0 {
			return errors.New("Must have an Author")
		}
		return nil
	}

	rattingSelect.PlaceHolder = rattings[0]
	rattingSelect.Selected = rattings[0]

	onLoanCheck := widget.NewCheck(
		"", 
		nil,
	)

	onLoanCheck.OnChanged = func (checked bool) {
		if checked {
			u.Emiter.Emit(NewOnLoanEvent, book)
			if book.IsOnLoan {
				onLoanCheck.Checked = true
				onLoanCheck.Refresh()
			}
		} else {
			dialog.ShowConfirm("Unloan Book", "Are you sure?", 
				func(b bool){
					if !b {
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

	Dialog := dialog.NewForm("New Book", "Add", "Cancel", f,
		func (ok bool) {
			if ok {
				fmt.Println("Yes")
				book.Title = titleEntry.Text
				book.Author = authorEntry.Text
				book.Genre = genreEntry.Text
				book.Ratting = rattingSelect.Selected
				fmt.Printf("%#v\n", book)
			} else {
				fmt.Println("No")
			}
		}, 
		u.Window,
	)
	height := Dialog.MinSize().Height
	size := fyne.NewSize(DialogWidth, 0)
	size.Height = height
	Dialog.Resize(size)
	return Dialog
}
