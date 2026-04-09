package command

import (
	repo "github.com/dubbersthehoser/mayble/internal/repository"
)

type Command interface{
	Execute() error
	Undo() error
}

type 
