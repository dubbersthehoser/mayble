package gui

import (
	//"log"
	"image/color"
	//"log"
	//"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	_ "fyne.io/fyne/v2/canvas"

	myapp "github.com/dubbersthehoser/mayble/internal/app"
)


/*
  I have don't know how I'm going to implement the data to canvas object entrys with the list widget.
  How am I going to hook the data to the edit button toggle that is with in the canvas entry object?
  What will I do when the edit button gets turned off, how is it going to update loaded data, how is
  it going to update that data to the tabel?
  When the labels gets wrapped how I'm I going to update the list object's height to be render with the new hight?
*/

type BooksTabel struct {
	Body         fyne.CanvasObject
	Header       fyne.CanvasObject
	List         *widget.List
	ListObjs     []fyne.CanvasObject
	BookTitles   []binding.String
	BookAuthors  []binding.String
	BookGenres   []binding.String
	BookRattings []binding.String
	BookSavable  []bool
}
func NewBooksTabel() *BooksTabel {
	b := &BooksTabel{
		ListObjs:     []fyne.CanvasObject{},
		BookTitles:   []binding.String{},
		BookAuthors:  []binding.String{},
		BookGenres:   []binding.String{},
		BookRattings: []binding.String{},
		BookSavable:  []bool{},
	}
	b.List = widget.NewList(
		func() int {return len(b.BookTitles)}, // Length
		func() { // CreateItem
			o := b.AddBookEntry()
			b.ListObjs = append(b.ListObjs, o)
			return o
			}),          
		func(id int, o fyne.CanvasObject) { // UpdateItem
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry).Bind(b.BookTitles[id])
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).Bind(b.BookTitles[id])
		})
	b.List.HideSeparators = true
	b.List.OnSelected = func(i int) {} // Disable selection highlighting
	b.List.OnUnselected = func(id int) {
		
		o := b.List.Objects[id]
		
		// Book Title
		o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry).Disable() // Entry
		o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry).Hide()
		o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).Show()    // Label

		// Book Author
		o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Entry).Disable() // Entry
		o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Entry).Hide()
		o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Label).Show()    // Label

		// Book Genre
		o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Entry).Disable() // Entry
		o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Entry).Hide()
		o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Label).Show()    // Label

		// Book Ratting
		o.(*fyne.Container).Objects[3].(*fyne.Container).Objects[0].(*widget.Entry).Disable() // Entry
		o.(*fyne.Container).Objects[3].(*fyne.Container).Objects[0].(*widget.Entry).Hide()
		o.(*fyne.Container).Objects[3].(*fyne.Container).Objects[1].(*widget.Label).Show()    // Label

		// Edit Button 
		o.(*fyne.Container).Objects[3].(*fyne.Container).Objects[0].(*widget.Entry).Hide()

		// Delete Button 
		editDisableBtn.Show()
		editEnabledBtn.Hide()
	}
	b.InitBooksHeader()
	return b
}

func (b *BooksTabel) InitBooksHeader() {
	style := fyne.TextStyle{
		Bold: true,
	}
	align := fyne.TextAlignCenter
	fields := []fyne.CanvasObject{
		widget.NewLabelWithStyle("Title", align, style),
		widget.NewLabelWithStyle("Author", align, style),
		widget.NewLabelWithStyle("Genre", align, style),
		widget.NewLabelWithStyle("Ratting", align, style),
		widget.NewLabelWithStyle("On Loan", align, style),
		container.New(layout.NewGridLayout(2), 
			widget.NewLabelWithStyle("Actions", align, style),
			widget.NewButtonWithIcon("", theme.ContentAddIcon(), 
				func() {
					b.BookTitles = append(b.BookTitles, binding.NewString())
					b.BookAuthors = append(b.BookAuthors, binding.NewString())
					b.BookGenres = append(b.BookGenres, binding.NewString())
					r := binding.NewString()
					r.Set("TBR")
					b.BookRattings = append(b.BookRattings, r)
					b.BookSavable = append(b.BookSavable, false)
					b.List.Refresh()
					}),
		),
	}
	b.Header = container.New(layout.NewGridLayout(len(fields)), fields...)
}

func (b *BooksTabel) AddBookEntry() fyne.CanvasObject {

	bookTitle := binding.NewString()
	bookAuthor := binding.NewString()
	bookGenre  := binding.NewString()
	bookRatting := binding.NewString()

	titleLabel := widget.NewLabelWithData(bookTitle)
	authorLabel := widget.NewLabelWithData(bookAuthor)
	genreLabel := widget.NewLabelWithData(bookGenre)
	rattingLabel := widget.NewLabelWithData(bookRatting)

	titleLabel.Wrapping = fyne.TextWrapWord
	authorLabel.Wrapping = fyne.TextWrapWord
	genreLabel.Wrapping = fyne.TextWrapWord
	rattingLabel.Wrapping = fyne.TextWrapWord

	titleLabel.Show()
	authorLabel.Show()
	genreLabel.Show()
	rattingLabel.Show()

	titleEntry :=  widget.NewEntry()
	authorEntry := widget.NewEntry()
	genreEntry :=  widget.NewEntry()
	rattingEntry := widget.NewSelectWithData([]string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}, bookRatting)

	titleEntry.Bind(bookTitle)
	authorEntry.Bind(bookAuthor)
	genreEntry.Bind(bookGenre)

	titleEntry.Disable()
	titleEntry.Hide()
	authorEntry.Disable()
	authorEntry.Hide()
	genreEntry.Disable()
	genreEntry.Hide()
	rattingEntry.Disable()
	rattingEntry.Hide()

	green := color.RGBA{0x00, 0xff, 0x00, 0xff}
	editEnabledBtn := widget.NewButtonWithIcon("", NewEditIcon(green, "EditEnable"), nil)
	
	editDisableBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)

	editEnabledBtn.Hide()

	editBtn := container.New(layout.NewStackLayout(), editEnabledBtn, editDisableBtn)

	editEnabledBtn.OnTapped = func() {
		titleEntry.Disable()
		authorEntry.Disable()
		genreEntry.Disable()
		rattingEntry.Disable()

		titleLabel.Show()
		authorLabel.Show()
		genreLabel.Show()
		rattingLabel.Show()
		
		titleEntry.Hide()
		authorEntry.Hide()
		genreEntry.Hide()
		rattingEntry.Hide()

		editDisableBtn.Show()
		editEnabledBtn.Hide()
	}
	editDisableBtn.OnTapped = func() {
		titleEntry.Enable()
		authorEntry.Enable()
		genreEntry.Enable()
		rattingEntry.Enable()

		titleLabel.Hide()
		authorLabel.Hide()
		genreLabel.Hide()
		rattingLabel.Hide()
		
		titleEntry.Show()
		authorEntry.Show()
		genreEntry.Show()
		rattingEntry.Show()

		editDisableBtn.Hide()
		editEnabledBtn.Show()
	}


	fields := []fyne.CanvasObject{
		container.New(layout.NewStackLayout(), titleEntry, titleLabel),
		container.New(layout.NewStackLayout(), authorEntry, authorLabel),
		container.New(layout.NewStackLayout(), genreEntry, genreLabel),
		container.New(layout.NewStackLayout(), rattingEntry, rattingLabel),
		widget.NewButtonWithIcon("", theme.CheckButtonIcon(), func(){return}),
		container.New(layout.NewGridLayout(2),
			editBtn,
			widget.NewButtonWithIcon("", NewDeleteIcon(), func(){return})),
	}

	entry := container.New(layout.NewGridLayout(len(fields)), fields...)
	return entry
}

func Run() {
	a := app.New()
	window := a.NewWindow(myapp.AppName)

	searchBy := widget.NewSelectWithData([]string{"Title", "Author", "Genre", "Ratting", "On Loan"}, binding.NewString())
	search := container.New(layout.NewGridLayout(2), searchBy, widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func(){return}))

	tabel := NewBooksTabel()

	top := container.New(layout.NewGridLayout(1), search, tabel.Header)
	//scroll := container.NewVScroll(tabel.Box)

	maincon := container.New(layout.NewBorderLayout(top, nil, nil, nil), top, tabel.List)

	window.SetContent(maincon)
	window.ShowAndRun()
}
