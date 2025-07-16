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


func (u *UIState) NewOnLoanDialog(loan *Loan) dialog.Dialog {
	if loan == nil {
		loan = &Loan{}
	}
	dateEntry := widget.NewDateEntry()
	nameEntry := widget.NewEntry()
	d := dialog.NewForm("Add Loaned Book", "Add", "Cancel", 
		[]*widget.FormItem{
			widget.NewFormItem("Name", nameEntry),
			widget.NewFormItem("Date", dateEntry),
		}, 
		func(b bool){
			if b {
				loan.Date = *dateEntry.Date
				loan.Name = nameEntry.Text
			}
		}, 
		u.Window,
	)
	return d
}

func (u *UIState) NewBookDialog(book *Book) dialog.Dialog {
	
	rattings := GetRattingStrings()

	titleEntry := widget.NewEntry()
	authorSelect := widget.NewSelectEntry(u.UniqueAuthors)
	genreEntry := widget.NewSelectEntry(u.UniqueGenres)
	rattingSelect := widget.NewSelect(rattings, nil)

	titleEntry.Validator = func(s string) error {
		if len(s) == 0 {
			return errors.New("Must have a Title")
		}
		return nil
	}

	authorSelect.Validator = func(s string) error {
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
			var loan *Loan = nil
			ld := u.NewOnLoanDialog(loan)
			ld.Show()
			if loan != nil {
				onLoanCheck.Checked = true
				onLoanCheck.Refresh()
			}
		} else {
			dialog.ShowConfirm("Remove Loaning", "Are you sure?", 
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
			"Title", 
			titleEntry,
		),
		widget.NewFormItem(
			"Author", 
			authorSelect,
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
		func (b bool) {
			if b {
				fmt.Println("Yes")
			} else {
				fmt.Println("No")
			}
		}, 
		u.Window,
	)
	height := Dialog.MinSize().Height
	size := fyne.NewSize(400, 0)
	size.Height = height
	Dialog.Resize(size)
	return Dialog
}

func (u *UIState) OpenNewBookForm() {
	b := &Book{}
	d := u.NewBookDialog(b)
	d.Show()
}
