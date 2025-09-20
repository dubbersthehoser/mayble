package view

import (
	"fyne.io/fyne/v2"
	
	"github.com/dubbersthehoser/mayble/internal/gui/controller"
)

type FunkView struct {
	controller *controller.Master
	View       fyne.CanvasObject
}

func NewFunkView(control *controller.Master) FunkView {
	f := FunkView{
		controller: control,
	}
	obj := f.BookEdit(f.NewCreateForm())
	f.View = obj
	return f
}

func (f *FunkView) Update() {
	f.View.Refresh()
}

func (f *FunkView) NewCreateForm() controller.BookForm {
	return f.controller.BookEditor.NewCreateForm()
}
func (f *FunkView) SubmitForm(form controller.BookForm) {
	f.controller.BookEditor.Submit(&form)
}
