package view

import (
	"github.com/dubbersthehoser/mayble/internal/viewmodel"
)

type Command interface {
	Call()
}

type Invoker struct {
	vm *viewmodel.Window
}

func Run(cmd Command)

