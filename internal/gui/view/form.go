package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/container"
	
	"github.com/dubbersthehoser/mayble/internal/gui/viewmodel"
)

func BookForm(vm *viewmodel.BookForm) fyne.CanvasObject {
	submit := widget.NewButton("Submit", vm.Submit)
	cancel := widget.NewButton("Cancel", vm.Cancel)
	message := widget.NewLabel("")

	//vm.Error.AddListener(binding.NewDataListener(func() {
	//	msg, _ := vm.Error.Get()
	//	message.SetText(msg)
	//}))

	//vm.Success.AddListener(binding.NewDataListener(func() {
	//	msg, _ := vm.Success.Get()
	//	message.SetText(msg)
	//}))

	loanCheck := widget.NewCheckWithData("On Loan", vm.IsLoaned)
	readCheck := widget.NewCheckWithData("Is Read", vm.IsRead)

	titleEntry := widget.NewEntryWithData(vm.Title)
	authorEntry := widget.NewEntryWithData(vm.Author)
	genreEntry := widget.NewEntryWithData(vm.Genre)

	titleEntry.SetPlaceHolder("Title...")
	authorEntry.SetPlaceHolder("Author...")
	genreEntry.SetPlaceHolder("Genre...")

	vm.Valid.AddListener(binding.NewDataListener(func() {
		ok, _ := vm.Valid.Get()
		if ok {
			message.Importance = widget.SuccessImportance
		} else {
			message.Importance = widget.DangerImportance
		}
	}))

	return container.New(layout.NewVBoxLayout(), 
		titleEntry,
		authorEntry,
		genreEntry,
		loanCheck,
		newLoanForm(vm.IsLoaned, vm.Date, vm.Borrower),
		readCheck,
		newReadForm(vm.IsRead, vm.Completed, vm.Rating), 
		message,
		submit,
		cancel,
	)
}

func newLoanForm(isLoaned binding.Bool, date binding.String, borrower binding.String) *fyne.Container {
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

func newReadForm(isRead binding.Bool, completed binding.String, rating binding.String) *fyne.Container {
	ratingEntry := widget.NewSelect([]string{}, nil)
	completedEntry := widget.NewDateEntry()

	ratingEntry.Bind(rating)
	completedEntry.Bind(completed)

	ratingEntry.PlaceHolder = "Rating"
	completedEntry.SetPlaceHolder("Completed DD/MM/YYYY")

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




