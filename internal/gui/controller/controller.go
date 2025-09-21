package controller

import (

	"github.com/dubbersthehoser/mayble/internal/core"
	"github.com/dubbersthehoser/mayble/internal/storage"
)

type Controller struct {
	Core        *core.Core
	BookList    *BookList
	BookEditor  *BookEditor
	Selector    *Selector
}

func NewContorller(core *core.Core) *Controller {
	var c Contorller
	c.Core = core
	c.BookList = NewBookList(&c)
	c.BookEditor = NewBookEditor(&c)
	c.Selector = NewSelector(&c)
	return &c
}

