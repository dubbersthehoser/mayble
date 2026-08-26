package view

import (
	"log"
)

type EventBus struct {
	handlers map[string][]func(any)
}

func NewEventBus() *EventBus {
	eb := &EventBus{
		handlers: make(map[string][]func(any)),
	}
	return eb
}

func (eb *EventBus) Subscribe(name string, h func(any)) {
	handlers, ok := eb.handlers[name]
	if !ok {
		handlers = make([]func(any), 0)
	}
	handlers = append(handlers, h)
	eb.handlers[name] = handlers
}

func (eb *EventBus) Notify(ev any) {

	var name string
	switch ev.(type) {
	case BodyChanged:
		name = "BodyChanged"
	case ColumnHiddenChanged:
		name = "ColumnHiddenChanged"
	default:
		log.Println("eventbus: notify: type not found")
		return
	}

	handlers, ok := eb.handlers[name]
	if !ok {
		log.Printf("eventbus: notify %s: name not found", name)
		return
	}
	for _, h := range handlers {
		h(ev)
	}
}

type BodyChanged struct {
	Body int
}

type ColumnHiddenChanged struct {
	ID      bool
	LoanSet bool
	ReadSet bool
}
