package gui

import (
	//"log"
	//"image/color"
	//"log"
	_"fmt"
	//"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	_"fyne.io/fyne/v2/data/binding"
	_"fyne.io/fyne/v2/theme"
	_"fyne.io/fyne/v2/driver/desktop"
	_"fyne.io/fyne/v2/dialog"
	_"fyne.io/fyne/v2/canvas"

	"github.com/dubbersthehoser/mayble/internal/controler"
	"github.com/dubbersthehoser/mayble/internal/event"
)

type UI struct {
	Window    fyne.Window
	Emiter    *event.EventEmiter
	VM        *controler.VM
}

func NewUI(window fyne.Window, vm *controler.VM) *UI {
	u := &UI{
		Emiter: event.NewEventEmiter(),
		Window: window,
		VM: vm,
	}
	return u
}

func Run() {
	a := app.New()
	vm := &controler.VM{}
	window := a.NewWindow("alpha")
	window.Resize(fyne.NewSize(800, 500))
	UI := NewUI(window, vm)

	//ctrlN := &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
	//ctrlS := &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	//ctrlD := &desktop.CustomShortcut{KeyName: fyne.KeyD, Modifier: fyne.KeyModifierControl}
	//ctrlE := &desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: fyne.KeyModifierControl}
	//window.Canvas().AddShortcut(ctrlN, func (shortcut fyne.Shortcut){
	//	fmt.Println("Control+N: Create New book")
	//})

	header := UI.NewHeaderComp()
	body := UI.NewBookTableComp()
	mainComp := container.New(layout.NewBorderLayout(header, nil, nil, nil), header, body)

	window.SetContent(mainComp)

	window.ShowAndRun()
}
