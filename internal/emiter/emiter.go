package emiter

import (
	"errors"
)


type Emiter struct {
	event map[string][]func(any)
}

func NewEmiter() *Emiter {
	return &Emiter{
		event: make(map[string][]func(any), 0),
	}
}

func (e *Emiter) OnEvent(key string, handler func(any)) {
	_, ok := e.event[key]
	if !ok {
		e.event[key] = make([]func(any), 0)
	}
	e.event[key] = append(e.event[key], handler)
}

func (e *Emiter) Emit(key string, data any) error {
	handlers, ok := e.event[key]
	if !ok {
		return errors.New("emiter: key not found")
	}
	for _, handle := range handlers {
		handle(data)
	}
	return nil
}

