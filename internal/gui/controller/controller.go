package controller

import (
	"github.com/dubbersthehoser/mayble/internal/app"
)

type Controller struct {
	App         app.Mayble
	BookList    *BookList
	BookEditor  *BookEditor
}

func New(a app.Mayble) *Controller {
	var c Controller
	c.App = a
	c.BookList = NewBookList(app.BookLoaning(a))
	c.BookEditor = NewBookEditor(app.BookLoaning(a))
	return &c
}

