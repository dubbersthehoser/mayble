package controller

import (
	"github.com/dubbersthehoser/mayble/internal/app"
	"github.com/dubbersthehoser/mayble/internal/emiter"
)

type Controller struct {
	App        app.API
	BookLoaner app.BookLoaning
	Redoer     app.Redoable
	Undoer     app.Undoable
	Importer   app.Importable
	Saver      app.Savable

	List      *BookLoanList
	Searcher  *BookLoanSearcher
	Editer    *BookEditer

	Broker    *emiter.Broker
}

func New(a app.API) *Controller {
	var c Controller
	c.Broker = &emiter.Broker{}
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

	c.List = NewBookLoanList(c.Broker, app.BookLoaning(a))
	c.Searcher = NewBookLoanSearcher(c.Broker, &c.List.list)
	c.Editer = NewBookEditer(c.Broker, app.BookLoaning(a))
}
