package main

import (

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/gui/view"
	//"github.com/dubbersthehoser/mayble/internal/gui/viewmodel"
)


func main() {
	a := app.New()
	window := a.NewWindow("New Mayble")
	window.Resize(fyne.NewSize(900, 600))
	window.CenterOnScreen()


	content := view.NewMainUI()

	window.SetContent(content)
	window.ShowAndRun()
}
