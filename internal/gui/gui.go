package gui

import (
	//"log"
	//"image/color"
	//"log"
	"fmt"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/dialog"
	_ "fyne.io/fyne/v2/canvas"

	myapp "github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/event"
)


/*
  I have don't know how I'm going to implement the data to canvas object entrys with the list widget.
  How am I going to hook the data to the edit button toggle that is with in the canvas entry object?
  What will I do when the edit button gets turned off, how is it going to update loaded data, how is
  it going to update that data to the tabel?
  When the labels gets wrapped how I'm I going to update the list object's height to be render with the new hight?
*/

type BookState uint
const (
	BSOnLoan BookState = 1 << iota 
	BSNew
	BSUpdate
	BSShow
)

type BooksTabel struct {
	Body         fyne.CanvasObject
	Header       fyne.CanvasObject

	BookTitles   []binding.String
	BookAuthors  []binding.String
	BookGenres   []binding.String
	BookRattings []binding.String
	BookStates   []BookState

	List         *widget.List
	BookShow     []int // Indexes to Bindings; Used for List.
}
func NewBooksTabel() *BooksTabel {
	b := &BooksTabel{
		BookTitles:   []binding.String{},
		BookAuthors:  []binding.String{},
		BookGenres:   []binding.String{},
		BookRattings: []binding.String{},
		BookStates:   []BookState{},
		BookShow:     []int{},
	}
	b.List = widget.NewList(
		func() int { // Length
			return len(b.BookShow)
		},
		func() fyne.CanvasObject { // CreateItem

			titleLabel := widget.NewLabel("")
			authorLabel := widget.NewLabel("")
			genreLabel := widget.NewLabel("")
			rattingLabel := widget.NewLabel("")

			titleLabel.Wrapping = fyne.TextWrapWord
			authorLabel.Wrapping = fyne.TextWrapWord
			genreLabel.Wrapping = fyne.TextWrapWord
			rattingLabel.Wrapping = fyne.TextWrapWord

			titleLabel.Truncation = fyne.TextTruncateEllipsis
			authorLabel.Truncation = fyne.TextTruncateEllipsis
			genreLabel.Truncation = fyne.TextTruncateEllipsis
			rattingLabel.Truncation = fyne.TextTruncateEllipsis

			titleLabel.Selectable = true
			authorLabel.Selectable = true
			genreLabel.Selectable = true
			rattingLabel.Selectable = true

			fields := []fyne.CanvasObject{
				titleLabel,
				authorLabel,
				genreLabel,
				rattingLabel,
				container.New(layout.NewPaddedLayout()),
				container.New(layout.NewPaddedLayout()),
				//widget.NewButtonWithIcon("", theme.CheckButtonIcon(), func(){return}), // this is for On Loaned
				//container.New(layout.NewGridLayout(2),
				//	editBtn,
				//	widget.NewButtonWithIcon("", NewDeleteIcon(), func(){return})),
			}
			entry := container.New(layout.NewGridLayout(len(fields)), fields...)
			return entry
		},          
		func(id int, o fyne.CanvasObject) { // UpdateItem
			index := b.BookShow[id]
			o.(*fyne.Container).Objects[0].(*widget.Label).Bind(b.BookTitles[index])
			o.(*fyne.Container).Objects[1].(*widget.Label).Bind(b.BookAuthors[index])
			o.(*fyne.Container).Objects[2].(*widget.Label).Bind(b.BookGenres[index])
			o.(*fyne.Container).Objects[3].(*widget.Label).Bind(b.BookRattings[index])
		})
	b.List.HideSeparators = false
	b.List.OnSelected = func(i int) {}   // todo
	b.List.OnUnselected = func(id int) {} // todo
	b.InitBooksHeader()
	return b
}

func (b *BooksTabel) AddNewBook(title, author, genre, ratting string) int {

	bTitle := binding.NewString()
	bAuthor := binding.NewString()
	bGenre := binding.NewString()
	bRatting := binding.NewString()

	bTitle.Set(title)
	bAuthor.Set(author)
	bGenre.Set(genre)
	bRatting.Set(ratting)
	
	index := len(b.BookTitles)
	b.BookTitles = append(b.BookTitles, bTitle)
	b.BookAuthors = append(b.BookAuthors, bAuthor)
	b.BookGenres = append(b.BookGenres, bGenre)
	b.BookRattings = append(b.BookRattings, bRatting)
	b.BookStates = append(b.BookStates, BSNew)

	return index
}
func (b *BooksTabel) ShowBook(index int) {
	state := b.BookStates[index]
	state = state | BSShow
	b.BookStates[index] = state
}
func (b *BooksTabel) HideBook(index int) {
	state := b.BookStates[index]
	state = state | BSShow
	b.BookStates[index] = state
}
func (b *BooksTabel) UpdateShowList() {
	list := []int{}
	for index, state := range b.BookStates {
		if (state & BSShow != 0) {
			list = append(list, index)
		}
	}
	// No need for preforments
	b.BookShow = list
	b.List.Refresh()
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
	}
	bottom := container.New(layout.NewGridLayout(len(fields)), fields...)

	fields = []fyne.CanvasObject{
			widget.NewButtonWithIcon("", theme.ContentAddIcon(), 
				func() {
					index := b.AddNewBook("Placeholder", "Placeholder", "Placeholder", "TBR")
					b.ShowBook(index)
					b.UpdateShowList()
				}),
	}

	top := container.New(layout.NewGridLayout(len(fields)), fields...)
	b.Header = container.New(layout.NewVBoxLayout(), top, bottom)
}

func (b *BooksTabel) AddBookEntry() fyne.CanvasObject {

	//bookTitle := binding.NewString()
	//bookAuthor := binding.NewString()
	//bookGenre  := binding.NewString()
	//bookRatting := binding.NewString()

	//bookTitle.Set("Placeholder")
	//bookAuthor.Set("Placeholder")
	//bookGenre.Set("Placeholder")
	//bookRatting.Set("Placeholder")

	titleLabel := widget.NewLabel("")
	authorLabel := widget.NewLabel("")
	genreLabel := widget.NewLabel("")
	rattingLabel := widget.NewLabel("")

	titleLabel.Wrapping = fyne.TextWrapWord
	authorLabel.Wrapping = fyne.TextWrapWord
	genreLabel.Wrapping = fyne.TextWrapWord
	rattingLabel.Wrapping = fyne.TextWrapWord

	//rattingEntry := widget.NewSelectWithData([]string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}, bookRatting)

	//green := color.RGBA{0x00, 0xff, 0x00, 0xff}
	//editEnabledBtn := widget.NewButtonWithIcon("", NewEditIcon(green, "EditEnable"), nil)
	//editDisableBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)

	//editEnabledBtn.Hide()

	//editBtn := container.New(layout.NewStackLayout(), editEnabledBtn, editDisableBtn)


	fields := []fyne.CanvasObject{
		titleLabel,
		authorLabel,
		genreLabel,
		rattingLabel,
		//widget.NewButtonWithIcon("", theme.CheckButtonIcon(), func(){return}), // this is for On Loaned
		//container.New(layout.NewGridLayout(2),
		//	editBtn,
		//	widget.NewButtonWithIcon("", NewDeleteIcon(), func(){return})),
	}
	entry := container.New(layout.NewGridLayout(len(fields)), fields...)
	return entry
}




func GetRattingStrings() []string {
	return []string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}
}

type UIState struct {
	Window fyne.Window
	DataHasChanged bool
	Emiter *event.EventEmiter
	BookSelected  int
	BookOrderedBy string
	BookList []string // TODO add book data
	UniqueGenres []string
	UniqueAuthors []string
}
func NewUIState(window fyne.Window) *UIState {
	u := &UIState{
		Emiter: event.NewEventEmiter(),
		Window: window,
	}
	return u
}

func (u *UIState) OpenNewBookForm() {

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
			dialog.ShowForm("Add Loaned Book", "Add", "Cancel", []*widget.FormItem{}, 
				func(b bool){
					if !b {
						onLoanCheck.Checked = false
						onLoanCheck.Refresh()
					}
				}, 
				u.Window,
			)
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
	Dialog.Show()

}

const (
	SaveButtonClicked    string = "SaveButtonClicked"
	ChangedOrderBy              = "ChangedOrderBy"
	ChangedOrderByAsc           = "ChangedOrderByAcs"
	ChangedOrderByDesc          = "ChangedOrderByDesc"
	ChangedSearchBy             = "ChangedSearchBy"
	ChangedSearch               = "ChangedSearch"
	NewBookButtonClicked        = "NewBookButtonClicked"
)

func Run() {

	a := app.New()
	window := a.NewWindow(myapp.AppName)

	UI := NewUIState(window)

	top:= UI.NewHeaderComp()
	body := UI.NewBookTableComp()
	mainComp := container.New(layout.NewBorderLayout(top, nil, nil, nil), top, body)
	window.SetContent(mainComp)
	window.ShowAndRun()
}
