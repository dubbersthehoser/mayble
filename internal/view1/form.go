package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	
	"github.com/dubbersthehoser/mayble/internal/viewmodel1"
)
func newEdit(vm *viewmodel.Window) fyne.CanvasObject {
	return nil
}

func newCreate(vm *viewmodel.Window) fyne.CanvasObject {
	return nil
}

func newBookForm(vm *viewmodel.Window, submit func()) fyne.CanvasObject {

	loanCheck := widget.NewCheck("Is on loan.", vm.updateForm.SetLoaned)
	readCheck := widget.NewCheckWithData("Has been completed.", vm.IsRead)

	submit := NewEnterButton("Submit", submit)

	submit.Alignment = widget.ButtonAlignLeading

	bookEntry := newBookEntry(vm.Title, vm.Author, vm.Genre, vm.Genres)

	top := container.NewVBox(
		bookEntry,
		loanCheck,
		newLoanEntry(vm.IsLoaned, vm.Date, vm.Borrower),
		readCheck,
		newReadEntry(vm.IsRead, vm.Completed, vm.Rating),
		container.NewHBox(add, submit, limit),
	)

	return top
}

func newBookEntry(title, author, genre binding.String, uniqueGenres *viewmodel.UniqueGenres) *fyne.Container {

	titleEntry := widget.NewEntryWithData(title)
	authorEntry := widget.NewEntryWithData(author)

	genres := uniqueGenres.Get()
	genreEntry := widget.NewSelectEntry(genres)
	genreEntry.Bind(genre)

	uniqueGenres.AddListener(binding.NewDataListener(func() {
		genres := uniqueGenres.Get()
		genreEntry.SetOptions(genres)
	}))

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

func newLoanEntry(isLoaned binding.Bool, date binding.String, borrower binding.String) *fyne.Container {
	dateEntry := widget.NewDateEntry()
	nameEntry := widget.NewEntry()

	dateEntry.Bind(date)
	nameEntry.Bind(borrower)

	nameEntry.SetPlaceHolder("Borrower...")
	dateEntry.SetPlaceHolder("DD/MM/YYYY")

	c := container.NewVBox(
		nameEntry,
		dateEntry,
	)

	ok, _ := isLoaned.Get()
	if !ok {
		dateEntry.Disable()
		nameEntry.Disable()
	}

	isLoaned.AddListener(binding.NewDataListener(func() {
		ok, _ := isLoaned.Get()
		if ok {
			dateEntry.Enable()
			nameEntry.Enable()
		} else {
			dateEntry.Disable()
			nameEntry.Disable()
		}
	}))
	return c
}

func newReadEntry(isRead binding.Bool, completed binding.String, rating binding.String) *fyne.Container {
	ratingEntry := widget.NewSelectWithData(viewmodel.Ratings(), rating)
	completedEntry := widget.NewDateEntry()

	ratingEntry.Bind(rating)
	completedEntry.Bind(completed)

	rattingStrings := viewmodel.Ratings()
	ratingEntry.PlaceHolder = rattingStrings[0]
	completedEntry.SetPlaceHolder("DD/MM/YYYY")

	c := container.NewVBox(
		completedEntry,
		ratingEntry,
	)

	ok, _ := isRead.Get()
	if !ok {
		ratingEntry.Disable()
		completedEntry.Disable()
	}

	isRead.AddListener(binding.NewDataListener(func() {
		ok, _ := isRead.Get()
		if ok {
			ratingEntry.Enable()
			completedEntry.Enable()
		} else {
			ratingEntry.Disable()
			completedEntry.Disable()
		}
	}))
	return c
}
