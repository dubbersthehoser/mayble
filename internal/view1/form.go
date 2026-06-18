package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
	
	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
)
func newEdit(vm *viewmodel.Window) fyne.CanvasObject {
	return newBookForm(vm, "Update", vm.Form.OnUpdate)
}

func newCreate(vm *viewmodel.Window) fyne.CanvasObject {
	return newBookForm(vm, "Create", vm.Form.OnCreate)
}

func newBookForm(vm *viewmodel.Window, label string, submit func()) fyne.CanvasObject {

	loanCheck := widget.NewCheckWithData("Is on loan.", vm.Form.Fyne.IsLoaned)
	readCheck := widget.NewCheckWithData("Has been completed.", vm.Form.Fyne.IsCompleted)

	submitBtn := NewEnterButton(label, submit)
	submitBtn.Alignment = widget.ButtonAlignLeading

	closeBtn := NewEnterButton("Cancel", func() {
		vm.Form.Reset()
		vm.Body.Set(viewmodel.BodyTable)
	})

	bookEntry := newBookEntry(vm)

	top := container.NewVBox(
		bookEntry,
		loanCheck,
		newLoanEntry(vm),
		readCheck,
		newReadEntry(vm),
		container.NewHBox(submitBtn, closeBtn),
	)

	return top
}

func newBookEntry(vm *viewmodel.Window) *fyne.Container {

	titleEntry := widget.NewEntryWithData(vm.Form.Fyne.Title)
	authorEntry := widget.NewEntryWithData(vm.Form.Fyne.Author)

	genres := vm.UniqueGenres.Genres()
	genreEntry := widget.NewSelectEntry(genres)
	genreEntry.Bind(vm.Form.Fyne.Genre)
	vm.UniqueGenres.AddListener(func() {
		genreEntry.SetOptions(vm.UniqueGenres.Genres())
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
	dateEntry.Bind(vm.Form.Fyne.LoanedAt)
	nameEntry := widget.NewEntry()
	nameEntry.Bind(vm.Form.Fyne.Borrower)

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
	}

	vm.Form.Fyne.IsLoaned.AddListener(binding.NewDataListener(update))
	update()

	return c
}

func newReadEntry(vm *viewmodel.Window) *fyne.Container {
	ratingEntry := widget.NewSelectWithData(viewmodel.Ratings(), vm.Form.Fyne.Rating)
	completedEntry := widget.NewDateEntry()
	completedEntry.Bind(vm.Form.Fyne.CompletedAt)

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
	}

	vm.Form.Fyne.IsCompleted.AddListener(binding.NewDataListener(update))
	update()
	return c
}
