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
	//"fyne.io/fyne/v2/canvas"

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


func ChangeSVGIconFG(svg []byte, c color.RGBA) {
	red := fmt.Sprintf("%02x", c.R)
	green := fmt.Sprintf("%02x", c.G)
	blue := fmt.Sprintf("%02x", c.B)
	valueStart := -1
	for i, c := range svg {
		if c == '"' && valueStart == -1 {
			valueStart = i+2
		} else if c == '"' && valueStart > -1 {
			value := svg[valueStart:i]
			if string(value) == "f3f3f3" {
				r := value[0:2]
				g := value[2:4]
				b := value[4:6]
				r[0] = red[0]
				r[1] = red[1]
				g[0] = green[0]
				g[1] = green[1]
				b[0] = blue[0]
				b[1] = blue[1]
			}
			valueStart = -1
		}
	}
}

func NewDeleteIcon() fyne.Resource {
	icon := theme.DeleteIcon()
	data := icon.Content()
	ChangeSVGIconFG(data, color.RGBA{0xff, 0x00, 0x00, 0x00})
	myIcon := &CustomIcon{
		MyName: "MyDeleteIcon",
		MyContent: data,
	}
	return myIcon
}
func NewEditIcon() fyne.Resource {
	data := theme.DocumentCreateIcon().Content()
	ChangeSVGIconFG(data, color.RGBA{0x00, 0x00, 0xff, 0x00})
	myIcon := &CustomIcon{
		MyName: "MyEditIcon",
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

func NewBookEntry() *fyne.Container {
	fields := []fyne.CanvasObject{
		widget.NewEntry(),
		widget.NewEntry(),
		widget.NewEntry(),
		widget.NewSelectWithData([]string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}, binding.NewString()),
		widget.NewButtonWithIcon("", theme.CheckButtonIcon(), func(){return}),
		container.New(layout.NewGridLayout(2),
			widget.NewButtonWithIcon("", NewEditIcon(), func(){return}),
			widget.NewButtonWithIcon("", NewDeleteIcon(), func(){return}),
		)}
	fields[0].(*widget.Entry).SetPlaceHolder("")
	fields[1].(*widget.Entry).SetPlaceHolder("")
	fields[2].(*widget.Entry).SetPlaceHolder("")
	fields[3].(*widget.Select).PlaceHolder = "TBR"
	//fields[4].(*widget.Entry).SetPlaceHolder("Lounded")

	return container.New(layout.NewGridLayout(len(fields)), fields...)
}

func Run() {
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
