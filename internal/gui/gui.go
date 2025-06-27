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

type BooksTabel struct {
	Body         fyne.CanvasObject
	Header       fyne.CanvasObject
	Box          fyne.CanvasObject
	Entries      []fyne.CanvasObject
	BookTitles   []binding.String
	BookAuthors  []binding.String
	BookGenres   []binding.String
	BookRattings []binding.String
	BookSavable  []bool
}
func NewBooksTabel() *BooksTabel {
	b := &BooksTabel{
		Entries: []fyne.CanvasObject{},
		BookTitles: []binding.String{},
		BookAuthors: []binding.String{},
		BookGenres: []binding.String{},
		BookRattings: []binding.String{},
		BookSavable: []bool{},
	}
	b.InitBooksHeader()
	b.Box = container.New(layout.NewVBoxLayout(), b.Entries...)
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
			widget.NewButtonWithIcon("", theme.ContentAddIcon(), b.AddBookEntry),
			//widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func(){return}),
		),
	}
	b.Header = container.New(layout.NewGridLayout(len(fields)), fields...)
}

func (b *BooksTabel) AddBookEntry() {

	bookTitle := binding.NewString()
	bookAuthor := binding.NewString()
	bookGenre  := binding.NewString()
	bookRatting := binding.NewString()

	b.BookTitles = append(b.BookTitles, bookTitle)
	b.BookAuthors = append(b.BookAuthors, bookAuthor)
	b.BookGenres = append(b.BookGenres, bookGenre)
	b.BookRattings = append(b.BookRattings, bookRatting)
	b.BookSavable = append(b.BookSavable, false)

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
	bookRatting.Set("TBR")

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
	b.Entries = append(b.Entries, entry)
	b.Box.(*fyne.Container).Add(entry)
}

func Run() {
	a := app.New()
	window := a.NewWindow(myapp.AppName)

	searchBy := widget.NewSelectWithData([]string{"Title", "Author", "Genre", "Ratting", "On Loan"}, binding.NewString())
	search := container.New(layout.NewGridLayout(2), searchBy, widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func(){return}))

	tabel := NewBooksTabel()

	top := container.New(layout.NewGridLayout(1), search, tabel.Header)
	scroll := container.NewVScroll(tabel.Box)

	maincon := container.New(layout.NewBorderLayout(top, nil, nil, nil), top, scroll)

	window.SetContent(maincon)
	window.ShowAndRun()
}
