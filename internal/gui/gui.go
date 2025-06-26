package gui

import (
	//"log"
	//"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/binding"
	//"fyne.io/fyne/v2/canvas"

	myapp "github.com/dubbersthehoser/mayble/internal/app"
)


func NewBookEntry() *fyne.Container {
	fields := []fyne.CanvasObject{
		widget.NewEntry(),
		widget.NewEntry(),
		widget.NewEntry(),
		widget.NewSelectWithData([]string{"TBR", "⭐", "⭐⭐", "⭐⭐⭐", "⭐⭐⭐⭐", "⭐⭐⭐⭐⭐"}, binding.NewString()),
		widget.NewEntry(),
	}
	fields[0].(*widget.Entry).SetPlaceHolder("Title")
	fields[1].(*widget.Entry).SetPlaceHolder("Author")
	fields[2].(*widget.Entry).SetPlaceHolder("Genre")
	fields[3].(*widget.Select).PlaceHolder = "Ratting"
	fields[4].(*widget.Entry).SetPlaceHolder("Lounded")

	return container.New(layout.NewGridLayout(len(fields)), fields...)

}


func Run() {
	a := app.New()
	window := a.NewWindow(myapp.AppName)

	entry_1 := NewBookEntry()
	entry_2 := NewBookEntry()
	content := container.New(layout.NewVBoxLayout(), entry_1, entry_2)

	window.SetContent(content)
	window.ShowAndRun()
	
}
