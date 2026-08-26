package view

import (
	"fmt"
)

type CommandBus struct {
	handlers map[string]func(any) error
}

func NewCommandBus() *CommandBus {
	cb := &CommandBus{
		handlers: make(map[string]func(any) error),
	}
	return cb
}

func (cb *CommandBus) Regester(name string, h func(any) error) {
	cb.handlers[name] = h
}

func (cb *CommandBus) Dispatch(command any) error {
	var name string
	switch command.(type) {
	case OpenDialog:
		name = "OpenDialog"
	case ToggleHiddenColumn:
		name = "ToggleHiddenColumn"
	default:
		return fmt.Errorf("commandbus: dispatch: unregister command type")
	}

	h, ok := cb.handlers[name]
	if !ok {
		return fmt.Errorf("commandbus: dispatch %s: command not found", name)
	}
	h(command)
	return nil
}

// Commands

type OpenDialog struct {
	Name string
}

type ToggleHiddenColumn struct {
	Column string
}
