package table

import (
	"fmt"
	"sync"
)

type eventBus struct {

	// I'm worried that not having a mutex may bite me.

	mu sync.Mutex
	handlers map[string][]func(any)
}

func newEventBus() *eventBus {
	eb := &eventBus{
		handlers: make(map[string][]func(any)),
	}
	return eb
}

func (eb *eventBus) Subscribe(name string, h func(any)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers, ok := eb.handlers[name]
	if !ok {
		handlers = make([]func(any), 0)
	}

	handlers = append(handlers, h)
	eb.handlers[name] = handlers
}

func (eb *eventBus) Notify(event any) error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	var name string
	switch event.(type) {
	case eventSelected:
		name = "eventSelected"
	case eventHiddenColumn:
		name = "eventHiddenColumn"
	case eventSearchBy:
		name = "eventSearchBy"
	default:
		return fmt.Errorf("event_bus: event type not found")
	}

	handlers, ok := eb.handlers[name]
	if !ok {
		return fmt.Errorf("event_bus %s: name not found", name)
	}
	for _, h := range handlers {
		h(event)
	}
	return nil
}

// Event Types

type eventSelected struct {
	has     bool
	point   Point
	id      int64
	version int64
}

type eventSort struct {
	column  string
	asc     bool
	version int64
}

type eventHiddenColumn struct {
	label  string
	hidden bool
}

type eventSearchBy struct {
	column int
}
