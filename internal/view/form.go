package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

func NewCreateBookForm(vm *viewmodel.CreateBookForm) fyne.CanvasObject {

	loanCheck := widget.NewCheckWithData("On Loan", vm.IsLoaned)
	readCheck := widget.NewCheckWithData("Is Read", vm.IsRead)
	submit := widget.NewButton("Submit", vm.Submit)
	add := widget.NewButton("Add Submission", vm.AddSubmission)
	submit.Alignment = widget.ButtonAlignLeading
	add.Alignment = widget.ButtonAlignLeading

	bookEntry := newBookEntry(vm.Title, vm.Author, vm.Genre, vm.Genres)

	top := container.NewVBox(
		bookEntry,
		loanCheck,
		newLoanEntry(vm.IsLoaned, vm.Date, vm.Borrower),
		readCheck,
		newReadEntry(vm.IsRead, vm.Completed, vm.Rating), 
		container.NewBorder(nil, nil, add, submit, add, submit),
	)

	//return container.New(layout.NewVBoxLayout(), 
	return container.NewBorder(top, nil, nil, nil,
		top,
		newSubmitionList(vm.SubmissionList()),
	)
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
	
	list := container.NewStack(container.NewVScroll(container.NewStack(content)))
	return list
}
