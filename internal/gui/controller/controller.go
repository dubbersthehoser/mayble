package controller

import (
	"github.com/dubbersthehoser/mayble/internal/core"
)

type Controller struct {
	Core        *core.Core
	BookList    *BookList
	BookEditor  *BookEditor
}

func New(core *core.Core) *Controller {
	var c Controller
	c.Core = core
	c.BookList = NewBookList(core)
	c.BookEditor = NewBookEditor(core)
	return &c
}

