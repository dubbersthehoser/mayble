package gui

import (
	//"log"
	"image/color"
	//"log"
	"fmt"

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

type CustomIcon struct {
	MyContent []byte
	MyName string
}
func (c *CustomIcon) Content() []byte {
	return c.MyContent
}
func (c *CustomIcon) Name() string {
	return c.MyName
}


func ChangeSVGIconFG(svg []byte, c color.RGBA) []byte {
	data := []byte{}
	red := fmt.Sprintf("%02x", c.R)
	green := fmt.Sprintf("%02x", c.G)
	blue := fmt.Sprintf("%02x", c.B)
	valueStart := -1
	for i, c := range svg {
		if c == '"' && valueStart == -1 {
			valueStart = i+1
			data = append(data, c)
		} else if c == '"' && valueStart > -1 {
			value := svg[valueStart:i]
			if string(value) == "#f3f3f3" {
				data = append(data, '#')
				data = append(data, red[0])
				data = append(data, red[1])
				data = append(data, green[0])
				data = append(data, green[1])
				data = append(data, blue[0])
				data = append(data, blue[1])
			} else {
				for _, c := range value {
					data = append(data, c)
				}
			}
			valueStart = -1
		}
		if valueStart == -1 {
			data = append(data, c)
		}
	}
	return data
}

func NewDeleteIcon() fyne.Resource {
	icon := theme.DeleteIcon()
	data := icon.Content()
	data = ChangeSVGIconFG(data, color.RGBA{0xff, 0x00, 0x00, 0x00})
	myIcon := &CustomIcon{
		MyName: "MyDeleteIcon",
		MyContent: data,
	}
	return myIcon
}
func NewEditIcon(color color.RGBA, name string) fyne.Resource {
	data := theme.DocumentCreateIcon().Content()
	data = ChangeSVGIconFG(data, color)
	myIcon := &CustomIcon{
		MyName: name,
		MyContent: data,
	}
	return myIcon
}

func NewBooksHeader() *fyne.Container {
	fields := []fyne.CanvasObject{
		//widget.NewButtonWithIcon("Title", theme.MoveDownIcon(), func(){return}),
		widget.NewLabel("Title"),
		widget.NewLabel("Author"),
		widget.NewLabel("Genre"),
		widget.NewLabel("Ratting"),
		widget.NewLabel("On Loan"),
		container.New(layout.NewGridLayout(2), 
			widget.NewLabel("Actions"),
			widget.NewButtonWithIcon("", theme.ContentAddIcon(), func(){return}),
			//widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func(){return}),
		),
	}

	return container.New(layout.NewGridLayout(len(fields)), fields...)
}


type BookEntryWidget struct {
	widget.BaseWidget
	object fyne.CanvasObject
}
func (b *BookEntryWidget)CreateRenderer() fyne.WidgetRenderer {
	if b == nil {
		return nil
	}
	if w, ok := b.object.(fyne.Widget); ok {
		return w.CreateRenderer()
	}
	return widget.NewSimpleRenderer(b.object)
}
func NewBookEntryWidget(o fyne.CanvasObject) *BookEntryWidget {
	b := &BookEntryWidget{object: o}
	b.ExtendBaseWidget(b)
	return b
}
	

func NewBookEntry() *fyne.Container {

	bookTitle := binding.NewString()
	bookTitle.Set("The Cat")

	titleLabel := widget.NewLabelWithData(bookTitle)
	titleLabel.Show()

	titleEntry :=  widget.NewEntry()
	authorEntry := widget.NewEntry()
	genreEntry :=  widget.NewEntry()
	rattingsEntry := widget.NewSelectWithData([]string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}, binding.NewString())

	titleEntry.Bind(bookTitle)

	titleEntry.Disable()
	titleEntry.Hide()
	authorEntry.Disable()
	genreEntry.Disable()
	rattingsEntry.Disable()


	green := color.RGBA{0x00, 0xff, 0x00, 0xff}
	editEnabledBtn := widget.NewButtonWithIcon("", NewEditIcon(green, "EditEnable"), nil)
	
	editDisableBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)

	editEnabledBtn.Hide()

	editBtn := container.New(layout.NewStackLayout(), editEnabledBtn, editDisableBtn)

	editEnabledBtn.OnTapped = func() {
		titleEntry.Disable()
		authorEntry.Disable()
		genreEntry.Disable()
		rattingsEntry.Disable()

		titleLabel.Show()
		titleEntry.Hide()

		editDisableBtn.Show()
		editEnabledBtn.Hide()
	}
	editDisableBtn.OnTapped = func() {
		titleEntry.Enable()
		authorEntry.Enable()
		genreEntry.Enable()
		rattingsEntry.Enable()

		titleLabel.Hide()
		titleEntry.Show()

		editDisableBtn.Hide()
		editEnabledBtn.Show()
	}


	fields := []fyne.CanvasObject{
		container.New(layout.NewStackLayout(), titleEntry, titleLabel),
		authorEntry,
		genreEntry,
		rattingsEntry,
		widget.NewButtonWithIcon("", theme.CheckButtonIcon(), func(){return}),
		container.New(layout.NewGridLayout(2),
			editBtn,
			widget.NewButtonWithIcon("", NewDeleteIcon(), func(){return})),
	}
	fields[3].(*widget.Select).PlaceHolder = "TBR"

	return container.New(layout.NewGridLayout(len(fields)), fields...)
}

func Run() {
	//data := theme.DocumentCreateIcon().Content()
	//fmt.Println(string(data))
	//data = ChangeSVGIconFG(data, color.RGBA{0xff, 0x00, 0x00, 0x00})
	//fmt.Println(string(data))
	a := app.New()
	window := a.NewWindow(myapp.AppName)

	searchBy := widget.NewSelectWithData([]string{"Title", "Author", "Genre", "Ratting", "On Loan"}, binding.NewString())
	header  := NewBooksHeader()

	search := container.New(layout.NewGridLayout(2), searchBy, widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func(){return}))

	top := container.New(layout.NewGridLayout(1), search, header)

	entry_1 := NewBookEntry()
	entry_2 := NewBookEntry()
	entries := container.New(layout.NewVBoxLayout(), entry_1, entry_2, NewBookEntry(), NewBookEntry(), NewBookEntry())
	scroll := container.NewVScroll(entries)

	maincon := container.New(layout.NewBorderLayout(top, nil, nil, nil), top, scroll)

	window.SetContent(maincon)
	window.ShowAndRun()
	
}
