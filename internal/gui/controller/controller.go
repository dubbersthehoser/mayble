package controller

import (
	"github.com/dubbersthehoser/mayble/internal/app"
)

type Controller struct {
	App        app.Mayble
	BookLoaner app.BookLoaning
	Redoer     app.Redoable
	Undoer     app.Undoable
	Importer   app.Importable
	Saver      app.Savable
	List      *BookList
	Editor *BookEditor
}

func New(a app.Mayble) *Controller {
	var c Controller
	c.App = a

	c.BookLoaner = a
	c.Redoer     = a
	c.Undoer     = a
	c.Importer   = a
	c.Saver      = a

	c.List = NewBookList(app.BookLoaning(a))
	c.Editor = NewBookEditor(app.BookLoaning(a))
	return &c
}

