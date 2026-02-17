package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func NewCreateBookForm(vm *viewmodel.CreateBookForm) fyne.CanvasObject {

	submit := widget.NewButton("Submit", vm.Submit)
	add := widget.NewButton("Add Submission", vm.AddSubmission)

	submit.Alignment = widget.ButtonAlignLeading
	add.Alignment = widget.ButtonAlignLeading

	loanCheck := widget.NewCheckWithData("On Loan", vm.IsLoaned)
	readCheck := widget.NewCheckWithData("Is Read", vm.IsRead)

	titleEntry := widget.NewEntryWithData(vm.Title)
	authorEntry := widget.NewEntryWithData(vm.Author)
	genreEntry := widget.NewEntryWithData(vm.Genre)

	titleEntry.SetPlaceHolder("Title...")
	authorEntry.SetPlaceHolder("Author...")
	genreEntry.SetPlaceHolder("Genre...")

	return container.New(layout.NewVBoxLayout(), 
		titleEntry,
		authorEntry,
		genreEntry,
		loanCheck,
		newLoanForm(vm.IsLoaned, vm.Date, vm.Borrower),
		readCheck,
		newReadForm(vm.IsRead, vm.Completed, vm.Rating), 
		container.NewBorder(nil, nil, add, submit, add, submit),
		widget.NewLabel(""),
		newSubmitionList(vm.SubmissionList()),
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
	ratingEntry := widget.NewSelectWithData(viewmodel.Ratings(), rating)
	completedEntry := widget.NewDateEntry()

	ratingEntry.Bind(rating)
	completedEntry.Bind(completed)

	ratingEntry.PlaceHolder = viewmodel.Ratings()[0]
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


func newSubmitionList(sl *viewmodel.SubmissionList) fyne.CanvasObject {	
	content := container.NewVBox()
	update := func() {
		content.RemoveAll()
		for i := range sl.Length() {
			v := sl.Get(i)
			del := widget.NewButtonWithIcon(
				"",
				theme.DeleteIcon(),
				func(id int) func() {
					return func() {
						sl.Remove(id)
					}
				}(i),
			)
			edt := widget.NewButtonWithIcon(
				"",
				theme.DocumentCreateIcon(),
				func(id int) func() {
					return func() {
						sl.Edit(id)
					}
				}(i),
			)

			del.Importance = widget.DangerImportance
			edt.Importance = widget.SuccessImportance

			btns := container.NewHBox(edt, del)

			object := container.NewBorder(nil, nil, nil, btns, btns, widget.NewLabel(v), )
			content.Add(object)
		}

	}
	sl.AddListener(binding.NewDataListener(func() {
		update()
	}))
	update()
	
	list := container.NewHScroll(container.NewStack(content))
	return container.NewStack(list)
}
