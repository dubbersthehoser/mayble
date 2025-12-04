package controller

import (
	"github.com/dubbersthehoser/mayble/internal/app"
)

type Controller struct {
	App        app.API
	BookLoaner app.BookLoaning
	Redoer     app.Redoable
	Undoer     app.Undoable
	Importer   app.Importable
	Saver      app.Savable

	List      *BookList
	Editor    *BookEditor
}

func New(a app.API) *Controller {
	var c Controller
	c.SetApp(a)
	return &c
}

func (c *Controller) SetApp(a app.API) {
	c.App = a

	c.BookLoaner = a
	c.Redoer     = a
	c.Undoer     = a
	c.Importer   = a
	c.Saver      = a

	c.List = NewBookList(app.BookLoaning(a))
	c.Editor = NewBookEditor(app.BookLoaning(a))
}
