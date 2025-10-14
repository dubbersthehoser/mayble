package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
)

const (
	OnSave string = "ON_SAVE"
	OnSort        = "ON_SORT"
	OnSelected    = "ON_SELECTED"
	OnUnselected    = "ON_UNSELECTED"
)

type emiter struct {
	items map[string][]func()
}

func (e *emiter) init() {
	if e.items == nil {
		e.items = make(map[string][]func())
	}
}

func (e *emiter) On(key string, do func()) {
	e.init()
	list, ok := e.items[key]
	if !ok {
		list = make([]func(), 0)
	}
	list = append(list, do)
	e.items[key] = list
}
func (e *emiter) Emit(key string) {
	e.init()
	list, ok := e.items[key]
	if !ok {
		return
	}
	for _, fn := range list {
		fn()
	}
}

type FunkView struct {
	controller *controller.Controller
	View       fyne.CanvasObject
	emiter     emiter
}

func NewFunkView(control *controller.Controller) (FunkView, error) {
	f := FunkView{
		controller: control,
		emiter: emiter{},
	}
	obj, err := f.BookEdit()
	if err != nil {
		return f, err
	}
	_ = obj

	f.View = f.Body()
	return f, nil
}

func (f *FunkView) Body() fyne.CanvasObject {
	topBar := f.TopBar()
	table := f.Table()
	body := container.New(layout.NewBorderLayout(topBar, nil, nil, nil), topBar, table)
	return body
}


func (f *FunkView) Update() {
	f.View.Refresh()
}

func (f *FunkView) DisplayError(err error) {
}
