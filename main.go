package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/dubbersthehoser/mayble/internal/view"
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)


func main() {
	a := app.New()
	window := a.NewWindow("New Mayble")
	window.Resize(fyne.NewSize(900, 600))
	window.CenterOnScreen()
	uiVM := viewmodel.NewMainUI()
	content := view.NewMainUI(window, uiVM)
	window.SetContent(content)
	window.ShowAndRun()
}
