package controller

import (

	"github.com/dubbersthehoser/mayble/internal/core"
)

type Controller struct {
	Core        *core.Core
	BookList    *BookList
	BookEditor  *BookEditor
}

func NewContorller(core *core.Core) *Master {
	var c Contorller
	c.Core = core
	c.BookList = NewBookList(&c)
	c.BookEditor = NewBookEditor(&c)
	return &c
}
