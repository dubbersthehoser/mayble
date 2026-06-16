package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	
	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
)
func newEdit(vm *viewmodel.Window) fyne.CanvasObject {
	return newBookForm(vm, "Update", vm.Form.OnUpdate)
}

func newCreate(vm *viewmodel.Window) fyne.CanvasObject {
	return newBookForm(vm, "Create", vm.Form.OnCreate)
}

func newBookForm(vm *viewmodel.Window, label string, submit func()) fyne.CanvasObject {


	loanCheck := widget.NewCheck("Is on loan.", vm.Form.SetLoaned)
	readCheck := widget.NewCheck("Has been completed.", vm.Form.SetCompleted)

	submitBtn := NewEnterButton(label, submit)

	submitBtn.Alignment = widget.ButtonAlignLeading

	bookEntry := newBookEntry(vm)

	top := container.NewVBox(
		bookEntry,
		loanCheck,
		newLoanEntry(vm),
		readCheck,
		newReadEntry(vm),
		container.NewHBox(submitBtn),
	)

	return top
}

func newBookEntry(vm *viewmodel.Window) *fyne.Container {

	titleEntry := widget.NewEntry()
	titleEntry.OnChanged = vm.Form.SetTitle
	authorEntry := widget.NewEntry()
	authorEntry.OnChanged = vm.Form.SetAuthor

	genres := vm.UniqueGenres.Genres()
	genreEntry := widget.NewSelectEntry(genres)
	genreEntry.OnChanged = vm.Form.SetGenre
	vm.UniqueGenres.AddListener(func() {
		genreEntry.SetOptions(vm.UniqueGenres.Genres())
	})

	vm.Form.AddListener(func() {
		titleEntry.SetText(vm.Form.GetTitle())
		authorEntry.SetText(vm.Form.GetAuthor())
	})

	titleEntry.SetPlaceHolder("Title...")
	authorEntry.SetPlaceHolder("Author...")
	genreEntry.SetPlaceHolder("Genre...")

	c := container.NewVBox(
		titleEntry,
		authorEntry,
		genreEntry,
	)
	return c

}

func newLoanEntry(vm *viewmodel.Window) *fyne.Container {
	dateEntry := widget.NewDateEntry()
	dateEntry.OnChanged = vm.Form.SetLoanedAt
	nameEntry := widget.NewEntry()
	nameEntry.OnChanged = vm.Form.SetBorrower

	nameEntry.SetPlaceHolder("Borrower...")
	dateEntry.SetPlaceHolder("DD/MM/YYYY")

	c := container.NewVBox(
		nameEntry,
		dateEntry,
	)

	update := func() {
		if vm.Form.IsLoaned() {
			dateEntry.Enable()
			nameEntry.Enable()
		} else {
			dateEntry.Disable()
			nameEntry.Disable()
		}

		dateEntry.Date = vm.Form.GetLoanedAt()
		dateEntry.Refresh()

		nameEntry.SetText(vm.Form.GetBorrower())
	}
	vm.Form.AddListener(update)
	update()

	return c
}

func newReadEntry(vm *viewmodel.Window) *fyne.Container {
	ratingEntry := widget.NewSelect(viewmodel.Ratings(), vm.Form.SetRating)
	completedEntry := widget.NewDateEntry()
	completedEntry.OnChanged = vm.Form.SetCompletedAt

	rattingStrings := viewmodel.Ratings()

	ratingEntry.PlaceHolder = rattingStrings[0]
	completedEntry.SetPlaceHolder("DD/MM/YYYY")

	c := container.NewVBox(
		completedEntry,
		ratingEntry,
	)

	update := func() {
		if vm.Form.IsCompleted() {
			ratingEntry.Enable()
			completedEntry.Enable()
		} else {
			ratingEntry.Disable()
			completedEntry.Disable()
		}

		// Don't use .SetSelected will create inf recursion.
		ratingEntry.Selected = vm.Form.GetRating()
		ratingEntry.Refresh()
		completedEntry.Date = vm.Form.GetCompletedAt()
	}

	vm.Form.AddListener(update)
	update()

	return c
}
