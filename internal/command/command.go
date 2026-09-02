package command

import (
	"fmt"
)

type Command any

type Handler func(Command) error

type CommandBus struct {
	handlers map[string]Handler
}

func NewCommandBus() *CommandBus {
	cb := &CommandBus{
		handlers: make(map[string]Handler),
	}
	return cb
}

// Register command to its handler.
func (cb *CommandBus) Register(c Command, fn Handler) {
	cb.handlers[fmt.Sprintf("%T", c)] = fn
}

// Dispatch registered command.
func (cb *CommandBus) Dispatch(c Command) error {
	handler, ok := cb.handlers[fmt.Sprintf("%T", c)]
	if !ok {
		return fmt.Errorf("dispatch %T: command not found", c)
	}
	return handler(c)
}
